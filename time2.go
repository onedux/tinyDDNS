package main

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"strings"
	"net"
	"os"
)

//curl -X GET -H "Authorization: sso-key 9jL9GABGRBx_xqY6VfDgW2bKUZaHS4sng:xqbBaBVo6kGy6v37H4yMA" "https://api.godaddy.com/v1/domains/available?domain=onedux.com"


func main() {
	const TIME_INTERVAL int = 30*60
	const DOMAIN_NAME string = "itts.io"
	const DOMAIN_SUBNAME string = "www"
	const DOMAIN_TYPE string = "A"
	const DOMAIN_TTL int = 600
	const GoDaddy_Key string = "9jL9GABGRBx_xqY6VfDgW2bKUZaHS4sng"
	const GoDaddy_Sec string = "xqbBaBVo6kGy6v37H4yMA"

	fmt.Println("Starting the application..."+"\n")
	for true  {
		ip1 :=   getPublicIP2()
		ip2 :=   getDNSRecordIP2(DOMAIN_SUBNAME + "." + DOMAIN_NAME)
		// ip2 = "0.0.0.0"
		fmt.Println("PublicIP : " + ip1+"\n")
		fmt.Println("DNSRecord : " + DOMAIN_SUBNAME + "." + DOMAIN_NAME + ":" + ip2+"\n")
		if ip1 == ip2 {
			fmt.Println("No Need For Update Domian Record "+"\n")
			time.Sleep(time.Duration(TIME_INTERVAL) * time.Second)
			continue
		}else{
			fmt.Print("Ready To Update Domian Record ..."+"\n")

		}
		setDNSRecordIP2(ip1,DOMAIN_TTL,DOMAIN_NAME,DOMAIN_TYPE,DOMAIN_SUBNAME,GoDaddy_Key,GoDaddy_Sec)
		time.Sleep(time.Duration(TIME_INTERVAL) * time.Second)
	}
}

func getPublicIP2() string {
	var ip = "0.0.0.0"
	response, err := http.Get("http://ipv4.icanhazip.com/")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		ip = string(data)

	}
	return strings.TrimSpace(ip)
}

//func getDNSRecordIP1(domainName string) string {
//	var ip = "0.0.0.0" <span> </span>domain := domainName <span> </span>ipAddr, err := net.ResolveIPAddr("ip", domain) <span> </span>if err != nil { <span> </span>fmt.Fprintf(os.Stderr, "Err: %s", err.Error()) <span> </span>return ip <span> </span>} <span> </span>ip = ipAddr.IP.String()
	// <span> </span>fmt.Println(ip)
	//return strings.TrimSpace(ip)
//}
func getDNSRecordIP2(domainName string) string{
pTCPAddr, err := net.ResolveTCPAddr("tcp", "www.itts.io:80")
if err != nil {
fmt.Fprintf(os.Stderr, "Err: %s", err.Error())
return ""
}

	fmt.Fprintf(os.Stdout, "www.itts.io:80 IP: %s PORT: %d", pTCPAddr.IP.String(), pTCPAddr.Port)
	return pTCPAddr.IP.String()
}



func setDNSRecordIP2(ip string, ttl int, domain string, domainType string, name string, key string, sec string) {
	url := fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/%s/%s",domain,domainType,name)
	data :=fmt.Sprintf("[{ \"data\": \"%s\", \"ttl\": %v, \"priority\": 0, \"weight\": 1 }]",ip,ttl)
	//生成client 参数为默认
	client := &http.Client{}
	//提交请求
	reqest, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		panic(err)
	}
	reqest.Header.Add("content-type", "application/json")
	reqest.Header.Add("Accept", "application/json")
	ssokey :=fmt.Sprintf("sso-key %s:%s",key,sec)
	reqest.Header.Add("Authorization", ssokey)
	//处理返回结果
	response, _ := client.Do(reqest)
	//返回的状态码
	status := response.StatusCode
	if status == 200 {
		fmt.Println("Done!")
	}else{
		fmt.Print(status)
		fmt.Println(" Error Hapened!")
	}

}