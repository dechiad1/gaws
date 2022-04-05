package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func CalculateCidrRange(ip string, cidr string) (ipRange string) {
	//TODO: real cidr calculation - for now just make it a 24
	parts := strings.Split(ip, ".")
	parts[3] = "0"
	ipRange = strings.Join(parts, ".") + "/24"
	return ipRange
}

func GetPortRange(portRange string) (int64, int64) {
	if portRange == "22" {
		return 22, 22
	}
	r := strings.Split(portRange, "-")
	if len(r) != 2 {
		fmt.Println("portRange must be in the format of ##-##. ie 8000-8080")
		panic("invalid port range")
	}
	fp, err := strconv.ParseInt(r[0], 10, 64) //string to parse, numerical base, size int (int64)
	tp, err := strconv.ParseInt(r[1], 10, 64)
	if err != nil {
		fmt.Printf("%s, %s are not ints\n", r[0], r[1])
		panic(err)
	}
	return fp, tp
}

type ipify struct {
	Ip string
}

func GetPublicIp() string {
	req, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		fmt.Println("could not reach api.ipify.org")
		os.Exit(1)
	}
	body, _ := ioutil.ReadAll(req.Body)
	var ipaddr ipify
	json.Unmarshal(body, &ipaddr)

	return ipaddr.Ip
}
