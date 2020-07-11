package pkg

import (
	"github.com/magiconair/properties/assert"
	"testing"
	"web-crawler/internal"
)

var statsCounter = NewStatsCounter()

func Test_statsCounterImpl_GetSnapshot(t *testing.T) {
	statsCounter := NewStatsCounter()
	statsCounter.RecordFetchSuccess(1)
	statsCounter.RecordFetchSuccess(1)
	statsCounter.RecordFetchSuccess(1)
	statsCounter.RecordParseFailure(1)
	statsCounter.RecordParseSuccess(1)
	statsCounter.RecordFetchFailure(1)
	gotSnapShot := statsCounter.GetSnapshot()
	wantSnapShot := internal.NewCrawlerStats(3, 1, 1, 1, 0)
	assert.Equal(t, gotSnapShot, wantSnapShot)
}

func Test_statsCounterImpl_RecordFetchSuccess(t *testing.T) {
	statsCounter.RecordFetchSuccess(1)
	statsCounter.RecordFetchSuccess(1)
	snapShot := statsCounter.GetSnapshot()
	got := snapShot.GetSucessFetchCount()
	var want int64 = 2
	assert.Equal(t, got, want)
}

func Test_statsCounterImpl_RecordFetchFailure(t *testing.T) {
	statsCounter.RecordFetchFailure(1)
	statsCounter.RecordFetchFailure(1)
	snapShot := statsCounter.GetSnapshot()
	got := snapShot.GetFailedFetchCount()
	var want int64 = 2
	assert.Equal(t, got, want)
}

func Test_statsCounterImpl_RecordParseSuccess(t *testing.T) {
	statsCounter.RecordParseSuccess(1)
	statsCounter.RecordParseSuccess(1)
	snapShot := statsCounter.GetSnapshot()
	got := snapShot.GetSucessParseCount()
	var want int64 = 2
	assert.Equal(t, got, want)
}

func Test_statsCounterImpl_RecordParseFailure(t *testing.T) {
	statsCounter.RecordParseFailure(1)
	statsCounter.RecordParseFailure(1)
	snapShot := statsCounter.GetSnapshot()
	got := snapShot.GetFailedParseCount()
	var want int64 = 2
	assert.Equal(t, got, want)
}