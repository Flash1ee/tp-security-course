package repository

import (
	"fmt"
	"net/http"

	"http-proxy/internal/proxy/models"
)

func FormRequestData(r *http.Request, dump []byte) *models.Request {
	req := &models.Request{
		Method: r.Method,
		Path:   r.URL.Path,
		Raw:    string(dump),
	}
	getParams := models.Map{}
	for key, value := range r.URL.Query() {
		getParams[key] = getValue(value)
	}
	req.GetParams = getParams

	headers := models.Map{
		"Host": r.Host,
	}

	for key, value := range r.Header {
		if key == "Cookie" {
			continue
		}
		headers[key] = getValue(value)
	}
	req.Headers = headers

	cookies := models.Map{}

	for _, value := range r.Cookies() {
		cookies[value.Name] = value.Value
	}
	req.Cookies = cookies

	postParams := models.Map{}
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	for key, value := range r.PostForm {
		postParams[key] = getValue(value)
	}
	req.PostParams = postParams

	return req
}
func FormResponseData(response *http.Response, body string) *models.Response {
	if response == nil {
		return nil
	}
	res := &models.Response{
		Code:    response.StatusCode,
		Message: response.Status,
	}

	headers := models.Map{}

	for key, value := range response.Header {
		if key == "Cookie" {
			continue
		}
		headers[key] = getValue(value)
	}
	res.Headers = headers
	res.Body = body

	return res

}
func getValue(value []string) interface{} {
	if len(value) == 1 {
		return value[0]
	}
	return value
}
