package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func GetRequestBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var requestBody T
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	return &requestBody, nil
}

func HttpGet[T any](url string) (T, error) {
	var data T
	resp, err := http.Get(url)
	if err != nil {
		return data, err
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, err
	}
	resp.Body.Close()
	return data, nil
}
