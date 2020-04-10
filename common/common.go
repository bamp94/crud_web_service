package common

import (
	"errors"
	"net/http"
	"regexp"
)

const (
	MethodGet    = "get"
	MethodPost   = "ins"
	MethodPut    = "upd"
	MethodDelete = "del"
)

var (
	ErrInternal = errors.New("Внутренняя ошибка сервера, повторите попытку позже или обратитесь к системному администратору")

	Methods = map[string]string{
		http.MethodGet:    MethodGet,
		http.MethodPost:   MethodPost,
		http.MethodPut:    MethodPut,
		http.MethodDelete: MethodDelete,
	}
)

// FindUrlPathObjects finds parts of path with object name and id
func FindUrlPathObjects(urlPath string) [][]string {
	re := regexp.MustCompile(`[0-9a-zA-Zа-яА-Я]*/[0-9]*`)
	return re.FindAllStringSubmatch(urlPath, -1)
}
