# Gocker
**Gocker** is a minimalist, straightforward time and pay tracker built in Go. Designed specifically for freelance contract work; built to help accurately record time spent to bill.

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

### Phase 1: The CLI & Internal Logic
- [x] Initialize project scaffolding.
- [x] Define SQLite schema for `sessions`, `tags` and `entry_tags`.
- [ ] Implement the `internal/engine` interface.
- [ ] Impliment basic Cobra commands `gocker start`, `gocker stop`, `gocker status`.

### Phase 2: Finance & Reporting
- [ ] Implement the `finance` package.
- [ ] Add `gocker report` command to show daily/weekly earnings.
- [ ] Status Command: A "Live" dashboard in the terminal showing current session duration, active tags, and real-time earnings.
- [ ] Implement CSV export logic for pay period cross-referencing.
- [ ] Add `gocker edit` to manually adjust timestamps.

### Phase 3: The Native GUI
- [ ] Implement Fyne for a consistent native UI across Windows, Mac, and Linux.
- [ ] Implement a "Live Timer" dashboard.
- [ ] Implement a table-based manager interface to easily edit timestamps, adjust hourly rates, or re-tag previous sessions.

### Phase 4: Polish & Quality of Life
- [ ] Cross-Platform Paths: Implement os.UserConfigDir to ensure the .db file is stored correctly on Windows, Linux, and macOS.
- [ ] Implement a --back or --offset flag for the CLI to account for forgotten start times.
- [ ] Add a rounding toggle for "Billable" vs. "Actual" time
- [ ] Add Automated backup logic to automatically mirror the gocker.db to a backup directory on session completion.
- [ ] Add theme support for light and dark
