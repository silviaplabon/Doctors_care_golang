package handlers

import (
	"doctors_care/src/models"
	"doctors_care/src/utils"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

func GetADoctor(c echo.Context) error {
	var wg sync.WaitGroup
	chnl := make(chan models.ResponseGetSpecificData)
	wg.Add(1)
	go utils.GetSpecificData(chnl, &wg,"doctors","PRd5l3ld6mMpLvqtv4bc");
	resp:=<-chnl
	go func() {
		wg.Wait()
		close(chnl)
		fmt.Println("POST Channel Closed")
	}()

	return c.JSON(http.StatusOK, resp)
}