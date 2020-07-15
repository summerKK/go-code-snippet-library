package timer

import "time"

func GetNowTime() time.Time {
	return time.Now()
}

func GetCalculateTime(currentTime time.Time, d string) (time.Time, error) {
	duration, err := time.ParseDuration(d)
	if err != nil {
		return time.Time{}, err
	}

	return currentTime.Add(duration), nil
}
