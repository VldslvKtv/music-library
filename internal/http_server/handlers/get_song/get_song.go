package get_song

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"music_library/internal/http_server/lib/utils"
	"music_library/internal/http_server/storage"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
)

// GetText представляет интерфейс для получения текста песни.
// @Description Интерфейс для получения текста песни.
type GetText interface {
	// GetSong получает текст песни по имени группы и имени песни.
	// @Description Получение текста песни по имени группы и имени песни.
	// @Param group string "Имя группы"
	// @Param song string "Имя песни"
	// @return string "Текст песни"
	// @return error "Ошибка выполнения"
	GetSong(group string, song string) (string, error)
}

// Response представляет структуру ответа с текстом песни и информацией о куплетах.
// @Description Структура ответа с текстом песни и информацией о куплетах.
type Response struct {
	Verses           []string `json:"verses"`
	MaxVersesPerPage int      `json:"maxVersesPerPage"`
	TotalPages       int      `json:"totalPages"`
	CurrentPage      int      `json:"currentPage"`
}

// New создает новый обработчик для получения текста песни с пагинацией по куплетам (метод GET).
// @Summary Получение текста песни
// @Description Получение текста песни с пагинацией по куплетам (метод GET).
// @ID get-song
// @Produce json
// @Param group query string true "Имя группы"
// @Param song query string true "Имя песни"
// @Param page query int false "Номер страницы"
// @Param pageSize query int false "Размер страницы"
// @Success 200 {object} Response
// @Failure 400 {object} map[string]string "group and song parameters are required"
// @Failure 500 {object} map[string]string "failed to get song"
// @Router /get_data/text [get]
func New(log *slog.Logger, getText GetText) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http_server.handlers.get_song.New"

		log.Info(fmt.Sprintf("op: %s", op))

		group := r.URL.Query().Get("group")
		song := r.URL.Query().Get("song")
		if group == "" || song == "" {
			utils.RenderCommonErr(errors.New("group and song parameters are required"),
				log, w, r, "group and song parameters are required", 400)
			return
		}

		pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1
		}

		songData, err := getText.GetSong(group, song)
		if err != nil {
			if errors.Is(err, storage.ErrSongNotFound) {
				utils.RenderCommonErr(err, log, w, r, "song not found", 500)
				return
			}
			utils.RenderCommonErr(err, log, w, r, "failed to get song", 500)
			return
		}

		songData = strings.ReplaceAll(songData, "\\n\\n", "\n\n")
		verses := strings.Split(songData, "\n\n")
		totalVerses := len(verses)
		totalPages := int(math.Ceil(float64(totalVerses) / float64(pageSize)))

		if page > totalPages {
			page = totalPages
		}

		start := (page - 1) * pageSize
		end := start + pageSize
		if end > totalVerses {
			end = totalVerses
		}

		paginatedVerses := verses[start:end]

		response := Response{
			Verses:           paginatedVerses,
			MaxVersesPerPage: pageSize,
			TotalPages:       totalPages,
			CurrentPage:      page,
		}

		log.Info("song get")

		render.JSON(w, r, response)
	}
}
