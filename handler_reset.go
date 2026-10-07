package main

import "net/http"

func (cfg *apiConfig) resetHandler(writer http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Only devs are allowed to reset database!"))
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.Reset(req.Context())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Reset fail: " + err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Database reset and hits zeroed"))
}