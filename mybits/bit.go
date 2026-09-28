package mybits

// OnesCount32 counts how many "1" in x
func OnesCount32(x uint32) int {
	count := 0
	for x > 0 {
		x &= (x - 1)
		count++
	}
	return count
}

// Len32 calculate how many bits needed to represent x in binary scale
func Len32(x uint32) int {
	if x == 0 {
		return 0
	}
	size := 0
	for x > 1 {
		x >>= 1
		size++
	}
	return size + 1
}

// RotateLeft32 rotate binary form of x to left (k%32) time
func RotateLeft32(x uint32, k int) uint32 {
	y := uint(k) % 32
	return x<<y | x>>(32-y)
}

// Reverse32 reverse the binary form of x
func Reverse32(x uint32) uint32 {
	var res uint32
	for range 32 {
		res = (res << 1) | (x & 1)
		x >>= 1
	}
	return res
}

// ReverseButes32 reverse x by the chunk of 8 bits(1 byte)
func ReverseBytes32(x uint32) uint32 {
	var res uint32
	for range 4 {
		res = (res << 8) | (x & 0xFF) //255, which is '11111111'
		x >>= 8
	}
	return res
}
