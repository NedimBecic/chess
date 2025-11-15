package handler

import (
	"chess/internal/handlers"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	handlers.HandlePostMove(w, r)
}

