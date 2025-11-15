import { initializeBoard } from "./js/board-init.js";
import { initializeUIControls } from "./js/ui-controls.js";
import { startGame } from "./js/api.js";

document.addEventListener("DOMContentLoaded", function () {
  initializeBoard();
  initializeUIControls();

  startGame().catch((error) => {
    console.error("Failed to initialize game:", error);
  });
});
