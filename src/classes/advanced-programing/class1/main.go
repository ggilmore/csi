package main

import "unsafe"

func StringLength(s string) int {
	var ptr unsafe.Pointer

	lenPtr := (*int)(
		unsafe.Pointer(
			uintptr(unsafe.Pointer(&s)) + unsafe.Sizeof(ptr),
		),
	)

	return *lenPtr
}

func StructField(p point) int {
	yPtr := (*int)(
		unsafe.Pointer(
			uintptr(unsafe.Pointer(&p)) + unsafe.Offsetof(p.Y)),
	)

	return *yPtr
}

type point struct {
	X int
	Y int
}
