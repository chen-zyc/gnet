package http

import (
	"reflect"
	"testing"
)

func TestParseRange(t *testing.T) {
	cases := []struct {
		s              string
		size           int64
		expectedRanges []Range
		expectedErr    error
	}{
		{s: "", size: 10, expectedRanges: nil, expectedErr: nil},
		{s: "0-10", size: 10, expectedRanges: nil, expectedErr: ErrInvalidRange},         // no 'bytes='
		{s: "bytes=0,10", size: 10, expectedRanges: nil, expectedErr: ErrInvalidRange},   // no '-'
		{s: "bytes=-abc", size: 10, expectedRanges: nil, expectedErr: ErrInvalidRange},   // no start, invalid end
		{s: "bytes=-5", size: 10, expectedRanges: []Range{{Start: 5, Length: 5}}},        // no start
		{s: "bytes=-15", size: 10, expectedRanges: []Range{{Start: 0, Length: 10}}},      // no start, end > size
		{s: "bytes=abc-15", size: 10, expectedRanges: nil, expectedErr: ErrInvalidRange}, // invalid start
		{s: "bytes=3-", size: 10, expectedRanges: []Range{{Start: 3, Length: 7}}},        // has start, no end
		{s: "bytes=3-abc", size: 10, expectedRanges: nil, expectedErr: ErrInvalidRange},  // has start, invalid end
		{s: "bytes=3-7", size: 10, expectedRanges: []Range{{Start: 3, Length: 5}}},       // valid
		{s: "bytes=3-17", size: 10, expectedRanges: []Range{{Start: 3, Length: 7}}},      // end > size
		{s: "bytes=", size: 10, expectedRanges: nil},                                     // no ranges
		{s: "bytes=10-15", size: 10, expectedRanges: nil, expectedErr: ErrNoOverlap},     // no overlap
		{s: "bytes=1-3, 5-7", size: 10, expectedRanges: []Range{ // multiple ranges
			{Start: 1, Length: 3},
			{Start: 5, Length: 3},
		}},
		{s: "bytes=1-3, 5-17", size: 10, expectedRanges: []Range{ // multiple ranges
			{Start: 1, Length: 3},
			{Start: 5, Length: 5},
		}},
		{s: "bytes=1-3, 15-17", size: 10, expectedRanges: []Range{ // multiple ranges
			{Start: 1, Length: 3},
		}},
	}

	for _, c := range cases {
		ranges, err := ParseRange(c.s, c.size)
		if c.expectedErr != nil {
			if c.expectedErr != err {
				t.Fatalf("expected err: %v, but got %v", c.expectedErr, err)
			}
			continue
		}
		if err != nil {
			t.Fatalf("expected nil err, but got %v", err)
		}
		if !reflect.DeepEqual(c.expectedRanges, ranges) {
			t.Fatalf("expected %#v, but got %#v", c.expectedRanges, ranges)
		}
	}
}

func TestRange_ContentRange(t *testing.T) {
	cases := []struct {
		r        Range
		size     int64
		expected string
	}{
		{r: Range{Start: 0, Length: 10}, size: 10, expected: "bytes 0-9/10"},
		{r: Range{Start: 0, Length: 15}, size: 10, expected: "bytes 0-14/10"},
		{r: Range{Start: 4, Length: 1}, size: 10, expected: "bytes 4-4/10"},
	}
	for _, c := range cases {
		got := c.r.ContentRange(c.size)
		if got != c.expected {
			t.Fatalf("expected %s, but got %s", c.expected, got)
		}
	}
}
