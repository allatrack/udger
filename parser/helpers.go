package parser

import (
	"strings"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"errors"
	"time"
	"log"
	"net"
	"strconv"
	"encoding/binary"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func findDataWithVersion(ua string, data []rexData, withVersion bool) (idx int64, value string, err error) {
	defer func() {
		if r := recover(); r != nil {
			idx, value, err = findData(ua, data, false)
		}
	}()

	idx, value, err = findData(ua, data, withVersion)

	return idx, value, err
}

func cleanRegex(r string) string {
	if strings.HasSuffix(r, "/si") {
		r = r[:len(r)-3]
	}
	if strings.HasPrefix(r, "/") {
		r = r[1:]
	}

	return r
}

func findData(ua string, data []rexData, withVersion bool) (idx int64, value string, err error) {
	for i := 0; i < len(data); i++ {
		r := data[i].RegexCompiled
		matcher := r.MatcherString(ua, 0)
		if !matcher.MatchString(ua, 0) {
			continue
		}

		if withVersion && matcher.Present(1) {
			return data[i].ID, matcher.GroupString(1), nil
		}

		return data[i].ID, "", nil
	}

	return -1, "", nil
}

func findRegex(ua string, data []rexData, withVersion bool) (regex string, value string, err error) {
	for i := 0; i < len(data); i++ {
		r := data[i].RegexCompiled
		matcher := r.MatcherString(ua, 0)
		if !matcher.MatchString(ua, 0) {
			continue
		}

		if withVersion && matcher.Present(1) {
			return data[i].Regex, matcher.GroupString(1), nil
		}

		return data[i].Regex, "", nil
	}

	return "", "", nil
}


func getCompiledReqObject(id int64, regString string) (d rexData, err error) {
	d.ID = id
	d.Regex = cleanRegex(regString)
	r, _err := pcre.Compile(d.Regex, pcre.CASELESS)
	if _err != nil {
		return d, errors.New(_err.String())
	}
	d.RegexCompiled = r
	return
}

func getIp6array(parsedIp net.IP) (ip6AsIntArray []int64) {

	ipAsString := parsedIp.String()
	if ipAsString == "" {
		return
	}
	ipAsArray := strings.SplitN(ipAsString, ":", 8)
	if ipAsArray == nil{
		return
	}

	for _, ipPart := range ipAsArray {
		if ipPart == "" {
			ip6AsIntArray = append(ip6AsIntArray, 0)
			continue
		}
		n, err := strconv.ParseInt(ipPart, 16, 32)
		if err != nil {
			panic(err)
		}
		ip6AsIntArray = append(ip6AsIntArray, n)
	}

	resultingArrayLen := len(ip6AsIntArray)
	if resultingArrayLen < 8 {
		lenToAdd := 8 - resultingArrayLen
		for i := 0; i < lenToAdd; i++ {
			ip6AsIntArray = append(ip6AsIntArray, 0)
		}
	}

	return
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func getIpVersion(ip string) (int, net.IP){
	_ip := net.ParseIP(ip)
	if _ip.To4() != nil {
		return net.IPv4len, _ip
	}

	if _ip.To16() != nil {
		return net.IPv6len, _ip
	}

	return 0, nil
}