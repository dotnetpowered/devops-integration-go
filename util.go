package main

import (
	"net"
	"strings"
)

const emptyString string = ""

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func getIpAddr(dnsName string) (string, string) {

	if !strings.Contains(dnsName, ".") {
		dnsName = dnsName + ".us.xeohealth.com"
	}
	ip, err := net.LookupHost(dnsName)
	_ = err
	ipAddr := ""
	if len(ip) > 0 {
		ipAddr = ip[0]
	}
	return dnsName, ipAddr
}
