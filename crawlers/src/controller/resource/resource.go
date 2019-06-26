package resource

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

// Config ...
type Config struct {
	ElasticServerIP string
	DBPort          string
}

// Data ...
type Data struct {
	Title            string
	Description      string
	Platform         string
	VideoID          string
	Host             string
	Status           string
	Thumbnails       string
	Published        string
	Tags             string
	GeneralTag       string
	Timestamp        string
	Language         string
	ViewCount        int
	Viewers          int
	VideoURL         string
	VideoEmbedded    string
	ChatRoomEmbedded string
	Channel          string
}

// GetFreePort ...
func GetFreePort() (port int, err error) {
	ln, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		return -1, err
	}
	port = ln.Addr().(*net.TCPAddr).Port
	err = ln.Close()
	return
}

// GetService ...
func GetService(port int) (*selenium.Service, error) {
	opts := []selenium.ServiceOption{
		// Enable fake XWindow session.
		// selenium.StartFrameBuffer(),
		//selenium.Output(os.Stderr), // Output debug information to STDERR
	}

	// Enable debug info.
	// selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService("chromedriver", port, opts...)

	return service, err
}

// GetRemote ...
func GetRemote(port int) (selenium.WebDriver, error) {
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{

		Path: "",
		Args: []string{
			"--headless",
			"--start-maximized",
			"--window-size=1200x600",
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
			"--disable-gpu",
			"--disable-impl-side-painting",
			"--disable-extensions",
			"--disable-gpu-sandbox",
			"--disable-accelerated-2d-canvas",
			"--test-type=ui",
		},
	}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", port))

	return wd, err
}

// SendToElasticByHTTPPost ...
func SendToElasticByHTTPPost(data Data, config Config) {
	url := "http://" + config.ElasticServerIP + ":17777/add"

	dataByte, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataByte))
	HandleError(err, "[resource] NewRequest err")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	// handleError(err, "[http.go] send data err")
	if err != nil {
		log.Println("[resource] send data err ", err)
	} else {
		defer resp.Body.Close()
	}
}

// GetLocalTime ...
func GetLocalTime() string {
	t := time.Now().UTC().Format("2006-01-02T15:04:05-0700")
	return t
}

// HandleError ...
func HandleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Printf(message+": %v", err.Error())
	}
}
