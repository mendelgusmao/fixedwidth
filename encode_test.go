package fixedwidth

import (
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type expected struct {
	begin, end int
	value      []byte
}

type line struct {
	I  int       `fixed:"%05d"`
	S1 string    `fixed:"%-10.10s"`
	S2 string    `fixed:"%10.10s"`
	B  bool      `fixed:"%s,n,y"`
	T1 time.Time `fixed:"20060102"`
	T2 time.Time `fixed:"20060102150405"`
	T3 time.Time `fixed:"20060102,zero"`
}

func TestEncode(t *testing.T) {
	tm := time.Date(2017, 05, 31, 15, 16, 17, 0, time.Local)

	data := []struct {
		description   string
		input         interface{}
		expected      []expected
		expectedError error
	}{
		{
			description: "it should marshal an object",
			input: line{
				I:  123,
				S1: "foobar",
				S2: "foobaz",
				B:  true,
				T1: tm,
				T2: tm,
			},
			expected: []expected{
				{0, 5, []byte("00123")},
				{5, 15, []byte("foobar    ")},
				{15, 25, []byte("    foobaz")},
				{25, 26, []byte("y")},
				{26, 34, []byte("20170531")},
				{34, 48, []byte("20170531151617")},
				{48, 56, []byte("00000000")},
			},
		},
	}

	for i, item := range data {
		output, err := Encode(item.input)

		if err != item.expectedError {
			t.Error(err)
		}

		t.Log("output:", string(output))

		for p, piece := range item.expected {
			value := output[piece.begin:piece.end]

			if !reflect.DeepEqual(value, piece.value) {
				t.Logf("failed test piece %d of item %d", p, i)
				t.Log("expected:", spew.Sdump(piece.value))
				t.Log("got:", spew.Sdump(value))
				t.Fatal()
			}

			t.Log("test piece", p, "passed")
		}
	}
}
