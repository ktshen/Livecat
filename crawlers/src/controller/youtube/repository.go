package youtube

import (
	"crawlers/controller/resource"

	"github.com/tebeka/selenium"
)

// Repository represent the youtube's repository contract
type Repository interface {
	Init(config resource.Config, url string)
	ParseStreamURLHTML(wd selenium.WebDriver, videoURLs []string) (datas []resource.Data, err error)
}
