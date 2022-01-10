package main

import "io"

type Record struct {
	rType    uint16
	class    uint16
	TTL      int32
	rdLength uint16
	rData    string
}

func Encode(r Record, w io.Writer) error {
	return nil
}

func (r Record) ResourceType() ResourceType {
	return ResourceType(r.rType)
}

func (r Record) ResourceClass() ResourceClass{
	return ResourceClass(r.class)
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

type ResourceClass int

const (
	ResourceClassIN ResourceClass = iota + 1
	ResourceClassCS
	ResourceClassCH
	ResourceClassHS
)
