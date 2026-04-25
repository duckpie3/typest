package typing

import "time"

type wpmDataPoint struct {
	Time float64
	Wpm  int
}

type TestStats struct {
	Characters  int
	Wpm         int
	WpmData     []wpmDataPoint
	Greatestwpm int
	startTime   time.Time
	ElapsedTime float64
}
