package parser

import (
	"reflect"
	"testing"

	"github.com/QingGo/teeworlds-master-cache/datatype"
)

func TestParseServerInfo(t *testing.T) {
	var parseTests = []struct {
		in       []byte                // input
		expected []datatype.ServerAddr // expected result
	}{
		{in: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 108, 105, 115, 50, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 164, 132, 46, 180, 32, 112, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 164, 132, 46, 180, 32, 111},
			expected: []datatype.ServerAddr{
				{
					IP:   "164.132.46.180",
					Port: 8304,
				},
				{
					IP:   "164.132.46.180",
					Port: 8303,
				},
			},
		},
	}
	for _, tt := range parseTests {
		actual := ParseServerInfo(tt.in)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("parseServerInfo(%+v) = %+v; expected %+v", tt.in, actual, tt.expected)
		}
	}
}

func TestParseIPListToBytes(t *testing.T) {
	var parseTests = []struct {
		in       []datatype.ServerAddr // input
		expected [][]byte              // expected result
	}{
		{
			in: []datatype.ServerAddr{
				{
					IP:   "164.132.46.180",
					Port: 8304,
				},
				{
					IP:   "164.132.46.180",
					Port: 8303,
				},
			},
			expected: [][]byte{{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 108, 105, 115, 50, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 164, 132, 46, 180, 32, 112, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 164, 132, 46, 180, 32, 111}},
		},
	}
	for _, tt := range parseTests {
		actual := ParseIPListToBytes(tt.in)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("parseServerInfo(%+v) = %+v; expected %+v", tt.in, actual, tt.expected)
		}
	}
}
