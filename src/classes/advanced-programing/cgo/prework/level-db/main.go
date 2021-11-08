package main

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -l leveldb
#include <leveldb/c.h>
#include <stdlib.h>
#include <stdio.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {

	// make db
	opts := C.leveldb_options_create()
	C.leveldb_options_set_create_if_missing(opts, 1)
	C.leveldb_options_set_error_if_exists(opts, C.uchar(1))

	name := C.CString("abc")
	errptr := C.CString("")

	db := C.leveldb_open(opts, name, &errptr)
	defer C.free(unsafe.Pointer(db))

	fmt.Printf("%v\n", db)
	fmt.Printf("%+v\n", opts)
	C.free(unsafe.Pointer(opts))

	entries := map[string]string{
		"A": "1",
		"B": "2",
		"C": "3",
	}
	errptr_two := C.CString("")
	woptions := C.leveldb_writeoptions_create()
	for k, v := range entries {
		C.leveldb_put(
			db, woptions, C.CString(k), C.size_t(len(k)),
			C.CString(v), C.size_t(len(v)), &errptr_two,
		)

	}

	options := C.leveldb_readoptions_create()
	for k, expected := range entries {
		lengthOfValue := new(C.size_t)
		result := C.leveldb_get(
			db, options, C.CString(k), C.ulong(len(k)), lengthOfValue, &errptr,
		)
		size := int(*lengthOfValue)
		fmt.Printf("expected: %s, got: %v for key: %s\n", expected, C.GoString(result)[:size], k)
	}

	// call it a night
	// zzzzz
}
