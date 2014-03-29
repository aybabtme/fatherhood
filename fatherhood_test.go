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

const smallJSON = ` {"class": "array", "references": ["hello", "bye", "lollll"], "line": 1662, "fd": -1, "shared": true, "flags": {"wbprotected": false, "old": -765, "marked": "probably"} } `

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
	dec := fatherhood.NewDecoder(bf)

	err := dec.EachMember(func(member string) error {
		switch member {
		case "class":
			return dec.ReadString(&got.Class)
		case "line":
			return dec.ReadUint64(&got.Line)
		case "fd":
			return dec.ReadInt(&got.Fd)
		case "shared":
			return dec.ReadBool(&got.Shared)
		case "references":
			return dec.EachValue(func(t fatherhood.JSONType) error {
				if t != fatherhood.String {
					return fmt.Errorf("unexpected type, %d", t)
				}
				var ref *string
				if err := dec.ReadString(ref); err != nil {
					return err
				}
				got.References = append(got.References, *ref)
				return nil
			})
		case "flags":
			return dec.EachMember(func(flagmember string) error {
				switch flagmember {
				case "wbprotected":
					return dec.ReadBool(&got.Flags.WbProtected)
				case "old":
					return dec.ReadInt(&got.Flags.Old)
				case "marked":
					return dec.ReadString(&got.Flags.Marked)
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
