from selenium import webdriver
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, NoSuchElementException
from selenium.webdriver.chrome.options import Options
import xml.etree.ElementTree as ET
from urllib3.exceptions import InsecureRequestWarning
import threading
import datetime
import requests
from requests.exceptions import RequestException
from bs4 import BeautifulSoup
import time

requests.packages.urllib3.disable_warnings(category=InsecureRequestWarning)

DRIVER_PATH = "./chromedriver"
GMAIL_EMAIL = "dlivecat@gmail.com"
GMAIL_PASS = "sTerr#21p"
PUPPETEER_URL = "http://140.115.153.208:8080/material/png_Crawler"
ELASTIC_CREATE_SERVER_URL = "http://120.126.16.88:17777/add"


class XMLCrawler(threading.Thread):
    def __init__(self, url, cookies):
        super().__init__()
        self.url = url
        self.daemon = True
        self.cookies = cookies

    def run(self):
        try:
            resp = requests.get(self.url, cookies=self.cookies, verify=False)
            self.root = ET.fromstring(resp.content)
            self.parse()
        except RequestException as e:
            print(e)

    def parse(self):
        for it in self.root.iter("item"):
            if not self.validate(it):
                continue
            data = {"title": self.get_title(it.find("title").text),
                    "host": self.get_host(),
                    "videourl": it.find("link").text,
                    "description": it.find("description").text.rsplit('<br>')[-1],
                    "published": datetime.datetime.strptime(it.find("pubDate").text, "%a, %d %b %Y %H:%M:%S +0000")\
                                                  .strftime("%Y-%m-%dT%H:%M:%SZ"),
                    "platform": "Facebook",
                    "thumbnails": it.find("enclosure").attrib["url"],
                    "timestamp": datetime.datetime.now().strftime("%Y-%m-%dT%H:%M:%SZ"),
                    }
            self.store(data)

    def store(self, data):
        # collection = client["Crawler"]["Livestreams"]
        # collection.find_one_and_update(
        #     {"videourl": data["videourl"]},
        #     {"$set":data},
        #     upsert=True,
        # )

        try:
            resp = requests.post(ELASTIC_CREATE_SERVER_URL, data=data)
            print(datetime.datetime.now(), " POST: ", data["title"], " ", resp)
        except Exception as e:
            print(e)

    
    def validate(self, it):
        _title = it.find("title")
        _thumbnails = it.find("enclosure")
        if _title is None or _thumbnails is None or not "Video" in _title.text :
            return False
        link = it.find("link").text
        if "youtu" in link or not "facebook" in link:
            return False
        published =  datetime.datetime.strptime(it.find("pubDate").text, "%a, %d %b %Y %H:%M:%S +0000")
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


class FacebookRssFetcher:
    def __init__(self):
        chrome_options = Options()
        chrome_options.add_argument('--headless')
        chrome_options.add_argument("--no-sandbox")
        chrome_options.add_argument("--disable-dev-shm-usage")
        self.driver = webdriver.Chrome(DRIVER_PATH, chrome_options=chrome_options)
        self.targets = set()
        self.workers = {}
        self.cookies = {}

    def site_login(self):
        print(datetime.datetime.now(), " Try login...")

        self.driver.delete_all_cookies()
        self.driver.get("https://fbrss.com/login")
        self.driver.find_element_by_id("email").send_keys(GMAIL_EMAIL)
        self.driver.find_element_by_id("pass").send_keys(GMAIL_PASS)
        self.driver.find_element_by_id("loginbutton").click()
        WebDriverWait(self.driver, 5).until(
            EC.title_is("Facebook to RSS - Export your facebook profile and page as RSS and Atom feed")
        )
        for c in self.driver.get_cookies():
            if "user" in c.values():
                self.cookies = {"user": c["value"]}
                break


    def close_driver(self):
        print("Closing driver...")
        self.driver.quit()

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

    def refresh_targets(self):
        self.driver.refresh()
        self.build_targets()

    def process_pages(self):
        # Refresh targets every day
        counter = 0
        while counter < 2880:
            for url in self.targets.copy():
                c = XMLCrawler(url, self.cookies)
                c.start()
                self.workers[url] = c
                self.targets.remove(url)

            for url in list(self.workers):
                if not self.workers[url].is_alive():
                    self.workers.pop(url)
                    self.targets.add(url)
            time.sleep(30)
            counter += 1

    def run(self):
        success = False
        while not success:
            try:
                self.site_login()
                self.build_targets()
                success = True
                print(datetime.datetime.now(), " Login Success!")
            except (AttributeError, NoSuchElementException):
                print(datetime.datetime.now(), " Can't login correctly, because sign in too often. Wait for an hour and try again.")
                time.sleep(60 * 60)

        try:
            while True:
                try:
                    self.process_pages()
                    self.refresh_targets()

                except (TimeoutException, AttributeError):
                    print(datetime.datetime.now(), " Can't login correctly, because sign in too often. Wait for thirty minutes and try again.")
                    time.sleep(60 * 30)
                    self.site_login()

        except KeyboardInterrupt:
            print("Forced Stop.")

        finally:
            self.close_driver()


f = FacebookRssFetcher()
f.run()


