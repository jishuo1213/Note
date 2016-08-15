package util

import (
	"fmt"
)

/*
Add test
*/
func Add(a int, b int) int {
	return a + b
}

/*
PrintTypeAndValue test
*/
func PrintTypeAndValue(a interface{}) {
	fmt.Printf("%T , %v \n", a, a)
}
