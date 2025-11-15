import { COLOR } from "./constants.js";
import {
  getCurrentMoveColor,
  getMoveNumber,
  incrementMoveNumber,
  addToMoveHistory,
} from "./state.js";
import {
  getSquareNotation,
  parseSquareId,
  getPieceNotation,
} from "./notation.js";

export function addMoveToHistory(fromSquare, toSquare, pieceImg) {
  const fromCoords = parseSquareId(fromSquare.id);
  const toCoords = parseSquareId(toSquare.id);

  if (!fromCoords || !toCoords) {
    return;
  }

  const fromNotation = getSquareNotation(fromCoords.row, fromCoords.col);
  const toNotation = getSquareNotation(toCoords.row, toCoords.col);
  const pieceNotation = getPieceNotation(pieceImg);
  const isWhiteMove = getCurrentMoveColor() === COLOR.WHITE;

  const moveNotation = `${pieceNotation}${fromNotation}-${toNotation}`;

  if (isWhiteMove) {
    const moveItem = document.createElement("div");
    moveItem.className = "move-item";
    moveItem.innerHTML = `<span class="move-number">${getMoveNumber()}.</span><span class="move-white">${moveNotation}</span>`;
    document.getElementById("moves-list").appendChild(moveItem);
    addToMoveHistory({
      move: moveNotation,
      color: COLOR.WHITE,
      moveNumber: getMoveNumber(),
    });
    incrementMoveNumber();
  } else {
    const movesList = document.getElementById("moves-list");
    const lastMove = movesList.lastElementChild;
    if (lastMove && lastMove.classList.contains("move-item")) {
      lastMove.innerHTML += ` <span class="move-black">${moveNotation}</span>`;
      addToMoveHistory({
        move: moveNotation,
        color: COLOR.BLACK,
        moveNumber: getMoveNumber() - 1,
      });
    } else {
      const moveItem = document.createElement("div");
      moveItem.className = "move-item";
      moveItem.innerHTML = `<span class="move-black">${moveNotation}</span>`;
      document.getElementById("moves-list").appendChild(moveItem);
      addToMoveHistory({
        move: moveNotation,
        color: COLOR.BLACK,
        moveNumber: 0,
      });
    }
  }

  const movesContainer = document.getElementById("moves-list");
  movesContainer.scrollTop = movesContainer.scrollHeight;
}

export function clearMoveHistory() {
  const movesList = document.getElementById("moves-list");
  if (movesList) {
    movesList.innerHTML = "";
  }
}
