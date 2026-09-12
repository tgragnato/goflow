package utils

import (
	"testing"
	"time"
)

func TestBatchMute(t *testing.T) {
	t.Parallel()

	tm := time.Date(2023, time.November, 10, 23, 0, 0, 0, time.UTC)
	bm := BatchMute{
		batchTime:     tm,
		resetInterval: time.Second * 10,
		max:           5,
	}

	for range 20 {
		tm = tm.Add(time.Second)
		t.Log(bm.increment(1, tm))
	}

}

func TestBatchMuteZero(t *testing.T) {
	t.Parallel()

	tm := time.Date(2023, time.November, 10, 23, 0, 0, 0, time.UTC)
	bm := BatchMute{
		batchTime:     tm,
		resetInterval: time.Second * 10,
		max:           0,
	}

	for range 20 {
		tm = tm.Add(time.Second)
		t.Log(bm.increment(1, tm))
	}

}

func TestBatchMuteInterval(t *testing.T) {
	t.Parallel()

	tm := time.Date(2023, time.November, 10, 23, 0, 0, 0, time.UTC)
	bm := BatchMute{
		batchTime:     tm,
		resetInterval: 0,
		max:           5,
	}

	for range 20 {
		tm = tm.Add(time.Second)
		t.Log(bm.increment(1, tm))
	}

}
