package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Data types
const (
	NONE   = 0
	BOOL   = 1
	DOUBLE = 2
	STRING = 3
	LIST   = 4
	DICT   = 5
)

type FModData struct {
	Version FVersion
	Data    interface{}
}

func (d *FModData) Decode(r io.Reader) {
	fr := FactorioReader{Reader: r}
	d.Version.Read(fr)
	d.Data = fr.Tree()
}

type FVersion struct {
	Major, Minor, Path, Dev uint16
}

func (ver *FVersion) Read(r FactorioReader) {
	r.Val(&ver.Major)
	r.Val(&ver.Minor)
	r.Val(&ver.Path)
	r.Val(&ver.Dev)
}

var le = binary.LittleEndian

type FactorioReader struct {
	io.Reader
}

func (r FactorioReader) Val(val interface{}) {
	err := binary.Read(r.Reader, le, val)
	if err != nil {
		panic(err)
	}
}

func (r FactorioReader) Bytes(len uint) []byte {
	ret := make([]byte, len)
	_, err := io.ReadFull(r.Reader, ret)
	if err != nil {
		panic(err)
	}
	return ret
}

func (r FactorioReader) Tree() interface{} {
	// Type of embedded data
	var kind byte
	r.Val(&kind)
	if r.Bool() {
		// No idea how to deal with it, but it should not appear in mod settings anyways
		panic("Legacy any-type flag detected")
	}
	switch kind {
	case BOOL:
		return r.Bool()

	case DOUBLE:
		var d float64
		r.Val(&d)
		return d

	case STRING:
		return r.String()

	case LIST:
		return r.List()

	case DICT:
		return r.Dict()

	default:
		panic(fmt.Sprintf("Unknown type %v", kind))
	}
}

func (r FactorioReader) List() []interface{} {
	var len uint32
	r.Val(&len)
	list := make([]interface{}, len, len)
	for i := uint32(0); i < len; i++ {
		list[i] = r.Tree()
	}
	return list
}

func (r FactorioReader) Dict() map[string]interface{} {
	var len uint32
	r.Val(&len)
	dict := make(map[string]interface{})
	for i := uint32(0); i < len; i++ {
		key := r.String()
		dict[key] = r.Tree()
	}
	return dict
}

func (r FactorioReader) Bool() bool {
	var b byte
	r.Val(&b)
	return b != 0
}

// Optimized uint https://wiki.factorio.com/Data_types#Space_Optimized
func (r FactorioReader) OUint() uint32 {
	var opt byte
	r.Val(&opt)
	if opt < 255 {
		return uint32(opt)
	}
	var ret uint32
	r.Val(&ret)
	return ret
}

func (r FactorioReader) String() string {
	// Whole byte for empty flag and optimized length after it. Makes sence %)
	empty := r.Bool()
	if empty {
		return ""
	}
	len := r.OUint()
	bytes := r.Bytes(uint(len))
	return string(bytes)
}
