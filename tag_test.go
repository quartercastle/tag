package tag_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/quartercastle/tag"
)

func Example() {
	t, err := tag.Parse(`json:"host" env:"SERVER_HOST" default:"localhost"`)

	if err != nil {
		panic(err)
	}

	fmt.Println(t["env"])
	// Output: SERVER_HOST
}

func ExampleMerge() {
	t1 := tag.Tag{
		"env": "TESTING",
	}

	t2 := tag.Tag{
		"env": "HELLO",
	}

	t := tag.Merge(t1, t2)
	fmt.Println(t)
	// Output: map[env:HELLO]
}

func ExampleParse() {
	t, err := tag.Parse(`env:"SERVER_HOST" default:"localhost"`)
	fmt.Println(t, err)
	// Output: map[default:localhost env:SERVER_HOST] <nil>
}

func ExampleTag_StructTag() {
	t, _ := tag.Parse(`env:"SERVER_HOST"`)
	st := t.StructTag()
	fmt.Println(st)
	// Output: env:"SERVER_HOST"
}

func TestParse(t *testing.T) {
	tag, _ := tag.Parse(`env:"SERVER_HOST" default:"localhost"`)

	if v, ok := tag["env"]; !ok || v != "SERVER_HOST" {
		t.Errorf("the key env is not equal to SERVER_HOST; got %s", v)
	}

	if v, ok := tag["default"]; !ok || v != "localhost" {
		t.Errorf("the key default is not equal to localhost; got %s", v)
	}
}

func TestParseErrors(t *testing.T) {
	_, err := tag.Parse(`invalid syntax`)
	if err != tag.ErrInvalidSyntax {
		t.Errorf("did not return error %s; got %s", tag.ErrInvalidSyntax, err)
	}

	cases := []reflect.StructTag{
		`:value`,
		`"":"value"`,
	}

	for _, c := range cases {
		_, err = tag.Parse(c)
		if err != tag.ErrInvalidKey {
			t.Errorf("did not return error %s; got %s", tag.ErrInvalidKey, err)
		}
	}

	cases = []reflect.StructTag{
		`key:value`,
		`key:value"`,
		`key:"value`,
		`key:"value\"`,
		`key:\"value"`,
		`key: ""`,
	}

	for _, c := range cases {
		_, err := tag.Parse(c)
		if err != tag.ErrInvalidValue {
			t.Errorf("did not return error %s; got %s", tag.ErrInvalidValue, err)
		}
	}

	_, err = tag.Parse(`key:"value", other:"value"`)
	if err != tag.ErrInvalidSeparator {
		t.Errorf("did not return error %s; got %s", tag.ErrInvalidSeparator, err)
	}
}
