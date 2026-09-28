# World Cup 2022 CLI Dashboard

A full-screen terminal dashboard for browsing World Cup 2022 matches, group tables, knockout brackets, lineups, events, and player stats — built in Go with Bubble Tea.

```
┌──────────────────────────────────────────────────────────────────────────────┐
│  Match navigation strip (paginated, selected match highlighted)              │
│  Venue            Match status / LIVE minute            Date                 │
│                                                                              │
│   Home team column        Big block-digit scoreline       Away team column    │
│   • True-color flag       2 - 0                           • True-color flag   │
│   • Lineup table          Match events (goals/cards/subs) • Lineup table      │
│     (# Player ⚽ ■ ■)                                                          │
│                                                                              │
│  ◄/a/h prev match   ►/d/l/space next match   q/ctrl+c quit                    │
│  Group table (group stage)  — or —  Knockout bracket (1/8 → Final)            │
│  API: local | spinner Refreshing... / ❌ err / Enjoy the match. | Last synced │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Features

- **Match browser** — step through all 64 matches with paginated nav (`ui/nav`).
- **Match detail** (`ui/match`) — venue, local kickoff time, LIVE minute, block-digit score (`ui/bigtext`), team columns with true-color pixel flags (`ui/flags`), lineups with per-player goals / yellow / red counts, and a chronological event feed (goals, penalties, own goals, cards, substitutions; canceled/VAR events struck through).
- **Group tables** (`ui/group`) — MP / W / D / L / GF / GA / GD / Pts with team color chips, shown automatically when the selected match is a group-stage match.
- **Knockout bracket** (`ui/bracket`) — ASCII bracket from Round of 16 through the Final, shown automatically for knockout matches. Requires the full 64-match dataset.
- **Player stats** (`ui/playerstats`) — goals, yellow and red cards aggregated per team from match events.
- **Status bar** (`ui/statusbar`) — data-source name, refresh spinner, last-sync time, error state.
- **Auto-refresh** — data refetch is scheduled every 1 minute (`refreshInterval` in `ui/dashboard.go`).
- **Responsive guards** — centered loading / error / empty states, plus a minimum-size gate (160 cols × 50 rows).

## Requirements

- Go 1.27.1+ (see `go.mod`)
- A true-color terminal with UTF-8 and emoji support (flags, ■ cards, ⬤ pagination, 🏆)
- Minimum terminal size: **160 columns × 50 rows**

## Quick start

```powershell
# from the repo root
go mod tidy
go run .
```

Build a binary:

```powershell
go build -o worldcupdashboard .
.\worldcupdashboard
```

Run tests / vet:

```powershell
go test ./...
go vet ./...
```

## Controls

| Key | Action |
|---|---|
| `→`, `d`, `l`, `Space` | Next match |
| `←`, `a`, `h` | Previous match |
| `q`, `ctrl+c`, `esc` | Quit |

Selection is sticky: once you navigate manually (`matchIndexChanged`), auto-refreshes will not yank the cursor back to the live match unless your index goes out of range. Initial selection prefers the first `Live` match, then the first `Scheduled` match, otherwise the last match.

## Data source

The app ships with an embedded offline dataset — no network calls:

- `data/models.go` — `Match`, `GroupTable`, `Event`, `Player`, `TeamInfo`, `Stage` (`Group`, `1/8`, `1/4`, `1/2`, `3rd`, `Final`), `Status` (`Scheduled`, `Live`, `Finished`).
- `data/teams.go` — all 32 teams: FIFA code → name, group (A–H), two display colors.
- `data/local/local.go` — `local.Client` (`Name() == "local"`) serving `GroupTables()` and `SortedMatches()` from embedded JSON.
- `data/local/groups.go`, `data/local/matches.go` — the embedded group and match payloads.

### Use a different source (e.g. live API)

Implement the small interface in `ui/fetch.go` and pass it to `ui.NewDashboard`:

```go
type dataFetcher interface {
    GroupTables() ([]data.GroupTable, error)
    SortedMatches() ([]data.Match, error)
    Name() string
}

// main.go
dashboard := ui.NewDashboard(&local.Client{}) // swap this
```

`dataFetchCmd` fetches both payloads, indexes groups by letter, sorts/keeps match order, and derives `playerStatsByTeam` via `playerstats.PlayerStatsByTeam`.

## Architecture

Elm-style Bubble Tea app (`charm.land/bubbletea/v2`):

- `main.go` — creates `ui.NewDashboard(&local.Client{})`, runs `tea.NewProgram`, fullscreen via `View.AltScreen = true`.
- `ui/dashboard.go` — root model: `Init` (spinner tick + fetch), `Update` (fetch results, 1-min `intervalRefreshMsg`, `KeyMsg`, `spinner.TickMsg`, `WindowSizeMsg`), `View` (composes nav / match / help / group-or-bracket / statusbar with Lip Gloss).
- `ui/fetch.go` — fetch commands and messages (`dataFetchMsg`, `dataFetchErrMsg`).
- `ui/keymap.go` — `help.KeyMap` implementation for the help bubble.

| Path | Purpose |
|---|---|
| `data/` | Domain models + team registry |
| `data/local/` | Embedded offline client |
| `ui/` | Root dashboard model, fetch logic, keymap |
| `ui/nav/` | Paginated match strip + pagination dots |
| `ui/match/` | Match detail view (score, events, lineups) |
| `ui/group/` | Group standings table |
| `ui/bracket/` | Knockout bracket renderer |
| `ui/bigtext/` | Block-digit font (`0-9`, `-`, `?`) for the scoreline |
| `ui/flags/` | Per-country 14×25 true-color pixel flags |
| `ui/playerstats/` | Goals/cards aggregation from events |
| `ui/statusbar/` | Bottom status bar |

## Tech stack

- `charm.land/bubbletea/v2` — TUI runtime (Model / Update / View, commands, alt-screen)
- `charm.land/bubbles/v2` — `spinner`, `help`, `key` bubbles
- `charm.land/lipgloss/v2` — layout, colors, borders, light/dark adaptation via `LightDark(HasDarkBackground(...))`

## Configuration tweaks

- Refresh cadence: `refreshInterval` in `ui/dashboard.go` (default `1 * time.Minute`).
- Size gate: `minWidth` / `minHeight` in `dashboard.View()` (default 160 × 50).
- Spinner style: `spinner.Globe` in `ui.NewDashboard`.
- Score font: glyphs in `ui/bigtext/font.go` (`chars` + `charHeight` must stay in sync or `NewBigText` panics).

## Troubleshooting

- **"Need at least 160 columns and 50 rows"** — enlarge the terminal (or shrink the font) and the dashboard will render on the next resize event.
- **"Initializing..." stuck** — the app hasn't received a `WindowSizeMsg` yet; focus/resize the terminal once.
- **Flags or block digits look broken** — use a UTF-8, true-color font with emoji coverage (Windows Terminal, WezTerm, Ghostty, Kitty, iTerm2 all work well).
- **Bracket says "want 64 matches"** — the bracket view only renders with the complete 64-match tournament dataset.

## Development notes

- `ui/nav` has unit tests (`nav_test.go`); other view packages are pure string renderers and easy to test the same way.
- Event-type string constants live in `data/models.go` — note the existing `"Yello Card"` spelling is load-bearing (datasets use it), so keep it in sync with your data.

## License

No license file is currently included in this repository. Add one (e.g. MIT) before distributing binaries.
