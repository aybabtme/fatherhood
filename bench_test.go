// Taken and adapted from the encoding/json package of
// Go's standard library.

// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Large data benchmark.
// The JSON data is a summary of agl's changes in the
// go, webkit, and chromium open source projects.
// We benchmark converting between the JSON form
// and in-memory data structures.

package fatherhood

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

type codeResponse struct {
	Tree     *codeNode `json:"tree"`
	Username string    `json:"username"`
}

type codeNode struct {
	Name     string      `json:"name"`
	Kids     []*codeNode `json:"kids"`
	CLWeight float64     `json:"cl_weight"`
	Touches  int         `json:"touches"`
	MinT     int64       `json:"min_t"`
	MaxT     int64       `json:"max_t"`
	MeanT    int64       `json:"mean_t"`
}

var codeJSON []byte
var codeReadJSON io.Reader
var codeStruct codeResponse

func Unmarshal(data []byte, code *codeResponse) error {

	read := bytes.NewReader(data)

	var (
		decodeNodeArr  func(*Decoder, interface{}, JSONType) error
		decodeNode     func(*Decoder, interface{}, string) error
		decodeResponse func(*Decoder, interface{}, string) error
	)

	decodeResponse = func(dec *Decoder, r interface{}, member string) error {
		resp := r.(*codeResponse)
		switch member {
		case "username":
			return dec.ReadString(&resp.Username)
		case "tree":
			resp.Tree = &codeNode{}
			return dec.EachMember(resp.Tree, decodeNode)
		}
		return fmt.Errorf("unsupported member %s", member)
	}

	decodeNode = func(dec *Decoder, n interface{}, member string) error {
		node := n.(*codeNode)
		switch member {
		case "name":
			return dec.ReadString(&node.Name)
		case "cl_weight":
			return dec.ReadFloat64(&node.CLWeight)
		case "touches":
			return dec.ReadInt(&node.Touches)
		case "min_t":
			return dec.ReadInt64(&node.MinT)
		case "max_t":
			return dec.ReadInt64(&node.MaxT)
		case "mean_t":
			return dec.ReadInt64(&node.MeanT)
		case "kids":
			node.Kids = make([]*codeNode, 0)
			return dec.EachValue(&node.Kids, decodeNodeArr)
		}
		return fmt.Errorf("unsupported member %s", member)
	}

	decodeNodeArr = func(dec *Decoder, a interface{}, t JSONType) error {
		arr := a.(*[]*codeNode)
		switch t {
		case Object:
			node := &codeNode{}
			err := dec.EachMember(node, decodeNode)
			*arr = append(*arr, node)
			return err
		}
		return fmt.Errorf("unsupported type %#v", t)
	}

	return NewDecoder(read).EachMember(code, decodeResponse)
}

func codeInit() {
	f, err := os.Open("testdata/code.json.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(gz)
	if err != nil {
		panic(err)
	}

	codeJSON = data

	if err := Unmarshal(codeJSON, &codeStruct); err != nil {
		panic(err)
	}

	// Encode it back with stdlib
	if data, err = json.Marshal(&codeStruct); err != nil {
		panic("marshal code.json: " + err.Error())
	}

	if !bytes.Equal(data, codeJSON) {
		println("different lengths", len(data), len(codeJSON))
		for i := 0; i < len(data) && i < len(codeJSON); i++ {
			if data[i] != codeJSON[i] {
				println("re-marshal: changed at byte", i)
				println("orig: ", string(codeJSON[i-10:i+10]))
				println("new: ", string(data[i-10:i+10]))
				break
			}
		}
		panic("re-marshal code.json: different result")
	}
}

func BenchmarkCodeUnmarshal(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	for i := 0; i < b.N; i++ {
		var r codeResponse
		if err := Unmarshal(codeJSON, &r); err != nil {
			b.Fatal("Unmmarshal:", err)
		}
	}
	b.SetBytes(int64(len(codeJSON)))
}

func BenchmarkCodeUnmarshalReuse(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	var r codeResponse
	for i := 0; i < b.N; i++ {
		if err := Unmarshal(codeJSON, &r); err != nil {
			b.Fatal("Unmmarshal:", err)
		}
	}
}
