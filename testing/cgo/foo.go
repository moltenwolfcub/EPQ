package main

// #cgo CFLAGS: -I${SRCDIR}
// #cgo LDFLAGS: -L${SRCDIR}/foo.c
// #include <foo.h>
import "C"

func main() {
	age := 3
	year := int(C.getSchoolYear(C.int(age)))
	println("School year:", year)
}
