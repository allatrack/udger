package main

import (
	"log"
	udger "github.com/peshkov3/udger-go"
)

func main() {
	
	// Example of Slow load - fast execution udgerSF client
	udgerSF, err := udger.NewSF("./udger.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	userAgent, err := udgerSF.ParseUa("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	if err != nil {
		log.Fatal(err)
	}

	//log.Printf("userAgent %+v", userAgent)
	log.Printf("Class %s", userAgent.Class)
	log.Printf("ClassCode %s", userAgent.ClassCode)
	log.Printf("Name %s", userAgent.Name)
	log.Printf("Version %s", userAgent.Version)

	ipAddress := udgerSF.ParseIp(`101.0.64.0`)
	log.Printf("Ip %s", ipAddress.Ip.String)
	log.Printf("Classification %s", ipAddress.Classification.String)
	log.Printf("ClassificationCode %s", ipAddress.ClassificationCode.String)
	log.Printf("Hostname %s", ipAddress.Hostname.String)
	log.Printf("Version %s", ipAddress.Version.String)
	log.Printf("Country %s", ipAddress.Country.String)

	// more examples see in udger_slow_fast_test.go

	// Example of Slow load - fast execution udgerSF client
	udgerFS, err := udger.NewFS("./udger.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	userAgent, err = udgerFS.ParseUa("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("User Agent %s", userAgent.String)
	log.Printf("Class %s", userAgent.Class)
	log.Printf("ClassCode %s", userAgent.ClassCode)
	log.Printf("User Agent Name %s", userAgent.Name)
	log.Printf("Version %s", userAgent.Version)
	log.Printf("VersionMajor %s", userAgent.VersionMajor)
	log.Printf("UpToDateCurrentVersion %s", userAgent.UpToDateCurrentVersion)

	ipAddress = udgerFS.ParseIp(`101.0.64.0`)
	log.Printf("Ip %s", ipAddress.Ip.String)
	log.Printf("Classification %s", ipAddress.Classification.String)
	log.Printf("ClassificationCode %s", ipAddress.ClassificationCode.String)
	log.Printf("Hostname %s", ipAddress.Hostname.String)
	log.Printf("Version %s", ipAddress.Version.String)
	log.Printf("Country %s", ipAddress.Country.String)

	// more examples see in udger_fast_slow_test.go
}
