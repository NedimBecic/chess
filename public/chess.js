import { initializeBoard } from "./js/board-init.js";
import { initializeUIControls } from "./js/ui-controls.js";
import { startGame } from "./js/api.js";
import { setCurrentFEN } from "./js/state.js";

document.addEventListener("DOMContentLoaded", function () {
  initializeBoard();
  initializeUIControls();

  startGame()
    .then((response) => {
      if (response && response.fen) {
        setCurrentFEN(response.fen);
      }
    })
    .catch((error) => {
      console.error("Failed to initialize game:", error);
    });
});
