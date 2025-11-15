package handler

import (
	"chess/handlers"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	handlers.HandleGetMoves(w, r)
}

