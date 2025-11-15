import { COLOR } from "./constants.js";
import { board, setPlayerViewColor, getPlayerViewColor } from "./state.js";

export function setBoardView(color) {
  if (!board) return;

  if (color === COLOR.BLACK) {
    board.classList.add("flipped-table");
  } else {
    board.classList.remove("flipped-table");
  }

  const rows = board.rows;
  for (let i = 0; i < rows.length; i++) {
    const cells = rows[i].cells;
    for (let j = 0; j < cells.length; j++) {
      let piece = cells[j].querySelector("img");
      if (piece != null) {
        if (color === COLOR.BLACK) {
          piece.classList.add("flipped-piece");
        } else {
          piece.classList.remove("flipped-piece");
        }
      }
    }
  }

  setPlayerViewColor(color);
}

export function getBoardView() {
  return getPlayerViewColor();
}
