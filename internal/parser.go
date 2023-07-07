package internal

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Event struct {
	ID    int
	Time  int
	Name  string
	Table int
}

var eventLen = []int{
	0: 3,
	1: 4,
	2: 3,
	3: 3,
}

var sample = `^[0-2][0-9]:[0-5][0-9] [1-4] [a-z0-9\-_]+$`

func ParseEvent(event string, tablesAmount int) (Event, bool) {
	var e Event

	parts := strings.Split(event, " ")
	if len(parts) < 3 {
		return e, false
	}
	re, err := regexp.Compile(sample)
	if err != nil {
		log.Fatal(err)
	}

	if !re.MatchString(strings.Join(parts[:3], " ")) {
		return e, false
	}

	id, _ := strconv.Atoi(parts[1])

	if id > len(eventLen) || len(parts) != eventLen[id-1] {
		return e, false
	}

	//if _, err := time.Parse("15:04", parts[0]); err != nil {
	//	return e, false
	//}

	e.ID = id
	var ok bool
	if e.Time, ok = ParseTimeToMinutes(parts[0]); !ok {
		return e, false
	}
	e.Name = parts[2]

	if id == 2 {
		e.Table, err = strconv.Atoi(parts[3])
		if err != nil || e.Table > tablesAmount || e.Table < 1 {
			return e, false
		}

	}

	return e, true
}

func ParseTimeToMinutes(time string) (int, bool) {
	parts := strings.Split(time, ":")
	hours, err := strconv.Atoi(parts[0])
	if err != nil || hours > 23 || hours < 0 {
		return 0, false
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil || minutes > 59 || minutes < 0 {
		return 0, false
	}
	return hours*60 + minutes, true
}

func ParseMinutesToTime(minutes int) string {
	return fmt.Sprintf("%02d:%02d", minutes/60, minutes%60)

}

func ParseInt(s string) (int, bool) {
	t, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return t, true
}
