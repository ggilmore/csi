//nolint
package main

import (
	"math"
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

func MapMax(m map[int]int) int {
	hashHeaderPtr := unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&m)))

	B := *(*uint8)(unsafe.Add(hashHeaderPtr, unsafe.Offsetof(myMap{}.B)))
	numBuckets := (1 << B)

	bucketBase := unsafe.Pointer(*(*uintptr)(unsafe.Add(hashHeaderPtr, unsafe.Offsetof(myMap{}.buckets))))

	max := math.MinInt
	for i := 0; i < numBuckets; i++ {
		bucketMax := maxBucket(unsafe.Add(bucketBase, uintptr(i)*unsafe.Sizeof(myBucket{})))
		if bucketMax > max {
			max = bucketMax
		}
	}

	return max
}

func maxBucket(bucketPtr unsafe.Pointer) int {
	topHashSize := unsafe.Sizeof(uint8(1))
	intSize := unsafe.Sizeof(1)

	max := math.MinInt

	for i := 0; i < 8; i++ {
		topHashPtr := unsafe.Add(bucketPtr, unsafe.Offsetof(myBucket{}.topHashes)+uintptr(i)*topHashSize)
		topHash := *(*uint8)(topHashPtr)

		if topHash == 0 {
			return max
		}

		valuePtr := unsafe.Add(
			bucketPtr,
			unsafe.Offsetof(myBucket{}.values)+uintptr(i)*intSize,
		)

		value := *(*int)(valuePtr)

		if value > max {
			max = value
		}
	}

	overflowPtr := unsafe.Add(bucketPtr, unsafe.Offsetof(myBucket{}.overflow))
	if overflowPtr != nil {
		overflowMax := maxBucket(unsafe.Pointer(*(*uintptr)(overflowPtr)))
		if overflowMax > max {
			max = overflowMax
		}
	}

	return max
}

type myBucket struct {
	topHashes [8]uint8
	keys      [8]int
	values    [8]int

	overflow unsafe.Pointer
}

// copied from src/runtime/map.go
type myMap struct {
	// Note: the format of the hmap is also encoded in cmd/compile/internal/reflectdata/reflect.go.
	// Make sure this stays in sync with the compiler's definition.
	count     int // # live cells == size of map.  Must be first (used by len() builtin)
	flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash seed

	buckets unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.

	// other fields are irrelevant for this
}
