package main

import "testing"

func TestSanitizeInputIDsList(t *testing.T) {
	// setup:
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},
		{input: "\n", expected: ""},
		{input: "\t\n ", expected: ""},
		{input: "\t\n id \t,id \t\n", expected: "id,id"},
		{input: "\t\n 4084d7a0-bbc2-4bb9-aebe-a1b70643c823 \t,82423bcb-ea16-48c9-bafb-9e77592cabb7 \t\n", expected: "4084d7a0-bbc2-4bb9-aebe-a1b70643c823,82423bcb-ea16-48c9-bafb-9e77592cabb7"},
	}

	for _, tc := range testCases {
		// execute:
		res := sanitizeInputIDsList(tc.input)

		// verify:
		if res != tc.expected {
			t.Fatalf(`expected "%s" but was "%s"`, tc.expected, res)
		}
	}
}
