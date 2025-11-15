import {
  getIsPieceMoving,
  setIsPieceMoving,
  getInitialSquare,
  setInitialSquare,
  clearInitialSquare,
  getMovingPiece,
  setMovingPiece,
  clearMovingPiece,
  toggleMoveColor,
  getCurrentMoveColor,
} from "./state.js";
import { addMoveToHistory } from "./move-history.js";
import { setBoardView } from "./board-view.js";
import { COLOR } from "./constants.js";
import { makeMove as apiMakeMove } from "./api.js";
import { parseSquareId, getSquareNotation } from "./notation.js";
import { showPromotionDialog } from "./promotion.js";
import { setPiece } from "./board-manipulation.js";

export function handleSquareClick(square) {
  if (square === getInitialSquare()) {
    setIsPieceMoving(false);
    clearMovingPiece();
    clearInitialSquare();
    return;
  }

  if (!getIsPieceMoving()) {
    if (square.innerHTML.trim() === "") {
      return;
    }
    setMovingPiece(square.innerHTML);
    setInitialSquare(square);
    setIsPieceMoving(true);
  } else {
    attemptMove(getInitialSquare(), square);
  }
}

async function attemptMove(fromSquare, toSquare) {
  const fromCoords = parseSquareId(fromSquare.id);
  const toCoords = parseSquareId(toSquare.id);

  if (!fromCoords || !toCoords) {
    resetMoveState();
    return;
  }

  const fromNotation = getSquareNotation(fromCoords.row, fromCoords.col);
  const toNotation = getSquareNotation(toCoords.row, toCoords.col);

  const pieceImg = fromSquare.querySelector("img");
  const pieceSrc = pieceImg ? pieceImg.getAttribute("src") : "";
  const isPawn = pieceSrc.includes("P.svg");
  const isWhite = pieceSrc.includes("w");
  const isKing = pieceSrc.includes("K.svg");
  const needsPromotion =
    isPawn &&
    ((isWhite && toCoords.row === 1) || (!isWhite && toCoords.row === 8));

  const toSquareWasEmpty = toSquare.innerHTML.trim() === "";

  let promotion = null;
  if (needsPromotion) {
    promotion = await showPromotionDialog(isWhite);
    if (!promotion) {
      resetMoveState();
      return;
    }
  }

  try {
    const response = await apiMakeMove(fromNotation, toNotation, promotion);

    if (response && response.success) {
      await applyMove(
        fromSquare,
        toSquare,
        pieceImg,
        fromCoords,
        toCoords,
        needsPromotion,
        promotion,
        isKing,
        isPawn,
        toSquareWasEmpty
      );

      if (pieceImg) {
        addMoveToHistory(fromSquare, toSquare, pieceImg);
      }

      toggleMoveColor();

      const playerColorSelect = document.getElementById("player-color");
      const selectedColor =
        playerColorSelect.value === "white" ? COLOR.WHITE : COLOR.BLACK;
      setBoardView(selectedColor);
    }
  } catch (error) {
    console.error("Move failed:", error);
  } finally {
    resetMoveState();
  }
}

async function applyMove(
  fromSquare,
  toSquare,
  pieceImg,
  fromCoords,
  toCoords,
  needsPromotion,
  promotion,
  isKing,
  isPawn,
  toSquareWasEmpty
) {
  const movingPieceContent = getMovingPiece();
  const pieceSrc = pieceImg ? pieceImg.getAttribute("src") : "";
  const isWhite = pieceSrc.includes("w");
  const fromNotation = getSquareNotation(fromCoords.row, fromCoords.col);
  const toNotation = getSquareNotation(toCoords.row, toCoords.col);

  if (isKing) {
    const colDiff = Math.abs(toCoords.col - fromCoords.col);
    if (colDiff === 2) {
      const isKingside = toCoords.col > fromCoords.col;
      handleCastling(fromCoords.row, fromCoords.col, isKingside);
      return;
    }
  }

  if (isPawn && toSquareWasEmpty) {
    const colDiff = Math.abs(toCoords.col - fromCoords.col);
    if (colDiff === 1) {
      handleEnPassant(fromSquare, toSquare, fromNotation, toNotation);
      return;
    }
  }

  if (needsPromotion && promotion) {
    const colorPrefix = isWhite ? "w" : "b";
    const promoType = `${colorPrefix}${promotion.toUpperCase()}`;
    toSquare.innerHTML = "";
    const newPiece = setPiece(promoType, toCoords.row, toCoords.col);
    toSquare.appendChild(newPiece);
    fromSquare.innerHTML = "";
  } else {
    toSquare.innerHTML = movingPieceContent;
    fromSquare.innerHTML = "";
  }
}

function handleCastling(kingRow, kingCol, isKingside) {
  const kingFrom = document.getElementById(`${kingRow}-${kingCol}`);
  const kingToCol = isKingside ? kingCol + 2 : kingCol - 2;
  const kingTo = document.getElementById(`${kingRow}-${kingToCol}`);

  const rookFromCol = isKingside ? 8 : 1;
  const rookToCol = isKingside ? kingCol + 1 : kingCol - 1;
  const rookFrom = document.getElementById(`${kingRow}-${rookFromCol}`);
  const rookTo = document.getElementById(`${kingRow}-${rookToCol}`);

  if (kingFrom && kingTo && rookFrom && rookTo) {
    kingTo.innerHTML = kingFrom.innerHTML;
    kingFrom.innerHTML = "";
    rookTo.innerHTML = rookFrom.innerHTML;
    rookFrom.innerHTML = "";
  }
}

function handleEnPassant(fromSquare, toSquare, fromNotation, toNotation) {
  const movingPieceContent = getMovingPiece();

  const file = toNotation[0];
  const sourceRank = fromNotation[1];
  const capturedSquareNotation = file + sourceRank;

  const fileChar = file.charCodeAt(0) - "a".charCodeAt(0) + 1;
  const ranks = ["8", "7", "6", "5", "4", "3", "2", "1"];
  const row = ranks.indexOf(sourceRank) + 1;
  const col = fileChar;

  const capturedSquare = document.getElementById(`${row}-${col}`);
  if (capturedSquare) {
    capturedSquare.innerHTML = "";
  }

  toSquare.innerHTML = movingPieceContent;
  fromSquare.innerHTML = "";
}

function resetMoveState() {
  setIsPieceMoving(false);
  clearMovingPiece();
  clearInitialSquare();
}
