package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func main(){
	url := "http://blog.f1ight.top/"
	resp, err := http.Get(url)
	if err != nil{
		log.Fatal(err)
	}
	bodyReader := resp.Body
	defer resp.Body.Close()

	bodyBytes,_ := ioutil.ReadAll(bodyReader)
	r, _ := regexp.Compile(`<title>(\s*.*\s*)</title>`)
	title := strings.TrimSpace(r.FindStringSubmatch(string(bodyBytes))[1])

	fmt.Println(title)
}
