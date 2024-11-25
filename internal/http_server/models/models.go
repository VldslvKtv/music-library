package models

import (
	"time"
)

type Data struct {
	SongAndGroup
	SongDetails
}

type SongAndGroup struct {
	Group string `json:"group" validate:"required"`
	Song  string `json:"song" validate:"required"`
}

type SongDetails struct {
	ReleaseDate CustomTime `json:"releaseDate"`
	Text        string     `json:"text"`
	Link        string     `json:"link"`
}

// CustomTimeFormat определяет формат даты, используемый для маршалинга и демаршалинга JSON.
const CustomTimeFormat = "02.01.2006"

// CustomTime расширяет структуру time.Time, предоставляя методы для работы с JSON.
type CustomTime struct {
	time.Time
}

// Реализует интерфейс json.Unmarshaler для CustomTime.
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

// Реализует интерфейс json.Marshaler для CustomTime.
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ct.Time.Format(CustomTimeFormat) + `"`), nil
}

// Возвращает строковое представление CustomTime в заданном формате даты.
func (ct CustomTime) String() string {
	return ct.Time.Format(CustomTimeFormat)
}
