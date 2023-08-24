package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Post struct {
	Success bool `json:"success"`
}

var email string
var message string
var name string

var Token string = "6169685035:AAEgNi4pC5gARzCiMlvkDTFIEOOClD6wHB0"
var ChatId string = "-1001888191995"

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func getUrl() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", Token)
}

func SendMessage(text string) (bool, error) {
	// Global variables
	var err error
	var response *http.Response

	// Send the message
	url := fmt.Sprintf("%s/sendMessage", getUrl())
	body, _ := json.Marshal(map[string]string{
		"chat_id": ChatId,
		"text":    text,
	})
	response, err = http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return false, err
	}

	// Close the request at the end
	defer response.Body.Close()

	// Body
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	// Log
	fmt.Printf("Message '%s' was sent\n", text)
	fmt.Printf("Response JSON: %s\n", string(body))

	// Return
	return true, nil
}

func post_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	r.ParseForm()
	email = r.Form["email"][0]
	message = r.Form["message"][0]
	name = r.Form["name"][0]
	solution := r.Form["frc-captcha-solution"][0]
	fmt.Println(solution)

	data := url.Values{}
	data.Add("solution", solution)
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
	// Log the request body
	//bodyString := string(body)
	//log.Print(bodyString)
	// Unmarshal result
	post := Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

	log.Println(post.Success)

	fmt.Println(r.Header.Get("Referer"))

	if post.Success == true {
		text := fmt.Sprintf("Name:\t%s\nE-mail:\t%s\nMessage:\t%s", name, email, message)
		SendMessage(text)

		url := fmt.Sprintf("%s/success", r.Header.Get("Referer"))

		http.Redirect(w, r, url, 302)
	} else {
		url := fmt.Sprintf("%s/ops", r.Header.Get("Referer"))

		http.Redirect(w, r, url, 302)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/post", post_handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
