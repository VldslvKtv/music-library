package update_song

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

	"github.com/go-playground/validator"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// UpdateSong представляет интерфейс для обновления данных песни.
// @Description Интерфейс для обновления данных песни.
type UpdateSong interface {
	// PatchSong изменяет данные песни по ID.
	// @Description Изменение данных песни по ID.
	// @Param idSong int ID песни
	// @Param data models.Data данные песни
	// @return error ошибка выполнения
	PatchSong(idSong int, data models.Data) error
}

// New создает новый обработчик для изменения данных песни (метод PATCH).
// @Summary Изменение данных песни
// @Description Изменение данных песни по ID.
// @ID update-song
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param data body models.Data true "Данные песни"
// @Success 200 {object} map[string]string "ok"
// @Failure 400 {object} map[string]string "failed to decode req-body"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /update/{id} [patch]
func New(log *slog.Logger, updateSong UpdateSong) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http_server.handlers.update_song.New"

		log.Info(fmt.Sprintf("op: %s", op))

		id, err := utils.CheckID(chi.URLParam(r, "id"))
		if err != nil {
			utils.RenderCommonErr(err, log, w, r, "invalid ID", 400)
			return
		}

		var req models.Data

		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			utils.RenderCommonErr(err, log, w, r, "failed to decode req-body", 400)
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validatorErr := err.(validator.ValidationErrors)
			log.Error("invalid request", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validatorErr)) // делаем ответ более читаемым
			return
		}

		log.Debug("request body decoded", slog.Any("request", req))

		err = updateSong.PatchSong(id, req)
		if err != nil {
			if errors.Is(err, storage.ErrGroupExists) {
				utils.RenderCommonErr(err, log, w, r, "group already exists", 500)
				return
			} else if errors.Is(err, storage.ErrSongExists) {
				utils.RenderCommonErr(err, log, w, r, "song already exists", 500)
				return
			}
			log.Error("error update data", logger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal server error"))
			return
		}

		log.Info("Song is update")
		render.JSON(w, r, resp.OK())
	}
}
