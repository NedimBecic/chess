package handlers

import (
	"chess/chess"
	"encoding/json"
	"log"
	"net/http"
)

var board chess.Board

func GetBoard() *chess.Board {
	return &board
}

func SetBoard(b chess.Board) {
	board = b
}

type GameState struct {
	FEN        string   `json:"fen"`
	MoveCount  int      `json:"move_count"`
	LegalMoves []string `json:"legal_moves"`
}

type MoveRequest struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Promotion string `json:"promotion,omitempty"`
	FEN       string `json:"fen,omitempty"`
}

type MoveResponse struct {
	Success    bool     `json:"success"`
	Message    string   `json:"message,omitempty"`
	FEN        string   `json:"fen,omitempty"`
	LegalMoves []string `json:"legal_moves,omitempty"`
}

func HandleStartGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	board = *chess.NewBoardFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	legalMoves := chess.GenerateAllLegalMoves(&board)
	moveStrings := make([]string, len(legalMoves))

	for i, move := range legalMoves {
		moveStrings[i] = move.ToString()
	}

	response := GameState {
		FEN: board.ToFEN(),
		MoveCount: len(legalMoves),
		LegalMoves: moveStrings,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding reponse: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func HandleGetMoves(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	legalMoves := chess.GenerateAllLegalMoves(&board)
	moveStrings := make([]string, len(legalMoves))
	for i, move := range legalMoves {
		moveStrings[i] = move.ToString()
	}

	response := GameState{
		FEN:        board.ToFEN(),
		MoveCount:  0,
		LegalMoves: moveStrings,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func HandlePostMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var moveReq MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&moveReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if moveReq.FEN != "" {
		board = *chess.NewBoardFromFEN(moveReq.FEN)
	} else {
		board = *chess.NewBoardFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	}

	var move chess.Move
	var err error
	if moveReq.Promotion != "" {
		move, err = board.ValidateMove(moveReq.From, moveReq.To, moveReq.Promotion)
	} else {
		move, err = board.ValidateMove(moveReq.From, moveReq.To)
	}
	if err != nil {
		response := MoveResponse{
			Success: false,
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	board.MakeMove(move)

	legalMoves := chess.GenerateAllLegalMoves(&board)
	moveStrings := make([]string, len(legalMoves))
	for i, m := range legalMoves {
		moveStrings[i] = m.ToString()
	}

	response := MoveResponse{
		Success:    true,
		Message:    "Move applied successfully",
		FEN:        board.ToFEN(),
		LegalMoves: moveStrings,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

