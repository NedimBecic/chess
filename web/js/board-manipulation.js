import {
  BLACK_PIECES_STARTING_RANK,
  WHITE_PIECES_STARTING_RANK,
  COLOR,
} from "./constants.js";
import { board, resetGameState } from "./state.js";
import { setBoardView } from "./board-view.js";

function dragStart(ev) {
  ev.dataTransfer.setData("text", ev.target.id);
}

export function setPiece(pieceType, row, col) {
  let piece = document.createElement("img");
  piece.setAttribute("src", `pieces/${pieceType}.svg`);
  piece.setAttribute("alt", `${pieceType}`);
  piece.setAttribute("id", `${pieceType}-${row}-${col}`);
  piece.addEventListener("dragstart", dragStart);
  return piece;
}

export function resetBoard() {
  if (!board) return;

  const rows = board.rows;
  for (let i = 0; i < rows.length; i++) {
    const cells = rows[i].cells;
    for (let j = 0; j < cells.length; j++) {
      cells[j].innerHTML = "";
    }
  }

  for (let i = 1; i <= 8; i++) {
    for (let j = 1; j <= 8; j++) {
      const cell = document.getElementById(`${i}-${j}`);
      if (!cell) continue;

      if (i == 1) {
        cell.append(setPiece(BLACK_PIECES_STARTING_RANK[j - 1], i, j));
      }
      if (i == 2) {
        cell.append(setPiece("bP", i, j));
      }
      if (i == 7) {
        cell.append(setPiece("wP", i, j));
      }
      if (i == 8) {
        cell.append(setPiece(WHITE_PIECES_STARTING_RANK[j - 1], i, j));
      }
    }
  }

  resetGameState();

  const playerColorSelect = document.getElementById("player-color");
  const selectedColor =
    playerColorSelect.value === "white" ? COLOR.WHITE : COLOR.BLACK;
  setBoardView(selectedColor);
}
