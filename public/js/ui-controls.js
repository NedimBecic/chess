import { COLOR } from "./constants.js";
import { setBoardView } from "./board-view.js";
import { resetBoard } from "./board-manipulation.js";
import { clearMoveHistory } from "./move-history.js";
import { startGame } from "./api.js";

export function initializeUIControls() {
  const playerColorSelect = document.getElementById("player-color");
  if (playerColorSelect) {
    playerColorSelect.addEventListener("change", function () {
      const selectedColor = this.value === "white" ? COLOR.WHITE : COLOR.BLACK;
      setBoardView(selectedColor);
    });
    const selectedColor =
      playerColorSelect.value === "white" ? COLOR.WHITE : COLOR.BLACK;
    setBoardView(selectedColor);
  }

  const resetBoardBtn = document.getElementById("reset-board-btn");
  if (resetBoardBtn) {
    resetBoardBtn.addEventListener("click", function () {
      resetBoard();
      clearMoveHistory();
    });
  }

  const newGameBtn = document.getElementById("new-game-btn");
  if (newGameBtn) {
    newGameBtn.addEventListener("click", async function () {
      try {
        await startGame();
        resetBoard();
        clearMoveHistory();
      } catch (error) {
        console.error("Failed to start game:", error);
        alert("Failed to start game. Please try again.");
      }
    });
  }

  const moveDepthSelect = document.getElementById("move-depth");
  if (moveDepthSelect) {
    moveDepthSelect.addEventListener("change", function () {
      console.log("Move depth changed to:", this.value);
    });
  }

  const engineEnabledCheckbox = document.getElementById("engine-enabled");
  if (engineEnabledCheckbox) {
    engineEnabledCheckbox.addEventListener("change", function () {
      console.log("Engine enabled:", this.checked);
    });
  }
}
