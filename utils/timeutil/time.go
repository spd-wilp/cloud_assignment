package timeutil

import "time"

// FindTimeBoundOfPreviousDay returns
//   - start and endtime of previous day [00.00.00-23.59.59]
//   - error, if provided date is not a valid unix timestamp
func FindTimeBoundOfPreviousDay(curTime time.Time) (int64, int64, error) {
	prevDay := curTime.Add(-24 * time.Hour)
	st := prevDay.Truncate(24 * time.Hour)
	et := st.Add(24 * time.Hour).Add(-1 * time.Second)
	return st.Unix(), et.Unix(), nil
}
