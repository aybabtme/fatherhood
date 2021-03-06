fatherhood
==========
[![Build Status](https://drone.io/github.com/aybabtme/fatherhood/status.png)](https://drone.io/github.com/aybabtme/fatherhood/latest)

`fatherhood` is a JSON stream decoding library wrapping the
[`megajson`](https://github.com/benbjohnson/megajson) scanner.

It offers a very ugly API in exchange for speed and no code generation.

## Performance

Its speed is equivalent to [`megajson`](https://github.com/benbjohnson/megajson),
since it uses the same scanner.  All kudos to [Ben Johnson](https://github.com/benbjohnson), not me.

|    package    |      ns/op |  MB/s |
|:-------------:|:----------:|:-----:|
|  fatherhood   | 52'156'933 | 37.20 |
|   megajson    | 53'557'744 | 36.23 |
| encoding/json | 98'061'899 | 19.79 |


## Godocs?

[Godoc](http://godoc.org/github.com/aybabtme/fatherhood)!

## Usage

The general idea of `fatherhood` goes like this:

* get a decoder.
* iterate over the values you need.
* extract them manually.

Getting a decoder:
```go
dec := fatherhood.NewDecoder(r)
```

Then you should know what's in the stream you are reading:
```go
err := dec.EachMember(&obj, objVisitor) // decodes objects
err := dec.EachValue(&arr, arrVisitor)  // decodes arrays
err := dec.ReadTypeX(&typeX)            // decodes strings, bool, floats, ints, etc
```

When you decode an object, you must provide a visitor `func` that will be invoked for each member of the object:

```go
obj := &objType{} // make sure obj is not a nil pointer!
err := dec.EachMember(&obj, objVisitor)

func objVisitor(dec *Decoder, o interface{}, member string) error {
  obj := o.(*objType)
  switch member {
  case "key1":
    return dec.ReadInt(obj.Key1)
    // and so on
  }
  // You MUST discard members that you choose not to decode.
  return dec.Discard()
}
```

Similarly, when you decode an array, you must provide a visitor `func` that will be invoked for each element in the array:

```go
arr := make([]arrType, 0) // make sure arr is not nil!
err := dec.EachValue(&arr, objVisitor)

func arrVisitor(dec *fatherhood.Decoder, a interface{}, t fatherhood.JSONType) error {
  arr := a.(*[]arrType)
  switch t {
  case fatherhood.Object:
    obj := &objType{}
    err := dec.EachMember(obj, decodeNode)
    *arr = append(*arr, obj)
    return err
  }
  // You MUST discard values that you choose not to decode.
  return dec.Discard()
}
```

__All this pointer stuff is tricky and easy to mess up.  Make sure you test
your things carefully.__

## Example

Extracted from the benchmark code. Be cautious, messed up code ahead:

```go
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
    // or dec.Discard() to carry on
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
```

### Details

Performance values taken from running benchmarks on the following revisions:

* `megajson`, on `master`:

```
benbjohnson/megajson $ git rev-parse HEAD
533c329f8535e121708a0ee08ea53bda5edfbe79
```

* `fatherhood`, on `master`:

```
aybabtme/fatherhood $ git rev-parse HEAD
5cfd87c089e3829a28c9cfcd8993370bf787ffa1
```

* `encoding/json`, on `release-branch.go1.2`:

```
encoding/json $ hg summary
parent: 18712:0ddbdc3c7ce2 go1.2.1 release
 go1.2.1
branch: release-branch.go1.2
```

## Why use this...

### Instead of `megajson`?

`megajson` uses code generation to create decoders/encoders for your types,
`fatherhood` doesn't.

Some combinations of types aren't working in megajson but they work with
`fatherhood`. For instance, I wrote `fatherhood` because `megajson` didn't
decode objects containing `[]uint64`.

### Instead of `encoding/json`?

`fatherhood` is faster than the standard library decoder.

## Why use `megajson` instead of `fatherhood`?

`megajson` offers an encoder, while `fatherhood` only decodes.

Regarding code generation, `megajson` gives you drop in codecs.
Meanwhile, the `fatherhood` API is ugly and painful to use.


## Why use `encoding/json` instead of `fatherhood`?

The standard library offers a much nicer API. You should always prefer
`encoding/json` to `fatherhood` unless JSON decoding speed becomes a problem.
