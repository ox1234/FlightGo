package main

import (
	"fmt"
	"pentestplatform/gather"
)

func main(){
	basicScanner := gather.NewBasicScanner()
	basicScanner.Set("www.xidian.edu.cn", "asdfasdf")
	basicScanner.Set("ehall.xidian.edu.cn", "kalsjdflk")
	basicScanner.Set("akjsdkf.xidian.edu.cn", "klajsdlkfkasdjf")
	basicScanner.DoGather()
	fmt.Println("complete")
}
