package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"crud_web_service/common"
	"crud_web_service/model"

	"github.com/labstack/echo/v4"
)

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
			healthcheck["DB"] = ok
		}
		mx.Unlock()
		wg.Done()
	}()

	wg.Wait()
	return ctx.JSON(http.StatusOK, echo.Map{"result": healthcheck})
}

func (c *Controller) getDBResult(ctx echo.Context) error {
	if common.Methods[ctx.Request().Method] == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": errWrongHTTPMethod})
	}
	urlPath := strings.TrimLeft(ctx.Request().URL.Path, "/"+c.config.DB.Endpoint+"/v1")
	pathParamsSlices := common.FindUrlPathObjects(urlPath)
	var objects []model.Object
	for _, pathParams := range pathParamsSlices {
		for _, object := range pathParams {
			objectInfo := strings.Split(object, "/")
			if len(objectInfo) == 0 || objectInfo[0] == "" {
				continue
			}
			if len(objectInfo) == 1 || objectInfo[1] == "" {
				objects = append(objects, model.Object{Name: objectInfo[0], ID: "0"})
				continue
			}
			objects = append(objects, model.Object{Name: objectInfo[0], ID: objectInfo[1]})
		}
	}
	params := model.QueryParameters{
		Method:  common.Methods[ctx.Request().Method],
		Objects: objects,
	}
	if ctx.Request().Method == http.MethodPut || ctx.Request().Method == http.MethodPost {
		jsonMap := make(map[string]interface{})
		err := json.NewDecoder(ctx.Request().Body).Decode(&jsonMap)
		if err != nil {
			return err
		}
		jsonString, _ := json.Marshal(jsonMap)
		params.Body = string(jsonString)
	}
	response, err := c.app.GetDBResult(params)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, response)
}
