package udpproxy

import (
	"regexp"
	"strconv"
)

var AddressPattern = regexp.MustCompile(`(?P<host>\d{1-3}(\.\d{1-3}){3})?:?(?P<port>\d+)`)

func ParseAddress(address string) (host string, port int){
	match := AddressPattern.FindStringSubmatch(address)
	result := make(map[string]string)
	for i, name := range AddressPattern.SubexpNames() {
		if i !=0 && name != "" {
			result[name] = match[i]
		}
	}
	host, exist := result["host"]
	if !exist {
		host = "127.0.0.1"
	}
	port, _ = strconv.Atoi(result["port"])
	return host, port
}
