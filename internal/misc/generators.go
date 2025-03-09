package misc

import (
	"math/rand"
	"strings"
	"time"
)

func RandString(min int, max int) string {
	var sb strings.Builder
	for range rand.Intn(max-min) + min {
		sb.WriteRune(rune(65 + rand.Intn(25)))
	}

	return sb.String()
}

func RanDate() time.Time {
	min := time.Date(2022, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2030, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
