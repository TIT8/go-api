package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

const Token string = "6169685035:AAEgNi4pC5gARzCiMlvkDTFIEOOClD6wHB0"
const ChatId string = "-1001888191995"

type Response struct {
	Response string `json:"response"`
}

type Post struct {
	Success bool `json:"success"`
}

type Request struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Message string `json:"message"`
	Captcha string `json:"frc-captcha-solution"`
}

func getUrl() string {

	return fmt.Sprintf("https://api.telegram.org/bot%s", Token)

}

func SendMessage(text string) (bool, error) {

	var err error
	var response *http.Response

	url := fmt.Sprintf("%s/sendMessage", getUrl())
	body, err := json.Marshal(map[string]string{
		"chat_id": ChatId,
		"text":    text,
	})
	if err != nil {
		return false, err
	}

	response, err = http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return false, err
	}

	defer response.Body.Close()

	return true, nil

}

func handler_get(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Ciao %s\n", r.Host)

}

func handler_post(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	content_type, err := regexp.MatchString("multipart/form-data", r.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Content-Type is multipart/form-data?  %v\n", content_type)

	if content_type {
		err = r.ParseMultipartForm(100)
	} else {
		err = r.ParseForm()
	}
	if err != nil {
		log.Fatal(err)
	}

	req := Request{
		Email:   r.FormValue("email"),
		Name:    r.FormValue("name"),
		Message: r.FormValue("message"),
		Captcha: r.FormValue("frc-captcha-solution"),
	}
	//log.Println(req)

	data := url.Values{}
	data.Add("solution", req.Captcha)
	data.Add("secret", "A1UGN12VU21PUJCEUJNDHHTP0CD835IGMDUO3IS0JVBDUBUUVQJ584DPD1")
	data.Add("sitekey", "FCMTGCV10AMHV9QE")

	posturl := "https://api.friendlycaptcha.com/api/v1/siteverify"

	resp, err := http.PostForm(posturl, data)
	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Formatting body failed: %s", err)
		return
	}

	post := Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}
	log.Printf("Captcha validation: %v", post.Success)

	var ok string
	if post.Success {

		text := fmt.Sprintf("New message!\n\nName:\t%s\n\nE-mail:\t%s\n\nMessage:\t%s\n", req.Name, req.Email, req.Message)
		result, err := SendMessage(text)
		if !result {
			log.Printf("Error sending telegram messag, %s\n", err)
			w.WriteHeader(http.StatusNotAcceptable)
		} else {
			ok = "ok"
		}

	} else {

		ok = "not ok"
		w.WriteHeader(http.StatusNotAcceptable)

	}

	response := Response{ok}
	res, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(res))

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}

func main() {

	http.HandleFunc("/", handler_get)
	http.HandleFunc("/post", handler_post)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
