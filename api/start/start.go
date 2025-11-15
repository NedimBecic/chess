package handler

import (
	"chess/internal/chess"
	"chess/internal/handlers"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	board := chess.NewBoardFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	handlers.SetBoard(*board)
	handlers.HandleStartGame(w, r)
}

