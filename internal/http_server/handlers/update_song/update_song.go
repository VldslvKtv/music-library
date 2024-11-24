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

// Request представляет структуру запроса с данными песни.
// @Description Структура запроса с данными песни.
type Request struct {
	song models.Data
}

// New создает новый обработчик для изменения данных песни (метод PATCH).
// @Summary Изменение данных песни
// @Description Изменение данных песни по ID.
// @ID update-song
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param data body Request true "Данные песни"
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
			utils.RenderCommonErr(log, w, r, "invalid ID", 400)
			return
		}

		var req Request
		err = render.DecodeJSON(r.Body, &req.song) // распарсим запрос
		if err != nil {
			utils.RenderCommonErr(log, w, r, "failed to decode req-body", 400)
			return // обязательно выйти тк render.JSON не прервет выполнение
		}

		log.Info("request body decoded", slog.Any("request", req.song))

		err = updateSong.PatchSong(id, req.song)
		if err != nil {
			if errors.Is(err, storage.ErrGroupExists) {
				utils.RenderCommonErr(log, w, r, "group already exists", 500)
				return
			} else if errors.Is(err, storage.ErrSongExists) {
				utils.RenderCommonErr(log, w, r, "song already exists", 500)
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
