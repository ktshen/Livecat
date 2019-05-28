package crawler

import (
	"fmt"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	portPiko       = 9000
	portSeventeen  = 9001
	portWatermelon = 9002
	portYouTube    = 9003

	chromedriverPath = "/home/user/chromedriver"
)

func getService(port int) *selenium.Service {
	opts := []selenium.ServiceOption{
		// Enable fake XWindow session.
		// selenium.StartFrameBuffer(),
		//selenium.Output(os.Stderr), // Output debug information to STDERR
	}

	// Enable debug info.
	// selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService("chromedriver", port, opts...)
	if err != nil {
		panic(err)
	}
	return service
}

func getRemote(port int) selenium.WebDriver {
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
			"--disable-gpu-sandbox",
			"--disable-accelerated-2d-canvas",
			"--test-type=ui",
		},
	}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	return wd
}
