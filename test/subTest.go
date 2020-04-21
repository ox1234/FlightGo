package main

import (
	"fmt"
	"pentestplatform/gather/subdomain"
)

func main() {
	subScanner := subdomain.NewSubDomainScanner()
	subScanner.Set("xidian.edu.cn")
	subScanner.DoGather()
	fmt.Println(subScanner.Report())
}
