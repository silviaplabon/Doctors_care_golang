package utils

import (
	"bytes"
	"context"
	"doctors_care/src/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/iterator"
)

var FBClient *firestore.Client

var Ctx context.Context

var Body []byte

func APICallWithSelfManagedBuffers(url string, method string, c echo.Context, chnl chan<- models.Response, wg *sync.WaitGroup) {

	defer wg.Done()

	var body []byte // This buffer could be acquired from a custom buffer pool

	var statusCode string

	var err error

	if method == "GET" || method == "DELETE" {

		fmt.Println(method + " " + url)

		req, err := http.NewRequest(method, url, nil)

		fmt.Println("1", err, req, statusCode)

		if err != nil {

			resp := models.Response{
				Success:    false,
				Message:    err.Error(),
				StatusCode: statusCode,
			}

			chnl <- resp

			return

		}

		// req.Header.Set("Accept", "application/json")
		// req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		req.Header.Set("Authorization", "")

		client := &http.Client{}
		resp, err := client.Do(req)

		fmt.Println("2", err, resp, statusCode)

		if err != nil {

			resp := models.Response{
				Success:    false,
				Message:    err.Error(),
				StatusCode: statusCode,
			}

			chnl <- resp

			return

		}

		defer resp.Body.Close()

		statusCode = resp.Status

		body, _ = ioutil.ReadAll(resp.Body)

	}

	if method == "POST" || method == "PUT" {

		b, err := ioutil.ReadAll(c.Request().Body)
		Body = b
		fmt.Println("b:", string(b))
		fmt.Println(string(b), "b data of post")
		if err != nil {

			resp := models.Response{
				Success:    false,
				Message:    err.Error(),
				StatusCode: statusCode,
			}

			chnl <- resp

			return

		}

		fmt.Println("method and url:", method, url)
		req, err := http.NewRequest(method, url, bytes.NewBuffer(b))

		if err != nil {

			resp := models.Response{
				Success:    false,
				Message:    err.Error(),
				StatusCode: statusCode,
			}
			fmt.Println(err, "err data")

			chnl <- resp

			return

		}

		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		req.Header.Set("Accept", "application/json")

		req.Header.Set("Authorization", "")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {

			resp := models.Response{
				Success:    false,
				Message:    err.Error(),
				StatusCode: statusCode,
			}

			chnl <- resp

			return

		}

		defer resp.Body.Close()

		statusCode = resp.Status

		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println(resp.Body)

	}

	fmt.Println("3", err, statusCode)

	if err != nil {

		resp := models.Response{
			Success:    false,
			Message:    err.Error(),
			StatusCode: statusCode,
		}

		chnl <- resp

		return

	}

	fmt.Println("4", err, statusCode)

	if statusCode != "200 OK" && statusCode != "204 No Content" && statusCode != "201 Created" {

		fmt.Println("5", err, statusCode)

		resp := models.Response{
			Success:    false,
			Message:    "Failure",
			StatusCode: statusCode,
		}

		chnl <- resp

		return

	}

	fmt.Println("6", err, statusCode)

	d := make(map[string]string)
	d["data"] = string(body)
	resp := models.Response{
		Success:    true,
		Message:    "Success",
		StatusCode: statusCode,
		Data:       d,
	}

	chnl <- resp

}

func GetSpecificData( chnl chan<- models.ResponseGetSpecificData, wg *sync.WaitGroup, databaseName string, id string) {
	defer wg.Done()
	dsnap, err := FBClient.Collection(databaseName).Doc(id).Get(Ctx)
	if err != nil {
		resp := models.ResponseGetSpecificData{
			Success:    false,
			Message:    err.Error(),
			StatusCode: "204 No Content",
		}
		chnl <- resp
		return
	} else {
		resp := models.ResponseGetSpecificData{
			Success:    true,
			Message:    "data",
			StatusCode: "200 OK",
			Data:       dsnap.Data(),
		}
		chnl <- resp
		return
	}
}

func GetAllData(c echo.Context, chnl chan<- models.ResponseGetAllData, wg *sync.WaitGroup, databaseName string) {
	defer wg.Done()
	arr := make([]interface{}, 0)
	iter := FBClient.Collection(databaseName).Documents(Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			resp := models.ResponseGetAllData{
				Success:    false,
				Message:    err.Error(),
				StatusCode: "500",
			}
			chnl <- resp
			return
		}
		arr = append(arr, doc.Data())
	}

	if len(arr) < 1 {
		resp := models.ResponseGetAllData{
			Success:    false,
			Message:    "No Record Found",
			StatusCode: "204 No Content",
		}
		chnl <- resp
		return
	} else {
		resp := models.ResponseGetAllData{
			Success:    true,
			Message:    "data",
			StatusCode: "200 OK",
			Data:       arr,
		}
		chnl <- resp
		return

	}
}

func UpdateSpecificData(c echo.Context, data map[string]interface{}, chnl chan<- models.ResponseGetSpecificData, wg *sync.WaitGroup, databaseName string, id string) {

	ref := FBClient.Collection(databaseName).Doc(id)
	_, err := ref.Set(Ctx, data, firestore.MergeAll)
	if err != nil {
		fmt.Println("An error has occurred:" + err.Error())
		resp := models.ResponseGetSpecificData{
			Success:    false,
			Message:    err.Error(),
			StatusCode: "500",
		}
		chnl <- resp
		return
	} else {
		resp := models.ResponseGetSpecificData{
			Success:    true,
			Message:    "data updated",
			StatusCode: "200 OK",
		}
		chnl <- resp
		return
	}
}
func DeleteSpecificData(c echo.Context, data map[string]interface{}, chnl chan<- models.ResponseGetSpecificData, wg *sync.WaitGroup, databaseName string, id string) {

	_, err := FBClient.Collection("merchants").Doc(id).Delete(Ctx)
	if err != nil {
		resp := models.ResponseGetSpecificData{
			Success:    false,
			Message:    err.Error(),
			StatusCode: "500",
		}
		chnl <- resp
		return
	} else {
		resp := models.ResponseGetSpecificData{
			Success:    true,
			Message:    "data deleted successfuly",
			StatusCode: "200 OK",
		}
		chnl <- resp
		return
	}
}
func AddSpecificData(c echo.Context, data map[string]interface{}, chnl chan<- models.ResponseGetSpecificData, wg *sync.WaitGroup, databaseName string, id string) {
	defer wg.Done()
	ref := FBClient.Collection("merchants").NewDoc()
	_, err := ref.Set(Ctx, data)
	if err != nil {
		resp := models.ResponseGetSpecificData{
			Success:    false,
			Message:    err.Error(),
			StatusCode: "204 No Content",
		}
		chnl <- resp
		return
	} else {
		resp := models.ResponseGetSpecificData{
			Success:    true,
			Message:    "data",
			StatusCode: "200 OK",
		}
		chnl <- resp
		return
	}
}
