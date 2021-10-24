package main

import "unsafe"

func StringLength(s string) int {
	var ptr unsafe.Pointer

	lenPtr := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + unsafe.Sizeof(ptr)))
	return *lenPtr
}
