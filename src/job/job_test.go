package job

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"testing"
)

func getMockJob() *Job {
	return &Job{
		Name:    "mock_job",
		Command: "bash -c 'date'",
		Owner:   "aj@ajvb.me",
		Retries: 2,
	}
}

func getMockJobWithSchedule(repeat int, scheduleTime time.Time, delay string) *Job {
	genericMockJob := getMockJob()

	parsedTime := scheduleTime.Format(time.RFC3339)
	scheduleStr := fmt.Sprintf("R%d/%s/%s", repeat, parsedTime, delay)
	genericMockJob.Schedule = scheduleStr

	return genericMockJob
}

func getMockJobWithGenericSchedule() *Job {
	fiveMinutesFromNow := time.Now().Add(
		time.Duration(time.Minute * 5),
	)
	return getMockJobWithSchedule(2, fiveMinutesFromNow, "P1DT10M10S")
}

func TestScheduleParsing(t *testing.T) {
	fiveMinutesFromNow := time.Now().Add(
		time.Duration(time.Minute * 5),
	)

	genericMockJob := getMockJobWithSchedule(2, fiveMinutesFromNow, "P1DT10M10S")

	genericMockJob.Init()

	assert.WithinDuration(
		t, genericMockJob.scheduleTime, fiveMinutesFromNow,
		time.Second, "The difference between parsed time and created time is to great.",
	)

}

var delayParsingTests = []struct {
	expected    time.Duration
	intervalStr string
}{
	{time.Duration(
		(time.Hour * 24) + (time.Second * 10) + (time.Minute * 10),
	), "P1DT10M10S"},
	{time.Duration(
		(time.Second * 10) + (time.Minute * 10),
	), "PT10M10S"},
	{time.Duration(
		(time.Hour * 24) + (time.Second * 10),
	), "P1DT10S"},
	{time.Duration(
		(time.Hour * 24 * 365) + (time.Hour * 24),
	), "P1Y1DT"},
}

func TestDelayParsing(t *testing.T) {
	testTime := time.Now().Add(
		time.Duration(time.Minute * 1),
	)

	for _, delayTest := range delayParsingTests {
		genericMockJob := getMockJobWithSchedule(1, testTime, delayTest.intervalStr)
		genericMockJob.Init()
		assert.Equal(t, delayTest.expected, genericMockJob.delayDuration.ToDuration(), "Parsed duration was incorrect")
	}
}

func TestJobInit(t *testing.T) {
	genericMockJob := getMockJobWithGenericSchedule()

	err := genericMockJob.Init()
	assert.Nil(t, err, "err should be nil")

	assert.NotEmpty(t, genericMockJob.Id, "Job.Id should not be empty")
	assert.NotEmpty(t, genericMockJob.jobTimer, "Job.jobTimer should not be empty")
}

func TestJobDisable(t *testing.T) {
	genericMockJob := getMockJobWithGenericSchedule()
	genericMockJob.Init()

	assert.False(t, genericMockJob.Disabled, "Job should start with disbaled == false")

	genericMockJob.Disable()
	assert.True(t, genericMockJob.Disabled, "Job.Disable() should set Job.Disabled to true")
	//TODO test genericMockJob.jobTimer is stopped
}

// TODO
func TestJobRun(t *testing.T) {
}

func TestOneOffJobs(t *testing.T) {
	j := getMockJob()

	assert.Equal(t, j.SuccessCount, uint(0))
	assert.Equal(t, j.ErrorCount, uint(0))
	assert.Equal(t, j.LastSuccess, time.Time{})
	assert.Equal(t, j.LastError, time.Time{})
	assert.Equal(t, j.LastAttemptedRun, time.Time{})

	j.Init()
	time.Sleep(time.Second)
	now := time.Now()

	assert.Equal(t, j.SuccessCount, uint(1))
	assert.WithinDuration(t, j.LastSuccess, now, 2*time.Second)
	assert.WithinDuration(t, j.LastAttemptedRun, now, 2*time.Second)
	assert.Equal(t, j.scheduleTime, time.Time{})
	assert.Nil(t, j.jobTimer)
}

// TODO
func TestDependentJobs(t *testing.T) {
}

// TODO
func TestMetaData(t *testing.T) {
}
