package megaplan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

var MP *MegaPlan
var client = &http.Client{
	Timeout: time.Second * 10,
}

func (mp *MegaPlan) doRequest(method, url string, body interface{}, response interface{}) error {
	jsonData, err := json.Marshal(body)
	req, err := http.NewRequest(method, mp.Url+url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header["AUTHORIZATION"] = []string{mp.getToken()}

	res, err := client.Do(req)
	resBody, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("megaplan(%d): %w\n\nResponse:\n\n %s", res.StatusCode, err, string(resBody))
	}

	defer res.Body.Close()
	json.Unmarshal(resBody, &response)

	return nil
}

func (mp *MegaPlan) Get(url string) *Response {
	var res Response
	if err := mp.doRequest("GET", url, nil, &res); err != nil {
		panic(err)
	}

	return &res
}

func (mp *MegaPlan) MustAuthenticateWithPassword(auth *AuthOpt) *MegaPlan {

	setField := func(w *multipart.Writer, fieldName string, value string) error {
		fw, err := w.CreateFormField(fieldName)
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, strings.NewReader(value))
		if err != nil {
			return err
		}

		return nil
	}

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	if err := setField(w, "username", auth.Username); err != nil {
		log.Fatalf("Error while authenticating megaplan: %v", err)
	}
	if err := setField(w, "password", auth.Password); err != nil {
		log.Fatalf("Error while authenticating megaplan: %v", err)
	}
	if err := setField(w, "grant_type", auth.GrantType); err != nil {
		log.Fatalf("Error while authenticating megaplan: %v", err)
	}
	w.Close()
	req, err := http.NewRequest("POST", mp.Url+auth.AccessTokenUrl, bytes.NewReader(body.Bytes()))
	if err != nil {
		log.Fatalf("Error while authenticating megaplan: %v", err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	res, err := client.Do(req)
	if res.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", res.StatusCode)
		log.Fatalf("Error while authenticating megaplan: %v", err)
	}
	token := Token{}
	json.NewDecoder(res.Body).Decode(&token)
	mp.Token = token
	return mp
}
