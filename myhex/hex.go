package myhex

const hex = "0123456789abcdef"

var hexSearch = map[byte]byte{
	'0': 0x0, '1': 0x1, '2': 0x2, '3': 0x3,
	'4': 0x4, '5': 0x5, '6': 0x6, '7': 0x7,
	'8': 0x8, '9': 0x9, 'a': 0xa, 'b': 0xb,
	'c': 0xc, 'd': 0xd, 'e': 0xe, 'f': 0xf,
	'A': 0xa, 'B': 0xb, 'C': 0xc, 'D': 0xd, 'E': 0xe, 'F': 0xf,
}

// Encode turn byte -> hex, 1 byte = 2 hex
func Encode(src []byte) []byte {
	res := make([]byte, len(src)*2)
	for i, b := range src {
		res[2*i] = hex[b>>4]     //high nibble, shift right that byte 4 positions
		res[2*i+1] = hex[b&0x0f] //low nibble, mask with 0000 1111
	}
	return res
}

// Decede do the opposite of Encode
func Decode(src []byte) ([]byte, bool) {
	if len(src)%2 != 0 {
		return nil, false
	}

	res := make([]byte, len(src)/2)
	for i := 0; i < len(src); i += 2 {
		high, ok := hexSearch[src[i]]
		if !ok {
			return nil, false
		}
		low, ok := hexSearch[src[i+1]]
		if !ok {
			return nil, false
		}
		res[i/2] = high<<4 | low
	}
	return res, true
}
