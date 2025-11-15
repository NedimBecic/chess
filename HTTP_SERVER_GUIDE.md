# HTTP Server Guide - Go net/http Package

This guide explains the Go `net/http` package concepts used in this chess game server.

## Project Structure

```
chess-game/
├── main.go                    # Main server entry point
├── go.mod                     # Go module definition
├── internal/
│   ├── chess/                # Chess game logic
│   │   ├── board.go          # Board representation
│   │   ├── move.go           # Move representation
│   │   └── moves.go          # Move generation
│   └── handlers/             # HTTP request handlers
│       └── api.go            # API endpoint handlers
└── web/                      # Frontend static files
    ├── chess.html
    ├── chess.js
    ├── style.css
    └── pieces/               # Piece images
```

## Key Concepts from net/http Package

### 1. http.Handler Interface

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

**Explanation:**

- `Handler` is an interface that any type can implement
- It requires one method: `ServeHTTP(w, r)`
- This method processes HTTP requests and writes responses
- Any type that implements this interface can handle HTTP requests

**Why it matters:**

- This is the foundation of Go's HTTP server
- Everything that handles requests implements this interface
- Allows for flexible request handling patterns

### 2. http.HandlerFunc Type

```go
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

**Explanation:**

- `HandlerFunc` is a type that adapts functions to the `Handler` interface
- It's a function type that matches the `ServeHTTP` signature
- Go automatically converts functions to `HandlerFunc` when needed
- This allows you to use regular functions as HTTP handlers

**Example:**

```go
func myHandler(w http.ResponseWriter, r *http.Request) {
    // Handle request
}

// This works because HandlerFunc implements Handler
http.HandleFunc("/path", myHandler)
```

**Why it matters:**

- Makes it easy to write handlers as simple functions
- No need to create types that implement the Handler interface
- The standard library handles the conversion automatically

### 3. http.ServeMux (Multiplexer)

```go
type ServeMux struct {
    // ... internal fields
}

func NewServeMux() *ServeMux
func (mux *ServeMux) Handle(pattern string, handler Handler)
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```

**Explanation:**

- `ServeMux` is an HTTP request router (multiplexer)
- It matches incoming requests against registered URL patterns
- Routes requests to the appropriate handler
- Patterns can be exact matches or prefixes (ending with `/`)

**Pattern Matching:**

- `/api/moves` - exact match
- `/api/` - prefix match (matches `/api/*`)
- Longer patterns take precedence

**Example:**

```go
mux := http.NewServeMux()
mux.HandleFunc("/api/moves", handleMoves)    // Exact match
mux.Handle("/api/", apiHandler)               // Prefix match
mux.Handle("/", fileServer)                   // Catch-all
```

**Why it matters:**

- Organizes routing logic
- Separates concerns (API vs frontend)
- Allows middleware chaining
- Standard way to route requests in Go

### 4. http.ResponseWriter Interface

```go
type ResponseWriter interface {
    Header() Header
    Write([]byte) (int, error)
    WriteHeader(statusCode int)
}
```

**Explanation:**

- `ResponseWriter` is an interface for writing HTTP responses
- `Header()` - returns the response headers map
- `Write([]byte)` - writes response body
- `WriteHeader(int)` - sets the status code

**Common Usage:**

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
w.Write([]byte(`{"status": "ok"}`))
```

**Why it matters:**

- Interface allows for different implementations
- Can be wrapped for middleware (logging, compression, etc.)
- Standard way to send responses

### 5. http.Request Structure

```go
type Request struct {
    Method string
    URL    *url.URL
    Header Header
    Body   io.ReadCloser
    // ... more fields
}
```

**Explanation:**

- `Request` contains all information about an HTTP request
- `Method` - HTTP method (GET, POST, etc.)
- `URL` - request URL with path, query, etc.
- `Header` - HTTP headers
- `Body` - request body (for POST/PUT)

**Common Usage:**

```go
method := r.Method                    // "GET", "POST", etc.
path := r.URL.Path                    // "/api/moves"
apiKey := r.Header.Get("API-KEY")     // Get header value
```

**Why it matters:**

- Contains all request information
- Read-only (safe to pass around)
- Standard way to access request data

### 6. http.FileServer

```go
func FileServer(root FileSystem) Handler
```

**Explanation:**

- `FileServer` creates a handler that serves files from a directory
- Takes a `FileSystem` (like `http.Dir`) as root
- Serves static files (HTML, CSS, JS, images)
- Handles directory listings and 404s

**Example:**

```go
fileServer := http.FileServer(http.Dir("./web"))
http.Handle("/", fileServer)
```

**Why it matters:**

- Simple way to serve static files
- Handles file system operations
- Production-ready file serving

### 7. http.Dir Type

```go
type Dir string

func (d Dir) Open(name string) (File, error)
```

**Explanation:**

- `Dir` is a string that implements the `FileSystem` interface
- Converts a directory path to a file system
- Used with `FileServer` to serve files

**Example:**

```go
fs := http.Dir("./web")        // Directory path
fileServer := http.FileServer(fs)  // Create file server
```

**Why it matters:**

- Simple way to serve files from a directory
- Type-safe file system access
- Standard library implementation

### 8. http.StripPrefix

```go
func StripPrefix(prefix string, h Handler) Handler
```

**Explanation:**

- `StripPrefix` removes a prefix from the request URL
- Useful when serving files from a subdirectory
- Allows serving files with different URL paths

**Example:**

```go
// Request: /api/moves
// After StripPrefix("/api", ...): /moves
mux.Handle("/api/", http.StripPrefix("/api", apiMux))
```

**Why it matters:**

- Allows clean URL structures
- Separates routing from file structure
- Common pattern in web servers

### 9. http.ListenAndServe

```go
func ListenAndServe(addr string, handler Handler) error
```

**Explanation:**

- `ListenAndServe` starts an HTTP server
- Listens on the specified address (e.g., `:8080`)
- Calls the handler for each request
- Blocks until the server stops

**Example:**

```go
err := http.ListenAndServe(":8080", mux)
if err != nil {
    log.Fatal(err)
}
```

**Why it matters:**

- Simple way to start a server
- Production-ready implementation
- Handles connection management

### 10. http.Error

```go
func Error(w ResponseWriter, error string, code int)
```

**Explanation:**

- `Error` writes an error response
- Sets the status code
- Writes the error message as the body
- Convenient for error handling

**Example:**

```go
http.Error(w, "Unauthorized", http.StatusUnauthorized)
// Equivalent to:
w.WriteHeader(http.StatusUnauthorized)
w.Write([]byte("Unauthorized"))
```

**Why it matters:**

- Convenient error handling
- Standard way to send errors
- Sets proper status codes

## Middleware Pattern

Middleware is a function that wraps a handler to add functionality.

```go
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do something before the handler
        // (e.g., check authentication)

        // Call the next handler
        next.ServeHTTP(w, r)

        // Do something after the handler
        // (e.g., log the response)
    })
}
```

**Why it matters:**

- Reusable functionality (logging, auth, CORS)
- Composable (chain multiple middlewares)
- Clean separation of concerns

## Request Flow in Our Server

1. **Client sends request** → `http://localhost:8080/api/moves`

2. **Server receives request** → `http.ListenAndServe` accepts connection

3. **Main ServeMux routes request** → Matches `/api/` pattern

4. **Middleware processes** → `apiKeyMiddleware` checks API key

5. **StripPrefix removes prefix** → `/api/moves` → `/moves`

6. **API ServeMux routes** → Matches `/moves` pattern

7. **Handler processes** → `HandleGetMoves` generates response

8. **Response sent** → JSON response written to client

## Common Patterns

### Pattern 1: Separate API and Frontend Routes

```go
apiMux := http.NewServeMux()
apiMux.HandleFunc("/moves", handleMoves)

frontendMux := http.NewServeMux()
frontendMux.Handle("/", http.FileServer(http.Dir("./web")))

mainMux := http.NewServeMux()
mainMux.Handle("/api/", apiMux)
mainMux.Handle("/", frontendMux)
```

**Why:**

- Clean separation of concerns
- Different middleware for API vs frontend
- Easier to test and maintain

### Pattern 2: Method Checking

```go
func handleMoves(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    // Handle GET request
}
```

**Why:**

- Enforce HTTP method restrictions
- Clear error messages
- Security best practice

### Pattern 3: JSON Encoding/Decoding

```go
// Decode request
var req MoveRequest
json.NewDecoder(r.Body).Decode(&req)

// Encode response
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(response)
```

**Why:**

- Standard way to handle JSON
- Type-safe serialization
- Built-in error handling

## Security Considerations

1. **API Key Authentication** - Middleware checks API key in headers
2. **Method Validation** - Handlers check HTTP methods
3. **Input Validation** - Validate and sanitize all inputs
4. **Error Handling** - Don't expose internal errors to clients
5. **CORS** - Add CORS headers if serving from different origins

## Testing the Server

1. **Start the server:**

   ```bash
   go run main.go
   ```

2. **Test frontend:**

   ```bash
   curl http://localhost:8080/
   ```

3. **Test API (with API key):**

   ```bash
   export API_KEY=your-secret-key
   curl -H "API-KEY: your-secret-key" http://localhost:8080/api/moves
   ```

4. **Test API (without API key):**
   ```bash
   curl http://localhost:8080/api/moves
   # Should return 401 Unauthorized
   ```

## Next Steps

1. **Add game state storage** - Store game states in memory or database
2. **Implement move validation** - Validate moves before applying
3. **Add CORS support** - Allow frontend to call API from browser
4. **Add logging** - Log requests and responses
5. **Add rate limiting** - Prevent abuse
6. **Add WebSocket support** - Real-time game updates

## References

- [Go net/http package documentation](https://pkg.go.dev/net/http)
- [Go by Example: HTTP Servers](https://gobyexample.com/http-servers)
- [Writing Web Applications in Go](https://go.dev/doc/articles/wiki/)
