package util

import (
	"fmt"
	"os"
)

var cwd string

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("")
	}
	fmt.Println(cwd)
}

/*
Add test
*/
func Add(a int, b int) int {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("")
	}
	fmt.Println(cwd)
	return a + b
}

/*
PrintTypeAndValue test
*/
func PrintTypeAndValue(a interface{}) {
	fmt.Println(cwd)
	fmt.Printf("%T , %v \n", a, a)
}
