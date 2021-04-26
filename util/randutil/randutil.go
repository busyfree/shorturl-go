package randutil

import (
	"math/rand"
)

func RandUint16(min, max int32) uint16 {
	if max > 65535 {
		return 65535
	}
	if min >= max || max == 0 {
		return uint16(max)
	}
	return uint16(min + rand.Int31n(max-min))
}
