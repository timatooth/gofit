package main

import (
	"fmt"
	"log"
	"bytes"
	"os"
	"strings"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"net/http"
  "net/url"
	"text/template"
)

type FitbitApi struct {
	ClientId, ClientSecret, RedirectUri, AuthorizeUri, Scopes string
	Auth FitbitAuth
}

type FitbitAuth struct {
		AccessToken string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		UserId string `json:"user_id"`
		TokenType string `json:"token_type"`
		Scope string `json:"scope"`
	}

func NewFitBitApi(ClientId string, ClientSecret string, RedirectUri string) FitbitApi {
	authorizeTmpl := "response_type=code&client_id={{.ClientId}}&redirect_uri={{.RedirectUri}}&scope={{.Scopes}}&expires_in=604800"
  tmpl, err := template.New("test").Parse(authorizeTmpl)
	api := FitbitApi{}
	api.ClientId = ClientId
	api.ClientSecret = ClientSecret
	api.RedirectUri = RedirectUri
	api.Scopes = "profile settings location heartrate activity weight sleep nutrition"
	if err != nil { panic(err) }
	var authorizeBuf bytes.Buffer
	err = tmpl.Execute(&authorizeBuf, api)
	if err != nil { panic(err) }
	api.AuthorizeUri = "https://www.fitbit.com/oauth2/authorize?" + url.PathEscape(authorizeBuf.String())
	return api
}

func (api *FitbitApi) EncodeBasicAuth() string{
	authstring := api.ClientId + ":" + api.ClientSecret
	return base64.StdEncoding.EncodeToString([]byte(authstring))
}

func (api *FitbitApi) GetProfile() string{
	req, _ := http.NewRequest("GET", "https://api.fitbit.com/1/user/-/profile.json", nil)
	req.Header.Set("Authorization", "Bearer " + api.Auth.AccessToken)
	res, _ := http.DefaultClient.Do(req)
	profiledata, _:= ioutil.ReadAll(res.Body)
	return string(profiledata)
}

func (api *FitbitApi) LoadAccessToken(code string){
	form := url.Values{}
	form.Add("clientId", api.ClientId)
	form.Add("grant_type", "authorization_code")
	form.Add("redirect_uri", api.RedirectUri)
	form.Add("code", code)
	req, err := http.NewRequest("POST", "https://api.fitbit.com/oauth2/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic " + api.EncodeBasicAuth())
	res, err := http.DefaultClient.Do(req)
	if err != nil { panic(err) }
	decoder := json.NewDecoder(res.Body)
	
	var auth FitbitAuth
	decerr := decoder.Decode(&auth)
	if decerr != nil { panic(decerr) }
	res.Body.Close()
	api.Auth = auth
}

func main() {
	mux := http.NewServeMux()
	
	api := NewFitBitApi(os.Getenv("FITBIT_CLIENT_ID"), os.Getenv("FITBIT_CLIENT_SECRET"), "http://localhost:3000/auth")
	
	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query()["code"][0]
		api.LoadAccessToken(code)
		fmt.Fprintf(w, api.GetProfile())
	})
	
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		fmt.Fprintf(w, "Visit: <a href=%q>%q</a>", api.AuthorizeUri, api.AuthorizeUri)
		
	})
	
	fmt.Println("Visit: " + api.AuthorizeUri)
	log.Fatal(http.ListenAndServe(":3000", mux))
}
