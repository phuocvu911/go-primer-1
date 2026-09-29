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

func ParseUint(s string, base int, bitSize int) (uint64, bool) {
	if base < 2 || base > 36 || bitSize < 0 || bitSize > 64 || len(s) == 0 {
		return 0, false
	}

	if bitSize == 0 {
		bitSize = 64
	}

	result := uint64(0)
	for i := 0; i < len(s); i++ {
		c := s[i]
		var digitValue int

		if c >= '0' && c <= '9' {
			digitValue = int(c - '0')
		} else if c >= 'a' && c <= 'z' {
			digitValue = int(c - 'a' + 10)
		} else if c >= 'A' && c <= 'Z' {
			digitValue = int(c - 'A' + 10)
		} else {
			return 0, false
		}

		if digitValue >= base {
			return 0, false
		}

		// Check for overflow, because the number wrap around when it become larger than Maxint64
		newResult := result*uint64(base) + uint64(digitValue)
		if newResult/uint64(base) != result || newResult%uint64(base) != uint64(digitValue) {
			return 0, false
		}
		result = newResult
	}

	// Check bitSize can hold the number
	if bitSize < 64 {
		maxVal := uint64(1) << uint(bitSize)
		if result >= maxVal {
			return 0, false
		}
	}

	return result, true
}

const maxInt = 1<<63 - 1

func Atoi(s string) (int, bool) {
	if len(s) == 0 {
		return 0, false
	}
	sign := s[0]
	isNeg := false
	if sign == '+' || sign == '-' {
		s = s[1:]
		if sign == '-' {
			isNeg = true
		}
	}
	res, ok := ParseUint(s, 10, 64)
	if !ok {
		return 0, false
	}

	if isNeg {
		if res > maxInt+1 {
			return 0, false
		}
		return int(-res), true
	}
	if res > maxInt {
		return 0, false
	}
	return int(res), true
}
