package chess

import (
	"fmt"
	"strings"
)

type Color int
type PieceType int

const (
	White Color = 1
	Black Color = 0
)

const (
	Pawn PieceType = iota
	Rook
	Knight
	Bishop
	Queen
	King
	NoPieceType
)

const (
	A1, B1, C1, D1, E1, F1, G1, H1 = 0, 1, 2, 3, 4, 5, 6, 7
	A2, B2, C2, D2, E2, F2, G2, H2 = 8, 9, 10, 11, 12, 13, 14, 15
	A3, B3, C3, D3, E3, F3, G3, H3 = 16, 17, 18, 19, 20, 21, 22, 23
	A4, B4, C4, D4, E4, F4, G4, H4 = 24, 25, 26, 27, 28, 29, 30, 31
	A5, B5, C5, D5, E5, F5, G5, H5 = 32, 33, 34, 35, 36, 37, 38, 39
	A6, B6, C6, D6, E6, F6, G6, H6 = 40, 41, 42, 43, 44, 45, 46, 47
	A7, B7, C7, D7, E7, F7, G7, H7 = 48, 49, 50, 51, 52, 53, 54, 55
	A8, B8, C8, D8, E8, F8, G8, H8 = 56, 57, 58, 59, 60, 61, 62, 63
)

type Board struct {
	WhitePawns   uint64
	WhiteKnights uint64
	WhiteBishops uint64
	WhiteRooks   uint64
	WhiteQueens  uint64
	WhiteKing    uint64

	BlackPawns   uint64
	BlackKnights uint64
	BlackBishops uint64
	BlackRooks   uint64
	BlackQueens  uint64
	BlackKing    uint64

	WhiteToMove                                                                          bool
	WhiteKingSideCastle, WhiteQueenSideCastle, BlackKingSideCastle, BlackQueenSideCastle bool
	EnPassantSquare                                                                      int
}

func (b *Board) PieceAt(square int) (PieceType, Color, bool) {
	squareBit := uint64(1) << uint64(square)

	if b.WhitePawns&squareBit != 0 {
		return Pawn, White, true
	}
	if b.WhiteKnights&squareBit != 0 {
		return Knight, White, true
	}
	if b.WhiteBishops&squareBit != 0 {
		return Bishop, White, true
	}
	if b.WhiteRooks&squareBit != 0 {
		return Rook, White, true
	}
	if b.WhiteQueens&squareBit != 0 {
		return Queen, White, true
	}
	if b.WhiteKing&squareBit != 0 {
		return King, White, true
	}

	if b.BlackPawns&squareBit != 0 {
		return Pawn, Black, true
	}
	if b.BlackKnights&squareBit != 0 {
		return Knight, Black, true
	}
	if b.BlackBishops&squareBit != 0 {
		return Bishop, Black, true
	}
	if b.BlackRooks&squareBit != 0 {
		return Rook, Black, true
	}
	if b.BlackQueens&squareBit != 0 {
		return Queen, Black, true
	}
	if b.BlackKing&squareBit != 0 {
		return King, Black, true
	}

	return 0, 0, false
}

func (b *Board) GetBitboard(pieceType PieceType, color Color) *uint64 {
	if color == White {
		switch pieceType {
		case Pawn:
			return &b.WhitePawns
		case Knight:
			return &b.WhiteKnights
		case Bishop:
			return &b.WhiteBishops
		case Rook:
			return &b.WhiteRooks
		case Queen:
			return &b.WhiteQueens
		case King:
			return &b.WhiteKing
		}
	} else {
		switch pieceType {
		case Pawn:
			return &b.BlackPawns
		case Knight:
			return &b.BlackKnights
		case Bishop:
			return &b.BlackBishops
		case Rook:
			return &b.BlackRooks
		case Queen:
			return &b.BlackQueens
		case King:
			return &b.BlackKing
		}
	}

	return nil
}

func (b *Board) WhitePieces() uint64 {
	return b.WhitePawns | b.WhiteRooks | b.WhiteKnights | b.WhiteBishops | b.WhiteQueens | b.WhiteKing
}

func (b *Board) BlackPieces() uint64 {
	return b.BlackPawns | b.BlackRooks | b.BlackKnights | b.BlackBishops | b.BlackQueens | b.BlackKing
}

func (b *Board) AllPieces() uint64 {
	return b.WhitePieces() | b.BlackPieces()
}

func NewBoardFromFEN(fen string) *Board {

	board := new(Board)
	parts := strings.Split(fen, " ")

	row := 8
	col := 1

	for i := 0; i < len(parts[0]); i++ {
		if parts[0][i] == '/' {
			row--
			col = 1
		} else if parts[0][i] >= '1' && parts[0][i] <= '8' {
			col += int(parts[0][i] - '0')
		} else {
			squareIndex := (row-1)*8 + (col - 1)

			switch parts[0][i] {
			case 'P':
				board.WhitePawns |= (1 << squareIndex)
			case 'p':
				board.BlackPawns |= (1 << squareIndex)
			case 'N':
				board.WhiteKnights |= (1 << squareIndex)
			case 'n':
				board.BlackKnights |= (1 << squareIndex)
			case 'B':
				board.WhiteBishops |= (1 << squareIndex)
			case 'b':
				board.BlackBishops |= (1 << squareIndex)
			case 'R':
				board.WhiteRooks |= (1 << squareIndex)
			case 'r':
				board.BlackRooks |= (1 << squareIndex)
			case 'Q':
				board.WhiteQueens |= (1 << squareIndex)
			case 'q':
				board.BlackQueens |= (1 << squareIndex)
			case 'K':
				board.WhiteKing |= (1 << squareIndex)
			case 'k':
				board.BlackKing |= (1 << squareIndex)
			}

			col++
		}
	}

	if len(parts) > 1 {
		board.WhiteToMove = (parts[1] == "w")
	}

	if len(parts) > 2 {
		board.WhiteKingSideCastle = strings.Contains(parts[2], "K")
		board.WhiteQueenSideCastle = strings.Contains(parts[2], "Q")
		board.BlackKingSideCastle = strings.Contains(parts[2], "k")
		board.BlackQueenSideCastle = strings.Contains(parts[2], "q")
	}

	if len(parts) > 3 && parts[3] != "-" {
		file := int(parts[3][0] - 'a')
		rank := int(parts[3][1] - '1')
		board.EnPassantSquare = rank*8 + file
	} else {
		board.EnPassantSquare = -1
	}

	return board
}
func (b *Board) MakeMove(move Move) UndoMoveInfo {
	undoInfo := UndoMoveInfo{
		PreviousEnPassantSquare: b.EnPassantSquare,
		WhiteKingSideCastle:     b.WhiteKingSideCastle,
		WhiteQueenSideCastle:    b.WhiteQueenSideCastle,
		BlackKingSideCastle:     b.BlackKingSideCastle,
		BlackQueenSideCastle:    b.BlackQueenSideCastle,
		CapturedPieceType:       NoPieceType,
	}

	from := move.From()
	to := move.To()
	flag := move.Flag()

	movingPiece, color, _ := b.PieceAt(int(from))
	undoInfo.MovingPieceType = movingPiece

	fromBit := uint64(1 << from)
	toBit := uint64(1 << to)
	fromToMask := fromBit | toBit

	var enemyColor Color
	if color == White {
		enemyColor = Black
	} else {
		enemyColor = White
	}

	if (flag & FlagCapture) != 0 {
		if flag == FlagEPCapture {
			undoInfo.CapturedPieceType = Pawn
			var capturedPawnSquare int
			if color == White {
				capturedPawnSquare = int(to) - 8
			} else {
				capturedPawnSquare = int(to) + 8
			}
			*b.GetBitboard(Pawn, enemyColor) &^= (uint64(1) << capturedPawnSquare)
		} else {
			capturedPiece, _, _ := b.PieceAt(int(to))
			undoInfo.CapturedPieceType = capturedPiece
			*b.GetBitboard(capturedPiece, enemyColor) &^= toBit
		}
	}

	b.EnPassantSquare = -1

	if flag >= FlagPromoKnight {
		*b.GetBitboard(Pawn, color) &^= fromBit
		var promoPiece PieceType
		switch flag & 0b1011 {
		case FlagPromoQueen:
			promoPiece = Queen
		case FlagPromoRook:
			promoPiece = Rook
		case FlagPromoBishop:
			promoPiece = Bishop
		case FlagPromoKnight:
			promoPiece = Knight
		}
		*b.GetBitboard(promoPiece, color) |= toBit
	} else {
		switch flag {
		case FlagKingCastle:
			*b.GetBitboard(King, color) ^= fromToMask
			if color == White {
				*b.GetBitboard(Rook, color) ^= (uint64(1) << H1) | (uint64(1) << F1)
			} else {
				*b.GetBitboard(Rook, color) ^= (uint64(1) << H8) | (uint64(1) << F8)
			}
		case FlagQueenCastle:
			*b.GetBitboard(King, color) ^= fromToMask
			if color == White {
				*b.GetBitboard(Rook, color) ^= (uint64(1) << A1) | (uint64(1) << D1)
			} else {
				*b.GetBitboard(Rook, color) ^= (uint64(1) << A8) | (uint64(1) << D8)
			}
		case FlagDoublePawn:
			*b.GetBitboard(Pawn, color) ^= fromToMask
			if color == White {
				b.EnPassantSquare = int(from) + 8
			} else {
				b.EnPassantSquare = int(from) - 8
			}
		default:
			*b.GetBitboard(movingPiece, color) ^= fromToMask
		}
	}

	if movingPiece == King {
		if color == White {
			b.WhiteKingSideCastle = false
			b.WhiteQueenSideCastle = false
		} else {
			b.BlackKingSideCastle = false
			b.BlackQueenSideCastle = false
		}
	}

	if from == H1 {
		b.WhiteKingSideCastle = false
	}
	if from == A1 {
		b.WhiteQueenSideCastle = false
	}
	if from == H8 {
		b.BlackKingSideCastle = false
	}
	if from == A8 {
		b.BlackQueenSideCastle = false
	}

	if to == H8 && undoInfo.CapturedPieceType == Rook {
		b.BlackKingSideCastle = false
	}
	if to == A8 && undoInfo.CapturedPieceType == Rook {
		b.BlackQueenSideCastle = false
	}
	if to == H1 && undoInfo.CapturedPieceType == Rook {
		b.WhiteKingSideCastle = false
	}
	if to == A1 && undoInfo.CapturedPieceType == Rook {
		b.WhiteQueenSideCastle = false
	}

	b.WhiteToMove = !b.WhiteToMove

	return undoInfo
}
func (b *Board) UndoMove(move Move, undoInfo UndoMoveInfo) {
	from := move.From()
	to := move.To()
	flag := move.Flag()

	var color Color
	if b.WhiteToMove {
		color = Black
	} else {
		color = White
	}

	fromBit := uint64(1 << from)
	toBit := uint64(1 << to)
	fromToMask := fromBit | toBit

	b.WhiteToMove = !b.WhiteToMove
	b.EnPassantSquare = undoInfo.PreviousEnPassantSquare
	b.WhiteKingSideCastle = undoInfo.WhiteKingSideCastle
	b.WhiteQueenSideCastle = undoInfo.WhiteQueenSideCastle
	b.BlackKingSideCastle = undoInfo.BlackKingSideCastle
	b.BlackQueenSideCastle = undoInfo.BlackQueenSideCastle

	if flag >= FlagPromoKnight {
		*b.GetBitboard(Pawn, color) |= fromBit

		var promoPiece PieceType
		switch flag &^ FlagCapture {
		case FlagPromoQueen:
			promoPiece = Queen
		case FlagPromoRook:
			promoPiece = Rook
		case FlagPromoBishop:
			promoPiece = Bishop
		case FlagPromoKnight:
			promoPiece = Knight
		}
		*b.GetBitboard(promoPiece, color) &^= toBit
	} else {
		switch flag {
		case FlagKingCastle:
			*b.GetBitboard(King, color) ^= fromToMask
			if color == White {
				*b.GetBitboard(Rook, color) ^= (uint64(1) << H1) | (uint64(1) << F1)
			} else {
				*b.GetBitboard(Rook, color) ^= (uint64(1) << H8) | (uint64(1) << F8)
			}
		case FlagQueenCastle:
			*b.GetBitboard(King, color) ^= fromToMask
			if color == White {
				*b.GetBitboard(Rook, color) ^= (uint64(1) << A1) | (uint64(1) << D1)
			} else {
				*b.GetBitboard(Rook, color) ^= (uint64(1) << A8) | (uint64(1) << D8)
			}
		default:
			*b.GetBitboard(undoInfo.MovingPieceType, color) ^= fromToMask
		}
	}

	if undoInfo.CapturedPieceType != NoPieceType {
		var enemyColor Color
		if color == White {
			enemyColor = Black
		} else {
			enemyColor = White
		}

		capturedSquare := to
		if flag == FlagEPCapture {
			if color == White {
				capturedSquare = to - 8
			} else {
				capturedSquare = to + 8
			}
		}
		*b.GetBitboard(undoInfo.CapturedPieceType, enemyColor) |= (uint64(1) << capturedSquare)
	}
}

func (b *Board) ToFEN() string {
	var fen strings.Builder
	
	for rank := 7; rank >= 0; rank-- {
		emptyCount := 0
		for file := 0; file < 8; file++ {
			square := rank*8 + file
			pieceType, color, present := b.PieceAt(square)
			
			if !present {
				emptyCount++
			} else {
				if emptyCount > 0 {
					fen.WriteString(fmt.Sprintf("%d", emptyCount))
					emptyCount = 0
				}
				
				var pieceChar byte
				switch pieceType {
				case Pawn:
					pieceChar = 'P'
				case Rook:
					pieceChar = 'R'
				case Knight:
					pieceChar = 'N'
				case Bishop:
					pieceChar = 'B'
				case Queen:
					pieceChar = 'Q'
				case King:
					pieceChar = 'K'
				}
				
				if color == Black {
					pieceChar = pieceChar - 'A' + 'a'
				}
				fen.WriteByte(pieceChar)
			}
		}
		if emptyCount > 0 {
			fen.WriteString(fmt.Sprintf("%d", emptyCount))
		}
		if rank > 0 {
			fen.WriteByte('/')
		}
	}
	
	if b.WhiteToMove {
		fen.WriteString(" w ")
	} else {
		fen.WriteString(" b ")
	}
	
	castling := ""
	if b.WhiteKingSideCastle {
		castling += "K"
	}
	if b.WhiteQueenSideCastle {
		castling += "Q"
	}
	if b.BlackKingSideCastle {
		castling += "k"
	}
	if b.BlackQueenSideCastle {
		castling += "q"
	}
	if castling == "" {
		castling = "-"
	}
	fen.WriteString(castling)
	
	fen.WriteString(" ")
	if b.EnPassantSquare >= 0 {
		file := b.EnPassantSquare % 8
		rank := b.EnPassantSquare / 8
		fen.WriteString(fmt.Sprintf("%c%c", 'a'+file, '1'+rank))
	} else {
		fen.WriteString("-")
	}
	
	fen.WriteString(" 0 1")
	
	return fen.String()
}

func (b *Board) ValidateMove(fromStr, toStr string, promotionPiece ...string) (Move, error) {
	fromSquare, err := ParseSquare(fromStr)
	if err != nil {
		return 0, fmt.Errorf("invalid 'from' square: %w", err)
	}

	toSquare, err := ParseSquare(toStr)
	if err != nil {
		return 0, fmt.Errorf("invalid 'to' square: %w", err)
	}

	legalMoves := GenerateAllLegalMoves(b)
	
	if len(promotionPiece) > 0 && promotionPiece[0] != "" {
		var promoFlag uint16
		switch promotionPiece[0] {
		case "q", "Q":
			promoFlag = FlagPromoQueen
		case "r", "R":
			promoFlag = FlagPromoRook
		case "b", "B":
			promoFlag = FlagPromoBishop
		case "n", "N":
			promoFlag = FlagPromoKnight
		default:
			return 0, fmt.Errorf("invalid promotion piece: %s", promotionPiece[0])
		}
		
		for _, move := range legalMoves {
			if move.From() == fromSquare && move.To() == toSquare {
				flag := move.Flag()
				if flag >= FlagPromoKnight && (flag&0b1011) == promoFlag {
					return move, nil
				}
			}
		}
	}
	
	move, found := FindMoveInList(fromSquare, toSquare, legalMoves)
	if !found {
		return 0, fmt.Errorf("move %s to %s is not legal", fromStr, toStr)
	}

	return move, nil
}
