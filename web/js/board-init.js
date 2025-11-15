import { setBoard } from "./state.js";
import { setPiece } from "./board-manipulation.js";
import {
  BLACK_PIECES_STARTING_RANK,
  WHITE_PIECES_STARTING_RANK,
} from "./constants.js";
import { handleSquareClick } from "./game-logic.js";
import { setupDragAndDrop } from "./drag-drop.js";

export function initializeBoard() {
  const boardElement = document.getElementById("board");
  if (!boardElement) {
    console.error("Board element not found");
    return;
  }

  setBoard(boardElement);

  for (let i = 1; i <= 8; i++) {
    let row = document.createElement("tr");
    boardElement.append(row);
    for (let j = 1; j <= 8; j++) {
      let cell = document.createElement("td");
      cell.setAttribute("draggable", "true");
      cell.setAttribute("id", `${i}-${j}`);
      row.append(cell);
      cell.onclick = function () {
        handleSquareClick(cell);
      };

      if ((i + j) % 2 == 0) {
        cell.classList.add("white-field");
      } else {
        cell.classList.add("black-field");
      }

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

  setupDragAndDrop();
}
