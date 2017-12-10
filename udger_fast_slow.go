package udger

import (
	"database/sql"
	"strings"
	"log"
	"fmt"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"net"
	"strconv"
	"os"
)

func NewFS(dbPath string) (*UdgerFastLoadSlowExec, error) {

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, err
	}
	uSlowFast := &UdgerFastLoadSlowExec{
		dbPath: dbPath,
	}

	return uSlowFast, nil
}

func (udger *UdgerFastLoadSlowExec) ParseUa(userAgentString string) (userAgent UserAgent, err error) {

	// defer timeTrack(time.Now(), "ParseUaFS")

	userAgent = UserAgent{}

	db, err := sql.Open("sqlite3", udger.dbPath)
	if err != nil {
		log.Panic(err)
		return userAgent, err
	}

	defer db.Close()

	if userAgentString != "" {

		userAgent.String = strings.TrimSpace(userAgentString)

		var (
			clientId         int64 = 0
			clientClassId    int64 = -1
			osId             int64 = 0
			deviceClassId    int64 = 0
			crawler          Crawler
			familyInfo       Info
			familyVendorInfo Info
		)

		rows, err := db.Query(fmt.Sprintf(CrawlersByUaSqlQuery, userAgent.String))
		if err != nil {
			panic(fmt.Sprintf("Error while CrawlersByUaSqlQuery: %s", err))
		}
		for rows.Next() {
			if err := rows.Scan(
				&crawler.Id,
				&crawler.Name,
				&crawler.Version,
				&crawler.VersionMajor,
				&crawler.LastSeen,
				&crawler.RespectRobotstxt,
				&familyInfo.Name,
				&familyInfo.Code,
				&familyInfo.Homepage,
				&familyInfo.Icon,
				&familyVendorInfo.Name,
				&familyVendorInfo.Code,
				&familyVendorInfo.Homepage,
				&crawler.CrawlerClassification,
				&crawler.CrawlerClassificationCode,
				&crawler.UaString,
			); err != nil {
				log.Print(Error{fmt.Sprintf("Error scanning CrawlersByUaSqlQuery: %s", err)})
			}
		}
		rows.Close()

		if crawler.Id.Int64 != 0 {
			clientClassId = CrawlerClassId
			userAgent.Class = CrawlerClass
			userAgent.ClassCode = CrawlerCode
			userAgent.Crawler = crawler
		}

		if crawler.Id.Int64 == 0 {
			rows, err := db.Query(BrowsersSqlQuery)
			if err != nil {
				panic(fmt.Sprintf("Error quering BrowsersSqlQuery: %s", err))
			}

			for rows.Next() {

				var browser Browser
				var regString string

				if err := rows.Scan(
					&browser.ClassId,
					&browser.Id,
					&regString,
					&browser.BrowserInfo.Name,
					&browser.BrowserInfo.Code,
					&browser.BrowserInfo.Homepage,
					&browser.BrowserInfo.Icon,
					&browser.BrowserInfo.IconBig,
					&browser.Engine,
					&browser.VendorInfo.Name,
					&browser.VendorInfo.Code,
					&browser.VendorInfo.Homepage,
					&browser.UptodateCurrentVersion,
					&browser.ClientClassification,
					&browser.ClientClassificationCode,
				); err != nil {
					panic(fmt.Sprintf("Error scunning BrowsersSqlQuery results: %s", err))
				}
				d, err := getCompiledReqObject(browser.Id, regString)
				if err != nil {
					panic(fmt.Sprintf("Error compiling req: %s %s", regString, err))
				}
				r := d.RegexCompiled
				matcher := r.MatcherString(userAgentString, 0)
				if !matcher.MatchString(userAgentString, 0) {
					continue
				}

				if matcher.Present(1) {
					userAgent.setVersion(matcher.GroupString(1), browser.BrowserInfo.Name.String)
					clientId = browser.Id
					clientClassId = browser.ClassId
					userAgent.Class = browser.ClientClassification
					userAgent.ClassCode = browser.ClientClassificationCode
					userAgent.UpToDateCurrentVersion = browser.UptodateCurrentVersion
					userAgent.FamilyInfo = browser.BrowserInfo
					userAgent.FamilyVendorInfo = browser.VendorInfo
					userAgent.FamilyInfo.Url = ToNullString(BrowserDetailUrl + browser.BrowserInfo.Name.String)
					userAgent.Engine = browser.Engine
					break
				}
			}
			rows.Close()
		}

		rows, err = db.Query(OsSqlQuery)
		if err != nil {
			panic(fmt.Sprintf("Error while OsSqlQuery: %s", err))
		}
		for rows.Next() {
			var operationSystem Os
			var regString string
			rows.Scan(
				&operationSystem.Id,
				&regString,
				&operationSystem.FamilyInfo.Name,
				&operationSystem.FamilyInfo.Code,
				&operationSystem.Info.Name,
				&operationSystem.Info.Code,
				&operationSystem.Info.Homepage,
				&operationSystem.Info.Icon,
				&operationSystem.Info.IconBig,
				&operationSystem.VendorInfo.Name,
				&operationSystem.VendorInfo.Code,
				&operationSystem.VendorInfo.Homepage,
			)
			d, err := getCompiledReqObject(operationSystem.Id, regString)
			if err != nil {
				panic(fmt.Sprintf("Error scanning OsSqlQuery: %s", err))
			}
			r := d.RegexCompiled
			matcher := r.MatcherString(userAgentString, 0)
			if !matcher.MatchString(userAgentString, 0) {
				continue
			}

			if matcher.Present(1) {
				osId = operationSystem.Id
				userAgent.OsInfo = operationSystem.Info
				userAgent.OsInfo.Url = ToNullString(OsDetailUrl + operationSystem.Info.Name.String)
				userAgent.OsFamilyInfo = operationSystem.FamilyInfo
				userAgent.OsFamilyVendorInfo = operationSystem.VendorInfo
				break
			}
		}
		rows.Close()

		//	// client_os_relation
		if osId == 0 && clientId != 0 {
			rows, err := db.Query(fmt.Sprintf(OsRelationSqlQuery, clientId))
			if err != nil {
				log.Printf("Error while OsRelationSqlQuery: %s", err)
			}
			for rows.Next() {
				rows.Scan(
					&osId,
					&userAgent.OsFamilyInfo.Name,
					&userAgent.OsFamilyInfo.Code,
					&userAgent.OsInfo.Name,
					&userAgent.OsInfo.Code,
					&userAgent.OsInfo.Homepage,
					&userAgent.OsInfo.Icon,
					&userAgent.OsInfo.IconBig,
					&userAgent.OsFamilyVendorInfo.Name,
					&userAgent.OsFamilyVendorInfo.Code,
					&userAgent.OsFamilyVendorInfo.Homepage,
				)
				break
			}
			rows.Close()
		}

		rows, err = db.Query(DevicesSqlQuery)
		if err != nil {
			panic(fmt.Sprintf("Error while DevicesSqlQuery: %s", err))
		}
		for rows.Next() {
			var device Device
			var regString string
			rows.Scan(
				&device.Id,
				&regString,
				&device.Info.Name,
				&device.Info.Code,
				&device.Info.Icon,
				&device.Info.IconBig,
			)
			d, err := getCompiledReqObject(device.Id, regString)
			if err != nil {
				panic(fmt.Sprintf("Error while req compilation: %s", err))
			}
			r := d.RegexCompiled
			matcher := r.MatcherString(userAgentString, 0)
			if !matcher.MatchString(userAgentString, 0) {
				continue
			}

			if matcher.Present(1) {
				deviceClassId = device.Id
				userAgent.DeviceInfo = device.Info
				userAgent.DeviceInfo.Url = ToNullString(DeviceDetailUrl + device.Info.Name.String)

				break
			}
		}

		if (deviceClassId == 0 && clientClassId != -1) {
			rows, err := db.Query(fmt.Sprintf(DeviceClassListSqlQuery))
			if err != nil {
				panic(fmt.Sprintf("Error while DeviceClassListSqlQuery: %s", err))
			}

			for rows.Next() {
				var device Device
				var classId int64
				rows.Scan(
					&classId,
					&device.Id,
					&device.Info.Name,
					&device.Info.Code,
					&device.Info.Icon,
					&device.Info.IconBig,
				)
				if clientClassId == classId {
					deviceClassId = device.Id
					userAgent.DeviceInfo = device.Info
					userAgent.DeviceInfo.Url = ToNullString(DeviceDetailUrl + device.Info.Name.String)
					break
				}
			}
			rows.Close()
		}

		// browser market name
		if userAgent.OsFamilyInfo.Code.String != "" {
			var (
				found     bool  = false
				deviceId  int64 = 0
				reqResult pcre.Regexp
			)
			rows, err := db.Query(DeviceCodeListSqlQuery)
			if err != nil {
				panic(fmt.Sprintf("Error while DeviceCodeListSqlQuery: %s", err))
			}
			for rows.Next() {
				var device DeviceCode
				var regString string
				rows.Scan(
					&device.Id,
					&device.OsFamilyCode,
					&device.OsCode,
					&regString,
				)
				if (userAgent.OsFamilyInfo.Code.String == device.OsFamilyCode && device.OsCode == "-all-") ||
					(userAgent.OsFamilyInfo.Code.String == device.OsFamilyCode && device.OsCode == userAgent.OsInfo.Code.String) {
					d, err := getCompiledReqObject(device.Id, regString)
					if err != nil {
						panic(fmt.Sprintf("Error while req compilation: %s", err))
					}
					reqResult = d.RegexCompiled
					found = true
					deviceId = device.Id
					break
				}

			}
			rows.Close()
			if found {
				matcher := reqResult.MatcherString(userAgent.String, 0)
				if matcher.Present(1) {
					rows, err := db.Query(fmt.Sprintf(BrandByRegexIdAndCodeSqlQuery, deviceId, matcher.GroupString(1)))
					if err != nil {
						log.Println("Error: ", err)
					}
					for rows.Next() {
						rows.Scan(
							&userAgent.DeviceMarketName,
							&userAgent.DeviceBrandInfo.Name,
							&userAgent.DeviceBrandInfo.Code,
							&userAgent.DeviceBrandInfo.Url,
							&userAgent.DeviceBrandInfo.Icon,
							&userAgent.DeviceBrandInfo.IconBig,
						)
					}
				}
			}
		}
	}

	return userAgent, nil
}

func (udger *UdgerFastLoadSlowExec) ParseIp(ip string) (ipAddress IpAddress) {

	if ip == "" {
		return
	}

	ipVersion, parsedIp := getIpVersion(ip)
	if ipVersion == 0 {
		return
	}

	var version string
	if ipVersion == net.IPv4len {
		version = IpV4
	}

	if ipVersion == net.IPv6len {
		version = IpV6
	}

	db, err := sql.Open("sqlite3", udger.dbPath)
	if err != nil {
		log.Println("sqlite3 error:", err)
		return
	}

	defer db.Close()

	rows, err := db.Query(IpSqlQuery)
	if err != nil {
		panic(fmt.Sprintf("Error while IpSqlQuery: %s", err))
	}
	found := false
	for rows.Next() {
		var _ipAddress IpAddress
		var ipString string
		if err := rows.Scan(
			&ipString,
			&_ipAddress.Crawler.Id,
			&_ipAddress.LastSeen,
			&_ipAddress.Hostname,
			&_ipAddress.Country,
			&_ipAddress.City,
			&_ipAddress.CountryCode,
			&_ipAddress.Classification,
			&_ipAddress.ClassificationCode,
			&_ipAddress.Crawler.Name,
			&_ipAddress.Crawler.Version,
			&_ipAddress.Crawler.VersionMajor,
			&_ipAddress.Crawler.LastSeen,
			&_ipAddress.Crawler.RespectRobotstxt,
			&_ipAddress.Crawler.FamilyInfo.Name,
			&_ipAddress.Crawler.FamilyInfo.Code,
			&_ipAddress.Crawler.FamilyInfo.Homepage,
			&_ipAddress.Crawler.FamilyInfo.Icon,
			&_ipAddress.Crawler.FamilyVendorInfo.Name,
			&_ipAddress.Crawler.FamilyVendorInfo.Code,
			&_ipAddress.Crawler.FamilyVendorInfo.Homepage,
			&_ipAddress.Crawler.CrawlerClassification,
			&_ipAddress.Crawler.CrawlerClassificationCode,
		); err != nil {
			panic(fmt.Sprintf("Error scanning IpSqlQuery results: %s", err))
		}

		if ip == ipString {
			ipAddress = _ipAddress
			found = true
			if ipAddress.ClassificationCode.String == CrawlerCode {
				botId := strconv.FormatInt(ipAddress.Crawler.Id.Int64, 10)
				ipAddress.Crawler.FamilyInfo.Url = ToNullString("https://udger.com/resources/ua-list/bot-detail?bot=" + ipAddress.Crawler.FamilyInfo.Name.String + "#id" + botId)
			}
		}
	}

	if !found {
		ipAddress.Classification = ToNullString("Unrecognized")
		ipAddress.ClassificationCode = ToNullString("unrecognized")
	}

	ipAddress.Ip.String = ip
	ipAddress.Version.String = version

	if ipAddress.Version.String == IpV4 {
		ipInt := ip2int(parsedIp)
		rows, err := db.Query(fmt.Sprintf(GetDataCenterV4SqlQuery, ipInt, ipInt))
		if err != nil {
			log.Println(err)
		}
		for rows.Next() {
			var dataCenterInfo Info
			rows.Scan(
				&dataCenterInfo.Name,
				&dataCenterInfo.Code,
				&dataCenterInfo.Homepage,
			)
			ipAddress.DataCenterInfo = dataCenterInfo
		}
	} else if ipAddress.Version.String == IpV6 {
		ip6AsIntArray := getIp6array(parsedIp)

		rows, err := db.Query(fmt.Sprintf(GetDataCenterV6SqlQuery,
			ip6AsIntArray[0],
			ip6AsIntArray[0],
			ip6AsIntArray[1],
			ip6AsIntArray[1],
			ip6AsIntArray[2],
			ip6AsIntArray[2],
			ip6AsIntArray[3],
			ip6AsIntArray[3],
			ip6AsIntArray[4],
			ip6AsIntArray[4],
			ip6AsIntArray[5],
			ip6AsIntArray[5],
			ip6AsIntArray[6],
			ip6AsIntArray[6],
			ip6AsIntArray[7],
			ip6AsIntArray[7],
		))
		if err != nil {
			log.Println(err)
		}
		for rows.Next() {
			var dataCenterInfo Info
			rows.Scan(
				&dataCenterInfo.Name,
				&dataCenterInfo.Code,
				&dataCenterInfo.Homepage,
			)
			ipAddress.DataCenterInfo = dataCenterInfo
		}
	} else {
		log.Panic("Unsupported ip, %s, version %s", ip, version)
	}

	return ipAddress
}
