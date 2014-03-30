fatherhood
==========
[![Build Status](https://drone.io/github.com/aybabtme/fatherhood/status.png)](https://drone.io/github.com/aybabtme/fatherhood/latest)

fatherhood is a JSON stream decoding library wrapping the scanner of
[`megajson`](https://github.com/benbjohnson/megajson).


Why use this and not megajson?
==============================

megajson uses code generation to create decoders/encoders for your types.
fatherhood doesn't.

Some combinaisons of types aren't working in megajson. They work with
fatherhood. For instance, I wrote fatherhood because megajson didn't

Why use megajson and not this?
==============================

megajson offers an encoder, fatherhood only decodes.

Aside for the code generation thing, megajson gives you drop in codecs.
Meanwhile, fatherhood's API is fugly and a pain to use.

Performance
===========

Speed is equivalent to [`megajson`](github.com/benbjohnson/megajson), since it
uses the same scanner.  All kudos to [Ben Johnson](github.com/benbjohnson),
not me.

|    package    |      ns/op |  MB/s |
|:-------------:|:----------:|:-----:|
|   megajson    | 53'557'744 | 36.23 |
|  fatherhood   | 53'717'435 | 36.12 |
| encoding/json | 98'061'899 | 19.79 |


Benchmark from the following revisions:

* megajson, on `master`:

```
benbjohnson/megajson $ git rev-parse HEAD
533c329f8535e121708a0ee08ea53bda5edfbe79
```

* fatherhood, on `master`:

```
aybabtme/fatherhood $ git rev-parse HEAD
5cfd87c089e3829a28c9cfcd8993370bf787ffa1
```

* encoding/json, on `release-branch.go1.2`:

```
encoding/json $ hg summary
parent: 18712:0ddbdc3c7ce2 go1.2.1 release
 go1.2.1
branch: release-branch.go1.2
```
