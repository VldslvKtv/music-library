package addsong

import (
	"errors"
	"fmt"
	"log/slog"
	"music_library/internal/http_server/lib/logger"
	resp "music_library/internal/http_server/lib/response"
	"music_library/internal/http_server/lib/utils"
	"music_library/internal/http_server/mocks"
	"music_library/internal/http_server/models"
	"music_library/internal/http_server/storage"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

// AddNewSong представляет интерфейс для добавления новой песни.
// @Description Интерфейс для добавления новой песни.
type AddNewSong interface {
	// CreateSong создает новую песню.
	// @Description Создание новой песни.
	// @Param song body models.Data true "Данные песни"
	// @return error ошибка выполнения
	CreateSong(song models.Data) error
}

// FullData представляет структуру запроса с данными песни.
// @Description Структура запроса с данными песни.
type FullData struct {
	Data models.Data
}

// RequestSongAndGroup представляет структуру запроса с данными группы и песни.
// @Description Структура запроса с данными группы и песни.
type RequestSongAndGroup struct {
	Song models.SongAndGroup
}

// RequestDetails представляет структуру запроса с данными о деталях песни.
// @Description Структура запроса с данными о деталях песни.
type RequestDetails struct {
	SongDetails models.SongDetails
}

// New создает новый обработчик для добавления новой песни (метод POST).
// @Summary Добавление новой песни
// @Description Добавление новой песни в формате JSON и запросе к внешнему API.
// @ID add-song
// @Accept json
// @Produce json
// @Param data body RequestSongAndGroup true "Данные песни"
// @Success 200 {object} FullData
// @Failure 400 {object} map[string]string "failed to decode req-body"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /add [post]
func New(log *slog.Logger, apiURL string, addSong AddNewSong) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http_server.handlers.add_song.New"

		log.Info(fmt.Sprintf("op: %s", op))

		var req RequestSongAndGroup

		err := render.DecodeJSON(r.Body, &req.Song) // распарсим запрос
		if err != nil {
			utils.RenderCommonErr(err, log, w, r, "failed to decode req-body", 400)
			return // обязательно выйти тк render.JSON не прервет выполнение
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validatorErr := err.(validator.ValidationErrors)
			log.Error("invalid request", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validatorErr)) // делаем ответ более читаемым
			return
		}

		// Для тестирования сделал мок:

		mocckClient := mocks.NewMockClient(`{ "releaseDate": 
		"16.07.2006", 
		"text": "Sample song lyrics\n\nMore lyrics...", 
		"link": "https://example.com" }`,
			http.StatusOK, nil)
		details, err := mocckClient.Get(fmt.Sprintf("%s/info?group=%s&song=%s", apiURL, req.Song.Group, req.Song.Song))

		// details, err := http.Get(fmt.Sprintf("%s/info?group=%s&song=%s", apiURL, req.Song.Group, req.Song.Song))
		if err != nil {
			log.Error("failed to get details", logger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal server error"))
			return
		}

		defer details.Body.Close()

		if details.StatusCode != http.StatusOK {
			if details.StatusCode == http.StatusBadRequest {
				log.Error("bad request (external api)", logger.Err(err))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.Error("bad request"))
				return
			} else {
				log.Error("internal server error (external api)", logger.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, resp.Error("internal server error"))
				return
			}
		}

		var detailsData RequestDetails

		err = render.DecodeJSON(details.Body, &detailsData.SongDetails) // распарсим запрос
		if err != nil {
			log.Error("failed to decode req-body (external api)", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode req"))
			return
		}

		log.Info("request body decoded 2", slog.Any("request", detailsData))

		full := FullData{
			Data: models.Data{
				SongAndGroup: req.Song,
				SongDetails:  detailsData.SongDetails,
			},
		}

		log.Info("request body decoded 3", slog.Any("request", full))

		err = addSong.CreateSong(full.Data)
		if err != nil {
			if errors.Is(err, storage.ErrGroupExists) {
				utils.RenderCommonErr(err, log, w, r, "group already exists", 500)
				return
			} else if errors.Is(err, storage.ErrSongExists) {
				utils.RenderCommonErr(err, log, w, r, "song already exists", 500)
				return
			}
			utils.RenderCommonErr(err, log, w, r, "failed to add song", 500)
			return
		}

		log.Info("Song is add")
		render.JSON(w, r, full)
	}
}
