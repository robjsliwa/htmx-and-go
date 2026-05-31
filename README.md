# Adventures in Go and HTMX

Source code for the **Adventures in Go and HTMX** blog series on [ShiftLeftAI](https://www.shiftleftai.dev/posts/adventures-in-go-htmx-1/), co-authored by Rob Sliwa and Pawan Tripathi.

The series builds a retro terminal-style text adventure game step-by-step, using Go for the server and HTMX for interactivity — no heavy JavaScript framework required.

## Blog Series

| Part | Topic | Link |
|------|-------|------|
| 1 | Server setup, Go templates, hypermedia architecture | [Adventures in Go and HTMX - Part 1](https://www.shiftleftai.dev/posts/adventures-in-go-htmx-1/) |
| 2 | HTMX `hx-boost` for SPA-like navigation | [Adventures in Go and HTMX - Part 2](https://www.shiftleftai.dev/posts/adventures-in-go-htmx-2/) |
| 3 | Form input, HTMX fragment-swapping, auto-reset, view transitions | [Adventures in Go and HTMX - Part 3](https://www.shiftleftai.dev/posts/adventures-in-go-htmx-3/) |

## Repository Structure

```
htmx-and-go/
├── part1/adv-htmx/    # Part 1 — basic server + Go templates (full page reloads)
├── part2/adv-htmx/    # Part 2 — adds hx-boost for smooth navigation
└── part3/adv-htmx/    # Part 3 — form input, fragment-swapping, view transitions
```

Each part is a standalone Go module you can run independently.

## Tech Stack

- **Go 1.22+** — standard library HTTP server with pattern-based routing
- **HTMX** — lightweight hypermedia library for partial page updates
- **Go `html/template`** — server-side HTML rendering
- **[Air](https://github.com/air-verse/air)** *(optional)* — live reload during development

## Getting Started

### Prerequisites

- Go 1.22 or later

### Run a Part

```bash
cd part1/adv-htmx   # or part2/adv-htmx or part3/adv-htmx
go run .
```

Open [http://localhost:4040](http://localhost:4040) in your browser.

### Development with Live Reload (optional)

```bash
go install github.com/air-verse/air@latest
cd part1/adv-htmx   # or part2/adv-htmx or part3/adv-htmx
air
```

## Key Concepts

- **Hypermedia architecture** — the server is the application; HTML is the API
- **Go 1.22 routing** — typed path parameters with `r.PathValue()`
- **`hx-boost`** — progressively enhances standard links to use AJAX, keeping full-page fallback when JS is unavailable
- **Server-driven state** — no client-side state management; rooms and game state live on the server

## License

MIT — see [LICENSE](LICENSE).
