package main

/*
#cgo CFLAGS:  -I/Users/ggilmore/dev/go/src/github.com/ggilmore/csi/src/classes/advanced-programing/cgo/prework/level-db/leveldb/include/
#cgo LDFLAGS: -L/Users/ggilmore/dev/go/src/github.com/ggilmore/csi/src/classes/advanced-programing/cgo/prework/level-db/leveldb/build -l leveldb

#include "leveldb/c.h"
#include <stdlib.h>
// leveldb_t* get_leveldb() { return calloc(1, sizeof(int)); }
int dumb_func() {
	leveldb_options_t _foo = leveldb_options_create();
	return 1;
}
*/
import "C"

import "fmt"

func main() {
	// var db *C.leveldb_t
	// opts := C.leveldb_options_create()
	// db = C.get_leveldb()
	// fmt.Printf("%v\n", opts)
	// fmt.Printf("%v\n", opts)
	x := C.dumb_func()
	fmt.Printf("x=%d\n", x)
}
