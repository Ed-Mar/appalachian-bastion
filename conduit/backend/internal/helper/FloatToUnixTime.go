package helper

import (
	"math"
	"time"
)

func FloatToUnixTime(f float64) time.Time {
	round, frac := math.Modf(f)
	return time.Unix(int64(round), int64(frac*1e9))
}
