package get_all_data

import (
	"fmt"
	"log/slog"
	"math"
	resp "music_library/internal/http_server/lib/response"
	"music_library/internal/http_server/lib/utils"
	"music_library/internal/http_server/models"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

// GetDataLibrary представляет интерфейс для получения данных библиотеки.
// @Description Интерфейс для получения данных библиотеки.
type GetDataLibrary interface {
	// GetData получает данные библиотеки с фильтрацией и пагинацией.
	// @Description Получение данных библиотеки с фильтрацией по всем полям и пагинацией (метод GET).
	// @Param filter map[string]interface{} "Фильтры для поиска"
	// @Param page int "Номер страницы"
	// @Param pageSize int "Размер страницы"
	// @return []models.Data "Массив данных песен"
	// @return error "Ошибка выполнения"
	GetData(filter map[string]interface{}, page int, pageSize int) ([]models.Data, error)
	// GetCountSongs получает общее количество песен с применением фильтров.
	// @Description Получение общего количества песен с применением фильтров.
	// @Param filter map[string]interface{} "Фильтры для поиска"
	// @return int "Общее количество песен"
	// @return error "Ошибка выполнения"
	GetCountSongs(filter map[string]interface{}) (int, error)
}

// Response представляет структуру ответа с данными песен и информацией о пагинации.
// @Description Структура ответа с данными песен и информацией о пагинации.
type Response struct {
	Songs       []models.Data `json:"songs"`
	MaxPageSize int           `json:"maxPageSize"`
	TotalPages  int           `json:"totalPages"`
	CurrentPage int           `json:"currentPage"`
	TotalSongs  int           `json:"totalSongs"`
}

// New создает новый обработчик для получения данных библиотеки.
// @Summary Получение данных библиотеки
// @Description Получение данных библиотеки с фильтрацией по всем полям и пагинацией (метод GET).
// @ID get-all-data
// @Produce json
// @Param group query string false "Имя группы"
// @Param song query string false "Имя песни"
// @Param releaseDate query string false "Дата релиза"
// @Param text query string false "Текст песни"
// @Param link query string false "Ссылка на песню"
// @Param page query int false "Номер страницы"
// @Param pageSize query int false "Размер страницы"
// @Success 200 {object} Response
// @Failure 500 {object} map[string]string "failed to get songs"
// @Router /get_data/songs [get]
func New(log *slog.Logger, getSongs GetDataLibrary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http_server.handlers.get_all_data.New"

		log.Info(fmt.Sprintf("op=%s", op))

		filter := getFilter(r, "group", "song", "releaseDate", "text", "link")

		totalSongs, err := getSongs.GetCountSongs(filter)
		if err != nil {
			utils.RenderCommonErr(err, log, w, r, "failed to get songs", 500)
			return
		}

		if totalSongs == 0 {
			render.JSON(w, r, resp.Empty("no songs"))
			return
		}

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1
		}
		pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		totalPages := int(math.Ceil(float64(totalSongs) / float64(pageSize)))

		if page > totalPages {
			page = totalPages
		}

		songs, err := getSongs.GetData(filter, page, pageSize)
		if err != nil {
			utils.RenderCommonErr(err, log, w, r, "failed to get songs", 500)
			return
		}

		response := Response{
			Songs:       songs,
			MaxPageSize: pageSize,
			TotalPages:  totalPages,
			CurrentPage: page,
			TotalSongs:  totalSongs,
		}

		log.Info("songs get")

		render.JSON(w, r, response)
	}
}

// Преобразование параметров запроса в map
func getFilter(r *http.Request, params ...string) map[string]interface{} {
	filter := make(map[string]interface{}, 5)
	for _, elem := range params {
		value := r.URL.Query().Get(elem)
		if value != "" {
			filter[elem] = value
		}
	}
	return utils.ChangeKeys(&filter)
}
