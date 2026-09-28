package myhex

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		// Normal cases
		{"single byte", []byte{0x00}, []byte("00")},
		{"single byte 0xFF", []byte{0xFF}, []byte("ff")},
		{"single byte 0x0F", []byte{0x0F}, []byte("0f")},
		{"single byte 0xF0", []byte{0xF0}, []byte("f0")},
		{"multiple bytes", []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F}, []byte("48656c6c6f")},
		{"all nibbles", []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}, []byte("0123456789abcdef")},

		// Edge cases
		{"empty input", []byte{}, []byte{}},
		{"max byte value", []byte{0xFF, 0xFF, 0xFF}, []byte("ffffff")},
		{"min byte value", []byte{0x00, 0x00, 0x00}, []byte("000000")},
		{"alternating high low", []byte{0xA5, 0x5A}, []byte("a55a")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Encode(tt.input)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("Encode(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
		valid    bool
	}{
		// Normal cases
		{"single byte lower", []byte("00"), []byte{0x00}, true},
		{"single byte upper", []byte("FF"), []byte{0xFF}, true},
		{"single byte mixed", []byte("aB"), []byte{0xAB}, true},
		{"multiple bytes lower", []byte("48656c6c6f"), []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F}, true},
		{"multiple bytes upper", []byte("48656C6C6F"), []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F}, true},
		{"full range", []byte("0123456789abcdef"), []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}, true},

		// Edge cases
		{"empty input", []byte{}, []byte{}, true},
		{"all zeros", []byte("000000"), []byte{0x00, 0x00, 0x00}, true},
		{"all ones", []byte("111111"), []byte{0x11, 0x11, 0x11}, true},
		{"all Fs", []byte("ffffff"), []byte{0xFF, 0xFF, 0xFF}, true},

		// Invalid cases - odd length
		{"odd length 1", []byte("a"), nil, false},
		{"odd length 3", []byte("abc"), nil, false},
		{"odd length 5", []byte("abcde"), nil, false},

		// Invalid cases - invalid characters
		{"invalid char lower g", []byte("0g"), nil, false},
		{"invalid char upper G", []byte("0G"), nil, false},
		{"invalid char g0", []byte("g0"), nil, false},
		{"invalid char G0", []byte("G0"), nil, false},
		{"invalid char gh", []byte("gh"), nil, false},
		{"invalid char GH", []byte("GH"), nil, false},
		{"invalid char space", []byte("48 65"), nil, false},
		{"invalid char space first", []byte(" 0"), nil, false},
		{"invalid char space second", []byte("0 "), nil, false},
		{"invalid char newline", []byte("48\n65"), nil, false},
		{"invalid char null byte", []byte("48\x0065"), nil, false},
		{"invalid char null first", []byte("\x000"), nil, false},
		{"invalid char tab", []byte("0\t"), nil, false},
		{"invalid char tab first", []byte("\t0"), nil, false},
		{"invalid char dash", []byte("0-"), nil, false},
		{"invalid char dash first", []byte("-0"), nil, false},
		{"invalid char underscore", []byte("0_"), nil, false},
		{"invalid char underscore first", []byte("_0"), nil, false},
		{"invalid char dot", []byte("0."), nil, false},
		{"invalid dot char", []byte(".0"), nil, false},
		{"invalid char z", []byte("z0"), nil, false},
		{"invalid char x first", []byte("x0"), nil, false},
		{"invalid char x second", []byte("0x"), nil, false},
		{"invalid char in middle", []byte("48656c6c6g"), nil, false},
		{"invalid char in middle upper", []byte("48656C6C6G"), nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := Decode(tt.input)
			if ok != tt.valid {
				t.Errorf("Decode(%v) validity = %v, want %v", tt.input, ok, tt.valid)
			}
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("Decode(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	// Test that encoding then decoding returns the original data
	inputs := [][]byte{
		nil,
		{},
		{0x00},
		{0xFF},
		{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF},
		{0xDE, 0xAD, 0xBE, 0xEF},
		[]byte("Hello, World!"),
		bytes.Repeat([]byte{0xAA}, 100),
		bytes.Repeat([]byte{0x55}, 1000),
	}

	for _, input := range inputs {
		encoded := Encode(input)
		decoded, ok := Decode(encoded)
		if !ok {
			t.Errorf("Round trip failed: Decode(Encode(%v)) returned invalid", input)
		}
		if !bytes.Equal(decoded, input) {
			t.Errorf("Round trip failed: Decode(Encode(%v)) = %v, want %v", input, decoded, input)
		}
	}
}

func TestDecodeEncodeRoundTrip(t *testing.T) {
	// Test that decoding then encoding returns the original hex (case may differ)
	// Note: Decode accepts both upper and lower case, Encode always produces lower case
	inputs := [][]byte{
		{},
		[]byte("00"),
		[]byte("FF"),
		[]byte("ff"),
		[]byte("Ff"),
		[]byte("48656c6c6f"),
		[]byte("48656C6C6F"),
		[]byte("0123456789ABCDEF"),
		[]byte("abcdefABCDEF"),
	}

	for _, input := range inputs {
		// Only test valid hex strings
		if len(input)%2 != 0 {
			continue
		}
		// Check if input is valid hex
		valid := true
		for _, b := range input {
			if _, ok := hexSearch[b]; !ok {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}

		decoded, ok := Decode(input)
		if !ok {
			t.Errorf("Decode(%v) returned invalid for valid hex", input)
			continue
		}
		reencoded := Encode(decoded)
		// Compare lowercase versions since Encode always produces lowercase
		if !bytes.Equal(reencoded, bytes.ToLower(input)) {
			t.Errorf("Encode(Decode(%v)) = %v, want %v", input, reencoded, bytes.ToLower(input))
		}
	}
}

func TestEncodeProperties(t *testing.T) {
	// Test that Encode produces the correct length output
	tests := []struct {
		inputLen    int
		expectedLen int
	}{
		{0, 0},
		{1, 2},
		{5, 10},
		{100, 200},
		{1000, 2000},
	}

	for _, tt := range tests {
		input := make([]byte, tt.inputLen)
		result := Encode(input)
		if len(result) != tt.expectedLen {
			t.Errorf("Encode length for input len %d = %d, want %d", tt.inputLen, len(result), tt.expectedLen)
		}
	}
}

func TestDecodeProperties(t *testing.T) {
	// Test that Decode produces the correct length output
	// Only valid even-length hex strings
	tests := []struct {
		inputLen    int
		expectedLen int
	}{
		{0, 0},
		{2, 1},
		{10, 5},
		{200, 100},
		{2000, 1000},
	}

	for _, tt := range tests {
		// Create a valid hex string of the specified length
		input := make([]byte, tt.inputLen)
		for i := 0; i < tt.inputLen; i += 2 {
			input[i] = 'a'
			input[i+1] = 'b'
		}
		result, ok := Decode(input)
		if !ok {
			t.Errorf("Decode(%v) returned invalid for valid hex", input)
			continue
		}
		if len(result) != tt.expectedLen {
			t.Errorf("Decode length for input len %d = %d, want %d", tt.inputLen, len(result), tt.expectedLen)
		}
	}
}
