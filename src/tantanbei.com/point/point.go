package xpoint

const (
	base      uint16 = 5
	MAX_TIMES uint16 = 4
)

func GetAddPoint(keepDay uint16) uint32 {
	if keepDay < MAX_TIMES {
		return uint32(base * keepDay)
	} else {
		return uint32(base * MAX_TIMES)
	}
}
