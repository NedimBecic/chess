export function showPromotionDialog(isWhite) {
  return new Promise((resolve) => {
    const modal = document.createElement("div");
    modal.className = "promotion-modal";
    modal.style.cssText = `
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background-color: rgba(0, 0, 0, 0.7);
      display: flex;
      justify-content: center;
      align-items: center;
      z-index: 1000;
    `;

    const dialog = document.createElement("div");
    dialog.style.cssText = `
      background-color: #2d2d2d;
      padding: 24px;
      border-radius: 8px;
      border: 2px solid #404040;
      display: flex;
      gap: 16px;
    `;

    const pieces = isWhite
      ? [
          { type: "wQ", label: "Queen" },
          { type: "wR", label: "Rook" },
          { type: "wB", label: "Bishop" },
          { type: "wN", label: "Knight" },
        ]
      : [
          { type: "bQ", label: "Queen" },
          { type: "bR", label: "Rook" },
          { type: "bB", label: "Bishop" },
          { type: "bN", label: "Knight" },
        ];

    pieces.forEach((piece) => {
      const button = document.createElement("button");
      button.style.cssText = `
        background: transparent;
        border: 2px solid #555;
        border-radius: 4px;
        padding: 8px;
        cursor: pointer;
        transition: all 0.2s;
      `;

      const img = document.createElement("img");
      img.src = `pieces/${piece.type}.svg`;
      img.alt = piece.label;
      img.style.cssText = `
        width: 48px;
        height: 48px;
        display: block;
      `;

      button.appendChild(img);
      button.title = piece.label;

      button.addEventListener("mouseenter", () => {
        button.style.borderColor = "#888";
        button.style.backgroundColor = "#3a3a3a";
      });

      button.addEventListener("mouseleave", () => {
        button.style.borderColor = "#555";
        button.style.backgroundColor = "transparent";
      });

      button.addEventListener("click", () => {
        const promoChar = piece.type[1].toLowerCase();
        document.body.removeChild(modal);
        resolve(promoChar);
      });

      dialog.appendChild(button);
    });

    modal.appendChild(dialog);
    document.body.appendChild(modal);
  });
}
