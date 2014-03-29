package fatherhood_test

import (
	"bytes"
	"fmt"
	"github.com/aybabtme/fatherhood"
	"log"
	"reflect"
	"testing"
)

//////////////////////
// Simple test

type object struct {
	Class      string
	References []string
	Line       uint64
	Fd         int
	Shared     bool
	Flags      flag
}

type flag struct {
	WbProtected bool
	Old         int
	Marked      string
}

const smallJSON = `{
    "class": "array",
    "line": 1662,
    "fd": -1,
    "shared": true,
    "flags": {
        "wbprotected": false,
        "old": -765,
        "marked": "probably"
    },
    "references": [
        "hello",
        "bye",
        "lollll"
    ]
}`

var want = object{
	Class:      "array",
	References: []string{"hello", "bye", "lollll"},
	Line:       1662,
	Fd:         -1,
	Shared:     true,
	Flags: flag{
		WbProtected: false,
		Old:         -765,
		Marked:      "probably",
	},
}

func TestParseSmallJSON(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	bf := bytes.NewBufferString(smallJSON)

	var got object

	err := fatherhood.NewDecoder(bf).EachMember(&got, func(dec *fatherhood.Decoder, dst interface{}, member string) error {
		obj := dst.(*object)
		switch member {
		case "class":
			return dec.ReadString(&obj.Class)
		case "line":
			return dec.ReadUint64(&obj.Line)
		case "fd":
			return dec.ReadInt(&obj.Fd)
		case "shared":
			return dec.ReadBool(&obj.Shared)
		case "references":
			return dec.EachValue(&obj.References, func(dec *fatherhood.Decoder, dst interface{}, t fatherhood.JSONType) error {
				arr := dst.(*[]string)
				if t != fatherhood.String {
					return fmt.Errorf("unexpected type, %d", t)
				}
				// can't pass pointer to naked string, need the string to be in a struct
				var ref struct {
					val string
				}
				if err := dec.ReadString(&ref.val); err != nil {
					return fmt.Errorf("reading reference string, %v", err)
				}
				*arr = append(*arr, ref.val)
				return nil
			})
		case "flags":
			return dec.EachMember(&obj.Flags, func(dec *fatherhood.Decoder, dst interface{}, flagmember string) error {
				innerObj := dst.(*flag)
				switch flagmember {
				case "wbprotected":
					return dec.ReadBool(&innerObj.WbProtected)
				case "old":
					return dec.ReadInt(&innerObj.Old)
				case "marked":
					return dec.ReadString(&innerObj.Marked)
				}
				return nil
			})
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v got %v", want, got)
	}
}
