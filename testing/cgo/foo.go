package main

// #cgo CFLAGS: -Ipwd
// #cgo LDFLAGS: foo.a
// #include <foo.h>
import "C"

func main() {
	age := 16
	year := int(C.getSchoolYear(C.int(age)))
	println("School year:", year)
}
