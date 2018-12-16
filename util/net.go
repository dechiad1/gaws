package util

import (
	"fmt"
	"os"
	"net/http"
  "io/ioutil"
	"encoding/json"
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
