package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenumber(t *testing.T) {
	tests := map[string]struct {
		input  string
		result string
	}{
		"un-numbered list": {
			input: `one
two
three
`,
			result: `1. one
2. two
3. three
`,
		},
		"mis-numbered list": {
			input: `3. one
3. two
1. three
`,
			result: `1. one
2. two
3. three
`,
		},
		"number in list entry text": {
			input: `1. one
f3. two
3. three
`,
			result: `1. one
2. f3. two
3. three
`,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var res bytes.Buffer
			r := strings.NewReader(test.input)

			// do the renumbering itself
			renumber(&res, r)

			// compare the results
			// convert from []byte
			got := string(res.Bytes())
			expected := test.result
			if got != expected {
				t.Fatalf("input:\n%v\ngot:\n%v\nwant:\n%v\n", test.input, got, expected)
			}
		})
	}
}
