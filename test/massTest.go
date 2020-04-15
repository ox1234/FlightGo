package main

import (
	"fmt"
	"pentestplatform/gather"
)

func main(){
	portScanner := gather.NewPortScanner()
	portScanner.Set("61.150.43.100")
	//portScanner.Set("120.77.152.169")
	portScanner.DoGather()
	fmt.Println(portScanner.Report())
}
