package udger

import (
	"sync"
	"fmt"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"log"
	"strings"
	"strconv"
	"net"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"runtime"
)

func NewSF(dbPath string) (*UdgerSlowLoadFastExec, error) {

	uSlowFast := &UdgerSlowLoadFastExec{
		BrowserList:     make(map[int64]Browser),
		OsList:          make(map[int64]Os),
		IpList:          make(map[string]IpAddress),
		CrawlerList:     make(map[string]Crawler),
		DeviceList:      make(map[string]Device),
		DeviceClassList: make(map[int64]Device),
		DeviceCodeList:  make(map[int64]DeviceCode),
	}

	var err error

	if _, err = os.Stat(dbPath); os.IsNotExist(err) {
		return uSlowFast, err
	}

	uSlowFast.db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return uSlowFast, err
	}

	defer uSlowFast.db.Close()
	uSlowFast.DbPath = dbPath

	err = uSlowFast.initSF()
	if err != nil {
		return uSlowFast, err
	}
	return uSlowFast, nil
}

func (udger *UdgerSlowLoadFastExec) initSF() error {

	//defer timeTrack(time.Now(), "udger init")

	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	errChannel := make(chan error, 1)
	finished := make(chan bool, 1)
	waitGroupLength := 7
	wg.Add(waitGroupLength)

	go func() {

		rows, err := udger.db.Query(BrowsersSqlQuery)
		if err != nil {
			errChannel <- Error{fmt.Sprintf("Error executing BrowsersSqlQuery: %s", err)}
			wg.Done()
			return
		}
		var regString string

		for rows.Next() {
			var browser Browser
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
				errChannel <- Error{fmt.Sprintf("Error scanning BrowsersSqlQuery row: %s", err)}
				break
			}
			d, err := getCompiledReqObject(browser.Id, regString)
			if err != nil {
				errChannel <- Error{fmt.Sprintf("Prepare regex error: %s", err)}
				break
			}
			udger.rexBrowsers = append(udger.rexBrowsers, d)
			udger.BrowserList[browser.Id] = browser
		}
		rows.Close()
		wg.Done()
	}()

	go func() {
		rows, err := udger.db.Query(OsSqlQuery)
		if err != nil {
			errChannel <- Error{fmt.Sprintf("Error executing OsSqlQuery: %s", err)}
			wg.Done()
			return
		}
		var regString string
		for rows.Next() {
			var operationSystem Os
			if err := rows.Scan(
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
			); err != nil {
				errChannel <- Error{fmt.Sprintf("Error scanning OsSqlQuery row: %s", err)}
				break
			}
			d, err := getCompiledReqObject(operationSystem.Id, regString)
			if err != nil {
				errChannel <- Error{fmt.Sprintf("Prepare regex error: %s", err)}
				break
			}
			udger.rexOS = append(udger.rexOS, d)
			udger.OsList[operationSystem.Id] = operationSystem
		}
		rows.Close()
		wg.Done()
	}()

	go func() {
		rows, err := udger.db.Query(DevicesSqlQuery)
		if err != nil {
			errChannel <- Error{fmt.Sprintf("Error executing DevicesSqlQuery: %s", err)}
			wg.Done()
			return
		}
		var regString string
		for rows.Next() {
			var device Device
			if err := rows.Scan(
				&device.Id,
				&regString,
				&device.Info.Name,
				&device.Info.Code,
				&device.Info.Icon,
				&device.Info.IconBig,
			); err != nil {
				errChannel <- Error{fmt.Sprintf("Error scanning DevicesSqlQuery row: %s", err)}
				break
			}
			d, err := getCompiledReqObject(device.Id, regString)
			if err != nil {
				errChannel <- Error{fmt.Sprintf("Prepare regex error: %s", err)}
				break
			}
			udger.rexDevices = append(udger.rexDevices, d)
			udger.DeviceList[ cleanRegex(regString)] = device
		}
		rows.Close()
		wg.Done()
	}()

	go func() {

		rows, err := udger.db.Query(fmt.Sprintf(DeviceClassListSqlQuery))
		if err != nil {
			errChannel <- Error{fmt.Sprintf("Error executing DeviceClassListSqlQuery: %s", err)}
			wg.Done()
			return
		}

		for rows.Next() {
			var device Device
			var classId int64
			if err := rows.Scan(
				&classId,
				&device.Id,
				&device.Info.Name,
				&device.Info.Code,
				&device.Info.Icon,
				&device.Info.IconBig,
			); err != nil {
				errChannel <- Error{fmt.Sprintf("Error scanning DeviceClassListSqlQuery row: %s", err)}
				break
			}
			udger.DeviceClassList[classId] = device
		}
		rows.Close()
		wg.Done()
	}()

	go func() {
		rows, err := udger.db.Query(DeviceCodeListSqlQuery)
		if err != nil {
			errChannel <- Error{fmt.Sprintf("Error executing DeviceCodeListSqlQuery: %s", err)}
			wg.Done()
			return
		}
		var index int64 = 0
		var regString string

		for rows.Next() {
			var device DeviceCode
			if err := rows.Scan(
				&device.Id,
				&device.OsFamilyCode,
				&device.OsCode,
				&regString,
			); err != nil {
				errChannel <- Error{fmt.Sprintf("Error scanning DeviceClassListSqlQuery row: %s", err)}
				break
			}
			var regex = cleanRegex(regString)
			regexCompiled, _err := pcre.Compile(regex, pcre.CASELESS)
			if _err != nil {
				errChannel <- Error{fmt.Sprintf("Prepare regex error: %s", _err.String())}
			}
			device.RegexCompiled = regexCompiled
			udger.DeviceCodeList[index] = device
			index++
		}
		rows.Close()
		wg.Done()
	}()

	go func() {

		rows, err := udger.db.Query(IpSqlQuery)
		if err != nil {
			errChannel <- Error{fmt.Sprintf("Error executing IpSqlQuery: %s", err)}
			wg.Done()
			return
		}
		var ipString string

		for rows.Next() {
			var ip IpAddress
			if err := rows.Scan(
				&ipString,
				&ip.Crawler.Id,
				&ip.LastSeen,
				&ip.Hostname,
				&ip.Country,
				&ip.City,
				&ip.CountryCode,
				&ip.Classification,
				&ip.ClassificationCode,
				&ip.Crawler.Name,
				&ip.Crawler.Version,
				&ip.Crawler.VersionMajor,
				&ip.Crawler.LastSeen,
				&ip.Crawler.RespectRobotstxt,
				&ip.Crawler.FamilyInfo.Name,
				&ip.Crawler.FamilyInfo.Code,
				&ip.Crawler.FamilyInfo.Homepage,
				&ip.Crawler.FamilyInfo.Icon,
				&ip.Crawler.FamilyVendorInfo.Name,
				&ip.Crawler.FamilyVendorInfo.Code,
				&ip.Crawler.FamilyVendorInfo.Homepage,
				&ip.Crawler.CrawlerClassification,
				&ip.Crawler.CrawlerClassificationCode,
			); err != nil {
				errChannel <- Error{fmt.Sprintf("Error scanning IpSqlQuery results: %s", err)}
				break
			}
			udger.IpList[ipString] = ip
		}

		rows.Close()
		wg.Done()
	}()
	go func() {
		rows, err := udger.db.Query(CrawlersSqlQuery)
		if err != nil {
			errChannel <- Error{fmt.Sprintf("Error executing CrawlersSqlQuery: %s", err)}
			wg.Done()
			return
		}
		for rows.Next() {
			var crawler Crawler
			var familyInfo Info
			var familyVendorInfo Info

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
				errChannel <- Error{fmt.Sprintf("Error scanning CrawlersSqlQuery: %s", err)}
				break
			}

			udger.CrawlerList[strings.TrimSpace(crawler.UaString)] = crawler
		}
		rows.Close()
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChannel:
		if err != nil {
			panic(err)
			return err
		}
	}

	return nil
}

func (udger *UdgerSlowLoadFastExec) ParseUa(userAgentString string) (userAgent UserAgent, err error) {

	//defer timeTrack(time.Now(), "ParseUa")

	userAgent = UserAgent{}
	udger.db, err = sql.Open("sqlite3", udger.DbPath)
	if err != nil {
		return userAgent, err
	}
	defer udger.db.Close()

	if userAgentString != "" {

		userAgent.String = strings.TrimSpace(userAgentString)

		var (
			clientId      int64 = 0
			clientClassId int64 = -1
			osId          int64 = 0
			deviceClassId int64 = 0
		)

		if crawler, ok := udger.CrawlerList[userAgent.String]; ok {

			userAgent.setVersion(crawler.Version, crawler.Name.String)

			clientClassId = CrawlerClassId

			userAgent.Class = CrawlerClass
			userAgent.ClassCode = CrawlerCode
			userAgent.FamilyInfo = crawler.FamilyInfo
			userAgent.FamilyVendorInfo = crawler.FamilyVendorInfo
			userAgent.FamilyInfo.Url = ToNullString(BotDetailUrl + crawler.FamilyInfo.Name.String + "#id" + strconv.FormatInt(crawler.Id.Int64, 10))
			userAgent.Crawler = crawler

		} else {

			// browser
			id, version, err := findDataWithVersion(userAgentString, udger.rexBrowsers, true)
			if err != nil {
				log.Println("Error findDataWithVersion", err)
				return userAgent, err
			}

			if browser, ok := udger.BrowserList[id]; ok {
				userAgent.setVersion(version, browser.BrowserInfo.Name.String)

				clientId = browser.Id
				clientClassId = browser.ClassId

				userAgent.Class = browser.ClientClassification
				userAgent.ClassCode = browser.ClientClassificationCode
				userAgent.UpToDateCurrentVersion = browser.UptodateCurrentVersion
				userAgent.FamilyInfo = browser.BrowserInfo
				userAgent.FamilyVendorInfo = browser.VendorInfo
				userAgent.FamilyInfo.Url = ToNullString(BrowserDetailUrl + browser.BrowserInfo.Name.String)
				userAgent.Engine = browser.Engine
			}

			// os
			id, _, err = findData(userAgentString, udger.rexOS, false)
			if err != nil {
				log.Println("Error findData", err)
				return userAgent, err
			}

			if clientOs, ok := udger.OsList[id]; ok {
				osId = clientOs.Id
				userAgent.OsInfo = clientOs.Info
				userAgent.OsInfo.Url = ToNullString(OsDetailUrl + clientOs.Info.Name.String)
				userAgent.OsFamilyInfo = clientOs.FamilyInfo
				userAgent.OsFamilyVendorInfo = clientOs.VendorInfo
			}

			// client_os_relation
			if osId == 0 && clientId != 0 {
				rows, err := udger.db.Query(fmt.Sprintf(OsRelationSqlQuery, clientId))
				if err != nil {
					log.Printf("Error while getting OS relation: %s", err)
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
			}

			// browser
			regex, _, err := findRegex(userAgentString, udger.rexDevices, false)
			if err != nil {
				log.Println("Error while browser recognition ", err)
				return userAgent, err
			}

			if device, ok := udger.DeviceList[regex]; ok {
				deviceClassId = device.Id
				userAgent.DeviceInfo = device.Info
				userAgent.DeviceInfo.Url = ToNullString(DeviceDetailUrl + device.Info.Name.String)
			}
			if (deviceClassId == 0 && clientClassId != -1) {
				if device, ok := udger.DeviceClassList[clientClassId]; ok {
					deviceClassId = device.Id
					userAgent.DeviceInfo = device.Info
					userAgent.DeviceInfo.Url = ToNullString(DeviceDetailUrl + device.Info.Name.String)
				}
			}

			// browser marketname
			if userAgent.OsFamilyInfo.Code.String != "" {
				var (
					found     bool  = false
					deviceId  int64 = 0
					reqResult pcre.Regexp
				)

				for _, deviceCode := range udger.DeviceCodeList {
					if (userAgent.OsFamilyInfo.Code.String == deviceCode.OsFamilyCode && deviceCode.OsCode == "-all-") ||
						(userAgent.OsFamilyInfo.Code.String == deviceCode.OsFamilyCode && deviceCode.OsCode == userAgent.OsInfo.Code.String) {
						reqResult = deviceCode.RegexCompiled
						found = true
						deviceId = deviceCode.Id
						break
					}
				}

				if found {
					matcher := reqResult.MatcherString(userAgent.String, 0)
					if matcher.Present(1) {
						rows, err := udger.db.Query(fmt.Sprintf(BrandByRegexIdAndCodeSqlQuery, deviceId, matcher.GroupString(1)))
						if err != nil {
							log.Println("Error quering BrandByRegexIdAndCodeSqlQuery: ", err)
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

			return userAgent, nil
		}
	}

	return userAgent, nil
}

func (udger *UdgerSlowLoadFastExec) ParseIp(ip string) (ipAddress IpAddress) {

	if ip == "" {
		return
	}

	ipVersion, parsedIp := getIpVersion(ip)

	var version string
	if ipVersion == net.IPv4len {
		version = IpV4
	}

	if ipVersion == net.IPv6len {
		version = IpV6
	}

	db, err := sql.Open("sqlite3", udger.DbPath)
	if err != nil {
		log.Println("Open sqlite3 DB error:", err)
		return
	}
	udger.db = db
	defer udger.db.Close()

	var found = false
	if _ipAddress, ok := udger.IpList[ip]; ok {
		ipAddress = _ipAddress
		found = true
		if ipAddress.ClassificationCode.String == CrawlerCode {
			botId := strconv.FormatInt(ipAddress.Crawler.Id.Int64, 10)
			ipAddress.Crawler.FamilyInfo.Url = ToNullString("https://udger.com/resources/ua-list/bot-detail?bot=" + ipAddress.Crawler.FamilyInfo.Name.String + "#id" + botId)
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
		rows, err := udger.db.Query(fmt.Sprintf(GetDataCenterV4SqlQuery, ipInt, ipInt))
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

		rows, err := udger.db.Query(fmt.Sprintf(GetDataCenterV6SqlQuery,
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
	}

	return ipAddress
}
