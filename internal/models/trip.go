package models

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

var weekdaysByName map[string]time.Weekday

func init() {
	weekdaysByName = make(map[string]time.Weekday)
	for d := time.Sunday; d <= time.Saturday; d++ {
		weekdaysByName[d.String()[:3]] = d
	}
}

type Weekdays map[time.Weekday]struct{}

func (w *Weekdays) UnmarshalJSON(i []byte) error {
	*w = make(map[time.Weekday]struct{})
	if len(i) == 0 {
		return nil
	}
	var rawWeekdays string
	err := json.Unmarshal(i, &rawWeekdays)
	if err != nil {
		return err
	}

	trimmed := strings.TrimSpace(rawWeekdays)
	if len(trimmed) == 0 {
		return nil
	}

	weekdayNames := strings.Split(trimmed, " ")
	for _, weekdayName := range weekdayNames {
		if weekday, ok := weekdaysByName[weekdayName]; ok {
			(*w)[weekday] = struct{}{}
		} else {
			return fmt.Errorf("invalid weekday '%s'", weekdayName)
		}
	}

	return nil
}

func (w Weekdays) MarshalJSON() ([]byte, error) {
	if len(w) == 0 {
		return json.Marshal("")
	}

	var days []int
	for weekday := range w {
		days = append(days, int(weekday))
	}

	sort.Ints(days)
	if days[0] == int(time.Sunday) {
		days = append(days[1:], days[0])
	}

	var buff bytes.Buffer
	for _, day := range days {
		if buff.Len() != 0 {
			buff.WriteString(" ")
		}
		buff.WriteString(time.Weekday(day).String()[:3])
	}

	return json.Marshal(buff.String())
}

func (w *Weekdays) Scan(src interface{}) error {
	scanned, ok := src.(int64)
	if !ok {
		return fmt.Errorf("database returned non int64 value: %+[1]v(%[1]T)", src)
	}

	byteArray := make([]byte, 8)
	binary.LittleEndian.PutUint64(byteArray, uint64(scanned))
	bitMask := byteArray[0]

	*w = make(map[time.Weekday]struct{})

	var bit byte
	for i := 0; i < 7; i++ {
		bit = 1 << i

		if (bitMask & bit) != 0 {
			(*w)[time.Weekday(i)] = struct{}{}
		}
	}

	return nil
}

func (w Weekdays) Value() (driver.Value, error) {
	if w == nil {
		return nil, errors.New("weekdays set is not initialized")
	}

	var bitmask int64
	for weekday := range w {
		bitmask |= 1 << byte(weekday)
	}

	return bitmask, nil
}

type Trip struct {
	Id          int64
	Origin      City
	Destination City
	Price       int64
	Dates       Weekdays
}
