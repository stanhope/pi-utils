package main

import (
       "fmt"
	"syscall"
	"encoding/json"
	"io/ioutil"
	"strings"
)

type Utsname syscall.Utsname

func uname() (*syscall.Utsname, error) {
	uts := &syscall.Utsname{}

	if err := syscall.Uname(uts); err != nil {
		return nil, err
	}
	return uts, nil
}


type SystemInfo struct {
	Sysname    string `json:"sysname"`
	Nodename   string `json:"nodename"`
	Release    string `json:"release"`
	Version    string `json:"version"`
	Machine    string `json:"machine"`
	Domainname string `json:"domain"`
	Serial     string `json:"serial"`
}

func convertToStringARM(x [65]uint8) (string) {
	i := 0
	for ; i < 65; i++  {
		if x[i] == 0 {
			break
		}
		
	}
	str := string(x[0:i])
	return str
}

func convertToString(x [65]int8) (string) {
	i := 0
        s := make([]byte, 65)
	for ; i < 65; i++  {
		if x[i] == 0 {
			break
		} else {
			s[i] = uint8(x[i])
		}
		
	}
	//// str := string(x[0:i])
	str := string(s[0:i])
	return str
}

func getInfo(uts *syscall.Utsname) (SystemInfo) {
	var info SystemInfo
	info.Sysname = convertToStringARM(uts.Sysname)
	info.Nodename = convertToStringARM(uts.Nodename)
	info.Release = convertToStringARM(uts.Release)
	info.Version = convertToStringARM(uts.Version)
	info.Machine = convertToStringARM(uts.Machine)
	info.Domainname = convertToStringARM(uts.Domainname)
	info.Serial = getCPUSerial()
	return info
}

func getCPUSerial() (serial string) {
	contents, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range(lines) {
		if strings.Index(line,"Serial") == 0 {
			fields := strings.Split(line, ":")
			serial = strings.TrimSpace(fields[1])
			return
		}
	}
	return
}

func main() {

	uname,err := uname()
	if err == nil {
		info := getInfo(uname)
		b, _ := json.Marshal(info)
		fmt.Printf("%s\n", string(b))
	} else {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}

