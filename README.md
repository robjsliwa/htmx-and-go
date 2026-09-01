# Adventures in Go and HTMX

Source code for the **Adventures in Go and HTMX** blog series on [ShiftLeftAI](https://www.shiftleftai.dev/posts/adventures-in-go-htmx-1/), co-authored by Rob Sliwa and Pawan Tripathi.

The series builds a retro terminal-style text adventure game step-by-step, using Go for the server and HTMX for interactivity — no heavy JavaScript framework required.

## Blog Series

| Part | Topic | Link |
|------|-------|------|
| 1 | Server setup, Go templates, hypermedia architecture | [Adventures in Go and HTMX - Part 1](https://www.shiftleftai.dev/posts/adventures-in-go-htmx-1/) |
| 2 | HTMX `hx-boost` for SPA-like navigation | [Adventures in Go and HTMX - Part 2](https://www.shiftleftai.dev/posts/adventures-in-go-htmx-2/) |
| 3 | Form input, HTMX fragment-swapping, auto-reset, view transitions | [Adventures in Go and HTMX - Part 3](https://www.shiftleftai.dev/posts/adventures-in-go-htmx-3/) |
| 4 | Inventory system, `hx-delete`, out-of-band swaps (`hx-swap-oob`), modular package refactor | [Adventures in Go and HTMX - Part 4](https://www.shiftleftai.dev/posts/adventures-in-go-htmx-4/) |
| 5 | Optimistic UI updates, equipment system, `hx-on` lifecycle events, toggle equip/unequip pattern | [Adventures in Go and HTMX - Part 5](https://www.shiftleftai.dev/posts/adventures-in-go-htmx-5/) |
| 6 | Active search spellbook, debounced `hx-trigger`, server-side filtering, `<dialog>` modal, loading indicators | [Adventures in Go and HTMX - Part 6](https://www.shiftleftai.dev/posts/adventures-in-go-htmx-6/) |
| 7 | Lazy loading room images, infinite scroll game log, `hx-trigger="intersect once"` | [Adventures in Go and HTMX - Part 7](https://www.shiftleftai.dev/posts/adventures-in-go-htmx-7/) |
| 8 | Server-generated SVG world map, polling with `hx-trigger="every 2s"`, wandering monster AI, `HX-Trigger` response header cascading updates | [Adventures in Go and HTMX - Part 8](https://www.shiftleftai.dev/posts/adventures-in-go-htmx-8/) |
| 9 | Data-driven world design with YAML, struct tags, fail-fast startup validation, map-based room exits | [Adventures in Go and HTMX - Part 9](https://www.shiftleftai.dev/posts/adventures-in-go-htmx-9/) |

## Repository Structure

```
htmx-and-go/
├── part1/adv-htmx/    # Part 1 — basic server + Go templates (full page reloads)
├── part2/adv-htmx/    # Part 2 — adds hx-boost for smooth navigation
├── part3/adv-htmx/    # Part 3 — form input, fragment-swapping, view transitions
├── part4/adv-htmx/    # Part 4 — inventory system, hx-delete, out-of-band swaps, modular packages
├── part5/adv-htmx/    # Part 5 — optimistic UI updates, equipment system, hx-on lifecycle events
├── part6/adv-htmx/    # Part 6 — active search spellbook, debounced hx-trigger, dialog modal
├── part7/adv-htmx/    # Part 7 — lazy-loaded room images, infinite scroll game log via intersect
├── part8/adv-htmx/    # Part 8 — server-generated SVG world map, polling, wandering monster AI
└── part9/adv-htmx/    # Part 9 — data-driven world via YAML, struct tags, fail-fast validation
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
cd part1/adv-htmx   # or part2/adv-htmx, part3/adv-htmx, part4/adv-htmx, part5/adv-htmx, part6/adv-htmx, part7/adv-htmx, part8/adv-htmx, part9/adv-htmx
go run .
```

Open [http://localhost:4040](http://localhost:4040) in your browser.

### Development with Live Reload (optional)

```bash
go install github.com/air-verse/air@latest
cd part1/adv-htmx   # or part2/adv-htmx, part3/adv-htmx, part4/adv-htmx, part5/adv-htmx, part6/adv-htmx, part7/adv-htmx, part8/adv-htmx, part9/adv-htmx
air
```

## Key Concepts

- **Hypermedia architecture** — the server is the application; HTML is the API
- **Go 1.22 routing** — typed path parameters with `r.PathValue()`
- **`hx-boost`** — progressively enhances standard links to use AJAX, keeping full-page fallback when JS is unavailable
- **Server-driven state** — no client-side state management; rooms and game state live on the server
- **`hx-swap-oob`** — out-of-band swaps let a single server response update multiple independent regions of the page
- **`hx-delete`** — maps HTTP DELETE to HTMX-powered removal of elements without page reloads
- **Optimistic UI updates** — the browser immediately reflects the expected outcome via `hx-on` lifecycle events (`beforeRequest`, `afterRequest`, `responseError`), then rolls back on server error
- **`hx-on` lifecycle events** — fine-grained client-side hooks that drive the optimistic pattern without custom JavaScript frameworks
- **Toggle equip pattern** — a single `ToggleEquip()` endpoint replaces separate equip/unequip routes; the server decides the outcome based on current state
- **Active search** — a debounced `hx-trigger="input changed delay:500ms, search"` sends server-side queries only after typing pauses, returning filtered HTML fragments instead of JSON
- **`<dialog>` modal** — the native HTML dialog element powers the spellbook UI with built-in overlay and focus management
- **`htmx-request` loading indicators** — HTMX automatically toggles this class during in-flight requests, enabling CSS-only loading spinners
- **Lazy loading** — the room image is deferred behind `hx-trigger="load"`, fetched only after the page renders so a skeleton placeholder shows first and the initial page stays fast
- **Infinite scroll** — a sentinel element at the top of the game log uses HTMX's built-in `hx-trigger="intersect once"` to fire when it scrolls into view, fetching the next page of older log entries and replacing itself with a new sentinel until history runs out
- **Server-generated SVG maps** — the world map is rendered server-side as an SVG string and swapped into the DOM, avoiding bitmap image assets entirely
- **Polling with `hx-trigger="every 2s"`** — the map polls the server on a fixed interval to simulate a "living world" without WebSockets or background goroutines
- **Monster AI** — a wandering goblin moves between rooms based on probability logic evaluated on each poll cycle, with its position tracked in the server-side game state
- **`HX-Trigger` response headers** — the server signals client-side events via response headers, letting one request (e.g. the map poll) cascade updates to other page regions (e.g. the game log) without extra round trips
- **Data-driven world design** — the game world moves from hardcoded Go data into `world.yaml`, letting level designers edit rooms and items without recompiling
- **YAML struct tags** — Go structs use `yaml:"..."` tags to map lowercase YAML keys onto exported Go fields
- **Fail-fast validation** — a startup validation pass checks referential integrity (e.g. room exits point to real rooms, items exist) and aborts with a descriptive error before the server starts if the data is invalid
- **Map-based room exits** — exits move from a structured list to `map[string]string`, simplifying data entry and requiring updated movement logic and template iteration

## License

MIT — see [LICENSE](LICENSE).
