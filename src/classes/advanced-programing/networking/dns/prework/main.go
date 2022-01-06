package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello, world!")
}

type Header struct {
	// A 16 bit identifier assigned by the program that
	// generates any kind of query.  This identifier is copied
	// the corresponding reply and can be used by the requester
	// to match up replies to outstanding queries.
	ID uint16

	// QR, Opcod, AA,TC, RD, RA
	optsA uint8

	// Z, Recode
	optsB uint8

	// Number of entries in the question section
	QDCount uint16

	// Number of records in the answer question
	ANCount uint16

	// Number of name records in the authority section
	NSCount uint16

	// Number of records in the additional records section
	ARCount uint16
}

func (h Header) RCode() ResponseCode {
	return ResponseCode(h.optsB & 0x7)
}

type ResponseCode int

const (
	ResponseNoError ResponseCode = iota
	ResponseFormatError
	ResponseServerFailure
	ResponseNameError
	ResponseNotImplemented
	ResponseRefused
)

type ResourceRecord struct {
}
