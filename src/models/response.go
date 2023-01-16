package models


type Response struct {
	Success    bool
	Message    string
	StatusCode string
	Data       map[string]string
}
type ResponseGetSpecificData struct {
	Success    bool
	Message    string
	StatusCode string
	Data      map[string]interface{}
}
type ResponseGetAllData struct {
	Success    bool
	Message    string
	StatusCode string
	Data      []interface{}
}