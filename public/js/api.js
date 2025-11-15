const API_BASE_URL = "/api";

async function apiRequest(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  const config = {
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    ...options,
  };

  try {
    const response = await fetch(url, config);
    let data;

    try {
      data = await response.json();
    } catch (jsonError) {
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      throw new Error("Invalid JSON response from server");
    }

    if (!response.ok) {
      const error = new Error(
        data.message || data.Message || `HTTP error! status: ${response.status}`
      );
      error.response = data;
      error.status = response.status;
      throw error;
    }

    return data;
  } catch (error) {
    if (error.response || error.status) {
      throw error;
    }
    console.error("API request failed:", error);
    throw new Error("Network error: Could not connect to server");
  }
}

export async function startGame() {
  return apiRequest("/start", {
    method: "POST",
    body: JSON.stringify({}),
  });
}

export async function getGameState() {
  return apiRequest("/moves", {
    method: "GET",
  });
}

export async function makeMove(from, to, promotion = null) {
  const body = { from, to };
  if (promotion) {
    body.promotion = promotion;
  }
  return apiRequest("/move", {
    method: "POST",
    body: JSON.stringify(body),
  });
}
