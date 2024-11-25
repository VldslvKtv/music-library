package update_song

import (
	"context"
	"encoding/json"
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
	// @Param ctx context.Context Контекст выполнения запроса
	// @Param idSong int ID песни
	// @Param data models.Data данные песни
	// @return error ошибка выполнения
	PatchSong(ctx context.Context, idSong int, data models.Data) error
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
// @Failure 400 {object} map[string]string "failed to decode req-body or any other errors"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /songs/{id} [patch]
func New(log *slog.Logger, updateSong UpdateSong) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http_server.handlers.update_song.New"
		ctx := r.Context()

		log.Info(fmt.Sprintf("op: %s", op))

		id, err := utils.CheckID(chi.URLParam(r, "id"))
		if err != nil {
			utils.RenderCommonErr(err, log, w, r, "invalid ID", 400)
			return
		}

		var bodyMap map[string]interface{}
		err = json.NewDecoder(r.Body).Decode(&bodyMap)
		if err != nil {
			utils.RenderCommonErr(err, log, w, r, "failed to decode req-body", 400)
			return
		}
		log.Debug("request body-map", slog.Any("req", bodyMap))

		if song, exists := bodyMap["song"]; exists && song == "" {
			utils.RenderCommonErr(errors.New("song cannot be empty"), log, w, r, "song cannot be empty", 400)
			return
		}
		if group, exists := bodyMap["group"]; exists && group == "" {
			utils.RenderCommonErr(errors.New("group cannot be empty"), log, w, r, "group cannot be empty", 400)
			return
		}

		bodyBytes, err := json.Marshal(bodyMap)
		if err != nil {
			utils.RenderCommonErr(err, log, w, r, "failed to marshal req-body", 500)
			return
		}

		var req models.Data
		err = json.Unmarshal(bodyBytes, &req)
		if err != nil {
			utils.RenderCommonErr(err, log, w, r, "failed to decode req-body", 400)
			return
		}

		log.Debug("request body decoded", slog.Any("request", req))

		err = updateSong.PatchSong(ctx, id, req)
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

		log.Info("Song is updated")
		render.JSON(w, r, resp.OK())
	}
}
