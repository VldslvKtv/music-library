package utils

import (
	"testing"
	"time"

	"music_library/internal/http_server/models"

	"github.com/stretchr/testify/assert"
)

func TestConvertStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    models.Data
		expected map[string]interface{}
	}{
		{
			name: "Полная структура",
			input: models.Data{
				SongAndGroup: models.SongAndGroup{
					Group: "Imagine Dragons",
					Song:  "Believer",
				},
				SongDetails: models.SongDetails{
					ReleaseDate: models.CustomTime{Time: time.Date(2006, 7, 16, 0, 0, 0, 0, time.UTC)},
					Text:        "First things first...",
					Link:        "https://example.com",
				},
			},
			expected: map[string]interface{}{
				"groups.name":  "Imagine Dragons",
				"songs.name":   "Believer",
				"release_date": "16.07.2006",
				"text":         "First things first...",
				"link":         "https://example.com",
			},
		},
		{
			name: "Пустая структура",
			input: models.Data{
				SongAndGroup: models.SongAndGroup{
					Group: "",
					Song:  "",
				},
				SongDetails: models.SongDetails{
					ReleaseDate: models.CustomTime{Time: time.Time{}},
					Text:        "",
					Link:        "",
				},
			},
			expected: map[string]interface{}{},
		},
		{
			name: "Частично заполненная структура",
			input: models.Data{
				SongAndGroup: models.SongAndGroup{
					Group: "Linkin Park",
					Song:  "",
				},
				SongDetails: models.SongDetails{
					ReleaseDate: models.CustomTime{Time: time.Date(2003, 3, 25, 0, 0, 0, 0, time.UTC)},
					Text:        "",
					Link:        "https://linkinpark.com",
				},
			},
			expected: map[string]interface{}{
				"groups.name":  "Linkin Park",
				"release_date": "25.03.2003",
				"link":         "https://linkinpark.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStruct(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestChangeKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "Стандартный случай",
			input: map[string]interface{}{
				"group":       "Imagine Dragons",
				"song":        "Believer",
				"releasedate": "16.07.2006",
			},
			expected: map[string]interface{}{
				"groups.name":  "Imagine Dragons",
				"songs.name":   "Believer",
				"release_date": "16.07.2006",
			},
		},
		{
			name: "Случай с другими ключами",
			input: map[string]interface{}{
				"artist":      "Linkin Park",
				"title":       "Numb",
				"releasedate": "25.03.2003",
				"album":       "Meteora",
			},
			expected: map[string]interface{}{
				"artist":       "Linkin Park",
				"title":        "Numb",
				"release_date": "25.03.2003",
				"album":        "Meteora",
			},
		},
		{
			name:     "Пустой map",
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			name: "Частичный map",
			input: map[string]interface{}{
				"group":       "Coldplay",
				"song":        "",
				"releasedate": "26.03.2004",
				"genre":       "Alternative Rock",
			},
			expected: map[string]interface{}{
				"groups.name":  "Coldplay",
				"songs.name":   "",
				"release_date": "26.03.2004",
				"genre":        "Alternative Rock",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ChangeKeys(&tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
