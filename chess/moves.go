package chess

import "math/bits"

var (
	NorthMasks [64]uint64
	SouthMasks [64]uint64
	EastMasks  [64]uint64
	WestMasks  [64]uint64

	NorthEastMasks [64]uint64
	NorthWestMasks [64]uint64
	SouthEastMasks [64]uint64
	SouthWestMasks [64]uint64

	PawnAttackMasks   [2][64]uint64
	KnightAttackMasks [64]uint64
	KingAttackMasks   [64]uint64
)

const NotAFile uint64 = 0xfefefefefefefefe
const NotHFile uint64 = 0x7f7f7f7f7f7f7f7f
const Rank3Mask uint64 = 0x0000000000ff0000
const Rank6Mask uint64 = 0x00ff000000000000
const NotABFile = 0xfcfcfcfcfcfcfcfc
const NotGHFile = 0x3f3f3f3f3f3f3f3f

func init() {
	for sq := 0; sq < 64; sq++ {
		NorthMasks[sq] = generateNorthMask(sq)
		SouthMasks[sq] = generateSouthMask(sq)
		EastMasks[sq] = generateEastMask(sq)
		WestMasks[sq] = generateWestMask(sq)

		NorthEastMasks[sq] = generateNorthEastMask(sq)
		NorthWestMasks[sq] = generateNorthWestMask(sq)
		SouthEastMasks[sq] = generateSouthEastMask(sq)
		SouthWestMasks[sq] = generateSouthWestMask(sq)

		PawnAttackMasks[White][sq] = generatePawnAttackMask(White, sq)
		PawnAttackMasks[Black][sq] = generatePawnAttackMask(Black, sq)
		KnightAttackMasks[sq] = generateKnightAttackMask(sq)
		KingAttackMasks[sq] = generateKingAttackMask(sq)
	}
}

func generateNorthMask(sq int) uint64 {
	var mask uint64 = 0
	for i := sq + 8; i < 64; i += 8 {
		mask |= (1 << i)
	}
	return mask
}

func generateSouthMask(sq int) uint64 {
	var mask uint64 = 0
	for i := sq - 8; i >= 0; i -= 8 {
		mask |= (1 << i)
	}
	return mask
}

func generateEastMask(sq int) uint64 {
	var mask uint64 = 0
	for i := sq + 1; i < 64; i++ {
		if i%8 == 0 {
			break
		}
		mask |= (1 << i)
	}
	return mask
}

func generateWestMask(sq int) uint64 {
	var mask uint64 = 0
	for i := sq - 1; i >= 0; i-- {
		if i%8 == 7 {
			break
		}
		mask |= (1 << i)
	}
	return mask
}

func generateNorthEastMask(sq int) uint64 {
	var mask uint64 = 0
	for i := sq + 9; i < 64; i += 9 {
		if (1 << (i - 9) & NotHFile) == 0 {
			break
		}
		mask |= (1 << i)
	}
	return mask
}

func generateNorthWestMask(sq int) uint64 {
	var mask uint64 = 0
	for i := sq + 7; i < 64; i += 7 {
		if (1 << (i - 7) & NotAFile) == 0 {
			break
		}
		mask |= (1 << i)
	}
	return mask
}

func generateSouthEastMask(sq int) uint64 {
	var mask uint64 = 0
	for i := sq - 7; i >= 0; i -= 7 {
		if i%8 == 0 {
			break
		}
		mask |= (1 << i)
	}
	return mask
}

func generateSouthWestMask(sq int) uint64 {
	var mask uint64 = 0
	for i := sq - 9; i >= 0; i -= 9 {
		if i%8 == 7 {
			break
		}
		mask |= (1 << i)
	}
	return mask
}

func generatePawnAttackMask(color Color, sq int) uint64 {
	targetSquareBitboard := uint64(1) << sq
	attackOrigins := uint64(0)

	if color == White {
		attackOrigins |= (targetSquareBitboard >> 9) & NotHFile
		attackOrigins |= (targetSquareBitboard >> 7) & NotAFile
	} else {
		attackOrigins |= (targetSquareBitboard << 7) & NotHFile
		attackOrigins |= (targetSquareBitboard << 9) & NotAFile
	}

	return attackOrigins
}

func generateKnightAttackMask(sq int) uint64 {
	targetSquareBitboard := uint64(1) << sq
	attacks := uint64(0)

	attacks |= (targetSquareBitboard & NotHFile) << 17
	attacks |= (targetSquareBitboard & NotAFile) << 15

	attacks |= (targetSquareBitboard & NotGHFile) << 10
	attacks |= (targetSquareBitboard & NotABFile) << 6

	attacks |= (targetSquareBitboard & NotGHFile) >> 6
	attacks |= (targetSquareBitboard & NotABFile) >> 10

	attacks |= (targetSquareBitboard & NotHFile) >> 15
	attacks |= (targetSquareBitboard & NotAFile) >> 17

	return attacks
}

func generateKingAttackMask(sq int) uint64 {
	targetSquareBitboard := uint64(1) << sq
	attacks := uint64(0)

	if sq+8 < 64 {
		attacks |= targetSquareBitboard << 8
	}
	if sq-8 >= 0 {
		attacks |= targetSquareBitboard >> 8
	}
	if sq%8 != 7 {
		attacks |= (targetSquareBitboard << 1) & NotAFile
	}
	if sq%8 != 0 {
		attacks |= (targetSquareBitboard >> 1) & NotHFile
	}

	if sq%8 != 7 && sq+9 < 64 {
		attacks |= (targetSquareBitboard << 9) & NotAFile
	}
	if sq%8 != 0 && sq+7 < 64 {
		attacks |= (targetSquareBitboard << 7) & NotHFile
	}
	if sq%8 != 7 && sq-7 >= 0 {
		attacks |= (targetSquareBitboard >> 7) & NotAFile
	}
	if sq%8 != 0 && sq-9 >= 0 {
		attacks |= (targetSquareBitboard >> 9) & NotHFile
	}

	return attacks
}
func GenerateRookMoves(rookSquare int, allOccupiedSquares uint64, friendlyPieces uint64) uint64 {
	var movesNorth, movesEast, movesSouth, movesWest uint64

	blockersNorth := allOccupiedSquares & NorthMasks[rookSquare]
	if blockersNorth != 0 {
		firstBlockerIndex := bits.TrailingZeros64(blockersNorth)
		movesNorth = NorthMasks[rookSquare] ^ NorthMasks[firstBlockerIndex]
	} else {
		movesNorth = NorthMasks[rookSquare]
	}

	blockersSouth := allOccupiedSquares & SouthMasks[rookSquare]
	if blockersSouth != 0 {
		firstBlockerIndex := 63 - bits.LeadingZeros64(blockersSouth)
		movesSouth = SouthMasks[rookSquare] ^ SouthMasks[firstBlockerIndex]
	} else {
		movesSouth = SouthMasks[rookSquare]
	}

	blockersWest := allOccupiedSquares & WestMasks[rookSquare]
	if blockersWest != 0 {
		firstBlockerIndex := 63 - bits.LeadingZeros64(blockersWest)
		movesWest = WestMasks[rookSquare] ^ WestMasks[firstBlockerIndex]
	} else {
		movesWest = WestMasks[rookSquare]
	}

	blockersEast := allOccupiedSquares & EastMasks[rookSquare]
	if blockersEast != 0 {
		firstBlockerIndex := bits.TrailingZeros64(blockersEast)
		movesEast = EastMasks[rookSquare] ^ EastMasks[firstBlockerIndex]
	} else {
		movesEast = EastMasks[rookSquare]
	}

	allPossibleMoves := movesNorth | movesEast | movesSouth | movesWest

	legalMoves := allPossibleMoves &^ friendlyPieces

	return legalMoves
}

func GenerateRookMovesDetailed(board *Board, color Color, allPieces uint64, friendlyPieces uint64) []Move {
	moves := make([]Move, 0, 256)
	var rookBitboard uint64
	var enemyColor Color

	if color == White {
		rookBitboard = board.WhiteRooks
	} else {
		rookBitboard = board.BlackRooks
	}

	if rookBitboard == 0 {
		return moves
	}

	if color == White {
		enemyColor = Black
	} else {
		enemyColor = White
	}

	for rookBitboard != 0 {
		fromSquare := bits.TrailingZeros64(rookBitboard)
		rookBitboard = rookBitboard & (rookBitboard - 1)

		destinations := GenerateRookMoves(fromSquare, allPieces, friendlyPieces)

		for destinations != 0 {
			toSquare := bits.TrailingZeros64(destinations)
			destinations = destinations & (destinations - 1)

			var flag uint16
			_, pieceColor, present := board.PieceAt(toSquare)

			if present && pieceColor == enemyColor {
				flag = FlagCapture
			} else {
				flag = FlagQuietMove
			}

			move := NewMove(uint16(fromSquare), uint16(toSquare), uint16(flag))
			moves = append(moves, move)
		}
	}

	return moves
}

func GenerateBishopMoves(bishopSquare int, allOccupiedSquares uint64, friendlyPieces uint64) uint64 {
	var movesNorthEast, movesNorthWest, movesSouthEast, movesSouthWest uint64

	blockersNorthEast := allOccupiedSquares & NorthEastMasks[bishopSquare]
	if blockersNorthEast != 0 {
		firstBlockerIndex := bits.TrailingZeros64(blockersNorthEast)
		movesNorthEast = NorthEastMasks[bishopSquare] ^ NorthEastMasks[firstBlockerIndex]
	} else {
		movesNorthEast = NorthEastMasks[bishopSquare]
	}

	blockersNorthWest := allOccupiedSquares & NorthWestMasks[bishopSquare]
	if blockersNorthWest != 0 {
		firstBlockerIndex := bits.TrailingZeros64(blockersNorthWest)
		movesNorthWest = NorthWestMasks[bishopSquare] ^ NorthWestMasks[firstBlockerIndex]
	} else {
		movesNorthWest = NorthWestMasks[bishopSquare]
	}

	blockersSouthEast := allOccupiedSquares & SouthEastMasks[bishopSquare]
	if blockersSouthEast != 0 {
		firstBlockerIndex := 63 - bits.LeadingZeros64(blockersSouthEast)
		movesSouthEast = SouthEastMasks[bishopSquare] ^ SouthEastMasks[firstBlockerIndex]
	} else {
		movesSouthEast = SouthEastMasks[bishopSquare]
	}

	blockersSouthWest := allOccupiedSquares & SouthWestMasks[bishopSquare]
	if blockersSouthWest != 0 {
		firstBlockerIndex := 63 - bits.LeadingZeros64(blockersSouthWest)
		movesSouthWest = SouthWestMasks[bishopSquare] ^ SouthWestMasks[firstBlockerIndex]
	} else {
		movesSouthWest = SouthWestMasks[bishopSquare]
	}

	allPossibleMoves := movesNorthEast | movesNorthWest | movesSouthEast | movesSouthWest

	return allPossibleMoves &^ friendlyPieces
}

func GenerateBishopMovesDetailed(board *Board, color Color, allPieces uint64, friendlyPieces uint64) []Move {
	moves := make([]Move, 0, 256)
	var bishopBitboard uint64

	if color == White {
		bishopBitboard = board.WhiteBishops
	} else {
		bishopBitboard = board.BlackBishops
	}

	if bishopBitboard == 0 {
		return moves
	}

	var enemyColor Color
	if color == White {
		enemyColor = Black
	} else {
		enemyColor = White
	}

	for bishopBitboard != 0 {
		fromSquare := bits.TrailingZeros64(bishopBitboard)
		bishopBitboard = bishopBitboard & (bishopBitboard - 1)

		destinations := GenerateBishopMoves(fromSquare, allPieces, friendlyPieces)

		for destinations != 0 {
			toSquare := bits.TrailingZeros64(destinations)
			destinations = destinations & (destinations - 1)

			var flag uint16
			_, pieceColor, present := board.PieceAt(toSquare)

			if present && pieceColor == enemyColor {
				flag = FlagCapture
			} else {
				flag = FlagQuietMove
			}

			move := NewMove(uint16(fromSquare), uint16(toSquare), flag)
			moves = append(moves, move)
		}
	}

	return moves
}

func GenerateQueenMoves(queenSquare int, allOccupiedSquares uint64, friendlyPieces uint64) uint64 {
	rookMoves := GenerateRookMoves(queenSquare, allOccupiedSquares, friendlyPieces)
	bishopMoves := GenerateBishopMoves(queenSquare, allOccupiedSquares, friendlyPieces)
	return rookMoves | bishopMoves
}

func GenerateQueenMovesDetailed(board *Board, color Color, allPieces uint64, friendlyPieces uint64) []Move {
	moves := make([]Move, 0, 256)
	var queenBitboard uint64

	if color == White {
		queenBitboard = board.WhiteQueens
	} else {
		queenBitboard = board.BlackQueens
	}

	if queenBitboard == 0 {
		return moves
	}

	var enemyColor Color
	if color == White {
		enemyColor = Black
	} else {
		enemyColor = White
	}

	for queenBitboard != 0 {
		fromSquare := bits.TrailingZeros64(queenBitboard)
		queenBitboard = queenBitboard & (queenBitboard - 1)

		destinations := GenerateQueenMoves(fromSquare, allPieces, friendlyPieces)

		for destinations != 0 {
			toSquare := bits.TrailingZeros64(destinations)
			destinations = destinations & (destinations - 1)

			var flag uint16
			_, pieceColor, present := board.PieceAt(toSquare)

			if present && pieceColor == enemyColor {
				flag = FlagCapture
			} else {
				flag = FlagQuietMove
			}

			move := NewMove(uint16(fromSquare), uint16(toSquare), flag)
			moves = append(moves, move)
		}
	}

	return moves
}

func GeneratePawnMoves(pawnSquare int, color Color, allOccupiedSquares uint64, enemyPieces uint64) uint64 {
	pawnBit := uint64(1) << pawnSquare
	var moves uint64

	if color == White {
		singlePush := (pawnBit << 8) &^ allOccupiedSquares
		moves |= singlePush

		if pawnSquare >= 8 && pawnSquare <= 15 {
			doublePush := (singlePush << 8) &^ allOccupiedSquares
			moves |= doublePush
		}

		if pawnSquare%8 != 7 {
			capturesNE := (pawnBit << 9) & enemyPieces
			moves |= capturesNE
		}
		if pawnSquare%8 != 0 {
			capturesNW := (pawnBit << 7) & enemyPieces
			moves |= capturesNW
		}
	} else {
		singlePush := (pawnBit >> 8) &^ allOccupiedSquares
		moves |= singlePush

		if pawnSquare >= 48 && pawnSquare <= 55 {
			doublePush := (singlePush >> 8) &^ allOccupiedSquares
			moves |= doublePush
		}

		if pawnSquare%8 != 0 {
			capturesSW := (pawnBit >> 9) & enemyPieces
			moves |= capturesSW
		}
		if pawnSquare%8 != 7 {
			capturesSE := (pawnBit >> 7) & enemyPieces
			moves |= capturesSE
		}
	}

	return moves
}
func GeneratePawnMovesDetailed(board *Board, color Color) []Move {
	moves := make([]Move, 0, 40)
	var pawnBitboard uint64
	var enemyPieces uint64

	if color == White {
		pawnBitboard = board.WhitePawns
		enemyPieces = board.BlackPieces()
	} else {
		pawnBitboard = board.BlackPawns
		enemyPieces = board.WhitePieces()
	}

	allPieces := board.AllPieces()

	for pawnBitboard != 0 {
		fromSquare := bits.TrailingZeros64(pawnBitboard)
		pawnBitboard &= pawnBitboard - 1

		destinations := GeneratePawnMoves(fromSquare, color, allPieces, enemyPieces)

		for destinations != 0 {
			toSquare := bits.TrailingZeros64(destinations)
			destinations = destinations & (destinations - 1)

			var flag uint16
			_, pieceColor, present := board.PieceAt(toSquare)

			if present && pieceColor != color {
				flag = FlagCapture
			} else {
				if color == White {
					if fromSquare >= 8 && fromSquare <= 15 && toSquare == fromSquare+16 {
						flag = FlagDoublePawn
					} else {
						flag = FlagQuietMove
					}
				} else {
					if fromSquare >= 48 && fromSquare <= 55 && toSquare == fromSquare-16 {
						flag = FlagDoublePawn
					} else {
						flag = FlagQuietMove
					}
				}
			}

			if color == White && toSquare >= 56 {
				moves = append(moves, NewMove(uint16(fromSquare), uint16(toSquare), FlagPromoQueen|flag))
				moves = append(moves, NewMove(uint16(fromSquare), uint16(toSquare), FlagPromoRook|flag))
				moves = append(moves, NewMove(uint16(fromSquare), uint16(toSquare), FlagPromoBishop|flag))
				moves = append(moves, NewMove(uint16(fromSquare), uint16(toSquare), FlagPromoKnight|flag))
			} else if color == Black && toSquare <= 7 {
				moves = append(moves, NewMove(uint16(fromSquare), uint16(toSquare), FlagPromoQueen|flag))
				moves = append(moves, NewMove(uint16(fromSquare), uint16(toSquare), FlagPromoRook|flag))
				moves = append(moves, NewMove(uint16(fromSquare), uint16(toSquare), FlagPromoBishop|flag))
				moves = append(moves, NewMove(uint16(fromSquare), uint16(toSquare), FlagPromoKnight|flag))
			} else {
				moves = append(moves, NewMove(uint16(fromSquare), uint16(toSquare), flag))
			}
		}
	}
	return moves
}

func GenerateKingMoves(kingSquare int, allOccupiedSquares uint64, friendlyPieces uint64) uint64 {
	return KingAttackMasks[kingSquare] &^ friendlyPieces
}

func GenerateKingMovesDetailed(board *Board, color Color, friendlyPieces uint64) []Move {
	moves := make([]Move, 0, 256)
	var kingBitboard uint64

	if color == White {
		kingBitboard = board.WhiteKing
	} else {
		kingBitboard = board.BlackKing
	}

	if kingBitboard == 0 {
		return moves
	}

	fromSquare := bits.TrailingZeros64(kingBitboard)
	allPieces := board.AllPieces()
	destinations := GenerateKingMoves(fromSquare, allPieces, friendlyPieces)

	var enemyColor Color
	if color == White {
		enemyColor = Black
	} else {
		enemyColor = White
	}

	for destinations != 0 {
		toSquare := bits.TrailingZeros64(destinations)
		destinations = destinations & (destinations - 1)

		var flag uint16
		_, pieceColor, present := board.PieceAt(toSquare)

		if present && pieceColor == enemyColor {
			flag = FlagCapture
		} else {
			flag = FlagQuietMove
		}

		move := NewMove(uint16(fromSquare), uint16(toSquare), flag)
		moves = append(moves, move)
	}

	return moves
}

func GenerateKnightMoves(knightSquare int, friendlyPieces uint64) uint64 {
	return KnightAttackMasks[knightSquare] &^ friendlyPieces
}

func GenerateKnightMovesDetailed(board *Board, color Color, friendlyPieces uint64) []Move {
	moves := make([]Move, 0, 256)
	var knightBitboard uint64

	if color == White {
		knightBitboard = board.WhiteKnights
	} else {
		knightBitboard = board.BlackKnights
	}

	if knightBitboard == 0 {
		return moves
	}

	var enemyColor Color
	if color == White {
		enemyColor = Black
	} else {
		enemyColor = White
	}

	for knightBitboard != 0 {
		fromSquare := bits.TrailingZeros64(knightBitboard)
		knightBitboard = knightBitboard & (knightBitboard - 1)

		destinations := GenerateKnightMoves(fromSquare, friendlyPieces)

		for destinations != 0 {
			toSquare := bits.TrailingZeros64(destinations)
			destinations = destinations & (destinations - 1)

			var flag uint16
			_, pieceColor, present := board.PieceAt(toSquare)

			if present && pieceColor == enemyColor {
				flag = FlagCapture
			} else {
				flag = FlagQuietMove
			}

			move := NewMove(uint16(fromSquare), uint16(toSquare), flag)
			moves = append(moves, move)
		}
	}

	return moves
}

func IsSquareAttacked(square int, attackerColor Color, board *Board) bool {
	var attackingPawns, attackingKnights, attackingRooks, attackingBishops, attackingQueens, attackingKing uint64

	if attackerColor == White {
		attackingPawns = board.WhitePawns
		attackingKnights = board.WhiteKnights
		attackingRooks = board.WhiteRooks
		attackingBishops = board.WhiteBishops
		attackingQueens = board.WhiteQueens
		attackingKing = board.WhiteKing
	} else {
		attackingPawns = board.BlackPawns
		attackingKnights = board.BlackKnights
		attackingRooks = board.BlackRooks
		attackingBishops = board.BlackBishops
		attackingQueens = board.BlackQueens
		attackingKing = board.BlackKing
	}

	if (PawnAttackMasks[attackerColor][square] & attackingPawns) != 0 {
		return true
	}
	if (KnightAttackMasks[square] & attackingKnights) != 0 {
		return true
	}
	if (KingAttackMasks[square] & attackingKing) != 0 {
		return true
	}

	allPieces := board.AllPieces()

	orthogonalAttackers := attackingRooks | attackingQueens
	if orthogonalAttackers != 0 {

		blockersNorth := allPieces & NorthMasks[square]
		if blockersNorth != 0 {
			firstBlockerIndex := bits.TrailingZeros64(blockersNorth)
			if ((uint64(1) << firstBlockerIndex) & orthogonalAttackers) != 0 {
				return true
			}
		} else {
			if (NorthMasks[square] & orthogonalAttackers) != 0 {
				return true
			}
		}

		blockersSouth := allPieces & SouthMasks[square]
		if blockersSouth != 0 {
			firstBlockerIndex := 63 - bits.LeadingZeros64(blockersSouth)
			if ((uint64(1) << firstBlockerIndex) & orthogonalAttackers) != 0 {
				return true
			}
		} else {
			if (SouthMasks[square] & orthogonalAttackers) != 0 {
				return true
			}
		}

		blockersEast := allPieces & EastMasks[square]
		if blockersEast != 0 {
			firstBlockerIndex := bits.TrailingZeros64(blockersEast)
			if ((uint64(1) << firstBlockerIndex) & orthogonalAttackers) != 0 {
				return true
			}
		} else {
			if (EastMasks[square] & orthogonalAttackers) != 0 {
				return true
			}
		}

		blockersWest := allPieces & WestMasks[square]
		if blockersWest != 0 {
			firstBlockerIndex := 63 - bits.LeadingZeros64(blockersWest)
			if ((uint64(1) << firstBlockerIndex) & orthogonalAttackers) != 0 {
				return true
			}
		} else {
			if (WestMasks[square] & orthogonalAttackers) != 0 {
				return true
			}
		}
	}

	diagonalAttackers := attackingBishops | attackingQueens
	if diagonalAttackers != 0 {

		blockersNorthEast := allPieces & NorthEastMasks[square]
		if blockersNorthEast != 0 {
			firstBlockerIndex := bits.TrailingZeros64(blockersNorthEast)
			if ((uint64(1) << firstBlockerIndex) & diagonalAttackers) != 0 {
				return true
			}
		} else {
			if (NorthEastMasks[square] & diagonalAttackers) != 0 {
				return true
			}
		}

		blockersNorthWest := allPieces & NorthWestMasks[square]
		if blockersNorthWest != 0 {
			firstBlockerIndex := bits.TrailingZeros64(blockersNorthWest)
			if ((uint64(1) << firstBlockerIndex) & diagonalAttackers) != 0 {
				return true
			}
		} else {
			if (NorthWestMasks[square] & diagonalAttackers) != 0 {
				return true
			}
		}

		blockersSouthEast := allPieces & SouthEastMasks[square]
		if blockersSouthEast != 0 {
			firstBlockerIndex := 63 - bits.LeadingZeros64(blockersSouthEast)
			if ((uint64(1) << firstBlockerIndex) & diagonalAttackers) != 0 {
				return true
			}
		} else {
			if (SouthEastMasks[square] & diagonalAttackers) != 0 {
				return true
			}
		}

		blockersSouthWest := allPieces & SouthWestMasks[square]
		if blockersSouthWest != 0 {
			firstBlockerIndex := 63 - bits.LeadingZeros64(blockersSouthWest)
			if ((uint64(1) << firstBlockerIndex) & diagonalAttackers) != 0 {
				return true
			}
		} else {
			if (SouthWestMasks[square] & diagonalAttackers) != 0 {
				return true
			}
		}
	}

	return false
}

func CanCastleKingSide(board *Board) bool {
	if board.WhiteToMove && !board.WhiteKingSideCastle {
		return false
	} else if !board.WhiteToMove && !board.BlackKingSideCastle {
		return false
	}

	var squaresBetween [2]int
	var kingPath [3]int
	if board.WhiteToMove {
		squaresBetween[0] = F1
		squaresBetween[1] = G1
		kingPath[0] = E1
		kingPath[1] = F1
		kingPath[2] = G1
	} else {
		squaresBetween[0] = F8
		squaresBetween[1] = G8
		kingPath[0] = E8
		kingPath[1] = F8
		kingPath[2] = G8
	}

	for _, square := range squaresBetween {
		_, _, populated := board.PieceAt(square)

		if populated {
			return false
		}
	}

	var enemyColor Color
	if board.WhiteToMove {
		enemyColor = Black
	} else {
		enemyColor = White
	}

	for _, square := range kingPath {
		if IsSquareAttacked(square, enemyColor, board) {
			return false
		}
	}

	return true
}

func CanCastleQueenSide(board *Board) bool {
	if board.WhiteToMove && !board.WhiteQueenSideCastle {
		return false
	} else if !board.WhiteToMove && !board.BlackQueenSideCastle {
		return false
	}

	var squaresBetween [3]int
	var kingPath [3]int
	if board.WhiteToMove {
		squaresBetween[0] = D1
		squaresBetween[1] = C1
		squaresBetween[2] = B1
		kingPath[0] = E1
		kingPath[1] = D1
		kingPath[2] = C1
	} else {
		squaresBetween[0] = D8
		squaresBetween[1] = C8
		squaresBetween[2] = B8
		kingPath[0] = E8
		kingPath[1] = D8
		kingPath[2] = C8
	}

	for _, square := range squaresBetween {
		_, _, populated := board.PieceAt(square)

		if populated {
			return false
		}
	}

	var enemyColor Color
	if board.WhiteToMove {
		enemyColor = Black
	} else {
		enemyColor = White
	}

	for _, square := range kingPath {
		if IsSquareAttacked(square, enemyColor, board) {
			return false
		}
	}

	return true
}

func GenerateCastlingMoves(board *Board, color Color) []Move {
	moves := make([]Move, 0, 2)

	if color == White {
		if CanCastleKingSide(board) {
			move := NewMove(E1, G1, FlagKingCastle)
			moves = append(moves, move)
		}

		if CanCastleQueenSide(board) {
			move := NewMove(E1, C1, FlagQueenCastle)
			moves = append(moves, move)
		}
	} else {
		if CanCastleKingSide(board) {
			move := NewMove(E8, G8, FlagKingCastle)
			moves = append(moves, move)
		}

		if CanCastleQueenSide(board) {
			move := NewMove(E8, C8, FlagQueenCastle)
			moves = append(moves, move)
		}
	}

	return moves
}

func GenerateEnPassantMoves(board *Board) []Move {
	moves := make([]Move, 0, 2)

	if board.EnPassantSquare == -1 {
		return moves
	}

	epSquare := uint16(board.EnPassantSquare)

	_, _, occupied := board.PieceAt(int(epSquare))
	if occupied {
		return moves
	}

	var capturedPawnSquare int

	if board.WhiteToMove {
		capturedPawnSquare = int(epSquare) - 8

		if (board.BlackPawns & (uint64(1) << capturedPawnSquare)) == 0 {
			return moves
		}

		attackerFromSW := int(epSquare) - 9
		if attackerFromSW >= 0 && (attackerFromSW%8) < 7 {
			if (board.WhitePawns & (uint64(1) << attackerFromSW)) != 0 {
				moves = append(moves, NewMove(uint16(attackerFromSW), epSquare, FlagEPCapture))
			}
		}

		attackerFromSE := int(epSquare) - 7
		if attackerFromSE >= 0 && (attackerFromSE%8) > 0 {
			if (board.WhitePawns & (uint64(1) << attackerFromSE)) != 0 {
				moves = append(moves, NewMove(uint16(attackerFromSE), epSquare, FlagEPCapture))
			}
		}
	} else {
		capturedPawnSquare = int(epSquare) + 8

		if (board.WhitePawns & (uint64(1) << capturedPawnSquare)) == 0 {
			return moves
		}

		attackerFromNW := int(epSquare) + 7
		if attackerFromNW < 64 && (attackerFromNW%8) < 7 {
			if (board.BlackPawns & (uint64(1) << attackerFromNW)) != 0 {
				moves = append(moves, NewMove(uint16(attackerFromNW), epSquare, FlagEPCapture))
			}
		}

		attackerFromNE := int(epSquare) + 9
		if attackerFromNE < 64 && (attackerFromNE%8) > 0 {
			if (board.BlackPawns & (uint64(1) << attackerFromNE)) != 0 {
				moves = append(moves, NewMove(uint16(attackerFromNE), epSquare, FlagEPCapture))
			}
		}
	}

	return moves
}
func GenerateAllLegalMoves(board *Board) []Move {
	legalMoves := make([]Move, 0, 256)
	allPossibleMoves := GenerateAllPossibleMoves(board)

	var currentColor Color

	if board.WhiteToMove {
		currentColor = White
	} else {
		currentColor = Black
	}

	for _, move := range allPossibleMoves {
		undoInfo := board.MakeMove(move)

		var currentKingBitboard uint64
		var enemyColor Color

		if currentColor == White {
			currentKingBitboard = board.WhiteKing
			enemyColor = Black
		} else {
			currentKingBitboard = board.BlackKing
			enemyColor = White
		}

		kingSquare := bits.TrailingZeros64(currentKingBitboard)

		if !IsSquareAttacked(kingSquare, enemyColor, board) {
			legalMoves = append(legalMoves, move)
		}

		board.UndoMove(move, undoInfo)
	}

	return legalMoves
}

func GenerateAllPossibleMoves(board *Board) []Move {
	moves := make([]Move, 0, 256)
	var color Color
	var friendlyPieces uint64

	if board.WhiteToMove {
		color = White
		friendlyPieces = board.WhitePieces()
	} else {
		color = Black
		friendlyPieces = board.BlackPieces()
	}

	allPieces := board.AllPieces()

	moves = append(moves, GeneratePawnMovesDetailed(board, color)...)
	moves = append(moves, GenerateKnightMovesDetailed(board, color, friendlyPieces)...)
	moves = append(moves, GenerateBishopMovesDetailed(board, color, allPieces, friendlyPieces)...)
	moves = append(moves, GenerateRookMovesDetailed(board, color, allPieces, friendlyPieces)...)
	moves = append(moves, GenerateQueenMovesDetailed(board, color, allPieces, friendlyPieces)...)
	moves = append(moves, GenerateKingMovesDetailed(board, color, friendlyPieces)...)
	moves = append(moves, GenerateCastlingMoves(board, color)...)
	moves = append(moves, GenerateEnPassantMoves(board)...)

	return moves
}
