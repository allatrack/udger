package parser

import (
	"database/sql"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"fmt"
	"strings"
	"log"
)

type UdgerSlowLoadFastExec struct {
	UserAgent UserAgent
	IpAddress IpAddress

	BrowserList     map[int64]Browser
	OsList          map[int64]Os
	DeviceList      map[string]Device
	DeviceClassList map[int64]Device
	DeviceCodeList  map[int64]DeviceCode
	IpList          map[string]IpAddress
	CrawlerList     map[string]Crawler

	rexBrowsers []rexData
	rexOS       []rexData
	rexDevices  []rexData

	db     *sql.DB
	DbPath string
}

type UdgerFastLoadSlowExec struct {
	dbPath string
}

type UserAgent struct {
	String                 string
	Name                   string
	Class                  string
	ClassCode              string
	Version                string
	VersionMajor           string
	UpToDateCurrentVersion string
	FamilyInfo             Info
	FamilyVendorInfo       Info
	Engine                 string
	OsInfo                 Info
	OsFamilyInfo           Info
	OsFamilyVendorInfo     Info
	DeviceInfo             Info
	DeviceMarketName       string
	DeviceBrandInfo        Info
	Crawler                Crawler
}

func (userAgent *UserAgent) setVersion(fullVersion interface{}, name string) {

	var version string
	switch v := fullVersion.(type) {
	default:
		log.Printf("Error while setting User Agent version: unexpected version type %T", v)
	case float64:
		version = fmt.Sprintf("%.2f", fullVersion)
	case string:
		version = fullVersion.(string)
	case sql.NullString:
		version = (fullVersion.(sql.NullString)).String
	}

	if version != "" {
		userAgent.Name = name + " " + version
		userAgent.Version = version
		userAgent.VersionMajor = string(strings.Split(string(version), ".")[0])
	} else {
		userAgent.Name = name
		userAgent.Version = ""
		userAgent.VersionMajor = ""
	}
}

type IpAddress struct {
	Ip                 sql.NullString
	Version            sql.NullString
	Classification     sql.NullString
	ClassificationCode sql.NullString
	Hostname           sql.NullString
	LastSeen           sql.NullString
	Country            sql.NullString
	CountryCode        sql.NullString
	City               sql.NullString
	DataCenterInfo     Info
	Crawler            Crawler
}

type Crawler struct {
	Id                        sql.NullInt64
	Name                      sql.NullString
	Version                   sql.NullString
	VersionMajor              sql.NullString
	FamilyInfo                Info
	FamilyVendorInfo          Info
	LastSeen                  sql.NullString
	CrawlerClassification     sql.NullString
	CrawlerClassificationCode sql.NullString
	RespectRobotstxt          sql.NullString
	UaString                  string
}

type Browser struct {
	ClassId                  int64
	Id                       int64
	BrowserInfo              Info
	VendorInfo               Info
	Engine                   string
	UptodateCurrentVersion   string
	ClientClassification     string
	ClientClassificationCode string
}

type Os struct {
	Id         int64
	Info       Info
	FamilyInfo Info
	VendorInfo Info
}

type Device struct {
	Id   int64
	Info Info
}

type DeviceCode struct {
	Id            int64
	OsFamilyCode  string
	OsCode        string
	RegexCompiled pcre.Regexp
}

type Info struct {
	Name     sql.NullString
	Code     sql.NullString
	Icon     sql.NullString
	IconBig  sql.NullString
	Url      sql.NullString
	Homepage sql.NullString
}

type rexData struct {
	ID            int64
	Regex         string
	RegexCompiled pcre.Regexp
}

func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

type Error struct {
	message string
}

func (e Error) Error() string {
	return e.message
}
