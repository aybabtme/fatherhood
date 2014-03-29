package fatherhood

import (
	"errors"
	"fmt"
	"github.com/benbjohnson/megajson/scanner"
	"io"
)

type Decoder struct {
	scan scanner.Scanner
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{scan: scanner.NewScanner(r)}
}

func (d *Decoder) ReadString(target *string) error              { return d.scan.ReadString(target) }
func (d *Decoder) ReadInt(target *int) error                    { return d.scan.ReadInt(target) }
func (d *Decoder) ReadInt64(target *int64) error                { return d.scan.ReadInt64(target) }
func (d *Decoder) ReadUint(target *uint) error                  { return d.scan.ReadUint(target) }
func (d *Decoder) ReadUint64(target *uint64) error              { return d.scan.ReadUint64(target) }
func (d *Decoder) ReadFloat32(target *float32) error            { return d.scan.ReadFloat32(target) }
func (d *Decoder) ReadFloat64(target *float64) error            { return d.scan.ReadFloat64(target) }
func (d *Decoder) ReadBool(target *bool) error                  { return d.scan.ReadBool(target) }
func (d *Decoder) ReadMap(target *map[string]interface{}) error { return d.scan.ReadMap(target) }

// EachMember iterates over the members of an object. Invoke the proper Read
// function to get the value back.
func (d *Decoder) EachMember(doFunc func(string) error) error {

	if tok, tokval, err := d.scan.Scan(); err != nil {
		return err
	} else if tok == scanner.TNULL {
		return nil
	} else if tok != scanner.TLBRACE {
		return fmt.Errorf("unexpected %s at %d: %s; expected '{'", scanner.TokenName(tok), d.scan.Pos(), string(tokval))
	}

	index := 0
	for {
		// Read in key.
		var key string
		tok, tokval, err := d.scan.Scan()
		if err != nil {
			return err
		} else if tok == scanner.TRBRACE {
			return nil
		} else if tok == scanner.TCOMMA {
			if index == 0 {
				return fmt.Errorf("unexpected comma at %d", d.scan.Pos())
			}
			if tok, tokval, err = d.scan.Scan(); err != nil {
				return err
			}
		}

		if tok != scanner.TSTRING {
			return fmt.Errorf("unexpected %s at %d: %s; expected '{' or string", scanner.TokenName(tok), d.scan.Pos(), string(tokval))
		}

		key = string(tokval)

		// Read in the colon.
		if tok, tokval, err := d.scan.Scan(); err != nil {
			return err
		} else if tok != scanner.TCOLON {
			return fmt.Errorf("unexpected %s at %d: %s; expected colon", scanner.TokenName(tok), d.scan.Pos(), string(tokval))
		}

		if err := doFunc(key); err != nil {
			return err
		}
		index++
	}
}

// EachValue iterates over and invokes doFunc at each value of an array. Invoke
// the proper Read function to get the value. The type of the value to read is
// specified by the JSONType given as argument.
func (d *Decoder) EachValue(doFunc func(JSONType) error) error {
	if tok, _, err := d.scan.Scan(); err != nil {
		return err
	} else if tok != scanner.TLBRACKET {
		return errors.New("expected '['")
	}

	// Loop over items.
	index := 0
	for {
		tok, tokval, err := d.scan.Scan()
		if err != nil {
			return err
		} else if tok == scanner.TRBRACKET {
			return nil
		} else if tok == scanner.TCOMMA {
			if index == 0 {
				return fmt.Errorf("unexpected comma in array at %d", d.scan.Pos())
			}
			if tok, tokval, err = d.scan.Scan(); err != nil {
				return err
			}
		}
		d.scan.Unscan(tok, tokval)
		if err := doFunc(toJSONType(tok)); err != nil {
			return err
		}
		index++
	}
}

// JSONType that can be received by
type JSONType uint8

const (
	String = iota
	Number
	Bool
	Null
	Object
	Array
)

func toJSONType(token int) JSONType {
	switch token {
	case scanner.TSTRING:
		return String
	case scanner.TNUMBER:
		return Number
	case scanner.TTRUE:
		return Bool
	case scanner.TFALSE:
		return Bool
	case scanner.TNULL:
		return Null
	case scanner.TLBRACE:
		return Object
	case scanner.TLBRACKET:
		return Array
	}
	panic(fmt.Sprintf("unexpected token type, %s (%d)", scanner.TokenName(token), token))
}
