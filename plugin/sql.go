package main

import "fmt"

type exploiter string
func (e exploiter) Exploit()(string, error){
	target := v[0].(string)
	fmt.Println("attack at " + target)
	msg := "success"
	return msg, nil
}

func (e exploiter) Describe()string{
	return "这是一个sql注入测试模块"
}

var FAttacker exploiter
