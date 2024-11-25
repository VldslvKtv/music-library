package add_song

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
	// @Param song body models.Data true "Данные песни"
	// @return error ошибка выполнения
	CreateSong(song models.Data) error
}

// New создает новый обработчик для добавления новой песни (метод POST).
// @Summary Добавление новой песни
// @Description Добавление новой песни в формате JSON и запросе к внешнему API.
// @ID add-song
// @Accept json
// @Produce json
// @Param song body models.SongAndGroup true "Данные песни"
// @Success 200 {object} models.Data
// @Failure 400 {object} map[string]string "failed to decode req-body"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /add [post]
func New(log *slog.Logger, apiURL string, addSong AddNewSong) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http_server.handlers.add_song.New"

		log.Info(fmt.Sprintf("op: %s", op))

		// Получаем название песни и группу
		var req models.SongAndGroup

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			utils.RenderCommonErr(err, log, w, r, "failed to decode req-body", 400)
			return
		}

		log.Debug("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validatorErr := err.(validator.ValidationErrors)
			log.Error("invalid request", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validatorErr)) // делаем ответ более читаемым
			return
		}

		// Для тестирования сделал мок:

		// mocckClient := mocks.NewMockClient(`{ "releaseDate":
		// "16.07.2006",
		// "text": "Sample song lyrics\n\nMore lyrics...",
		// "link": "https://example.com" }`,
		// 	http.StatusOK, nil)
		// details, err := mocckClient.Get(fmt.Sprintf("%s/info?group=%s&song=%s", apiURL, req.Group, req.Song))

		details, err := http.Get(fmt.Sprintf("%s/info?group=%s&song=%s", apiURL, req.Group, req.Song))
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

		// Получаем расширенную информацию о песне
		var detailsData models.SongDetails

		err = render.DecodeJSON(details.Body, &detailsData)
		if err != nil {
			log.Error("failed to decode req-body (external api)", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode req"))
			return
		}

		log.Debug("request body decoded (ext api)", slog.Any("request", detailsData))

		// Объединяем данные перед отправкой в БД
		fullData := models.Data{
			SongAndGroup: req,
			SongDetails:  detailsData,
		}

		log.Debug("full data", slog.Any("data", fullData))

		err = addSong.CreateSong(fullData)
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
		render.JSON(w, r, fullData)
	}
}
