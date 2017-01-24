package helper

import (
	"net/http"
	"testing"
)

func TestPaginationFromRequest(t *testing.T) {

	type testpair struct {
		url            string
		expectedOffset int
		expectedRows   int
	}

	// Test cases with disered output
	var tests = []testpair{
		{"http://localhost/", 1, 10},
		{"http://localhost/?offset=3", 3, 10},
		{"http://localhost/?offset=5&rows=4", 5, 4},
		{"http://localhost/?offset=abc&rows=12", 1, 12},
		{"http://localhost/?offset=-1", 1, 10},
		{"http://localhost/?rows=3", 1, 3},
		{"http://localhost/?offset=4&rows=abc", 4, 10},
		{"http://localhost/?rows=-1", 1, 10},
		{"", 1, 10},
	}

	// Run all tests
	for _, test := range tests {

		req, err := http.NewRequest("GET", test.url, nil)
		if err != nil {
			t.Fatal(err)
		}

		actualOffset, actualRows := PaginationFromRequest(req)

		if actualOffset != test.expectedOffset {
			t.Errorf("Expected %v but got %v", test.expectedOffset, actualOffset)
		}
		if actualRows != test.expectedRows {
			t.Errorf("Expected %v but got %v", test.expectedRows, actualRows)
		}
	}

}
