package tag_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kvartborg/tag"
)

func Example() {
	tags, err := tag.Parse(`json:"host" env:"SERVER_HOST" default:"localhost"`)
	fmt.Println(tags, err)
	// Output: map[default:localhost env:SERVER_HOST json:host] <nil>
}

func ExampleMerge() {
	t1 := tag.Tag{
		"env": "TESTING",
	}

	t2 := tag.Tag{
		"env": "HELLO",
	}

	tags := tag.Merge(t1, t2)
	fmt.Println(tags)
	// Output: map[env:HELLO]
}

func ExampleParse() {
	tags, err := tag.Parse(`env:"SERVER_HOST" default:"localhost"`)
	fmt.Println(tags, err)
	// Output: map[default:localhost env:SERVER_HOST] <nil>
}

func ExampleStructTag() {
	t, _ := tag.Parse(`env:"SERVER_HOST"`)
	st := t.StructTag()
	fmt.Println(st)
	// Output: env:"SERVER_HOST"
}

func TestParse(t *testing.T) {
	tags, _ := tag.Parse(`env:"SERVER_HOST" default:"localhost"`)

	if v, ok := tags["env"]; !ok || v != "SERVER_HOST" {
		t.Errorf("the key env does not SERVER_HOST; got %s", v)
	}

	if v, ok := tags["default"]; !ok || v != "localhost" {
		t.Errorf("the key default does not localhost; got %s", v)
	}
}

func TestParseErrors(t *testing.T) {
	_, err := tag.Parse(`invalid syntax`)
	if err != tag.ErrInvalidSyntax {
		t.Errorf("did not return error %s; got %s", tag.ErrInvalidSyntax, err)
	}

	_, err = tag.Parse(`:"value"`)
	if err != tag.ErrInvalidKey {
		t.Errorf("did not return error %s; got %s", tag.ErrInvalidKey, err)
	}

	cases := []reflect.StructTag{
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
	if err != tag.ErrInvalidSeperator {
		t.Errorf("did not return error %s; got %s", tag.ErrInvalidSeperator, err)
	}
}
