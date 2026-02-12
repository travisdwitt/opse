# OPSE — One Page Solo Engine (Now with more _stuff_)

A TUI implementation of [One Page Solo Engine](https://inflatablestudios.itch.io/one-page-solo-engine) (v1.6) by Karl Hendricks.
Plus some additional generators and tools I like to use when running my own solo RPG adventures.

---

## Features

- **OPSE v1.6** — every oracle, generator, focus table, and GM move from the original rules
- **Adventure journal** — write narrative text and engine results into a single scrollable log
- **Character voices** — attribute log entries to named characters with `/char`
- **Markdown output** — journals save as `.md` files with timestamps, readable anywhere
- **Save and resume** — reopen any adventure and pick up where you left off
- **Saved rolls** — persistent dice roll templates organized into folders
- **Autocomplete** — fuzzy-matching suggestions as you type
- **Built-in help** — 9-page reference covering rules, generators, and commands

---

## Quick Start

### Prerequisites

- [Go](https://go.dev/dl/) 1.22 or later

### Build and Run

```sh
git clone https://github.com/your-username/opse.git
cd opse
go build -o opse .
./opse
```

Or run directly without building a binary:

```sh
go run .
```

## Interface

### Layout

```
┌──────────────────────────────────────────────────────────────┐
│ OPSE — My Adventure                                          │
├────────────┬─────────────────────────────────────────────────┤
│  ORACLE    │                                                 │
│  1 Likely  │  14:32  Engine                                  │
│  2 Even    │  ╭─────────────────────────────────────╮        │
│  3 Unlikely│  │ Oracle: Yes/No (Even)               │        │
│  4 How     │  │  Answer: 5 → Yes                    │        │
│            │  │  Modifier: 6 → and...               │        │
│  FOCUS     │  │                                     │        │
│  5 Action  │  │  Result: Yes, and...                │        │
│  6 Detail  │  ╰─────────────────────────────────────╯        │
│  7 Topic   │                                                 │
│            │  14:33  User                                    │
│  SCENE     │  The door swings open to reveal a vast cavern.  │
│  8 Scene   │                                                 │
│  9 Random  ├─────────────────────────────────────────────────┤
│            │                                                 │
│  GM MOVES  │  > _                                            │
│  0 Pacing  │                                                 │
│  - Failure ├─────────────────────────────────────────────────┤
│            │ Tab: switch | 1-9: generators | ?: help         │
└────────────┴─────────────────────────────────────────────────┘
```

### Navigation

| Key | Action |
|---|---|
| `Tab` | Cycle focus: Input → Sidebar → Log |
| `Esc` | Return focus to Input |
| `j` / `k` or `↑` / `↓` | Scroll or navigate |
| `?` | Toggle help (from Sidebar or Log) |
| `Ctrl+S` | Save journal |
| `Ctrl+R` | Open saved rolls manager |
| `Ctrl+Q` | Save and quit |

### Shortcuts

| Key | Generator |
|---|---|
| `1` | Oracle: Yes/No (Likely) |
| `2` | Oracle: Yes/No (Even) |
| `3` | Oracle: Yes/No (Unlikely) |
| `4` | Oracle: How |
| `5` | Action Focus |
| `6` | Detail Focus |
| `7` | Topic Focus |
| `8` | Set the Scene |
| `9` | Random Event |
| `0` | Pacing Move |
| `-` | Failure Move |
| `=` | Generic Generator |

---

## Slash Commands

Type these in the Input area and press Enter. Autocomplete suggestions appear as you type — press `Tab` to complete.

### Dice & Randomizers

| Command | Description | Examples |
|---|---|---|
| `/roll NdS` | Roll N dice with S sides | `/roll 2d6`, `/roll 4d6k3` |
| `/roll NdS+M` | Roll with modifier | `/roll 1d20+5`, `/r 2d8-1` |
| `/r` | Shorthand for `/roll` | `/r 3d6` |
| `/flip [N]` | Flip coins | `/flip`, `/flip 5` |
| `/f` | Shorthand for `/flip` | `/f 3` |
| `/draw [N]` | Draw cards from the utility deck | `/draw`, `/draw 3` |
| `/card` | Same as `/draw` | `/card 2` |
| `/shuffle` | Reshuffle the utility deck | `/shuffle` |

### World Building

| Command | Description |
|---|---|
| `/dir [N]` | Random compass direction (4, 8, or 16 point) |
| `/weather` | Random weather (condition, temperature, wind) |
| `/color` | Random color |
| `/sound [CATEGORY]` | Random sound effect |
| `/scene` | Same as Set the Scene (`8`) |

Sound categories: `nature`, `urban`, `combat`, `social`, `mechanical`, `animal`, `weather`, `supernatural`, `domestic`, `musical`

### Character Voice

| Command | Description |
|---|---|
| `/char NAME TEXT` | Log entry attributed to a named character |
| `/char "FULL NAME" TEXT` | Use quotes for multi-word names |

Examples:
```
/char Elara I search the room carefully.
/char "Captain Pike" Set course for the nebula.
```

Entries appear in the log and saved markdown with the character's name instead of "User":
```
14:32  Elara
I search the room carefully.
```

---

## Generators

### The Oracle

**Yes/No** (`1`, `2`, `3`) — Ask a question, pick a likelihood:
- **Likely** — Yes on 3+ (d6)
- **Even** — Yes on 4+ (d6)
- **Unlikely** — Yes on 5+ (d6)

**How** (`4`) — For "how much?" or "how strong?" questions. Returns a scale from "Surprisingly lacking" to "Extraordinary."

### Focus Tables

Draw a card for open-ended inspiration. The rank gives the table entry, the suit adds a domain for interpretation.

- **Action Focus** (`5`) — What does it do?
- **Detail Focus** (`6`) — What kind of thing is it?
- **Topic Focus** (`7`) — What is this about?

### Suit Domains

| Suit | Domain | Examples |
|---|---|---|
| Clubs | Physical | Strength, damage, objects, terrain |
| Diamonds | Technical | Plans, devices, intelligence, craft |
| Spades | Mystical | Destiny, spells, omens, power |
| Hearts | Social | Loyalty, love, betrayal, reputation |

### Scene Management

- **Set the Scene** (`8`) — Generates a complication for each new situation. May trigger an altered scene with cascading effects.
- **Random Event** (`9`) — Draws two cards (Action Focus + Topic Focus) for an unexpected twist.

### GM Moves

- **Pacing Move** (`0`) — Use when there's a lull. Foreshadow trouble, reveal details, advance threats.
- **Failure Move** (`-`) — Use when PCs fail a check. Cause harm, present choices, reveal unwelcome truths.

### Complex Generators

Access these from the Sidebar:

- **Generic** (`=`) — Towns, factions, items, monsters — anything
- **Plot Hook** — Quest objective + adversary + reward
- **NPC** — Identity, goal, notable feature, attitude, conversation topic
- **Dungeon Theme** — Look and purpose of a dangerous location
- **Dungeon Room** — Location, encounter, object, exits
- **Hex** — Terrain, contents, features, events for wilderness exploration

---

## Saved Rolls

Press `Ctrl+R` to open the saved rolls manager. Create reusable rolls for whatever system you're playing.

| Key | Action |
|---|---|
| `j` / `k` | Navigate |
| `Enter` | Execute selected roll |
| `n` | Create new roll |
| `f` | Create folder |
| `d` | Delete selected |
| `←` / `→` | Collapse / expand folder |
| `Esc` | Close |

Saved rolls persist across sessions in `~/.config/opse/saved_rolls.json`.

---

## Journal Format

Adventures save as Markdown files in the current directory. The format is designed to look clean when viewed in any Markdown renderer:

```markdown
# The Fall of Blackspire

*Started: 2026-02-10*

---

*14:32 — User*

The party approaches the crumbling tower at dusk.

*14:33 — Engine*

> **Oracle (Yes/No, Even):** Yes, and...

*14:33 — User*

The gates stand open. Something drove the guards away.

*14:34 — Engine*

> **Set the Scene**
> - **Complication:** Hostile forces oppose you
> - **Altered Scene:** Unexpected NPCs are present

*14:35 — Elara*

We should proceed with caution.
```

---

### Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Bubbles](https://github.com/charmbracelet/bubbles)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)

---

## Credits

### One Page Solo Engine

**Created by Karl Hendricks — [Inflatable Studios](https://inflatablestudios.itch.io/one-page-solo-engine)**
One Page Solo Engine v1.6 is licensed under [CC-BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/). The oracle system, focus tables, generators, GM moves, and all game mechanics implemented in this application are based on the original OPSE rules.

---

## License

This software is released under the [MIT License](LICENSE).
