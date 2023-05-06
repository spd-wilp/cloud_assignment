package timeutil

import (
	"testing"
	"time"
)

func Test_FindTimeBoundOfPreviousDay(t *testing.T) {
	curTime := time.Unix(1683359006, 0)
	expectedSt := int64(1683244800)
	expectedEt := int64(1683331199)

	st, et, _ := FindTimeBoundOfPreviousDay(curTime)
	if st != expectedSt {
		t.Logf("st not correct, expected=%v received=%v", expectedSt, st)
		t.Fail()
	}
	if et != expectedEt {
		t.Logf("et not correct, expected=%v received=%v", expectedEt, et)
		t.Fail()
	}
}
