export function setupDragAndDrop() {
  const board = document.getElementById("board");
  if (!board) return;

  const rows = board.rows;
  for (let i = 0; i < rows.length; i++) {
    const cells = rows[i].cells;
    for (let j = 0; j < cells.length; j++) {
      cells[j].addEventListener("dragover", dragOver);
      cells[j].addEventListener("drop", drop);
    }
  }
}

function dragOver(ev) {
  ev.preventDefault();
}

function drop(ev) {
  ev.preventDefault();
  const data = ev.dataTransfer.getData("text");
  ev.target.appendChild(document.getElementById(data));
}
