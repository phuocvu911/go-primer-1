package mystrconv

import (
	"math"
	"testing"
)

func TestFormatUint(t *testing.T) {
	tests := []struct {
		name string
		in   uint64
		base int
		want string
		ok   bool
	}{
		{name: "zero in decimal", in: 0, base: 10, want: "0", ok: true},
		{name: "decimal", in: 123456789, base: 10, want: "123456789", ok: true},
		{name: "binary", in: 42, base: 2, want: "101010", ok: true},
		{name: "octal", in: 511, base: 8, want: "777", ok: true},
		{name: "hexadecimal", in: math.MaxUint64, base: 16, want: "ffffffffffffffff", ok: true},
		{name: "base 36", in: 35, base: 36, want: "z", ok: true},
		{name: "base below minimum", in: 10, base: 1, want: "", ok: false},
		{name: "base above maximum", in: 10, base: 37, want: "", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := FormatUint(tt.in, tt.base)
			if got != tt.want || ok != tt.ok {
				t.Errorf("FormatUint(%d, %d) = (%q, %v), want (%q, %v)", tt.in, tt.base, got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestFormatInt(t *testing.T) {
	tests := []struct {
		name string
		in   int64
		base int
		want string
		ok   bool
	}{
		{name: "zero", in: 0, base: 10, want: "0", ok: true},
		{name: "positive decimal", in: 123456789, base: 10, want: "123456789", ok: true},
		{name: "negative decimal", in: -123456789, base: 10, want: "-123456789", ok: true},
		{name: "negative hexadecimal", in: -255, base: 16, want: "-ff", ok: true},
		{name: "minimum int64", in: math.MinInt64, base: 10, want: "-9223372036854775808", ok: true},
		{name: "base below minimum", in: 10, base: 1, want: "", ok: false},
		{name: "base above maximum", in: 10, base: 37, want: "", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := FormatInt(tt.in, tt.base)
			if got != tt.want || ok != tt.ok {
				t.Errorf("FormatInt(%d, %d) = (%q, %v), want (%q, %v)", tt.in, tt.base, got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestParseUint(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		base    int
		bitSize int
		want    uint64
		ok      bool
	}{
		{name: "decimal", input: "123456789", base: 10, bitSize: 64, want: 123456789, ok: true},
		{name: "binary", input: "101010", base: 2, bitSize: 64, want: 42, ok: true},
		{name: "hexadecimal", input: "deadBEEF", base: 16, bitSize: 64, want: 0xDEADBEEF, ok: true},
		{name: "base zero", input: "0x2a", base: 0, bitSize: 64, want: 0, ok: false},
		{name: "bitSize zero", input: "6767", base: 10, bitSize: 0, want: 6767, ok: true},
		{name: "zero", input: "0", base: 10, bitSize: 8, want: 0, ok: true},
		{name: "maximum uint8", input: "255", base: 10, bitSize: 8, want: 255, ok: true},
		{name: "overflow", input: "256", base: 10, bitSize: 8, want: 0, ok: false},
		{name: "negative value", input: "-54", base: 10, bitSize: 64, want: 0, ok: false},
		{name: "carry a plus sign", input: "+67", base: 10, bitSize: 64, want: 0, ok: false},
		{name: "carry a sign in middle", input: "82+7", base: 10, bitSize: 64, want: 0, ok: false},
		{name: "has surrounding space", input: " 123 ", base: 10, bitSize: 64, want: 0, ok: false},
		{name: "invalid digit", input: "12g", base: 10, bitSize: 64, want: 0, ok: false},
		{name: "empty input", input: "", base: 10, bitSize: 64, want: 0, ok: false},
		{name: "invalid base", input: "10", base: 1, bitSize: 64, want: 0, ok: false},
		{name: "invalid bit size", input: "10", base: 10, bitSize: 100, want: 0, ok: false},
		{name: "value too large for bitSize", input: "7", base: 2, bitSize: 2, want: 0, ok: false},
		{name: "contains special char", input: "(666)", base: 10, bitSize: 10, want: 0, ok: false},
		{name: "should fail", input: "50000000000000000000", base: 10, bitSize: 64, want: 0, ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseUint(tt.input, tt.base, tt.bitSize)
			if got != tt.want || ok != tt.ok {
				t.Errorf("ParseUint(%q, %d, %d) = (%d, %v), want (%d, %v)", tt.input, tt.base, tt.bitSize, got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestAtoi(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
		ok    bool
	}{
		{name: "zero", input: "-0", want: 0, ok: true},
		{name: "positive", input: "12345", want: 12345, ok: true},
		{name: "negative", input: "-12345", want: -12345, ok: true},
		{name: "max negative", input: "-9223372036854775808", want: -9223372036854775808, ok: true},
		{name: "overflow negative", input: "-9223372036854775809", want: 0, ok: false},
		{name: "leading plus", input: "+42", want: 42, ok: true},
		{name: "leading whitespace is invalid", input: " 42", want: 0, ok: false},
		{name: "decimal point is invalid", input: "42.0", want: 0, ok: false},
		{name: "empty input", input: "", want: 0, ok: false},
		{name: "contain char", input: "123abc", want: 0, ok: false},
		{name: "sign in the middle", input: "123+789", want: 0, ok: false},
		{name: "has surrounding space", input: " 666 ", want: 0, ok: false},
		{name: "overflow possitive", input: "9223372036854775808", want: 0, ok: false},
		{name: "more than 1 sign", input: "+-123", want: 0, ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Atoi(tt.input)
			if got != tt.want || ok != tt.ok {
				t.Errorf("Atoi(%q) = (%d, %v), want (%d, %v)", tt.input, got, ok, tt.want, tt.ok)
			}
		})
	}
}
