package delete_song

import (
	"errors"
	"fmt"
	"log/slog"
	"music_library/internal/http_server/lib/logger"
	resp "music_library/internal/http_server/lib/response"
	"music_library/internal/http_server/lib/utils"
	"music_library/internal/http_server/storage"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// DeleteSong представляет интерфейс для удаления песни.
// @Description Интерфейс для удаления песни.
type DeleteSong interface {
	// DeleteSong удаляет песню по ID.
	// @Description Удаление песни по ID.
	// @Param idSong int ID песни
	// @return error ошибка выполнения
	DeleteSong(idSong int) error
}

// New создает новый обработчик для удаления песни (метод DELETE).
// @Summary Удаление песни
// @Description Удаление песни по ID.
// @ID delete-song
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {object} map[string]string "ok"
// @Failure 400 {object} map[string]string "invalid ID"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /delete/{id} [delete]
func New(log *slog.Logger, deleteSong DeleteSong) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http_server.handlers.delete_song.New"

		log.Info(fmt.Sprintf("op: %s", op))

		id, err := utils.CheckID(chi.URLParam(r, "id"))
		if err != nil {
			utils.RenderCommonErr(log, w, r, "invalid ID", 400)
			return
		}

		err = deleteSong.DeleteSong(id)
		if err != nil {
			if errors.Is(err, storage.ErrSongNotFound) {
				utils.RenderCommonErr(log, w, r, "song not found", 500)
				return
			}
			log.Error("error deleting", logger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal server error"))
			return
		}

		log.Info("Song is delete")
		render.JSON(w, r, resp.OK())
	}
}
