package main

import (
	"fmt"
	"time"
)

func temp() {
	fmt.Println("temp")
}
func main() {
	var a = [5]int{1, 2, 3, 4, 5}
	var r [5]int

	fmt.Println("original a =", a)

	for i, v := range a {
		if i == 0 {
			a[1] = 12
			a[2] = 13
		}
		r[i] = v
	}

	fmt.Println("after for range loop, r =", r)
	fmt.Println("after for range loop, a =", a)
}

type field struct {
	name string
}

func (p *field) print() {
	fmt.Println(p.name)
}

func main1() {
	data1 := []*field{{"one"}, {"two"}, {"three"}}
	data1[0].name = "modified"
	for _, v := range data1 {
		go v.print()
	}

	data2 := []field{{"four"}, {"five"}, {"six"}}

	for _, v := range data2 {
		go v.print()
	}
	time.Sleep(3 * time.Second)

}
