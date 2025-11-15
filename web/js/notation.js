export function getSquareNotation(row, col) {
  const files = ["a", "b", "c", "d", "e", "f", "g", "h"];
  const ranks = ["8", "7", "6", "5", "4", "3", "2", "1"];
  return files[col - 1] + ranks[row - 1];
}

export function getPieceNotation(pieceImg) {
  if (!pieceImg) return "";
  const src = pieceImg.getAttribute("src") || "";
  const pieceMatch = src.match(/\/([wb][KQRNBP])\.svg/);
  if (!pieceMatch) return "";

  const piece = pieceMatch[1];
  const type = piece[1];

  const pieceMap = {
    K: "K",
    Q: "Q",
    R: "R",
    B: "B",
    N: "N",
    P: "",
  };

  return pieceMap[type] || "";
}

export function parseSquareId(squareId) {
  const parts = squareId.split("-");
  if (parts.length !== 2) {
    return null;
  }
  return {
    row: parseInt(parts[0]),
    col: parseInt(parts[1]),
  };
}
