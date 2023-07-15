package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

const keyServerAddress = "KeyServerAddress"

func main() {
	myMux := http.NewServeMux()

	ctx, cancelCtx := context.WithCancel(context.Background())
	serverOne := &http.Server{
		Addr:    ":3333",
		Handler: myMux,
		BaseContext: func(listener net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddress, listener.Addr().String())
			return ctx
		}}
	serverTwo := &http.Server{
		Addr:    ":4444",
		Handler: myMux,
		BaseContext: func(listener net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddress, listener.Addr().String())
			return ctx
		}}

	fmt.Println("My Golang Test")
	myMux.HandleFunc("/", getRoot)
	myMux.HandleFunc("/hello", getHello)
	myMux.HandleFunc("/gustine", getGustine)

	go func() {
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error: Server Closed\n")
		} else if err != nil {
			fmt.Printf("Unexpected Server Error: %s\n", err)
		}
		cancelCtx()
	}()

	go func() {
		err := serverTwo.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error: Server Closed\n")
		} else if err != nil {
			fmt.Printf("Unexpected Server Error: %s\n", err)
		}
		cancelCtx()
	}()
	<-ctx.Done()
	//err := http.ListenAndServe(":3333", myMux)
	//if errors.Is(err, http.ErrServerClosed) {
	//	fmt.Printf("\n***** Server Closed ***\n")
	//} else if err != nil {
	//	fmt.Printf("\n***** Error Starting server ***\n")
	//	os.Exit(1)
	//}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	myContext := r.Context()
	fmt.Printf("%s : got / request\n", myContext.Value(keyServerAddress))
	io.WriteString(w, "This is my website!\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	myContext := r.Context()

	var hasFirst bool = r.URL.Query().Has("first") // returns whether there is a filter value for the given key
	var hasSecond bool = r.URL.Query().Has("second")
	first := r.URL.Query().Get("first") //returns empty string if no input from user or a filter value if provided
	second := r.URL.Query().Get("second")

	fmt.Printf("%s : got /hello request.\nHas First: %t = %s\nHas second: %t = %s\n", myContext.Value(keyServerAddress), hasFirst, first, hasSecond, second)
	io.WriteString(w, "Hello, HTTP!\n")
}

func getGustine(w http.ResponseWriter, r *http.Request) {
	myContext := r.Context()
	var data = r.Body
	var body Body
	fmt.Printf("Body Data: %s\n", data)

	err := json.NewDecoder(data).Decode(&body)

	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("Body Decoded: %s\n", body)
	fmt.Println("Body Decoded Name:", body.Name)

	fmt.Printf("%s : Got a /gustine Request\n", myContext.Value(keyServerAddress))
	io.WriteString(w, "My Name Is Gustine")

}

type Body struct {
	Name string `json:"name"`
	Page string `json:"page"`
}
