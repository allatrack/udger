package udger

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net"
)

func TestGetEmptyIpVersion(t *testing.T) {
	ip,_:=getIpVersion(``)
	assert.Equal(t, 0, ip)
}

func TestGetEmptyIp6Version(t *testing.T) {
	assert.Equal(t, 8, len(getIp6array(net.ParseIP(`::`))))
}

func TestGetIp6arrayLength(t *testing.T) {
	assert.Equal(t, 8, len(getIp6array(net.ParseIP(`FE80:CD00:0000:0CDE:1257:0000:211E:729C`))))
	assert.Equal(t, 8, len(getIp6array(net.ParseIP("::1"))))
}

func TestGetIp6array(t *testing.T) {
	assert.Equal(t, []int64{65152, 52480, 0, 3294, 4695, 0, 8478, 29340,}, getIp6array(net.ParseIP(`FE80:CD00:0000:0CDE:1257:0000:211E:729C`)))
	assert.Equal(t, []int64{0, 0, 1, 0, 0, 0, 0, 0,}, getIp6array(net.ParseIP("::1")))
}

func TestGetIp6Version(t *testing.T) {
	ip, _:= getIpVersion(`FE80:CD00:0000:0CDE:1257:0000:211E:729C`)
	assert.Equal(t, 16, ip)
	ip, _= getIpVersion(`FE80:CD00:0:CDE:1257:0:211E:729C`)
	assert.Equal(t, 16, ip)
}

func TestGetIp66LoopbackVersion(t *testing.T) {
	ip, _:= getIpVersion(`::1`)
	assert.Equal(t, 16, ip)
}

func TestGetIp4Version(t *testing.T) {
	ip, _:= getIpVersion(`0.0.0.0`)
	assert.Equal(t, 4, ip)
	ip, _= getIpVersion(`127.0.0.1`)
	assert.Equal(t, 4, ip)
}

