package main

import (
	"bufio"
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"net/url"
	"os"
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

func main() {
	// https://developers.google.com/identity/protocols/oauth2
	conf := &oauth2.Config{
		ClientID:     "1075913669393-qhnfn93sgoibq7d5l09o131btkgu3sqo.apps.googleusercontent.com",
		ClientSecret: "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",
		Endpoint:     google.Endpoint,
		RedirectURL:  "https://chat.openai.com/chat",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
	// Get the authorization url
	uri := conf.AuthCodeURL("state")

	// Print the auth url
	fmt.Printf("Visit the URL for the auth dialog: %v\n", uri)

	// Read the authorization code from the command line input
	fmt.Print("Enter the authorization code: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	link := strings.TrimSpace(scanner.Text())
	parsed, err := url.Parse(link)
	if err != nil {
		log.Fatalln(err)
	}
	code := parsed.Query().Get("code")
	if code == "" {
		fmt.Println("Authorization code not found in URI")
		return
	}

	ctx := context.Background()
	// Exchange the auth code for an access token
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatalln(err)
	}

	// use the access token to make protected resources on the site
	client := conf.Client(ctx, token)
	client.Jar = NewJar()
	req, err := http.NewRequest("GET", "https://chat.openai.com/chat", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	//
	fmt.Println(resp)

}
