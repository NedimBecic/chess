const blackPiecesStartingRank = [
  "bR",
  "bN",
  "bB",
  "bQ",
  "bK",
  "bB",
  "bN",
  "bR",
];
const whitePiecesStartingRank = [
  "wR",
  "wN",
  "wB",
  "wQ",
  "wK",
  "wB",
  "wN",
  "wR",
];

const COLOR = {
  WHITE: "White",
  BLACK: "Black",
};

let board = document.getElementById("board");
let currentMoveColor = COLOR.WHITE;
let isPieceMoving = false;
let movingPiece = null;
let initialSquare = null;

function setBoardView(color) {
  // const rows = board.rows;
  // for (let i = 0; i < 4; i++) {
  //   let currentCells = rows[i].cells;
  //   let swapCells = rows[7 - i].cells;
  //
  //   for (let j = 0; j < 8; j++) {
  //     let temp = swapCells[7 - j].innerHTML;
  //     swapCells[7 - j].innerHTML = currentCells[j].innerHTML;
  //     currentCells[j].innerHTML = temp;
  //   }
  // }
  if (color == COLOR.BLACK) {
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
        color == COLOR.BLACK
          ? piece.classList.add("flipped-piece")
          : piece.classList.remove("flipped-piece");
      }
    }
  }
}

function movePiece(square) {
  if (square == initialSquare) {
    return;
  }

  if (!isPieceMoving) {
    movingPiece = square.innerHTML;
    initialSquare = square;
    isPieceMoving = true;
  } else {
    square.innerHTML = movingPiece;
    initialSquare.innerHTML = "";
    isPieceMoving = false;
    currentMoveColor == COLOR.WHITE
      ? (currentMoveColor = COLOR.BLACK)
      : (currentMoveColor = COLOR.WHITE);
    setBoardView(currentMoveColor);
  }
}

function dragStart(ev) {
  ev.dataTransfer.setData("text", ev.target.id);
}

function dragOver(ev) {
  ev.preventDefault();
}

function drop(ev) {
  ev.preventDefault();
  const data = ev.dataTransfer.getData("text");
  ev.target.appendChild(document.getElementById(data));
}

function setPiece(pieceType, row, col) {
  let piece = document.createElement("img");
  piece.setAttribute("src", `pieces/${pieceType}.svg`);
  piece.setAttribute("alt", `${pieceType}`);
  piece.setAttribute("id", `${pieceType}-${row}-${col}`);
  piece.addEventListener("dragstart", dragStart);
  return piece;
}

for (let i = 1; i <= 8; i++) {
  let row = document.createElement("tr");
  board.append(row);
  for (let j = 1; j <= 8; j++) {
    let cell = document.createElement("td");
    cell.setAttribute("draggable", "true");
    cell.setAttribute("id", `${i}-${j}`);
    cell.addEventListener("dragover", dragOver);
    cell.addEventListener("drop", drop);
    row.append(cell);
    cell.onclick = function () {
      movePiece(cell);
    };

    if ((i + j) % 2 == 0) {
      cell.classList.add("white-field");
    } else {
      cell.classList.add("black-field");
    }

    if (i == 1) {
      cell.append(setPiece(blackPiecesStartingRank[j - 1], i, j));
    }

    if (i == 2) {
      cell.append(setPiece("bP", i, j));
    }

    if (i == 7) {
      cell.append(setPiece("wP", i, j));
    }

    if (i == 8) {
      cell.append(setPiece(whitePiecesStartingRank[j - 1], i, j));
    }
  }
}
