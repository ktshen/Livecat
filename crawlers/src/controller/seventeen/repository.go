package seventeen

import (
	"crawlers/controller/resource"

	"github.com/tebeka/selenium"
)

// Repository represent the seventeen's repository contract
type Repository interface {
	Init(config resource.Config)
	ParseStreamURLHTML(wd selenium.WebDriver, videoURLs []string) (datas []resource.Data, err error)
}
