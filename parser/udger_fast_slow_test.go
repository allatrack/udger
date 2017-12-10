package parser

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"log"
)


func TestParseUaFS(t *testing.T) {
	log.Println("Start TestParseUaFS")
	udger, err := NewSF("../udger.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	assert := assert.New(t)
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36";
	data, _ := udger.ParseUa(userAgent)

	assert.Equal("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36", data.String)
	assert.Equal("Browser", data.Class)
	assert.Equal("browser", data.ClassCode)
	assert.Equal("Chrome 39.0.2171.71", data.Name)
	assert.Equal("39.0.2171.71", data.Version)
	assert.Equal("39", data.VersionMajor)
	assert.Equal("61", data.UpToDateCurrentVersion)

	assert.Equal("Chrome", data.FamilyInfo.Name.String)
	assert.Equal("chrome", data.FamilyInfo.Code.String)
	assert.Equal("http://www.google.com/chrome/", data.FamilyInfo.Homepage.String)
	assert.Equal("Google Inc.", data.FamilyVendorInfo.Name.String)
	assert.Equal("google_inc", data.FamilyVendorInfo.Code.String)
	assert.Equal("https://www.google.com/about/company/", data.FamilyVendorInfo.Homepage.String)
	assert.Equal("chrome.png", data.FamilyInfo.Icon.String)
	assert.Equal("chrome_big.png", data.FamilyInfo.IconBig.String)
	assert.Equal("https://udger.com/resources/ua-list/browser-detail?browser=Chrome", data.FamilyInfo.Url.String)

	assert.Equal("WebKit/Blink", data.Engine)

	assert.Equal("OS X 10.9 Mavericks", data.OsInfo.Name.String)
	assert.Equal("osx_10_9", data.OsInfo.Code.String)
	assert.Equal("https://en.wikipedia.org/wiki/OS_X_Mavericks", data.OsInfo.Homepage.String)
	assert.Equal("macosx.png", data.OsInfo.Icon.String)
	assert.Equal("macosx_big.png", data.OsInfo.IconBig.String)
	assert.Equal("https://udger.com/resources/ua-list/os-detail?os=OS X 10.9 Mavericks", data.OsInfo.Url.String)

	assert.Equal("OS X", data.OsFamilyInfo.Name.String)
	assert.Equal("osx", data.OsFamilyInfo.Code.String)
	assert.Equal("Apple Computer, Inc.", data.OsFamilyVendorInfo.Name.String)
	assert.Equal("apple_inc", data.OsFamilyVendorInfo.Code.String)
	assert.Equal("http://www.apple.com/", data.OsFamilyVendorInfo.Homepage.String)

	assert.Equal("Desktop", data.DeviceInfo.Name.String)
	assert.Equal("desktop", data.DeviceInfo.Code.String)
	assert.Equal("desktop.png", data.DeviceInfo.Icon.String)
	assert.Equal("desktop_big.png", data.DeviceInfo.IconBig.String)
	assert.Equal("https://udger.com/resources/ua-list/device-detail?device=Desktop", data.DeviceInfo.Url.String)

	assert.Equal("Mac", data.DeviceMarketName)
	assert.Equal("apple", data.DeviceBrandInfo.Code.String)
	assert.Equal("Apple", data.DeviceBrandInfo.Name.String)
	assert.Equal("http://www.apple.com/", data.DeviceBrandInfo.Url.String)
	assert.Equal("apple.png", data.DeviceBrandInfo.Icon.String)
	assert.Equal("apple_big.png", data.DeviceBrandInfo.IconBig.String)

	// test crawler
	data, _ = udger.ParseUa(`192.comAgent`)

	assert.Equal("Crawler", data.Class)
	assert.Equal("crawler", data.ClassCode)
}

func TestParseIpFS(t *testing.T) {
	log.Println("Start TestParseIpFS")
	udger, err := NewSF("../udger.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	assert := assert.New(t)

	ipAddress := udger.ParseIp(`1.0.128.191`)
	assert.Equal("1.0.128.191", ipAddress.Ip.String)
	assert.Equal("Web proxy", ipAddress.Classification.String)
	assert.Equal("web_proxy", ipAddress.ClassificationCode.String)
	assert.Equal("node-5b.pool-1-0.dynamic.totbb.net", ipAddress.Hostname.String)
	assert.Equal("v4", ipAddress.Version.String)
	assert.Equal("Thailand", ipAddress.Country.String)

	ipAddress = udger.ParseIp(`101.0.64.0`)
	assert.Equal("Digital Pacific", ipAddress.DataCenterInfo.Name.String)
	assert.Equal("digital_pacific", ipAddress.DataCenterInfo.Code.String)
	assert.Equal("http://www.digitalpacific.com.au/", ipAddress.DataCenterInfo.Homepage.String)

	ipAddress = udger.ParseIp(`2001:e42:101::`)
	assert.Equal("Sakura Internet Inc.", ipAddress.DataCenterInfo.Name.String)
	assert.Equal("sakura_internet", ipAddress.DataCenterInfo.Code.String)
	assert.Equal("http://www.sakura.ne.jp/", ipAddress.DataCenterInfo.Homepage.String)
}

