package chess

import (
	"fmt"
)

type UndoMoveInfo struct {
	PreviousEnPassantSquare                                                              int
	WhiteKingSideCastle, WhiteQueenSideCastle, BlackKingSideCastle, BlackQueenSideCastle bool
	CapturedPieceType                                                                    PieceType
	MovingPieceType                                                                      PieceType
}

type Move uint16

const (
	FlagQuietMove   = 0b0000
	FlagDoublePawn  = 0b0001
	FlagKingCastle  = 0b0010
	FlagQueenCastle = 0b0011
	FlagCapture     = 0b0100
	FlagEPCapture   = 0b0101
	FlagPromoKnight = 0b1000
	FlagPromoBishop = 0b1001
	FlagPromoRook   = 0b1010
	FlagPromoQueen  = 0b1011
)

func NewMove(from uint16, to uint16, flag uint16) Move {
	return Move((from & 0x3f) | ((to & 0x3f) << 6) | (flag << 12))
}

func (move Move) From() uint16 {
	return uint16(move) & 0x3f
}

func (move Move) To() uint16 {
	return (uint16(move) >> 6) & 0x3f
}

func (move Move) Flag() uint16 {
	return uint16(move) >> 12
}

func (move Move) ToString() string {
	from := move.From()
	to := move.To()
	flag := move.Flag()
	fromStr := fmt.Sprintf("%c%c", 'a'+(from%8), '1'+(from/8))
	toStr := fmt.Sprintf("%c%c", 'a'+(to%8), '1'+(to/8))
	
	if flag >= FlagPromoKnight {
		var promoChar byte
		switch flag & 0b1011 {
		case FlagPromoQueen:
			promoChar = 'q'
		case FlagPromoRook:
			promoChar = 'r'
		case FlagPromoBishop:
			promoChar = 'b'
		case FlagPromoKnight:
			promoChar = 'n'
		}
		return fromStr + toStr + string(promoChar)
	}
	
	return fromStr + toStr
}

func ParseSquare(squareStr string) (uint16, error) {
	if len(squareStr) != 2 {
		return 0, fmt.Errorf("invalid square format: %s (expected format: e.g., 'e2')", squareStr)
	}
	file := squareStr[0] - 'a'
	rank := squareStr[1] - '1'
	if file > 7 || rank > 7 {
		return 0, fmt.Errorf("invalid square: %s (file or rank out of range)", squareStr)
	}
	square := uint16(rank)*8 + uint16(file)
	return square, nil
}

func FindMoveInList(fromSquare, toSquare uint16, moves []Move) (Move, bool) {
	for _, move := range moves {
		if move.From() == fromSquare && move.To() == toSquare {
			return move, true
		}
	}
	return 0, false
}
