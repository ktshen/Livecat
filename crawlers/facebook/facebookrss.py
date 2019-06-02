import threading
import datetime
import requests
from requests.exceptions import RequestException
import time
import json
import os
from json.decoder import JSONDecodeError
from bs4 import BeautifulSoup

from selenium import webdriver
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, NoSuchElementException
from selenium.webdriver.chrome.options import Options
import xml.etree.ElementTree as ET
from urllib3.exceptions import InsecureRequestWarning




requests.packages.urllib3.disable_warnings(category=InsecureRequestWarning)

DRIVER_PATH = "../chromedriver"
GMAIL_EMAIL = "dlivecat@gmail.com"
GMAIL_PASS = "bTerr#21p"
ELASTIC_CREATE_SERVER_URL = "http://120.126.16.88:17777/add"


def logfunc(*string):
    print(datetime.datetime.now(), end='  ')
    for i in string:
        print(i, end=' ')
    print()


class NotLoginStatusException(Exception):
    def __init__(self,*args,**kwargs):
        Exception.__init__(self,*args,**kwargs)


class XMLCrawler(threading.Thread):
    def __init__(self, url, cookies):
        super().__init__()
        self.url = url
        self.daemon = True
        self.cookies = cookies

    def store(self, data):
        try:
            resp = requests.post(ELASTIC_CREATE_SERVER_URL, data=data)
            logfunc("POST", "status:"+data["status"], data["title"], resp)
        except Exception as e:
            print(e)

    def validate(self, it):
        _title = it.find("title")
        _thumbnails = it.find("enclosure")
        if _title is None or _thumbnails is None or not "Video" in _title.text:
            return False
        # Empty Title
        if not self.get_title(_title.text):
            return False
        link = it.find("link").text
        if "youtu" in link or not "facebook" in link:
            return False
        published = datetime.datetime.strptime(it.find("pubDate").text, "%a, %d %b %Y %H:%M:%S +0000")
        if published < datetime.datetime.now() - datetime.timedelta(days=7):
            return False
        return True

    def get_title(self, title):
        if "Video - " == title[:8]:
            title = title[8:]
        title = title.split("<br>")[0]
        return title

    def get_host(self):
        _t = self.root.find("channel").find("title").text
        host = _t.split("FB-RSS feed for ")[-1]
        return host

    def test_video_live_status(self, url):
        try:
            r = requests.get(url)
        except Exception as e:
            logfunc(url, e)
            return 'live'

        if r.status_code == 404:
            return 'invalid'

        if "正在直播" in r.text:
            return 'live'
        else:
            return 'video'

    def parse(self):
        for it in self.root.iter("item"):
            if not self.validate(it):
                continue
            video_status = self.test_video_live_status(it.find("link").text)
            if not video_status == 'live':
                continue

            data = {"title": self.get_title(it.find("title").text),
                    "host": self.get_host(),
                    "videourl": it.find("link").text,
                    "description": it.find("description").text.rsplit('<br>')[-1],
                    "published": datetime.datetime.strptime(it.find("pubDate").text, "%a, %d %b %Y %H:%M:%S +0000") \
                        .strftime("%Y-%m-%dT%H:%M:%SZ"),
                    "platform": "Facebook",
                    "thumbnails": it.find("enclosure").attrib["url"],
                    "timestamp": datetime.datetime.now().strftime("%Y-%m-%dT%H:%M:%SZ"),
                    "status": video_status,
                    }
            self.store(data)

    def run(self):
        success = False
        counter = 0
        while not success and counter < 5:
            try:
                resp = requests.get(self.url, cookies=self.cookies, verify=False)
                self.root = ET.fromstring(resp.content.decode('utf-8'))
                self.parse()
                success = True
            except (RequestException, ET.ParseError) as e:
                logfunc(self.url, e)
            except Exception as e:
                logfunc(e)
            time.sleep(1)
            counter += 1


class FacebookRssFetcher:
    def __init__(self):
        chrome_options = Options()
        chrome_options.add_argument('--headless')
        chrome_options.add_argument("--no-sandbox")
        chrome_options.add_argument("--disable-dev-shm-usage")
        self.driver = webdriver.Chrome(DRIVER_PATH, options=chrome_options)
        self.targets = set()
        self.workers = {}
        self.user_cookie = {}
        self.user_cookie_filename = "user_cookie.txt"

    def site_login(self):
        logfunc("Try login...")
        success = False
        while not success:
            self.driver.delete_all_cookies()
            self.driver.get("https://fbrss.com/login")
            try:
                self.driver.find_element_by_id("email").send_keys(GMAIL_EMAIL)
                self.driver.find_element_by_id("pass").send_keys(GMAIL_PASS)
                self.driver.find_element_by_id("loginbutton").click()
            except NoSuchElementException as e:
                logfunc(e)
                logfunc("Can't log in, wait for an hour to try again")
                time.sleep(60*60)
                continue

            try:
                WebDriverWait(self.driver, 5).until(
                    EC.title_is("Facebook to RSS - Export your facebook profile and page as RSS and Atom feed")
                )
                if "home" in self.driver.current_url:
                    logfunc("Log in successfully")
                    success = True
            except TimeoutException as e:
                if "facebook.com" in self.driver.current_url:
                    logfunc("Wrong Username or Password!")
                    raise ValueError
                logfunc(e)
                logfunc("Can't log in, wait for an hour to try again")
                time.sleep(60 * 60)
                continue

        self.write_cookies_to_file()
        self.update_user_cookie()


    def update_user_cookie(self):
        for c in self.driver.get_cookies():
            if "user" in c.values():
                self.user_cookie = {"user": c["value"]}
                break

    def write_cookies_to_file(self):
        logfunc("Writing cookies to {}".format(self.user_cookie_filename))
        with open(self.user_cookie_filename, "w") as fp:
            json.dump(self.driver.get_cookies(), fp)

    def load_cookies_from_file(self):
        logfunc("Loading cookies from {}".format(self.user_cookie_filename))
        if not "fbrss.com" in self.driver.current_url:
            self.driver.get("https://fbrss.com")
        try:
            with open(self.user_cookie_filename, "r") as fp:
                cookies = json.load(fp)
        except FileNotFoundError:
            logfunc("Cannot find file: {}".format(self.user_cookie_filename))
            return
        except JSONDecodeError:
            os.remove(self.user_cookie_filename)
            return
        if type(cookies) is dict:
            self.driver.add_cookie(cookies)
        elif type(cookies) is list:
            for c in cookies:
                self.driver.add_cookie(c)
        self.update_user_cookie()


    def close_driver(self):
        print("Closing driver...")
        self.driver.quit()


    def test_is_login_status(self):
        self.driver.get("https://fbrss.com/home")
        try:
            self.driver.find_element_by_partial_link_text('Sign in')
        except NoSuchElementException:
            if "home" in self.driver.current_url:
                return True
        logfunc("Not log in yet.")
        return False


    def build_targets(self):
        self.targets = set()
        soup = BeautifulSoup(self.driver.page_source, 'html.parser')
        feeds = soup.find(id="the_feeds")
        for tr in feeds.find_all('tr'):
            tds = tr.find_all('td')
            if not len(tds):
                continue
            feed_url = tds[1].find('a').get('href')
            self.targets.add(feed_url)


    def process_pages(self):
        while True:
            if not self.test_is_login_status():
                raise NotLoginStatusException
            self.build_targets()
            for url in self.targets.copy():
                c = XMLCrawler(url, self.user_cookie)
                c.start()
                self.workers[url] = c
                time.sleep(1)

            while len(self.workers):
                for url in list(self.workers.keys()):
                    if not self.workers[url].is_alive():
                        self.workers.pop(url)

    def run(self):
        self.load_cookies_from_file()
        if not self.test_is_login_status():
            self.site_login()
        logfunc("Start Crawling!")
        try:
            while True:
                try:
                    self.process_pages()
                except (NotLoginStatusException, AttributeError):
                    try:
                        self.site_login()
                    except Exception as e:
                        raise e
        except KeyboardInterrupt:
            logfunc("Forced Stop.")
        except ValueError:
            pass
        except Exception as e:
            logfunc(e)
        finally:
            self.close_driver()

f = FacebookRssFetcher()
f.run()
