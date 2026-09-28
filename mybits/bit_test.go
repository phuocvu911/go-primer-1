package mybits

import (
	"math"
	"math/bits"
	"testing"
)

func TestOnesCount32(t *testing.T) {
	tests := []struct {
		name string
		in   uint32
		want int
	}{
		{name: "zero", in: 0, want: 0},
		{name: "one (lowest bit)", in: 1, want: 1},
		{name: "two (second bit)", in: 2, want: 1},
		{name: "three (two bits)", in: 3, want: 2},
		{name: "highest bit set (1<<31)", in: 1 << 31, want: 1},
		{name: "max uint32 (all 32 bits set)", in: math.MaxUint32, want: 32},
		{name: "alternating bits 0x55555555", in: 0x55555555, want: 16},
		{name: "alternating bits 0xAAAAAAAA", in: 0xAAAAAAAA, want: 16},
		{name: "alternating bytes 0x00FF00FF", in: 0x00FF00FF, want: 16},
		{name: "arbitrary value 0x12345678", in: 0x12345678, want: 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OnesCount32(tt.in)
			if got != tt.want {
				t.Errorf("OnesCount32(%#x) = %d; want %d", tt.in, got, tt.want)
			}
		})
	}

	t.Run("single bit at each position", func(t *testing.T) {
		for i := range 32 {
			val := uint32(1) << i
			if got := OnesCount32(val); got != 1 {
				t.Errorf("OnesCount32(1 << %d = %#x) = %d; want 1", i, val, got)
			}
		}
	})

	t.Run("powers of two minus 1", func(t *testing.T) {
		for i := 1; i <= 32; i++ {
			var val uint32
			if i == 32 {
				val = math.MaxUint32
			} else {
				val = (uint32(1) << i) - 1
			}
			if got := OnesCount32(val); got != i {
				t.Errorf("OnesCount32((1 << %d) - 1) = %d; want %d", i, got, i)
			}
		}
	})

	t.Run("cross check with math/bits", func(t *testing.T) {
		samples := []uint32{
			0, 1, 2, 0x7F, 0x80, 0xFF, 0x100, 0x7FFF, 0x8000, 0xFFFF,
			0x7FFFFFFF, 0x80000000, math.MaxUint32, 0xA5A5A5A5,
		}
		for _, x := range samples {
			want := bits.OnesCount32(x)
			if got := OnesCount32(x); got != want {
				t.Errorf("OnesCount32(%#x) = %d; want %d (from math/bits)", x, got, want)
			}
		}
	})
}

func TestLen32(t *testing.T) {
	tests := []struct {
		name string
		in   uint32
		want int
	}{
		{name: "zero", in: 0, want: 0},
		{name: "one", in: 1, want: 1},
		{name: "two", in: 2, want: 2},
		{name: "three", in: 3, want: 2},
		{name: "four", in: 4, want: 3},
		{name: "seven", in: 7, want: 3},
		{name: "eight", in: 8, want: 4},
		{name: "highest power of two minus 1", in: 1<<31 - 1, want: 31},
		{name: "highest power of two (1<<31)", in: 1 << 31, want: 32},
		{name: "max uint32", in: math.MaxUint32, want: 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Len32(tt.in)
			if got != tt.want {
				t.Errorf("Len32(%#x) = %d; want %d", tt.in, got, tt.want)
			}
		})
	}

	t.Run("powers of two and predecessor", func(t *testing.T) {
		for i := range 32 {
			val := uint32(1) << i
			if got := Len32(val); got != i+1 {
				t.Errorf("Len32(1 << %d) = %d; want %d", i, got, i+1)
			}
			if i > 0 {
				pre := val - 1
				if got := Len32(pre); got != i {
					t.Errorf("Len32((1 << %d) - 1) = %d; want %d", i, got, i)
				}
			}
		}
	})

	t.Run("cross check with math/bits", func(t *testing.T) {
		samples := []uint32{
			0, 1, 2, 3, 4, 15, 16, 255, 256, 65535, 65536,
			0x7FFFFFFF, 0x80000000, 0x80000001, math.MaxUint32,
		}
		for _, x := range samples {
			want := bits.Len32(x)
			if got := Len32(x); got != want {
				t.Errorf("Len32(%#x) = %d; want %d (from math/bits)", x, got, want)
			}
		}
	})
}

func TestRotateLeft32(t *testing.T) {
	tests := []struct {
		name string
		x    uint32
		k    int
		want uint32
	}{
		{name: "zero rotated by 0", x: 0, k: 0, want: 0},
		{name: "zero rotated by positive", x: 0, k: 15, want: 0},
		{name: "zero rotated by negative", x: 0, k: -15, want: 0},
		{name: "no rotation k=0", x: 0x12345678, k: 0, want: 0x12345678},
		{name: "full rotation k=32", x: 0x12345678, k: 32, want: 0x12345678},
		{name: "multiple full rotations k=64", x: 0x12345678, k: 64, want: 0x12345678},
		{name: "multiple full rotations k=96", x: 0x12345678, k: 96, want: 0x12345678},
		{name: "rotate left by 1", x: 1 << 16, k: 1, want: 1 << 17},
		{name: "rotate left by 1 with low bit", x: 1, k: 1, want: 2},
		{name: "rotate left by 31", x: 1, k: 31, want: 1 << 31},
		{name: "rotate left by 33 (equivalent to 1)", x: 1 << 31, k: 33, want: 1},
		{name: "negative rotation k=-1 (equivalent to rotate right 1)", x: 1, k: -1, want: 1 << 31},
		{name: "negative rotation k=-1 from most significant bit", x: 1 << 31, k: -1, want: 1 << 30},
		{name: "negative rotation k=-31 (equivalent to rotate left 1)", x: 1 << 31, k: -31, want: 1},
		{name: "negative full rotation k=-32", x: 0x12345678, k: -32, want: 0x12345678},
		{name: "negative multiple full rotations k=-64", x: 0x12345678, k: -64, want: 0x12345678},
		{name: "negative rotation k=-33 (equivalent to -1)", x: 0x00000001, k: -33, want: 0x80000000},
		{name: "all bits set rotated", x: math.MaxUint32, k: 13, want: math.MaxUint32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RotateLeft32(tt.x, tt.k)
			if got != tt.want {
				t.Errorf("RotateLeft32(%#08x, %d) = %#08x; want %#08x", tt.x, tt.k, got, tt.want)
			}
		})
	}

	t.Run("inverse property RotateLeft(RotateLeft(x, k), -k) == x", func(t *testing.T) {
		values := []uint32{0x12345678, 0x80000001, 0xAAAAAAAA, 0x55555555, 1, math.MaxUint32}
		shifts := []int{-65, -33, -32, -31, -7, -1, 0, 1, 7, 31, 32, 33, 65}
		for _, x := range values {
			for _, k := range shifts {
				rotated := RotateLeft32(x, k)
				restored := RotateLeft32(rotated, -k)
				if restored != x {
					t.Errorf("RotateLeft32(RotateLeft32(%#08x, %d), %d) = %#08x; want %#08x", x, k, -k, restored, x)
				}
			}
		}
	})

	t.Run("cross check with math/bits", func(t *testing.T) {
		values := []uint32{0, 1, 0x80000000, 0x12345678, 0xDEADBEEF, math.MaxUint32}
		for _, x := range values {
			for k := -64; k <= 64; k++ {
				want := bits.RotateLeft32(x, k)
				if got := RotateLeft32(x, k); got != want {
					t.Errorf("RotateLeft32(%#08x, %d) = %#08x; want %#08x (from math/bits)", x, k, got, want)
				}
			}
		}
	})
}

func TestReverse32(t *testing.T) {
	tests := []struct {
		name string
		in   uint32
		want uint32
	}{
		{name: "zero", in: 0, want: 0},
		{name: "max uint32 (all ones)", in: math.MaxUint32, want: math.MaxUint32},
		{name: "lowest bit set (1)", in: 1, want: 1 << 31},
		{name: "highest bit set (1<<31)", in: 1 << 31, want: 1},
		{name: "two outer bits (palindrome)", in: (1 << 31) | 1, want: (1 << 31) | 1},
		{name: "alternating bits 0x55555555", in: 0x55555555, want: 0xAAAAAAAA},
		{name: "alternating bits 0xAAAAAAAA", in: 0xAAAAAAAA, want: 0x55555555},
		{name: "lower 16 bits set", in: 0x0000FFFF, want: 0xFFFF0000},
		{name: "upper 16 bits set", in: 0xFFFF0000, want: 0x0000FFFF},
		{name: "arbitrary pattern 0x12345678", in: 0x12345678, want: 0x1E6A2C48},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Reverse32(tt.in)
			if got != tt.want {
				t.Errorf("Reverse32(%#08x) = %#08x; want %#08x", tt.in, got, tt.want)
			}
		})
	}

	t.Run("involution property Reverse32(Reverse32(x)) == x", func(t *testing.T) {
		values := []uint32{
			0, 1, 2, 0x12345678, 0x80000000, 0x80000001,
			0xAAAAAAAA, 0x55555555, 0xCAFEBABE, math.MaxUint32,
		}
		for _, x := range values {
			if got := Reverse32(Reverse32(x)); got != x {
				t.Errorf("Reverse32(Reverse32(%#08x)) = %#08x; want %#08x", x, got, x)
			}
		}
	})

	t.Run("single bit reversal across all 32 positions", func(t *testing.T) {
		for i := range 32 {
			in := uint32(1) << i
			want := uint32(1) << (31 - i)
			if got := Reverse32(in); got != want {
				t.Errorf("Reverse32(1 << %d) = %#08x; want %#08x", i, got, want)
			}
		}
	})

	t.Run("cross check with math/bits", func(t *testing.T) {
		samples := []uint32{
			0, 1, 2, 3, 0x55555555, 0xAAAAAAAA, 0x12345678,
			0x7FFFFFFF, 0x80000000, 0xDEADBEEF, math.MaxUint32,
		}
		for _, x := range samples {
			want := bits.Reverse32(x)
			if got := Reverse32(x); got != want {
				t.Errorf("Reverse32(%#08x) = %#08x; want %#08x (from math/bits)", x, got, want)
			}
		}
	})
}

func TestReverseBytes32(t *testing.T) {
	tests := []struct {
		name string
		in   uint32
		want uint32
	}{
		{name: "zero", in: 0, want: 0},
		{name: "all ones", in: math.MaxUint32, want: math.MaxUint32},
		{name: "lowest byte only 0x000000FF", in: 0x000000FF, want: 0xFF000000},
		{name: "highest byte only 0xFF000000", in: 0xFF000000, want: 0x000000FF},
		{name: "lowest bit set (1)", in: 1, want: 0x01000000},
		{name: "highest bit set (1<<31)", in: 1 << 31, want: 0x00000080},
		{name: "sequential bytes 0x12345678", in: 0x12345678, want: 0x78563412},
		{name: "sequential bytes 0x01020304", in: 0x01020304, want: 0x04030201},
		{name: "symmetric byte pattern 0x12343412", in: 0x12343412, want: 0x12343412},
		{name: "outer bytes set 0xFF0000FF", in: 0xFF0000FF, want: 0xFF0000FF},
		{name: "inner bytes set 0x00FFFF00", in: 0x00FFFF00, want: 0x00FFFF00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReverseBytes32(tt.in)
			if got != tt.want {
				t.Errorf("ReverseBytes32(%#08x) = %#08x; want %#08x", tt.in, got, tt.want)
			}
		})
	}

	t.Run("involution property ReverseBytes32(ReverseBytes32(x)) == x", func(t *testing.T) {
		values := []uint32{
			0, 1, 0x12345678, 0x01020304, 0x80000000,
			0xAABBCCDD, 0x12343412, math.MaxUint32,
		}
		for _, x := range values {
			if got := ReverseBytes32(ReverseBytes32(x)); got != x {
				t.Errorf("ReverseBytes32(ReverseBytes32(%#08x)) = %#08x; want %#08x", x, got, x)
			}
		}
	})

	t.Run("each byte position mapped to reverse byte position", func(t *testing.T) {
		bytes := []uint32{
			0xAA000000,
			0x00BB0000,
			0x0000CC00,
			0x000000DD,
		}
		expected := []uint32{
			0x000000AA,
			0x0000BB00,
			0x00CC0000,
			0xDD000000,
		}
		for i := range bytes {
			if got := ReverseBytes32(bytes[i]); got != expected[i] {
				t.Errorf("ReverseBytes32(%#08x) = %#08x; want %#08x", bytes[i], got, expected[i])
			}
		}
	})

	t.Run("cross check with math/bits", func(t *testing.T) {
		samples := []uint32{
			0, 1, 2, 0xFF, 0x100, 0x12345678, 0xDEADBEEF,
			0x80000000, 0x7FFFFFFF, math.MaxUint32,
		}
		for _, x := range samples {
			want := bits.ReverseBytes32(x)
			if got := ReverseBytes32(x); got != want {
				t.Errorf("ReverseBytes32(%#08x) = %#08x; want %#08x (from math/bits)", x, got, want)
			}
		}
	})
}
