package utils

import (
	"fmt"
	"music_library/internal/http_server/models"
	"testing"
	"time"
)

func TestConvertStruct(t *testing.T) {
	data := models.Data{
		SongAndGroup: models.SongAndGroup{
			Group: "Imagine Dragons",
			Song:  "Believer",
		},
		SongDetails: models.SongDetails{
			ReleaseDate: models.CustomTime{Time: time.Date(2006, 7, 16, 0, 0, 0, 0, time.UTC)},
			Text:        "First things first...",
			Link:        "https://example.com",
		},
	}

	result := ConvertStruct(data)
	fmt.Println(result)

}

func TestChangeKeys(t *testing.T) {
	data := models.Data{
		SongAndGroup: models.SongAndGroup{
			Group: "Imagine Dragons",
			Song:  "Believer",
		},
		SongDetails: models.SongDetails{
			ReleaseDate: models.CustomTime{Time: time.Date(2006, 7, 16, 0, 0, 0, 0, time.UTC)},
			Text:        "First things first...",
			Link:        "https://example.com",
		},
	}

	result := ConvertStruct(data)

	result = ChangeKeys(&result)
	fmt.Println(result)
}
