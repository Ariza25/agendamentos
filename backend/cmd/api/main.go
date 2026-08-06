package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := newLazyAppHandler()
	defer func() {
		if err := handler.Close(); err != nil {
			log.Printf("error closing container: %v", err)
		}
	}()

	addr := ":" + port
	log.Printf("backend listening on %s", addr)
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: resolveDurationEnv("HTTP_READ_HEADER_TIMEOUT", 5*time.Second),
		ReadTimeout:       resolveDurationEnv("HTTP_READ_TIMEOUT", 15*time.Second),
		WriteTimeout:      resolveDurationEnv("HTTP_WRITE_TIMEOUT", 2*time.Minute),
		IdleTimeout:       resolveDurationEnv("HTTP_IDLE_TIMEOUT", 60*time.Second),
		MaxHeaderBytes:    resolvePositiveIntEnv("HTTP_MAX_HEADER_BYTES", 1<<20),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func resolvePositiveIntEnv(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	value, err := strconv.Atoi(raw)
	if raw == "" || err != nil || value <= 0 {
		return fallback
	}
	return value
}

func resolveDurationEnv(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	if duration, err := time.ParseDuration(raw); err == nil && duration > 0 {
		return duration
	}
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return fallback
	}
	return time.Duration(seconds) * time.Second
}
