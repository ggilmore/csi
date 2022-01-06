package main

type Record struct {
	rType    uint16
	class    uint16
	ttl      int32
	rdLength uint16
	rData    string
}

type ResourceType int

const (
	ResourceTypeA ResourceType = iota + 1
	ResourceTypeNS
	ResourceTypeMD
	ResourceTypeMF
	ResourceTypeCNAME
	ResourceTypeSOA
	ResourceTypeMB
	ResourceTypeMG
	ResourceTypeMR
	ResourceTypeNULL
	ResourceTypeWKS
	ResourceTypePTR
	ResourceTypeHINFO
	ResourceTypeMINFO
	ResourceTypeMX
	ResourceTypeTXT
)
