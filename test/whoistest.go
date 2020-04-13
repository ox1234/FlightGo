package main

import (
	"fmt"
	"pentestplatform/gather"
)

func main(){

	whoisScanner := gather.NewWhoisScanner("xidian.edu.cn")
	whoisScanner.DoGather()
	json, err := whoisScanner.Report()
	if err != nil{

	}
	fmt.Println(json)
}
