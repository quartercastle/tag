# tag

[![Version](https://img.shields.io/github/release/kvartborg/tag.svg)](https://github.com/kvartborg/tag/releases)
[![Build Status](https://travis-ci.org/kvartborg/tag.svg?branch=master)](https://travis-ci.org/kvartborg/tag)
[![GoDoc](https://godoc.org/github.com/kvartborg/tag?status.svg)](https://godoc.org/github.com/kvartborg/tag)
[![Go Report Card](https://goreportcard.com/badge/github.com/kvartborg/tag)](https://goreportcard.com/report/github.com/kvartborg/tag)


The motivation behind this package is that the [`StructTag`](https://github.com/golang/go/blob/0377f061687771eddfe8de78d6c40e17d6b21a39/src/reflect/type.go#L1110)
implementation shipped with Go's standard library is very limited in
detecting a malformed StructTag and each time `StructTag.Get(key)` gets called,
it results in the `StructTag` being parsed agian. Another problem is that the
`StructTag` can not be easily manipulated because it is basically a string.
This package provides a way to parse the `StructTag` into a `Tag` map. This
allows fast lookups and easy manipulation of the key value pairs within the
`Tag`.

```go
// Example of struct using tags to append metadata to fields.
type Server struct {
	Host string `json:"host" env:"SERVER_HOST" default:"localhost"`
	Port int    `json:"port" env:"SERVER_PORT" default:"3000"`
}
```

### Install
```
go get github.com/kvartborg/tag
```

### Usage
```go
t, err := tag.Parse(`json:"host" env:"SERVER_HOST"`)

if err != nil {
  panic(err)
}

fmt.Println(t["json"])
```
See [godoc](https://godoc.org/github.com/kvartborg/tag) for full documentation.

### License
This project is licensed under the [MIT License](https://github.com/kvartborg/tag/blob/master/LICENSE).
