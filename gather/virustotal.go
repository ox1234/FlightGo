package gather

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net"
	"net/http"
	"pentestplatform/logger"
	"strings"
)

const VT_KEY string = "406b7ac8ae475758ea595ebc7e1cfc9eb81ae754bcef33df1811710d4bb5da0e"
const apiUrl string = "https://www.virustotal.com/api/v3/domains/%s/subdomains?limit=40"

type vtScanner struct {
	Domain string
	VtDomainSet []subDomain
}

func NewVtScanner() *vtScanner{
	return &vtScanner{}
}

func (v *vtScanner) Set(va ...interface{}){
	v.Domain = va[0].(string)
}

func (v *vtScanner) DoGather(){
	url := fmt.Sprintf(apiUrl, v.Domain)
	for{
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("x-apikey", VT_KEY)
		resp, err := (&http.Client{}).Do(req)
		if err != nil{
			logger.Red.Println(err)
		}
		bodyReader := resp.Body
		bodyBytes, _ := ioutil.ReadAll(bodyReader)
		body := string(bodyBytes)
		bodyReader.Close()
		domainSet, next := v.extractDomain(body)
		v.VtDomainSet = append(v.VtDomainSet, domainSet...)
		if next == ""{
			break
		}else{
			url = next
		}
	}
	logger.Green.Println("vt scan over!")
}

func (v *vtScanner) Report() (string, error){
	jsondata, err := json.Marshal(v)
	if err != nil{
		logger.Red.Println(err)
		return "", err
	}
	return string(jsondata), nil
}

func (v *vtScanner) extractDomain(jsondata string)([]subDomain,string){
	var domainSet []subDomain
	datajson := gjson.Get(jsondata, "data")
	if datajson.IsArray(){
		for _, element := range datajson.Array(){
			domainName := element.Get("id").String()
			lastDnsRecords := element.Get("attributes").Get("last_dns_records")
			var ipAddress string
			if lastDnsRecords.IsArray() && len(lastDnsRecords.Array()) > 0{
				ipAddress = lastDnsRecords.Array()[0].Get("value").String()
			}else{
				ips, _ := net.LookupHost(domainName)
				ipAddress = strings.Join(ips, ",")
			}
			domainSet = append(domainSet, subDomain{
				IPAddress: ipAddress,
				HostName:  domainName,
			})
		}
	}
	next := gjson.Get(jsondata, "links").Get("next").String()
	return domainSet, next
}


