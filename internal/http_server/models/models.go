package models

import (
	"time"
)

type Data struct {
	SongAndGroup
	SongDetails
}

type SongAndGroup struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type SongDetails struct {
	ReleaseDate CustomTime `json:"releaseDate"`
	Text        string     `json:"text"`
	Link        string     `json:"link"`
}

type CustomTime struct {
	time.Time
}

const CustomTimeFormat = "02.01.2006"

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]

	t, err := time.Parse(CustomTimeFormat, s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ct.Time.Format(CustomTimeFormat) + `"`), nil
}

func (ct CustomTime) String() string {
	return ct.Time.Format(CustomTimeFormat)
}
