// +build !appengine

package optimus

import

// "log"

"testing"

func TestOptimus(t *testing.T) {
	var encodeTests = []struct {
		prime  uint64
		mod    uint64
		rand   uint64
		input  uint64
		encode uint64
	}{
		{1580030173, 59260789, 1163945558, 15, 1103647397},
		{1580030173, 59260789, 1163945558, 99999, 1458223381},

		{2123809381, 1885413229, 146808189, 15, 1645575830},
		{2123809381, 1885413229, 146808189, 99999, 1124632518},

		{3, 715827883, 1234567890, 15, 1234567935},
		{3, 715827883, 1234567890, 99999, 1234342159},

		{837350711, 1701236871, 1234567890, 1024, 1782767314},
		{837350711, 1701236871, 1234567890, 1025, 472634341},
	}

	for i, tt := range encodeTests {
		o := New(tt.prime, tt.mod, tt.rand)

		enc := o.Encode(tt.input)
		if enc != tt.encode {
			t.Errorf("[%d] encode failed: expected %d, got %d", i, tt.encode, enc)
		}

		dec := o.Decode(enc)
		if dec != tt.input {
			t.Errorf("[%d] decode failed: expected %d, got %d", i, tt.input, dec)
		}
	}
}
