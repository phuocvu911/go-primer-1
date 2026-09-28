package mystrconv

const baseTable = "0123456789abcdefghijklmnopqrstuvwxyz"

func FormatUint(u uint64, base int) (string, bool) {
	if base < 2 || base > 36 {
		return "", false
	}
	if u == 0 {
		return "0", true
	}
	baseStr := baseTable[:base]
	res := ""
	for u > 0 {
		digit := baseStr[u%uint64(base)]
		res = string(digit) + res
		u /= uint64(base)
	}
	return res, true
}

func FormatInt(i int64, base int) (string, bool) {
	if base < 2 || base > 36 {
		return "", false
	}
	if i == 0 {
		return "0", true
	}
	isNeg := i < 0
	baseStr := baseTable[:base]
	res := ""
	u := uint64(0)
	if isNeg {
		u = uint64(-i)
	} else {
		u = uint64(i)
	}
	for u > 0 {
		digit := baseStr[u%uint64(base)]
		res = string(digit) + res
		u /= uint64(base)
	}
	if isNeg {
		res = "-" + res
	}
	return res, true
}

// func ParseUint(s string, base int, bitSize int) (uint64, bool)
// func Atoi(s string) (int, bool)
