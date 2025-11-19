# OPSE
An app for managing solo role-playing game logs using the [One Page Solo Engine](https://inflatablestudios.itch.io/one-page-solo-engine) rules.

## Installation
   ```bash
   go mod tidy
   ```
   ```bash
   go build -o opse
   ```
   ```bash
   ./opse
   ```

## Usage

### Creating a New Log
1. "New Game Log" from the main menu
2. Enter a title for your game log
3. Press Enter to create

### Loading a Log
1. "Load Game Log" from the main menu
2. Select a log file
3. Press Enter to load

### Log View
When viewing a log, you can:

- **n**: Add a new text entry
- **l**: Oracle (Yes/No - Likely)
- **e**: Oracle (Yes/No - Even)
- **u**: Oracle (Yes/No - Unlikely)
- **h**: Use Oracle (How)
- **1**: Generate Scene Complication
- **2**: Generate Altered Scene
- **3**: Generate Pacing Move
- **4**: Generate Failure Move
- **5**: Generate Random Event
- **6**: Generate Plot Hook
- **7**: Generate NPC
- **8**: Generate Generic Generator
- **9**: Generate Dungeon Crawler
- **0**: Generate Hex Crawler
- **Esc**: Open Menu

## File Format
Game logs are saved as YAML files in the current directory.

## License
This application implements the One Page Solo Engine rules, which are licensed under CC-BY-SA 4.0.
The app is licensed MIT - if you want to make it better or change it up to suit your solo sessions
better, please do (and send me a link so I can play too)

## Credits
One Page Solo Engine was created by Karl Hendricks.

