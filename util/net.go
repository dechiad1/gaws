package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Ipify struct {
	Ip string
}

func GetPublicIp() string {
	req, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		fmt.Println("could not reach api.ipify.org")
		os.Exit(1)
	}
	body, _ := ioutil.ReadAll(req.Body)
	var ipaddr Ipify
	json.Unmarshal(body, &ipaddr)

	return ipaddr.Ip
}
