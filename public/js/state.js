import { COLOR } from "./constants.js";

export let board = null;
export let currentMoveColor = COLOR.WHITE;
export let isPieceMoving = false;
export let movingPiece = null;
export let initialSquare = null;
export let moveHistory = [];
export let moveNumber = 1;
export let playerViewColor = COLOR.WHITE;
export let currentFEN = null;

export function setBoard(boardElement) {
  board = boardElement;
}

export function resetGameState() {
  currentMoveColor = COLOR.WHITE;
  moveHistory = [];
  moveNumber = 1;
  isPieceMoving = false;
  movingPiece = null;
  initialSquare = null;
}

export function setCurrentMoveColor(color) {
  currentMoveColor = color;
}

export function getCurrentMoveColor() {
  return currentMoveColor;
}

export function toggleMoveColor() {
  currentMoveColor =
    currentMoveColor === COLOR.WHITE ? COLOR.BLACK : COLOR.WHITE;
}

export function setMovingPiece(piece) {
  movingPiece = piece;
}

export function getMovingPiece() {
  return movingPiece;
}

export function clearMovingPiece() {
  movingPiece = null;
}

export function setInitialSquare(square) {
  initialSquare = square;
}

export function getInitialSquare() {
  return initialSquare;
}

export function clearInitialSquare() {
  initialSquare = null;
}

export function setIsPieceMoving(value) {
  isPieceMoving = value;
}

export function getIsPieceMoving() {
  return isPieceMoving;
}

export function setPlayerViewColor(color) {
  playerViewColor = color;
}

export function getPlayerViewColor() {
  return playerViewColor;
}

export function addToMoveHistory(move) {
  moveHistory.push(move);
}

export function getMoveHistory() {
  return moveHistory;
}

export function clearMoveHistory() {
  moveHistory = [];
}

export function incrementMoveNumber() {
  moveNumber++;
}

export function getMoveNumber() {
  return moveNumber;
}

export function resetMoveNumber() {
	moveNumber = 1;
}

export function setCurrentFEN(fen) {
	currentFEN = fen;
}

export function getCurrentFEN() {
	return currentFEN;
}

export function resetFEN() {
	currentFEN = null;
}
