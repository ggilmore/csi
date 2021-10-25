package main

import (
	"unsafe"
)

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

func SumSlice(xs []int) int {
	intSize := unsafe.Sizeof(1)
	ptrSize := unsafe.Sizeof(&intSize)

	sliceLen := *((*int)(
		unsafe.Pointer(uintptr(unsafe.Pointer(&xs)) + ptrSize),
	))

	sum := 0
	for i := 0; i < sliceLen; i++ {
		elem := *(*int)(
			unsafe.Pointer(
				uintptr(
					// follow the pointer to the addr of the first element
					*(*int)(unsafe.Pointer(&xs)),
				) +
					uintptr(i)*intSize,
			))

		sum += elem
	}

	return sum
}
