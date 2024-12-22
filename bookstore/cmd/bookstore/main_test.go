package main

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	arr := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sl := arr[3:7:9]

	var p map[string]int
	for _, v := range sl {
		fmt.Println(v)
	}
}
