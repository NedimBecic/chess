package main

import (
	"chess/internal/chess"
	"testing"
)

func perft(board *chess.Board, depth int) uint64 {
	if depth == 0 {
		return 1
	}

	moves := chess.GenerateAllLegalMoves(board)
	if depth == 1 {
		return uint64(len(moves))
	}

	var nodes uint64
	for _, move := range moves {
		undoInfo := board.MakeMove(move)
		nodes += perft(board, depth - 1)
		board.UndoMove(move, undoInfo)
	}

	return nodes
}

func TestPerftInitialBoard(t *testing.T) {
	tests := map[string]struct {
		depth int
		nodes uint64
	}{
		"depth1": {
			depth: 1,
			nodes: 20,
		},
		"depth2": {
			depth: 2,
			nodes: 400,
		},
		"depth3": {
			depth: 3,
			nodes: 8902,
		},
		"depth4": {
			depth: 4,
			nodes: 197281,
		},
		"depth5": {
			depth: 5,
			nodes: 4865609,
		},
		"depth6": {
			depth: 6,
			nodes: 119060324,
		},
		// "depth7": {
		// 	depth: 7,
		// 	nodes: 3195901860,
		// },
	}

	board := chess.NewBoardFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := perft(board, test.depth)
			if result != test.nodes {
				t.Errorf("Depth %d: expected %d nodes, got %d", test.depth, test.nodes, result)
			}
		})
	}
}

func TestKiwipetePosition(t *testing.T) {
	tests := map[string]struct {
		depth int
		nodes uint64
	} {
		"depth1": {
			depth: 1,
			nodes: 48,
		}, 
		"depth2": {
			depth: 2,
			nodes: 2039,
		},
		"depth3": {
			depth: 3,
			nodes: 97862,
		}, 
		"depth4": {
			depth: 4,
			nodes: 4085603,
		}, 
		"depth5": {
			depth: 5,
			nodes: 193690690,
		}, 
		// "depth6": {
		// 	depth: 6,
		// 	nodes: 8031647685,
		// },
	}

	board := chess.NewBoardFromFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -")

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := perft(board, test.depth) 
			if result != test.nodes {
				t.Errorf("Depth %d: expected %d nodes, got %d", test.depth, test.nodes, result)
			}
		})
	}
}

func TestPosition3(t *testing.T) {
	tests := map[string]struct {
		depth int
		nodes uint64
	} {
		"depth1": {
			depth: 1,
			nodes: 14,
		},
		"depth2": {
			depth: 2,
			nodes: 191,
		}, 
		"depth3": {
			depth: 3,
			nodes: 2812,
		},
		"depth4": {
			depth: 4,
			nodes: 43238,
		},
		"depth5": {
			depth: 5,
			nodes: 674624,
		},
		"depth6": {
			depth: 6,
			nodes: 11030083,
		},
		"depth7": {
			depth: 7,
			nodes: 178633661,
		},
	}

	board := chess.NewBoardFromFEN("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := perft(board, test.depth)
			if result != test.nodes {
				t.Errorf("Depth %d: expected %d nodes, got %d", test.depth, test.nodes, result)
			}
		})
	}
}

func TestPosition4EnPassant(t *testing.T) {
	tests := map[string]struct {
		depth int
		nodes uint64
	}{
		"depth1": {
			depth: 1,
			nodes: 6,
		},
		"depth2": {
			depth: 2,
			nodes: 264,
		},
		"depth3": {
			depth: 3,
			nodes: 9467,
		},
		"depth4": {
			depth: 4,
			nodes: 422333,
		},
		"depth5": {
			depth: 5,
			nodes: 15833292,
		},
		"depth6": {
			depth: 6,
			nodes: 706045033,
		},
	}

	board := chess.NewBoardFromFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := perft(board, test.depth)
			if result != test.nodes {
				t.Errorf("Depth %d: expected %d nodes, got %d", test.depth, test.nodes, result)
			}
		})
	}
}

func TestPosition5Castling(t *testing.T) {
	tests := map[string]struct {
		depth int
		nodes uint64
	}{
		"depth1": {
			depth: 1,
			nodes: 44,
		},
		"depth2": {
			depth: 2,
			nodes: 1486,
		},
		"depth3": {
			depth: 3,
			nodes: 62379,
		},
		"depth4": {
			depth: 4,
			nodes: 2103487,
		},
		"depth5": {
			depth: 5,
			nodes: 89941194,
		},
	}

	board := chess.NewBoardFromFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := perft(board, test.depth)
			if result != test.nodes {
				t.Errorf("Depth %d: expected %d nodes, got %d", test.depth, test.nodes, result)
			}
		})
	}
}

func TestPosition6Promotion(t *testing.T) {
	tests := map[string]struct {
		depth int
		nodes uint64
	}{
		"depth1": {
			depth: 1,
			nodes: 46,
		},
		"depth2": {
			depth: 2,
			nodes: 2079,
		},
		"depth3": {
			depth: 3,
			nodes: 89890,
		},
		"depth4": {
			depth: 4,
			nodes: 3894594,
		},
		"depth5": {
			depth: 5,
			nodes: 164075551,
		},
	}

	board := chess.NewBoardFromFEN("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10")

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := perft(board, test.depth)
			if result != test.nodes {
				t.Errorf("Depth %d: expected %d nodes, got %d", test.depth, test.nodes, result)
			}
		})
	}
}