package main

import (
	"fmt"
	"pentestplatform/gather"
)

func main(){
	rapidDnsScanner := gather.NewRapidDnsScanner()
	rapidDnsScanner.Set("xidian.edu.cn")
	rapidDnsScanner.DoGather()
	fmt.Println(rapidDnsScanner.Report())
}
