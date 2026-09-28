package mybits

func OnesCount32(x uint32) int {
	count := 0
	for x > 0 {
		if x&1 == 1 {
			count++
		}
		x >>= 1
	}
	return count
}
func Len32(x uint32) int
func RotateLeft32(x uint32, k int) uint32
func Reverse32(x uint32) uint32
func ReverseBytes32(x uint32) uint32
