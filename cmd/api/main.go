package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	origin, err := url.Parse("http://localhost:3000")
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(origin)
	authCallBackProxy := httputil.NewSingleHostReverseProxy(origin)

	authCallBackProxy.ModifyResponse = func(response *http.Response) error {
		//read to response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		response.Body = io.NopCloser(strings.NewReader(string(body)))

		var userData map[string]interface{}
		//parse to json
		if err := json.Unmarshal(body, &userData); err != nil {
			return err
		}

		fmt.Println("Parsed JSON:", userData)
		return nil
	}

	http.Handle("/sso-auth/{provider}/callback", authCallBackProxy)
	http.Handle("/sso-auth/{provider}", proxy)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
