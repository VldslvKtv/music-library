package utils

import (
	"fmt"
	"log/slog"
	resp "music_library/internal/http_server/lib/response"
	"music_library/internal/http_server/models"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/render"
)

func CheckID(param string) (int, error) {
	var value int
	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("error convert to int: %v", err)
	}
	return value, nil
}

func ConvertStruct(s models.Data) map[string]interface{} {
	newMap := make(map[string]interface{})
	v := reflect.ValueOf(s)

	// Проходим по всем полям структуры, включая вложенные структуры
	convertStructHelper(v, newMap)

	return ChangeKeys(&newMap)
}

func convertStructHelper(v reflect.Value, newMap map[string]interface{}) {
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := v.Type().Field(i)
			fieldName := strings.ToLower(fieldType.Name)

			// Если значение поля не нулевое, добавляем его в map
			if !isZero(field) {
				if field.Kind() == reflect.Struct {
					if field.Type() == reflect.TypeOf(models.CustomTime{}) {
						// Обработка CustomTime
						customTime := field.Interface().(models.CustomTime)
						newMap[fieldName] = customTime.Time.Format(models.CustomTimeFormat)
					} else {
						// Рекурсивно обрабатываем вложенные структуры
						convertStructHelper(field, newMap)
					}
				} else {
					newMap[fieldName] = field.Interface()
				}
			}
		}
	}
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Struct:
		return v.Interface() == reflect.Zero(v.Type()).Interface()
	default:
		return false
	}
}

func ChangeKeys(m *map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for key, value := range *m {
		switch key {
		case "group":
			newMap["groups.name"] = value
		case "song":
			newMap["songs.name"] = value
		case "releasedate":
			newMap["release_date"] = value
		default:
			newMap[key] = value
		}
	}
	return newMap
}

func RenderCommonErr(log *slog.Logger, w http.ResponseWriter, r *http.Request, text string, statusCode int) {

	log.Error(text)
	if statusCode == 400 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	render.JSON(w, r, resp.Error(text))
}
