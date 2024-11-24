package addsong

import (
	"errors"
	"fmt"
	"log/slog"
	"music_library/internal/http_server/lib/logger"
	resp "music_library/internal/http_server/lib/response"
	"music_library/internal/http_server/lib/utils"
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
	// @Param song models.Data данные песни
	// @return error ошибка выполнения
	CreateSong(song models.Data) error
}

// Request представляет структуру запроса с данными песни.
// @Description Структура запроса с данными песни.
type Request struct {
	Data models.Data
}

// New создает новый обработчик для добавления новой песни (метод POST).
// @Summary Добавление новой песни
// @Description Добавление новой песни в формате JSON и запросе к внешнему API.
// @ID add-song
// @Accept json
// @Produce json
// @Param data body Request true "Данные песни"
// @Success 200 {object} Request
// @Failure 400 {object} map[string]string "failed to decode req-body"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /add [post]
func New(log *slog.Logger, apiURL string, addSong AddNewSong) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http_server.handlers.add_song.New"

		log.Info(fmt.Sprintf("op: %s", op))

		var req Request

		err := render.DecodeJSON(r.Body, &req.Data.SongAndGroup) // распарсим запрос
		if err != nil {
			utils.RenderCommonErr(log, w, r, "failed to decode req-body", 400)
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
		// mocckClient := mocks.NewMockClient(`{ "releaseDate":
		//  "16.07.2006",
		//  "text": "Sample song lyrics\n\nMore lyrics...",
		//   "link": "https://example.com" }`,
		//  http.StatusOK, nil)
		// details, err := mocckClient.Get(fmt.Sprintf("%s/info?group=%s&song=%s", apiURL, req.Data.Group, req.Data.Song))

		details, err := http.Get(fmt.Sprintf("%s/info?group=%s&song=%s", apiURL, req.Data.Group, req.Data.Song))
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

		err = render.DecodeJSON(details.Body, &req.Data.SongDetails) // распарсим запрос
		if err != nil {
			log.Error("failed to decode req-body (external api)", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode req"))
			return
		}

		log.Info("request body decoded 2", slog.Any("request", req))

		err = addSong.CreateSong(req.Data)
		if err != nil {
			if errors.Is(err, storage.ErrGroupExists) {
				utils.RenderCommonErr(log, w, r, "group already exists", 500)
				return
			} else if errors.Is(err, storage.ErrSongExists) {
				utils.RenderCommonErr(log, w, r, "song already exists", 500)
				return
			}
			utils.RenderCommonErr(log, w, r, "failed to add song", 500)
			return
		}

		log.Info("Song is add")
		render.JSON(w, r, req)
	}
}
