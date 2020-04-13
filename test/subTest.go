package main

import (
	"fmt"
)

type test1 struct {
	aaa []test2
}

type test2 struct {
	Name string
}

func main(){
	fmt.Println()
	t1 := new(test1)
	for _, t := range t1.aaa{
		fmt.Println(t.Name)
	}
}
