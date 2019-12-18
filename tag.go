// The motivation behind this package is that the StructTag implementation shipped
// with Go's standard library is very limited in detecting a malformed StructTag
// and each time StructTag.Get(key) gets called, it results in the StructTag
// being parsed agian. Another problem is that the StructTag can not be
// easily manipulated because it is basically a string.
// This package provides a way to parse the StructTag into a Tag, which
// allows for fast lookups and easy manipulation of key value pairs within the
// a Tag.
//
// 	// Example of struct using tags to append metadata to fields.
// 	type Server struct {
//		Host string `json:"host" env:"SERVER_HOST" default:"localhost"`
//		Port int    `json:"port" env:"SERVER_PORT" default:"3000"`
//	}
package tag

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrInvalidSyntax    = errors.New("invalid syntax for key value pair")
	ErrInvalidKey       = errors.New("invalid key")
	ErrInvalidValue     = errors.New("invalid value, missing qoutes around value")
	ErrInvalidSeperator = errors.New("invalid seperator, key value pairs should be seperated by spaces")
)

// Tag is just a map of key value pairs.
type Tag map[string]string

// Merge multiple tags together into a single Tag.
// In case of duplicate keys, the last encountered key will overwrite the any existing.
func Merge(tags ...Tag) Tag {
	for _, t := range tags {
		for k, v := range t {
			tags[0][k] = v
		}
	}

	return tags[0]
}

// StructTag converts the Tag into a StructTag.
func (t Tag) StructTag() reflect.StructTag {
	var s string
	for k, v := range t {
		s += fmt.Sprintf(`%s:"%s" `, k, v)
	}
	return reflect.StructTag(strings.TrimSpace(s))
}

// Parse takes a StructTag and parses it into a Tag or returns an error.
// If the given string contains duplicate key value pairs the last pair
// will overwrite the previous in the map.
//
// The parsing logic is a slightly modified version of the StructTag.Lookup
// function from the reflect package included in the standard library.
// https://github.com/golang/go/blob/0377f061687771eddfe8de78d6c40e17d6b21a39/src/reflect/type.go#L1132
func Parse(st reflect.StructTag) (Tag, error) {
	tag := Tag{}

	for st != "" {
		i := 0
		for i < len(st) && st[i] == ' ' {
			i++
		}

		st = st[i:]
		if st == "" {
			break
		}

		i = 0
		for i < len(st) && st[i] > ' ' && st[i] != ':' && st[i] != '"' && st[i] != 0x7f {
			if st[i] == ',' {
				return tag, ErrInvalidSeperator
			}
			i++
		}

		if i == 0 {
			return tag, ErrInvalidKey
		}

		if i+1 >= len(st) || st[i] != ':' {
			return tag, ErrInvalidSyntax
		}

		if st[i+1] != '"' {
			return tag, ErrInvalidValue
		}

		key := string(st[:i])
		st = st[i+1:]

		i = 1
		for i < len(st) && st[i] != '"' {
			if st[i] == '\\' {
				i++
			}
			i++
		}

		if i >= len(st) {
			return tag, ErrInvalidValue
		}

		qvalue := string(st[:i+1])
		st = st[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			return tag, ErrInvalidValue
		}

		tag[key] = value
	}

	return tag, nil
}
