package controller

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type HealthCheck struct {
	Result struct {
		DB string `json:"DB,omitempty" example:"ok"`
	} `json:"result"`
}

// @Summary Проверка сервера
// @Description Запускает процесс проверки всех зависимостей, необходимых для корректной работы сервера
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @tags Служебные
// @Param endpoint path string true "Идентификатор слушателя" default(api)
// @Param version path int true "Версия API" default(1)
// @Success 200 {object} HealthCheck
// @Router /{endpoint}/v{version}/healthcheck [get]
func (c *Controller) healthcheck(ctx echo.Context) error {
	healthcheck := make(map[string]string)
	mx := sync.Mutex{}
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		err := c.app.PingDatabase()

		mx.Lock()
		if err != nil {
			healthcheck["DB"] = err.Error()
		} else {
			healthcheck["DB"] = "ok"
		}
		mx.Unlock()
		wg.Done()
	}()

	wg.Wait()
	return ctx.JSON(http.StatusOK, echo.Map{"result": healthcheck})
}
