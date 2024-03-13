package main

import "net/http"

func NewMux() http.Handler {
	mux := http.NewServeMux()
	//An endpoint to check if http server is running correctly
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	return mux
}
