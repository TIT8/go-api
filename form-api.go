package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

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

type GForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
	File    string `json:"file"`
}

func getUrl() string {

	return fmt.Sprintf("https://api.telegram.org/bot%s", os.Getenv("TELEGRAM_TOKEN"))

}

func SendMessage(text string) (bool, error) {

	var err error
	var response *http.Response

	url := fmt.Sprintf("%s/sendMessage", getUrl())
	body, err := json.Marshal(map[string]string{
		"chat_id": os.Getenv("CHAT_ID"),
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

	r.Body = http.MaxBytesReader(w, r.Body, 100)

	fmt.Fprintf(w, "Ciao da %s\nList of IP addresses received (including yours, if visible):\n%s\n", r.Host, r.Header.Get("X-Forwarded-For"))

}

func handler_sheet(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 5*1024)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	content_type, err := regexp.MatchString("application/json", r.Header.Get("Content-Type"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !content_type {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var data GForm
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	text := fmt.Sprintf("New message from Google Form!\n\nName:\t%s\n\nE-mail:\t%s\n\nMessage:\t%s\n\nFile link:\t%s", data.Name, data.Email, data.Message, data.File)
	result, err := SendMessage(text)
	if !result {
		log.Printf("Error sending telegram messag, %s\n", err)
		w.WriteHeader(http.StatusNotAcceptable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func handler_post(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 50*1024)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	content_type, err := regexp.MatchString("multipart/form-data", r.Header.Get("Content-Type"))
	if err != nil {
		log.Println(err)
	}
	log.Printf("Content-Type is multipart/form-data?  %v\n", content_type)

	if content_type {
		err = r.ParseMultipartForm(1000)
	} else {
		err = r.ParseForm()
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
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
	data.Add("secret", os.Getenv("CAPTCHA_SECRET"))
	data.Add("sitekey", os.Getenv("CAPTCHA_SITEKEY"))

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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	post := Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//log.Printf("Captcha validation: %v", post.Success)

	var ok string
	if post.Success {

		text := fmt.Sprintf("New message from HTML form!\n\nName:\t%s\n\nE-mail:\t%s\n\nMessage:\t%s\n", req.Name, req.Email, req.Message)

		result, err := SendMessage(text)
		if !result {
			log.Printf("Error sending telegram messag, %s\n", err)
			w.WriteHeader(http.StatusNotAcceptable)
			ok = "Valid CAPTCHA but error on Telegram request"
		} else {
			w.WriteHeader(http.StatusOK)
			ok = "Valid CAPTCHA and Telegram sent"
		}

	} else {
		ok = "CAPTCHA validation failed"
		w.WriteHeader(http.StatusNotAcceptable)
	}

	response := Response{ok}
	res, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}
	//log.Println(string(res))

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}

func main() {

	http.HandleFunc("/", handler_get)
	http.HandleFunc("/post", handler_post)
	http.HandleFunc("/sheet", handler_sheet)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
