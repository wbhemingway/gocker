# Gocker
**Gocker** is a minimalist, straightforward time and pay tracker built in Go. Designed specifically for freelance contract work, built to help accurately record time spent to bill.

## Key Features
* **Session Management:** Start and stop work sessions from the CLI with single-word commands or through the GUI.
* **Break Tracking:** Log breaks without stopping the main session for accurate billing.
* **Freelance Finance:** Automated tax estimates based on location.
* **Data Integrity:** Powered by **SQLite** and **sqlc** for type-safe, local-first persistence.
* **Native Performance:** Compiles to a zero-dependency Windows binary using Go drivers.

## Tech Stack
* **Language:** Go (Golang)
* **Database:** SQLite via `modernc.org/sqlite` (CGO-free)
* **Codegen:** `sqlc`
* **CLI Framework:** Cobra

---

## Project Roadmap

### Phase 1: The CLI Foundation
- [x] Initialize project scaffolding (Standard Go Layout).
- [ ] Define SQLite schema for `sessions` and `breaks`.
- [ ] Implement the `internal/engine` interface.
- [ ] Basic Cobra commands: `gocker start`, `gocker stop`, `gocker status`.

### Phase 2: Finance & Reporting
- [ ] Implement the `finance` package.
- [ ] Add `gocker report` command to show daily/weekly earnings.
- [ ] Implement CSV export logic for pay period cross-referencing.
- [ ] Add `gocker edit` to manually adjust timestamps.

### Phase 3: The Native GUI (Windows)
- [ ] Transition to a native Windows GUI.
- [ ] Implement a "Live Timer" dashboard.

### Phase 4: Polish & Quality of Life
**TBD**
