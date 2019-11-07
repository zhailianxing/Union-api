package strconvx

import "strconv"

func Itobs(i int) []byte {
	return formatBits(nil, uint64(i), i < 0, false)
}

func Uint32ToA(i uint32) string {
	return strconv.Itoa(int(i))
}

func Uint16ToA(i uint16) string {
	return strconv.Itoa(int(i))
}

const (
	digits   = "0123456789abcdefghijklmnopqrstuvwxyz"
	digits01 = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
	digits10 = "0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999"
)

var shifts = [len(digits) + 1]uint{
	1 << 1: 1,
	1 << 2: 2,
	1 << 3: 3,
	1 << 4: 4,
	1 << 5: 5,
}

// formatBits computes the string representation of u in the given base.
// If neg is set, u is treated as negative int64 value. If append_ is
// set, the string is appended to dst and the resulting byte slice is
// returned as the first result value; otherwise the string is returned
// as the second result value.
//
func formatBits(dst []byte, u uint64, neg, append_ bool) []byte {

	var a [64 + 1]byte // +1 for sign of 64bit value in base 2
	//	a := make([]byte, 65)
	i := len(a)

	if neg {
		u = -u
	}

	// common case: use constants for / and % because
	// the compiler can optimize it into a multiply+shift,
	// and unroll loop
	for u >= 100 {
		i -= 2
		q := u / 100
		j := uintptr(u - q*100)
		a[i+1] = digits01[j]
		a[i+0] = digits10[j]
		u = q
	}
	if u >= 10 {
		i--
		q := u / 10
		a[i] = digits[uintptr(u-q*10)]
		u = q
	}

	// u < base
	i--
	a[i] = digits[uintptr(u)]

	// add sign, if any
	if neg {
		i--
		a[i] = '-'
	}

	//	bs := make([]byte, len(a[i:]))
	//	copy(bs, a[i:])
	//	return bs

	return a[i:]
}

func AToUint32(str string) uint32 {
	i, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}

	return uint32(i)
}
