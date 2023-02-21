package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type StartConversationRequest struct {
	UserId string `json:"user_id"`
}

type StartConversationResponse struct {
	ConversationId string `json:"conversation_id"`
}

type Jar struct {
	sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.Lock()
	if _, ok := jar.cookies[u.Host]; ok {
		for _, c := range cookies {
			jar.cookies[u.Host] = append(jar.cookies[u.Host], c)
		}
	} else {
		jar.cookies[u.Host] = cookies
	}
	jar.Unlock()
}

func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

func NewJarClient() *http.Client {
	return &http.Client{
		Jar: NewJar(),
	}
}

// https://stackoverflow.com/questions/37037773/use-golang-to-login-to-private-site-and-pull-info

func fetch(w http.ResponseWriter, r *http.Request) {

	// create the client
	client := NewJarClient()

	// get the csrf token
	req, _ := http.NewRequest("GET", "https://chat.openai.com/auth/login", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Header)
	fmt.Println(resp.ContentLength)

}

/*doc, err := goquery.NewDocumentFromResponse(resp)
if err != nil {
	log.Fatal(err)
}
doc.
if val, ok := doc.Find(`head meta[name="csrf-token-value"]`).Attr("content"); ok {
	csrfToken = val
}*/

/*func main() {
	http.HandleFunc("/", fetch)
	http.ListenAndServe(":8080", nil)
}*/

func main() {
	// https://developers.google.com/identity/protocols/oauth2
	conf := &oauth2.Config{
		ClientID:     "1075913669393-a35h2pn996rdrblbfc07a5d4a6vmo82g.apps.googleusercontent.com",
		ClientSecret: "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",
		Endpoint:     google.Endpoint,
		RedirectURL:  "https://chat.openai.com/auth/login",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
	// Start the authorization flow
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser: %v\n", url)

	// Handle the callback request
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		token, err := conf.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Access token: ", token.AccessToken)
		fmt.Println("Refresh token: ", token.RefreshToken)
		fmt.Println("Type: ", token.TokenType)
		fmt.Println("Expiration: ", token.Expiry)

		// get the user's information
		httpClient := &http.Client{}
		access := strings.Join([]string{"Bearer", token.AccessToken}, " ")
		req, _ := http.NewRequest("GET", "https://api.openai.com/v1/users/me", nil)
		req.Header.Set("Authorization", access)
		resp, err := httpClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Println("made it")

		// parse the response
		fmt.Println("Content length: ", resp.ContentLength)
		fmt.Println("Status: ", resp.Status)
		fmt.Println("Status Code: ", resp.StatusCode)
		fmt.Println("Header", resp.Header)
		fmt.Println("Body: ")
		if resp.StatusCode == http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Println(string(bodyBytes))
		}
	})

	http.ListenAndServe(":8080", nil)

}
