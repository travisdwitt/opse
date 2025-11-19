package main

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateMainMenu state = iota
	stateNewLog
	stateLoadLog
	stateLogView
	stateAddEntry
	stateOracleYesNo
	stateOracleHow
	stateSceneComplication
	stateAlteredScene
	statePacingMove
	stateFailureMove
	stateRandomEvent
	statePlotHook
	stateNPCGenerator
	stateGenericGenerator
	stateDungeonCrawler
	stateHexCrawler
	stateViewGeneratorResult
	stateViewLicense
	stateViewOPSERawText
	stateViewOPSERawTextHeadings
)

type model struct {
	state               state
	currentLog          *GameLog
	logFilename         string
	menuIndex           int
	textInput           string
	generatorResult     string
	errorMsg            string
	likelihood          string
	availableLogs       []string
	logIndex            int
	scrollOffset        int
	width               int
	height              int
	headings            []Heading
	headingIndex        int
	headingsSidebarOpen bool
	headingSearch       string
	filteredHeadings    []Heading
	headingScrollOffset int
	logMenuOpen         bool
	logMenuIndex        int
	selectedWindow      string // "log" or "prompt"
	promptScrollOffset  int
	cursorPos           int // Cursor position in textInput
}

type Heading struct {
	Title string
	Line  int
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("3")).
			Padding(1, 2)

	menuItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	selectedMenuItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("3")).
				Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			PaddingTop(1)

	resultStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")). // Dark gray
			Padding(1, 2).
			Margin(1, 0)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	borderStyle = lipgloss.NewStyle().
			Padding(0)

	centeredWindowStyle = lipgloss.NewStyle().
				Width(60).
				Align(lipgloss.Left).
				Padding(1, 2)

	menuTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("3")).
			Align(lipgloss.Center).
			PaddingBottom(1)

	sidebarTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("21")). // Blue
				Align(lipgloss.Center).
				PaddingBottom(1)

	sidebarHeadingStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("250")). // Light gray
				Align(lipgloss.Center).
				PaddingBottom(1)

	sidebarShortcutStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("250")). // Light gray
				Bold(true)

	sidebarTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")) // Dark gray

	promptTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3")) // Yellow

	promptSymbolStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("21")) // Blue

	promptHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")) // Dark gray

	licensePopupStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")). // Dark gray
				Padding(1, 2).
				Width(98).
				Height(20)

	rawTextPopupStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")). // Dark gray
				Padding(1, 2).
				Width(98).
				Height(20)

	logViewStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")). // Dark gray
			Padding(1, 2)

	logViewStyleSelected = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("250")). // Light gray
				Padding(1, 2)

	logMenuPopupStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")). // Dark gray
				Padding(1, 2).
				Width(40).
				Height(8)

	logSidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")). // Dark gray
			Padding(1, 2)

	promptWindowStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")). // Dark gray
				Padding(1, 2).
				Height(5)

	promptWindowStyleSelected = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("250")). // Light gray
					Padding(1, 2).
					Height(5)

	resultPopupStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")). // Dark gray
				Padding(1, 2).
				Width(60).
				Height(12)

	oraclePopupStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")). // Dark gray
				Padding(1, 2).
				Width(40).
				Height(10)

	// Log entry color styles
	timestampStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6")) // Cyan

	generatorTypeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("2")). // Green
				Bold(true)

	oracleTypeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")). // Magenta
			Bold(true)

	userTypeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")). // Light gray
			Bold(true)

	commandNameStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")) // Dark gray
)

const opseRawText = `ONE PAGE SOLO ENGINE

Version: 1.6	CC-BY-SA 4.0

A minimal, all-in-one toolkit to play your favorite tabletop RPGs without a GM.

HOW TO PLAY

1	Create one or more characters using your chosen game system.

2	Roll a starting PLOT HOOK and a RANDOM EVENT, then SET THE SCENE.

3	Start asking the ORACLE questions.  Interpret the answers in context.

4	Play the game to overcome the challenges of the scene.

5	Use GM MOVES to move the action.

6	SET THE SCENE for the next thing you want your character to do.

USING PLAYING CARDS

This system uses a deck of playing cards to inspire answers.  Look up the rank in the appropriate table and combine with the SUIT DOMAIN below to determine the answer.  When you draw a Joker, shuffle the deck and add a RANDOM EVENT.

SUIT DOMAIN:

Clubs – Physical (appearance, existence)

Diamonds – Technical (mental, operation)

Spades – Mystical (meaning, capability)

Hearts – Social (personal, connection)

OPTIONAL: USE ONLY CARDS

When you would roll a d6, draw a card and use the rank divided by 2 (round down).  Discard Aces.

OPTIONAL: USE ONLY DICE

When you would draw a card, roll a d12 for the rank and a d4 for the suit.  On a 12, flip a coin to see if you use the Q or K.

SET THE SCENE

Describe where your character is and what they are trying to accomplish, then roll or choose a SCENE COMPLICATION.

SCENE COMPLICATION (D6):

1	Hostile forces oppose you

2	An obstacle blocks your way

3	Wouldn't it suck if…

4	An NPC acts suddenly

5	All is not as is seems

6	Things actually go as planned

Roll 1d6, on a 5+, it is an ALTERED SCENE.

ALTERED SCENE (D6):

1	A major detail of the scene is enhanced or somehow worse

2	The environment is different

3	Unexpected NPCs are present

4	Add a SCENE COMPLICATION

5	Add a PACING MOVE

6	Add a RANDOM EVENT

ORACLE (YES/NO)

When you need to ask a simple question, choose the likelihood and roll 2d6.

Answer (d6)

Likely	  Yes on 3+

Even	  Yes on 4+

Unlikely  Yes on 5+

Mod (d6)

1	but…

2-5	

6	and…

ORACLE (HOW)

When you need to know how big, good, strong, numerous, etc. something is.

Answer (D6):

1	Surprisingly lacking

2	Less than expected

3-4	About average

5	More than expected

6	Extraordinary

GM MOVES

When you need to advance the action, roll on the tables below and describe the results as the GM normally would.

Use a PACING MOVE when there is a lull in the action, or you think "what now?"  Use a FAILURE MOVE to move things forward when the PCs fail a check.

PACING MOVES (D6):

1	Foreshadow Trouble

2	Reveal a New Detail

3	An NPC Takes Action

4	Advance a Threat

5	Advance a Plot

6	Add a RANDOM EVENT to the scene

FAILURE MOVES (D6):

1	Cause Harm

2	Put Someone in a Spot

3	Offer a Choice

4	Advance a Threat

5	Reveal an Unwelcome Truth

6	Foreshadow Trouble

RANDOM EVENT

When you need to create a random event, draw the following.

What happens: ACTION FOCUS

Involving: TOPIC FOCUS

COMPLEX QUESTIONS

When you need to ask an open-ended question, try to find the most appropriate ORACLE (FOCUS) to use.  If the question is not sufficiently answered, add results from a second ORACLE (FOCUS).

ORACLE (FOCUS)

When you have a broad question or need to know details about something, draw on one of the tables below.  Remember to apply the SUIT DOMAIN when interpreting the result.

ACTION FOCUS (CARD):

What does it do?

2	Seek		9	Command

3	Oppose		T	Take

4	Communicate	J	Protect

5	Move		Q	Assist

6	Harm		K	Transform

7	Create		A	Deceive

8	Reveal			

DETAIL FOCUS (CARD):

What kind of thing is it?

2	Small		9	Unsavory

3	Large		T	Specialized

4	Old		J	Unexpected

5	New		Q	Exotic

6	Mundane		K	Dignified

7	Simple		A	Unique

8	Complex			

TOPIC FOCUS (CARD):

What is this about?

2	Current Need	9	Rumors

3	Allies		T	A Plot Arc

4	Community	J	Recent Events

5	History		Q	Equipment

6	Future Plans	K	A Faction

7	Enemies		A	The PCs

8	Knowledge			

				

 

ONE PAGE GENERATORS

Version: 1.6	CC-BY-SA 4.0

Content-neutral generators to aid your GM-less adventures in any setting.

GENERIC GENERATOR

Use this to generate towns, spaceships, factions, magic items, taverns, monsters, or anything else you can think of.

What it does: ACTION FOCUS

How it looks: DETAIL FOCUS

How significant: ORACLE (HOW)

PLOT HOOK GENERATOR

Use this to generate plot hooks, quests, or missions for the PCs to follow.

OBJECTIVE (D6):

1	Eliminate a threat

2	Learn the truth

3	Recover something valuable

4	Escort or deliver to safety

5	Restore something broken

6	Save an ally in peril

ADVERSARIES (D6):

1	A powerful organization

2	Outlaws

3	Guardians

4	Local inhabitants

5	Enemy horde or force

6	A new or recurring villain

REWARDS (D6):

1	Money or valuables

2	Money or valuables

3	Knowledge and secrets

4	Support of an ally

5	Advance a plot arc

6	A unique item of power

NPC GENERATOR

Use this to generate NPCs that may be encountered while playing.

IDENTITY (CARD):

2	Outlaw		9	Entertainer

3	Drifter		T	Adherent

4	Tradesman	J	Leader

5	Commoner	Q	Mystic

6	Soldier	K	Adventurer

7	Merchant	A	Lord

8	Specialist			

GOAL (CARD):

2	Obtain		9	Enrich Self

3	Learn		T	Avenge

4	Harm		J	Fulfill Duty

5	Restore		Q	Escape

6	Find		K	Create

7	Travel		A	Serve

8	Protect			

NOTABLE FEATURE (D6):

1	Unremarkable

2	Notable nature

3	Obvious physical trait

4	Quirk or mannerism

5	Unusual equipment

6	Unexpected age or origin

Draw a DETAIL FOCUS for the description of the notable feature.

CURRENT SITUATION

Attitude to PCs: ORACLE (HOW)

Conversation: TOPIC FOCUS

DUNGEON CRAWLER

Use this when exploring a dangerous location like a typical dungeon.

DUNGEON THEME:

How it looks: DETAIL FOCUS

How it is used: ACTION FOCUS

The first area always has 3 exits.  As you explore, roll once on each table below to create the new area.

LOCATION (D6):

1	Typical area

2	Transitional area

3	Living area or meeting place

4	Working or utility area

5	Area with a special feature

6	Location for a specialized purpose

ENCOUNTER (D6):

1-2	None

3-4	Hostile enemies

5	An obstacle blocks the way

6	Unique NPC or adversary

OBJECT (D6):

1-2	Nothing, or mundane objects

3	An interesting item or clue

4	A useful tool, key, or device

5	Something valuable

6	Rare or special item

TOTAL EXITS (D6):

1-2	Dead end

3-4	1 additional exit

5-6	2 additional exits

HEX CRAWLER

Use this to generate maps of larger areas.  Whenever the characters enter a hex, generate the TERRAIN and CONTENTS of all surrounding hexes, then roll an EVENT for the current hex.

REGION

Each hex is part of a region.  Define the three terrain types for the starting region (common, uncommon, and rare).  New regions might be discovered later.

TERRAIN (D6):

1-2	Same as current hex

3-4	Common terrain

5	Uncommon terrain

6	Rare terrain

CONTENTS (D6):

1-5	Nothing notable

6	Roll a FEATURE

FEATURES (D6):

1	Notable structure

2	Dangerous hazard

3	A settlement

4	Strange natural feature

5	New region (set new terrain types)

6	DUNGEON CRAWLER entrance

EVENT (D6):

1-4	None

5-6	RANDOM EVENT then SET THE SCENE

 

MORE INFORMATION

One Page Solo Engine was designed to be incredibly concise and minimalistic, but still have all the essential tools required to run a game without a GM.  The first two pages of this document are all that are required to play.  I understand some people might want to know a bit more, though, so here are some notes.

INTENDED AUDIENCE

This is really meant for people who are already familiar with RPGs and playing them solo.  Most of the tools assume you have already encountered similar concepts in other products.  If you're completely new to solo or GM-less gaming, check out some of the products in the Acknowledgements to get started.

DESIGN PHILOSOPHY

There are a great many excellent tools out there to run a solo RPG game.  I always found, however, that many of them were overly complicated.  You shouldn't have to read 15 pages of rules and make 10 dice rolls just to determine what the guards in a room are doing.

Also, many tools only provide part of what you need to actually play.  Some only answer questions, while others only provide narrative structure or generate random elements.  A complete oracle should do all these things. 

One Page Solo Engine was designed to provide every tool needed to run a solo game using any game system while using as few words as humanly possible.

GM MOVES

Though the GM Moves section is highly inspired by PbtA games, the system will work with virtually any tabletop RPG. The reason the PbtA framework was chosen is that it gamifies the role of the GM with discrete moves that can fit in a table.

How you use the GM Moves will depend on the game system you are playing.  If you're playing a PbtA game, it will be obvious when to use them because PbtA is built around the concept of partial success.  If you aren't, just use them when you want to move things along.

If you need more information on how to use the individual GM Moves, check out any PbtA game such as Dungeon World, Uncharted Worlds, or many others.

PACING MOVES

Pacing Moves should be used to fill in the gaps during those times when the players would normally look to the GM to see what happens next.  They represent the little prompting and extra details that a GM usually adds.  Try using one whenever you want to move the action forward.

FAILURE MOVES

Failure Moves represent setbacks or partial successes.  Maybe the roll failed, but the character still gets part of what they wanted or all of it with a cost.  These moves keep the action moving during failures and can be used in virtually any RPG system instead of just saying "no that failed".

NON-PBTA GAMES

When playing a Non-PbtA game, it is important to remember that not every failure should result in a GM Move. 

Sometimes the Spot check just fails because there was nothing there.  GM Moves should be used when a roll fails and there are consequences for failure, or the action needs to pick up.  Checking the room for secret doors?  Probably not.  Climbing a cliff in the rain to escape group of cultists?  Definitely.

THE POWER OF INTERPRETATION

Some solo RPG tools contain dozens of tables with hundreds of entries each.  The problem with these is that they are either thematically tuned to a certain genre of game, or they are so specific the results just don't make sense.

When using the One Page Solo Engine, remember that the answers are meant to inspire an idea that makes sense in the context of your game.  The answer should have meaning, not just be a random detail.  The result may be surprising, but it should always be logical.

Give all results meaning.  Embrace the unexpected.  Reject the nonsensical.

USING A DECK OF CARDS

Many people who try this system wonder why a deck of cards was chosen for randomization.  The reason is because a playing card carries more information than a die roll and the suits work well for applying a "domain" to the results.

I've seen this used to great success in systems like World vs Hero and decided to apply it to a generic solo engine.  Instead of having a huge table with every adjective you can pull from the dictionary, you have a smaller table with more general words and a domain that they can apply to.  This results in more interpretation and less guesswork about how "divinely slippery" could possibly apply to your current situation.

TIPS FOR BEST RESULTS

*	Ask mostly yes/no questions

*	Loose interpretations are okay

*	Always go with what's cool

*	If it doesn't make sense, try again

*	Use GM Moves to drive the action

*	Try group play with no GM, it's great

ACKNOWLEDGEMENTS

One Page Solo Engine was created by taking the things I liked from other solo tools, stripping them down to the bare bones, and then adding in a bit of the process I use for my own games.  It would not be possible without inspiration from:

*	Mythic (Tana Pigeon)

*	World vs Hero (John Fiore)

*	Conjecture, UNE (Zach Best)

*	Dungeon World (Koebel, LaTorra)

*	The Black Hack (David Black)

*	Maze Rats (Ben Milton)

*	The Lone Wolf Solo RPG community

ABOUT

Written by:

    Karl Hendricks – Inflatable Studios

    Email: support@inflatablestudios.dev

    Reddit: u/archon1024

Download at:

    https://inflatablestudios.itch.io/

License: 

    CC-BY-SA 4.0

    Free to use and adapt for your own works (even commercially) as long as you provide credit and share your changes.  See the link above for more details.`

// calculatePromptWrappedLines calculates how the prompt text wraps and returns
// the wrapped lines and the cursor's position in the wrapped lines
func (m model) calculatePromptWrappedLines() ([]string, int, int) {
	// Calculate available width (must match View calculation)
	sidebarWidth := 28
	spacing := 2
	estimatedLogWindowWidth := m.width - sidebarWidth - spacing - 4
	if estimatedLogWindowWidth < 50 {
		estimatedLogWindowWidth = 50
	}
	availableWidth := estimatedLogWindowWidth - 6
	if availableWidth < 10 {
		availableWidth = 10
	}

	// Build the prompt line with cursor
	var promptLine strings.Builder
	promptLine.WriteString("> ")
	if m.cursorPos >= len(m.textInput) {
		promptLine.WriteString(m.textInput)
		promptLine.WriteString("_")
	} else {
		beforeCursor := m.textInput[:m.cursorPos]
		afterCursor := m.textInput[m.cursorPos:]
		promptLine.WriteString(beforeCursor)
		promptLine.WriteString("_")
		promptLine.WriteString(afterCursor)
	}

	// Wrap the prompt line
	wrappedLine := lipgloss.NewStyle().
		Width(availableWidth).
		Render(promptLine.String())
	wrappedLineParts := strings.Split(wrappedLine, "\n")

	// Calculate cursor position in wrapped lines
	// We need to find which wrapped line contains the cursor
	cursorLine := 0
	cursorCol := 0
	promptPrefix := "> "
	prefixLen := len(promptPrefix)

	// Build text without cursor to calculate positions
	var textWithoutCursor strings.Builder
	textWithoutCursor.WriteString(promptPrefix)
	textWithoutCursor.WriteString(m.textInput)
	textWithoutCursorStr := textWithoutCursor.String()

	// Wrap text without cursor to find line breaks
	wrappedWithoutCursor := lipgloss.NewStyle().
		Width(availableWidth).
		Render(textWithoutCursorStr)
	wrappedWithoutCursorParts := strings.Split(wrappedWithoutCursor, "\n")

	// Find which line the cursor is on
	targetPos := prefixLen + m.cursorPos
	currentLineStart := 0
	for i, line := range wrappedWithoutCursorParts {
		lineDisplayLen := lipgloss.Width(line)
		if targetPos >= currentLineStart && targetPos < currentLineStart+lineDisplayLen {
			cursorLine = i
			cursorCol = targetPos - currentLineStart
			break
		}
		currentLineStart += lineDisplayLen
		// Account for line break
		if i < len(wrappedWithoutCursorParts)-1 {
			currentLineStart += 0 // No extra offset needed
		}
	}

	// If cursor is beyond all lines, put it on the last line
	if cursorLine >= len(wrappedLineParts) {
		cursorLine = len(wrappedLineParts) - 1
		if cursorLine < 0 {
			cursorLine = 0
		}
		cursorCol = lipgloss.Width(wrappedLineParts[cursorLine])
	}

	return wrappedLineParts, cursorLine, cursorCol
}

// moveCursorUpDown moves the cursor up or down through wrapped lines
func (m model) moveCursorUpDown(direction int, logWindowWidth int) (int, int) {
	// Calculate available width
	availableWidth := logWindowWidth - 6
	if availableWidth < 10 {
		availableWidth = 10
	}

	promptPrefix := "> "
	fullText := promptPrefix + m.textInput

	// Find current cursor line and column using the helper
	_, currentWrappedLine, currentCol := m.calculatePromptWrappedLines()

	// Wrap text without cursor to find line boundaries
	wrappedText := lipgloss.NewStyle().
		Width(availableWidth).
		Render(fullText)
	wrappedLines := strings.Split(wrappedText, "\n")

	// Move to adjacent line
	newWrappedLine := currentWrappedLine + direction
	if newWrappedLine < 0 {
		newWrappedLine = 0
	}
	if newWrappedLine >= len(wrappedLines) {
		newWrappedLine = len(wrappedLines) - 1
		if newWrappedLine < 0 {
			newWrappedLine = 0
		}
	}

	// Find character positions for line boundaries
	// Find start of new line
	newLineStartChar := len(promptPrefix)
	if newWrappedLine > 0 {
		// Binary search for the character position that starts this line
		low := len(promptPrefix)
		high := len(fullText)
		for low < high {
			mid := (low + high) / 2
			testPrefix := fullText[:mid]
			testWrapped := lipgloss.NewStyle().
				Width(availableWidth).
				Render(testPrefix)
			testLines := strings.Split(testWrapped, "\n")
			if len(testLines) > newWrappedLine+1 {
				high = mid
			} else {
				low = mid + 1
			}
		}
		newLineStartChar = low
	}

	// Find end of new line
	newLineEndChar := len(fullText)
	if newWrappedLine < len(wrappedLines)-1 {
		low := newLineStartChar
		high := len(fullText)
		for low < high {
			mid := (low + high) / 2
			testPrefix := fullText[:mid]
			testWrapped := lipgloss.NewStyle().
				Width(availableWidth).
				Render(testPrefix)
			testLines := strings.Split(testWrapped, "\n")
			if len(testLines) > newWrappedLine+1 {
				high = mid
			} else {
				low = mid + 1
			}
		}
		newLineEndChar = low
	}

	// Calculate target column (maintain column, clamp to line length)
	targetCol := currentCol
	newLineText := fullText[newLineStartChar:newLineEndChar]
	newLineDisplayWidth := lipgloss.Width(newLineText)
	if targetCol > newLineDisplayWidth {
		targetCol = newLineDisplayWidth
	}

	// Find character position at target column
	newCursorCharPos := newLineStartChar
	for testPos := newLineStartChar; testPos <= newLineEndChar; testPos++ {
		testPrefix := fullText[newLineStartChar:testPos]
		if lipgloss.Width(testPrefix) >= targetCol {
			newCursorCharPos = testPos
			break
		}
	}
	if newCursorCharPos > newLineEndChar {
		newCursorCharPos = newLineEndChar
	}

	// Convert to textInput cursor position
	newCursorPos := newCursorCharPos - len(promptPrefix)
	if newCursorPos < 0 {
		newCursorPos = 0
	}
	if newCursorPos > len(m.textInput) {
		newCursorPos = len(m.textInput)
	}

	// Calculate scroll to keep cursor visible
	promptHeight := 5
	contentHeight := promptHeight - 4
	maxVisibleLines := contentHeight - 1
	if maxVisibleLines < 1 {
		maxVisibleLines = 1
	}

	// Recalculate wrapped lines with new cursor position to find its line
	var newPromptLine strings.Builder
	newPromptLine.WriteString("> ")
	if newCursorPos >= len(m.textInput) {
		newPromptLine.WriteString(m.textInput)
		newPromptLine.WriteString("_")
	} else {
		newPromptLine.WriteString(m.textInput[:newCursorPos])
		newPromptLine.WriteString("_")
		newPromptLine.WriteString(m.textInput[newCursorPos:])
	}
	newWrappedPrompt := lipgloss.NewStyle().
		Width(availableWidth).
		Render(newPromptLine.String())
	newWrappedLines := strings.Split(newWrappedPrompt, "\n")

	// Find which line the cursor is on now
	cursorWrappedLine := 0
	newCursorCharPosInPrompt := len("> ") + newCursorPos
	charPos := 0
	for i := 0; i < len(newWrappedLines); i++ {
		// Find character count for this line
		lineStartChar := charPos
		for testChar := charPos; testChar <= len(m.textInput)+len("> ")+1; testChar++ {
			var testText strings.Builder
			testText.WriteString("> ")
			if newCursorPos < len(m.textInput) {
				testText.WriteString(m.textInput[:newCursorPos])
				testText.WriteString("_")
				testText.WriteString(m.textInput[newCursorPos:])
			} else {
				testText.WriteString(m.textInput)
				testText.WriteString("_")
			}
			testPrefix := testText.String()[:min(testChar, testText.Len())]
			testWrapped := lipgloss.NewStyle().
				Width(availableWidth).
				Render(testPrefix)
			testLines := strings.Split(testWrapped, "\n")
			if len(testLines) > i+1 {
				charPos = testChar
				break
			}
		}
		if newCursorCharPosInPrompt >= lineStartChar && newCursorCharPosInPrompt < charPos {
			cursorWrappedLine = i
			break
		}
	}

	// Add error message lines if any
	totalWrappedLines := len(newWrappedLines)
	if m.errorMsg != "" {
		errorWrapped := lipgloss.NewStyle().
			Width(availableWidth).
			Render(m.errorMsg)
		totalWrappedLines += len(strings.Split(errorWrapped, "\n"))
	}

	maxScroll := totalWrappedLines - maxVisibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}

	// Calculate scroll to keep cursor visible
	newScrollOffset := m.promptScrollOffset
	if cursorWrappedLine < newScrollOffset {
		newScrollOffset = cursorWrappedLine
	} else if cursorWrappedLine >= newScrollOffset+maxVisibleLines {
		newScrollOffset = cursorWrappedLine - maxVisibleLines + 1
	}

	if newScrollOffset < 0 {
		newScrollOffset = 0
	}
	if newScrollOffset > maxScroll {
		newScrollOffset = maxScroll
	}

	return newCursorPos, newScrollOffset
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// renderScrollBar renders a vertical scroll bar
func renderScrollBar(height int, scrollOffset int, totalLines int, visibleLines int) string {
	if totalLines <= visibleLines {
		// No scrolling needed, return empty scroll bar (spaces)
		var scrollBar strings.Builder
		for i := 0; i < height; i++ {
			scrollBar.WriteString(" ")
			if i < height-1 {
				scrollBar.WriteString("\n")
			}
		}
		return scrollBar.String()
	}

	maxScroll := totalLines - visibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}

	// Calculate scroll bar thumb position
	// thumbPosition is 0 at top, height-1 at bottom
	thumbPosition := 0
	if maxScroll > 0 {
		thumbPosition = (scrollOffset * (height - 1)) / maxScroll
		if thumbPosition >= height {
			thumbPosition = height - 1
		}
	}

	// Calculate thumb size (proportional to visible area)
	thumbSize := 1
	if height > 3 {
		// Make thumb at least 1 character, but proportional to visible/total ratio
		ratio := float64(visibleLines) / float64(totalLines)
		thumbSize = int(float64(height) * ratio)
		if thumbSize < 1 {
			thumbSize = 1
		}
		if thumbSize > height {
			thumbSize = height
		}
	}

	// Build scroll bar
	var scrollBar strings.Builder
	for i := 0; i < height; i++ {
		if i >= thumbPosition && i < thumbPosition+thumbSize {
			// Thumb (filled)
			scrollBar.WriteString("█")
		} else {
			// Track (empty)
			scrollBar.WriteString("░")
		}
		if i < height-1 {
			scrollBar.WriteString("\n")
		}
	}

	return scrollBar.String()
}

// calculateLogMaxScroll calculates the maximum scroll offset for the log view
func (m model) calculateLogMaxScroll() int {
	// Calculate max scroll based on log window dimensions (must match View calculation)
	sidebarHeight := m.height - 4
	if sidebarHeight < 15 {
		sidebarHeight = 15
	}
	promptHeight := 5
	logHeight := sidebarHeight - promptHeight - 2
	if logHeight < 10 {
		logHeight = 10
	}
	contentHeight := logHeight - 2       // Account for border only (top + bottom)
	maxVisibleLines := contentHeight - 1 // Allow 1 row margin from bottom
	if maxVisibleLines < 1 {
		maxVisibleLines = 1
	}

	// Calculate total lines in log content (must match View calculation)
	totalLines := 2 // Title + blank line
	for _, entry := range m.currentLog.Entries {
		entryText := formatLogEntryColored(entry)
		entryLines := strings.Split(entryText, "\n")
		totalLines += len(entryLines) + 1 // Add spacing between entries
	}
	if len(m.currentLog.Entries) == 0 {
		totalLines = 3 // Title + blank + message
	}

	maxScroll := totalLines - maxVisibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}
	return maxScroll
}

// formatLogEntryColored formats a log entry with colors
func formatLogEntryColored(entry LogEntry) string {
	// Color the timestamp
	timestamp := timestampStyle.Render(entry.Timestamp.Format("2006-01-02 15:04:05"))

	// Color the type based on entry type
	var typeColored string
	switch entry.Type {
	case "generator":
		typeColored = generatorTypeStyle.Render(entry.Type)
	case "oracle":
		typeColored = oracleTypeStyle.Render(entry.Type)
	case "user":
		typeColored = userTypeStyle.Render(entry.Type)
	default:
		typeColored = entry.Type
	}

	// Parse content to find command/generator name and color it
	// Format is typically: "Command Name: result" or "Oracle (Type): result"
	content := entry.Content
	colonIndex := strings.Index(content, ":")
	if colonIndex > 0 {
		// Extract the command/generator name (before the colon)
		commandName := strings.TrimSpace(content[:colonIndex])
		// Check if it's an Oracle with parentheses
		if strings.HasPrefix(commandName, "Oracle") {
			// For Oracle entries, color the whole "Oracle (Type)" part
			commandNameColored := commandNameStyle.Render(commandName)
			restOfContent := content[colonIndex:]
			return fmt.Sprintf("[%s] %s: %s%s", timestamp, typeColored, commandNameColored, restOfContent)
		} else {
			// For other entries, color just the command name
			commandNameColored := commandNameStyle.Render(commandName)
			restOfContent := content[colonIndex:]
			return fmt.Sprintf("[%s] %s: %s%s", timestamp, typeColored, commandNameColored, restOfContent)
		}
	}

	// If no colon found, just return with colored timestamp and type
	return fmt.Sprintf("[%s] %s: %s", timestamp, typeColored, content)
}

func (m model) Init() tea.Cmd {
	// Window size will be sent automatically by bubbletea
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.state == stateMainMenu {
				return m, tea.Quit
			}
			if m.currentLog != nil {
				m.state = stateLogView
				m.textInput = ""
				m.cursorPos = 0
				m.logMenuOpen = false
				m.selectedWindow = "prompt"
			} else {
				m.state = stateMainMenu
			}
			return m, nil

		case "esc":
			if m.state == stateMainMenu {
				return m, tea.Quit
			}
			// Let stateLogView handle Esc for popup menu
			if m.state == stateLogView {
				// Don't handle here, let the state-specific handler deal with it
				break
			}
			if m.currentLog != nil {
				m.state = stateLogView
				m.textInput = ""
				m.cursorPos = 0
				m.logMenuOpen = false
				m.selectedWindow = "prompt"
			} else {
				m.state = stateMainMenu
			}
			return m, nil
		}

		switch m.state {
		case stateMainMenu:
			return m.updateMainMenu(msg)
		case stateNewLog:
			return m.updateNewLog(msg)
		case stateLoadLog:
			return m.updateLoadLog(msg)
		case stateLogView:
			return m.updateLogView(msg)
		case stateAddEntry:
			return m.updateAddEntry(msg)
		case stateOracleYesNo:
			return m.updateOracleYesNo(msg)
		case stateViewGeneratorResult:
			return m.updateViewGeneratorResult(msg)
		case stateViewLicense:
			return m.updateViewLicense(msg)
		case stateViewOPSERawText:
			return m.updateViewOPSERawText(msg)
		}
	}

	return m, nil
}

func (m model) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.menuIndex > 0 {
				m.menuIndex--
			}
		case "down", "j":
			if m.menuIndex < 4 {
				m.menuIndex++
			}
		case "enter":
			switch m.menuIndex {
			case 0:
				m.state = stateNewLog
				m.textInput = ""
				m.errorMsg = ""
			case 1:
				m.state = stateLoadLog
				m.errorMsg = ""
				logs, err := ListLogs()
				if err != nil {
					m.errorMsg = fmt.Sprintf("Error listing logs: %v", err)
				} else {
					m.availableLogs = logs
					m.logIndex = 0
				}
			case 2:
				m.state = stateViewLicense
				m.scrollOffset = 0 // Reset scroll when opening license
				m.errorMsg = ""
			case 3:
				m.state = stateViewOPSERawText
				m.scrollOffset = 0
				m.headings = parseHeadings(opseRawText)
				m.headingsSidebarOpen = false
				m.headingSearch = ""
				m.headingScrollOffset = 0
				m.filteredHeadings = m.headings
				m.errorMsg = ""
			case 4:
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) updateNewLog(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if strings.TrimSpace(m.textInput) != "" {
				now := time.Now()
				m.currentLog = &GameLog{
					Title:     strings.TrimSpace(m.textInput),
					CreatedAt: now,
					Entries:   []LogEntry{},
				}
				m.logFilename = strings.ToLower(strings.ReplaceAll(m.currentLog.Title, " ", "_")) + ".yaml"
				m.state = stateLogView
				m.textInput = ""
				m.cursorPos = 0
				m.logMenuOpen = false
				m.selectedWindow = "prompt"
			}
		case "backspace":
			if len(m.textInput) > 0 {
				m.textInput = m.textInput[:len(m.textInput)-1]
			}
		default:
			if len(msg.Runes) > 0 {
				m.textInput += string(msg.Runes)
			}
		}
	}
	return m, nil
}

func (m model) updateLoadLog(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.logIndex > 0 {
				m.logIndex--
			}
		case "down", "j":
			if m.logIndex < len(m.availableLogs)-1 {
				m.logIndex++
			}
		case "enter":
			if len(m.availableLogs) > 0 && m.logIndex < len(m.availableLogs) {
				log, err := LoadLog(m.availableLogs[m.logIndex])
				if err != nil {
					m.errorMsg = fmt.Sprintf("Error loading log: %v", err)
				} else {
					m.currentLog = log
					m.logFilename = m.availableLogs[m.logIndex]
					m.state = stateLogView
					m.textInput = ""
					m.cursorPos = 0
					m.logMenuOpen = false
					m.errorMsg = ""
					m.selectedWindow = "prompt"
				}
			}
		}
	}
	return m, nil
}

func (m model) updateLogView(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Initialize selected window to prompt if not set
	if m.selectedWindow == "" {
		m.selectedWindow = "prompt"
	}

	// Clamp scrollOffset to valid bounds at the start of each update
	if m.currentLog != nil {
		maxScroll := m.calculateLogMaxScroll()
		if m.scrollOffset > maxScroll {
			m.scrollOffset = maxScroll
		}
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
	}

	// Handle popup menu input
	if m.logMenuOpen {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.logMenuOpen = false
			case "up", "k":
				if m.logMenuIndex > 0 {
					m.logMenuIndex--
				}
			case "down", "j":
				if m.logMenuIndex < 4 {
					m.logMenuIndex++
				}
			case "enter":
				switch m.logMenuIndex {
				case 0: // Save
					err := SaveLog(m.currentLog, m.logFilename)
					if err != nil {
						m.errorMsg = fmt.Sprintf("Error saving: %v", err)
					} else {
						m.errorMsg = "Log saved successfully!"
					}
					m.logMenuOpen = false
				case 1: // Load
					m.logMenuOpen = false
					m.state = stateLoadLog
					m.errorMsg = ""
					logs, err := ListLogs()
					if err != nil {
						m.errorMsg = fmt.Sprintf("Error listing logs: %v", err)
					} else {
						m.availableLogs = logs
						m.logIndex = 0
					}
				case 2: // View License
					m.logMenuOpen = false
					m.state = stateViewLicense
					m.scrollOffset = 0
					m.errorMsg = ""
				case 3: // View Raw Text
					m.logMenuOpen = false
					m.state = stateViewOPSERawText
					m.scrollOffset = 0
					m.headings = parseHeadings(opseRawText)
					m.headingsSidebarOpen = false
					m.headingSearch = ""
					m.filteredHeadings = m.headings
					m.headingIndex = 0
					m.errorMsg = ""
				case 4: // Exit OPSE
					return m, tea.Quit
				}
			}
		}
		return m, nil
	}

	// Normal log view input handling
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Parse command or add entry
			cmd := strings.TrimSpace(m.textInput)
			m.textInput = ""
			m.cursorPos = 0

			if cmd == "" {
				return m, nil
			}

			// Check for single character commands
			if len(cmd) == 1 {
				switch cmd {
				case "n":
					// Add entry mode - already handled by text input
					return m, nil
				case "l":
					yes, result := OracleYesNo("Likely")
					if yes {
						m.generatorResult = fmt.Sprintf("Oracle (Yes/No - Likely): Yes - %s", result)
					} else {
						m.generatorResult = fmt.Sprintf("Oracle (Yes/No - Likely): No - %s", result)
					}
					m.currentLog.AddEntry(m.generatorResult, "oracle")
					m.errorMsg = ""
					// Auto-scroll to bottom
					maxScroll := m.calculateLogMaxScroll()
					m.scrollOffset = maxScroll
					return m, nil
				case "e":
					yes, result := OracleYesNo("Even")
					if yes {
						m.generatorResult = fmt.Sprintf("Oracle (Yes/No - Even): Yes - %s", result)
					} else {
						m.generatorResult = fmt.Sprintf("Oracle (Yes/No - Even): No - %s", result)
					}
					m.currentLog.AddEntry(m.generatorResult, "oracle")
					m.errorMsg = ""
					// Auto-scroll to bottom
					maxScroll := m.calculateLogMaxScroll()
					m.scrollOffset = maxScroll
					return m, nil
				case "u":
					yes, result := OracleYesNo("Unlikely")
					if yes {
						m.generatorResult = fmt.Sprintf("Oracle (Yes/No - Unlikely): Yes - %s", result)
					} else {
						m.generatorResult = fmt.Sprintf("Oracle (Yes/No - Unlikely): No - %s", result)
					}
					m.currentLog.AddEntry(m.generatorResult, "oracle")
					m.errorMsg = ""
					m.scrollOffset = 9999 // Auto-scroll to bottom
					return m, nil
				case "1":
					result := SceneComplication()
					m.currentLog.AddEntry("Scene Complication: "+result, "generator")
					m.errorMsg = ""
					// Auto-scroll to bottom
					maxScroll := m.calculateLogMaxScroll()
					m.scrollOffset = maxScroll
					return m, nil
				case "2":
					roll := RollD6()
					if roll >= 5 {
						result := AlteredScene()
						m.currentLog.AddEntry("Altered Scene: "+result, "generator")
					} else {
						m.currentLog.AddEntry("Normal Scene", "generator")
					}
					m.errorMsg = ""
					// Auto-scroll to bottom
					maxScroll := m.calculateLogMaxScroll()
					m.scrollOffset = maxScroll
					return m, nil
				case "3":
					result := PacingMove()
					m.currentLog.AddEntry("Pacing Move: "+result, "generator")
					m.errorMsg = ""
					// Auto-scroll to bottom
					maxScroll := m.calculateLogMaxScroll()
					m.scrollOffset = maxScroll
					return m, nil
				case "4":
					result := FailureMove()
					m.currentLog.AddEntry("Failure Move: "+result, "generator")
					m.errorMsg = ""
					// Auto-scroll to bottom
					maxScroll := m.calculateLogMaxScroll()
					m.scrollOffset = maxScroll
					return m, nil
				case "5":
					result := RandomEvent()
					m.currentLog.AddEntry("Random Event: "+result, "generator")
					m.errorMsg = ""
					// Auto-scroll to bottom
					maxScroll := m.calculateLogMaxScroll()
					m.scrollOffset = maxScroll
					return m, nil
				case "6":
					result := PlotHook()
					m.currentLog.AddEntry("Plot Hook: "+result, "generator")
					m.errorMsg = ""
					// Auto-scroll to bottom
					maxScroll := m.calculateLogMaxScroll()
					m.scrollOffset = maxScroll
					return m, nil
				case "7":
					result := NPCGenerator()
					m.currentLog.AddEntry("NPC: "+result, "generator")
					m.errorMsg = ""
					// Auto-scroll to bottom
					maxScroll := m.calculateLogMaxScroll()
					m.scrollOffset = maxScroll
					return m, nil
				case "8":
					result := GenericGenerator()
					m.currentLog.AddEntry("Generic Generator: "+result, "generator")
					m.errorMsg = ""
					// Auto-scroll to bottom
					maxScroll := m.calculateLogMaxScroll()
					m.scrollOffset = maxScroll
					return m, nil
				case "9":
					result := DungeonCrawler()
					m.currentLog.AddEntry("Dungeon Crawler: "+result, "generator")
					m.errorMsg = ""
					// Auto-scroll to bottom
					maxScroll := m.calculateLogMaxScroll()
					m.scrollOffset = maxScroll
					return m, nil
				case "0":
					result := HexCrawler()
					m.currentLog.AddEntry("Hex Crawler: "+result, "generator")
					m.errorMsg = ""
					// Auto-scroll to bottom
					maxScroll := m.calculateLogMaxScroll()
					m.scrollOffset = maxScroll
					return m, nil
				case "h":
					result := OracleHow()
					m.currentLog.AddEntry("Oracle (How): "+result, "oracle")
					m.errorMsg = ""
					// Auto-scroll to bottom
					maxScroll := m.calculateLogMaxScroll()
					m.scrollOffset = maxScroll
					return m, nil
				case "s":
					err := SaveLog(m.currentLog, m.logFilename)
					if err != nil {
						m.errorMsg = fmt.Sprintf("Error saving: %v", err)
					} else {
						m.errorMsg = ""
					}
					return m, nil
				}
			}

			// Check for multi-character commands
			cmdLower := strings.ToLower(cmd)
			if strings.HasPrefix(cmdLower, "save") || cmdLower == "s" {
				err := SaveLog(m.currentLog, m.logFilename)
				if err != nil {
					m.errorMsg = fmt.Sprintf("Error saving: %v", err)
				} else {
					m.errorMsg = ""
				}
				return m, nil
			}

			// Otherwise, treat as a new entry
			m.currentLog.AddEntry(cmd, "user")
			m.errorMsg = ""
			m.textInput = ""
			m.cursorPos = 0
			// Auto-scroll to bottom to show new entry
			maxScroll := m.calculateLogMaxScroll()
			m.scrollOffset = maxScroll
			return m, nil

		case "left":
			// Move cursor left in prompt window
			if m.selectedWindow == "prompt" {
				if m.cursorPos > 0 {
					m.cursorPos--
				}
				return m, nil
			}
		case "right":
			// Move cursor right in prompt window
			if m.selectedWindow == "prompt" {
				if m.cursorPos < len(m.textInput) {
					m.cursorPos++
				}
				return m, nil
			}
		case "backspace":
			if m.selectedWindow == "prompt" {
				// Delete character at cursor position
				if m.cursorPos > 0 && m.cursorPos <= len(m.textInput) {
					m.textInput = m.textInput[:m.cursorPos-1] + m.textInput[m.cursorPos:]
					m.cursorPos--
				}
			} else {
				// Legacy behavior for other windows
				if len(m.textInput) > 0 {
					m.textInput = m.textInput[:len(m.textInput)-1]
				}
			}
		case "tab":
			// Switch between log and prompt windows
			if m.selectedWindow == "log" {
				m.selectedWindow = "prompt"
			} else {
				m.selectedWindow = "log"
			}
			return m, nil
		case "up", "k":
			// Scroll the currently selected window
			if m.selectedWindow == "log" {
				// Scroll log window
				if m.scrollOffset > 0 {
					m.scrollOffset--
				}
			} else {
				// Move cursor up through wrapped lines in prompt window
				sidebarWidth := 28
				spacing := 2
				logWindowWidth := m.width - sidebarWidth - spacing - 4
				if logWindowWidth < 50 {
					logWindowWidth = 50
				}
				newCursorPos, newScrollOffset := m.moveCursorUpDown(-1, logWindowWidth)
				m.cursorPos = newCursorPos
				m.promptScrollOffset = newScrollOffset
			}
			return m, nil
		case "down", "j":
			// Scroll the currently selected window
			if m.selectedWindow == "log" {
				// Scroll log window
				maxScroll := m.calculateLogMaxScroll()
				if m.scrollOffset < maxScroll {
					m.scrollOffset++
				}
			} else {
				// Move cursor down through wrapped lines in prompt window
				sidebarWidth := 28
				spacing := 2
				logWindowWidth := m.width - sidebarWidth - spacing - 4
				if logWindowWidth < 50 {
					logWindowWidth = 50
				}
				newCursorPos, newScrollOffset := m.moveCursorUpDown(1, logWindowWidth)
				m.cursorPos = newCursorPos
				m.promptScrollOffset = newScrollOffset
			}
			return m, nil
		case "esc":
			// Open popup menu
			m.logMenuOpen = true
			m.logMenuIndex = 0
		default:
			// Handle text input only when prompt window is selected
			if m.selectedWindow == "prompt" {
				if len(msg.Runes) > 0 {
					// Insert text at cursor position
					insertText := string(msg.Runes)
					if m.cursorPos >= len(m.textInput) {
						// Append to end
						m.textInput += insertText
						m.cursorPos = len(m.textInput)
					} else {
						// Insert in middle
						m.textInput = m.textInput[:m.cursorPos] + insertText + m.textInput[m.cursorPos:]
						m.cursorPos += len(insertText)
					}
					// Auto-scroll prompt to bottom when typing
					m.promptScrollOffset = 9999 // Will be clamped in View
				}
			}
		}
	}
	return m, nil
}

func (m model) updateAddEntry(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if strings.TrimSpace(m.textInput) != "" {
				m.currentLog.AddEntry(strings.TrimSpace(m.textInput), "user")
				m.textInput = ""
				m.cursorPos = 0
				m.state = stateLogView
				m.logMenuOpen = false
				m.selectedWindow = "prompt"
			}
		case "backspace":
			if len(m.textInput) > 0 {
				m.textInput = m.textInput[:len(m.textInput)-1]
			}
		default:
			if len(msg.Runes) > 0 {
				m.textInput += string(msg.Runes)
			}
		}
	}
	return m, nil
}

func (m model) updateOracleYesNo(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.menuIndex > 0 {
				m.menuIndex--
			}
		case "down", "j":
			if m.menuIndex < 2 {
				m.menuIndex++
			}
		case "enter":
			likelihoods := []string{"Likely", "Even", "Unlikely"}
			m.likelihood = likelihoods[m.menuIndex]
			yes, result := OracleYesNo(m.likelihood)
			if yes {
				m.generatorResult = fmt.Sprintf("Oracle (Yes/No - %s):\nYes - %s", m.likelihood, result)
			} else {
				m.generatorResult = fmt.Sprintf("Oracle (Yes/No - %s):\nNo - %s", m.likelihood, result)
			}
			m.currentLog.AddEntry(m.generatorResult, "oracle")
			m.state = stateViewGeneratorResult
		}
	}
	return m, nil
}

func (m model) updateViewGeneratorResult(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			m.state = stateLogView
			m.textInput = ""
			m.cursorPos = 0
			m.logMenuOpen = false
			m.selectedWindow = "prompt"
		}
	}
	return m, nil
}

func (m model) updateViewLicense(msg tea.Msg) (tea.Model, tea.Cmd) {
	licenseText := "CC-BY-SA 4.0\n\nCreative Commons Attribution-ShareAlike 4.0 International\n\nThis work is licensed under the Creative Commons Attribution-ShareAlike 4.0 International License. To view a copy of this license, visit http://creativecommons.org/licenses/by-sa/4.0/ or send a letter to Creative Commons, PO Box 1866, Mountain View, CA 94042, USA.\n\nYou are free to:\n- Share — copy and redistribute the material in any medium or format\n- Adapt — remix, transform, and build upon the material for any purpose, even commercially\n\nUnder the following terms:\n- Attribution — You must give appropriate credit, provide a link to the license, and indicate if changes were made.\n- ShareAlike — If you remix, transform, or build upon the material, you must distribute your contributions under the same license as the original."
	lines := strings.Split(licenseText, "\n")

	// Calculate maxVisibleLines dynamically based on terminal height
	popupHeight := m.height - 10
	if popupHeight < 10 {
		popupHeight = 10
	}
	contentHeight := popupHeight - 4
	maxVisibleLines := contentHeight - 1 // Subtract 1 for help text
	if maxVisibleLines < 1 {
		maxVisibleLines = 1
	}

	maxScroll := len(lines) - maxVisibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = stateMainMenu
			m.scrollOffset = 0
		case "up", "k":
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}
		case "down", "j":
			if m.scrollOffset < maxScroll {
				m.scrollOffset++
			}
		case "pgup":
			m.scrollOffset -= 10
			if m.scrollOffset < 0 {
				m.scrollOffset = 0
			}
		case "pgdown":
			m.scrollOffset += 10
			if m.scrollOffset > maxScroll {
				m.scrollOffset = maxScroll
			}
		}
	}
	return m, nil
}

func parseHeadings(text string) []Heading {
	lines := strings.Split(text, "\n")
	var headings []Heading

	skipPatterns := []string{
		"Version:", "Email:", "Reddit:", "Download at:", "License:",
		"Written by:", "ONE PAGE SOLO ENGINE", "ONE PAGE GENERATORS",
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Check if line is a heading (all caps, not empty, not indented)
		if trimmed != "" &&
			strings.ToUpper(trimmed) == trimmed &&
			len(trimmed) > 3 &&
			!strings.HasPrefix(line, "\t") &&
			!strings.HasPrefix(line, " ") {
			// Skip metadata lines
			skip := false
			for _, pattern := range skipPatterns {
				if strings.Contains(trimmed, pattern) {
					skip = true
					break
				}
			}
			// Skip if it's a numbered list item (starts with number and tab)
			if !skip && strings.Contains(line, "\t") {
				parts := strings.SplitN(line, "\t", 2)
				if len(parts) > 0 {
					var num int
					if _, err := fmt.Sscanf(parts[0], "%d", &num); err == nil {
						skip = true
					}
				}
			}
			if !skip {
				headings = append(headings, Heading{Title: trimmed, Line: i})
			}
		}
	}
	return headings
}

func fuzzyMatch(search, text string) bool {
	search = strings.ToLower(search)
	text = strings.ToLower(text)
	if search == "" {
		return true
	}
	searchIdx := 0
	for i := 0; i < len(text) && searchIdx < len(search); i++ {
		if text[i] == search[searchIdx] {
			searchIdx++
		}
	}
	return searchIdx == len(search)
}

func (m model) updateViewOPSERawText(msg tea.Msg) (tea.Model, tea.Cmd) {
	lines := strings.Split(opseRawText, "\n")

	// Calculate fixed window dimensions: 4 rows from top + 4 rows from bottom
	popupHeight := m.height - 8
	if popupHeight < 10 {
		popupHeight = 10
	}
	contentHeight := popupHeight - 4     // border (2) + padding (2)
	maxVisibleLines := contentHeight - 1 // Subtract 1 for help text
	if maxVisibleLines < 1 {
		maxVisibleLines = 1
	}

	maxScroll := len(lines) - maxVisibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}

	// If sidebar is open, handle sidebar input
	if m.headingsSidebarOpen {
		// Calculate available heading lines for scrolling
		popupHeight := m.height - 8
		if popupHeight < 10 {
			popupHeight = 10
		}
		sidebarContentHeight := popupHeight - 4
		availableHeadingLines := sidebarContentHeight - 3 // title (1) + search (1) + spacing (1)
		if availableHeadingLines < 1 {
			availableHeadingLines = 1
		}

		maxHeadingScroll := len(m.filteredHeadings) - availableHeadingLines
		if maxHeadingScroll < 0 {
			maxHeadingScroll = 0
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				// Close sidebar first, window stays open
				m.headingsSidebarOpen = false
				m.headingSearch = ""
			case "up", "k":
				if m.headingIndex > 0 {
					m.headingIndex--
				}
				// Update scroll offset to keep selected item visible
				if m.headingIndex < m.headingScrollOffset {
					m.headingScrollOffset = m.headingIndex
				}
				if m.headingScrollOffset < 0 {
					m.headingScrollOffset = 0
				}
			case "down", "j":
				if m.headingIndex < len(m.filteredHeadings)-1 {
					m.headingIndex++
				}
				// Update scroll offset to keep selected item visible
				if m.headingIndex >= m.headingScrollOffset+availableHeadingLines {
					m.headingScrollOffset = m.headingIndex - availableHeadingLines + 1
				}
				if m.headingScrollOffset > maxHeadingScroll {
					m.headingScrollOffset = maxHeadingScroll
				}
			case "enter":
				if len(m.filteredHeadings) > 0 && m.headingIndex < len(m.filteredHeadings) {
					// Jump to the heading - set it to appear at the top of the window
					headingLine := m.filteredHeadings[m.headingIndex].Line
					// Calculate maxVisibleLines to ensure heading appears at top
					popupHeight := m.height - 8
					if popupHeight < 10 {
						popupHeight = 10
					}
					contentHeight := popupHeight - 4
					maxVisibleLines := contentHeight - 1
					if maxVisibleLines < 1 {
						maxVisibleLines = 1
					}
					// Set scroll so heading appears at the top
					m.scrollOffset = headingLine
					// Ensure we don't scroll past the end
					maxScroll := len(strings.Split(opseRawText, "\n")) - maxVisibleLines
					if maxScroll < 0 {
						maxScroll = 0
					}
					if m.scrollOffset > maxScroll {
						m.scrollOffset = maxScroll
					}
					m.headingsSidebarOpen = false
					m.headingSearch = ""
				}
			case "backspace":
				if len(m.headingSearch) > 0 {
					m.headingSearch = m.headingSearch[:len(m.headingSearch)-1]
					// Re-filter headings
					m.filteredHeadings = []Heading{}
					for _, h := range m.headings {
						if fuzzyMatch(m.headingSearch, h.Title) {
							m.filteredHeadings = append(m.filteredHeadings, h)
						}
					}
					// Reset index and scroll if out of bounds
					if m.headingIndex >= len(m.filteredHeadings) {
						m.headingIndex = 0
					}
					if len(m.filteredHeadings) == 0 {
						m.headingIndex = 0
					}
					m.headingScrollOffset = 0
				}
			default:
				// Handle text input for search
				if len(msg.Runes) > 0 {
					m.headingSearch += string(msg.Runes)
					// Re-filter headings
					m.filteredHeadings = []Heading{}
					for _, h := range m.headings {
						if fuzzyMatch(m.headingSearch, h.Title) {
							m.filteredHeadings = append(m.filteredHeadings, h)
						}
					}
					// Reset index and scroll
					m.headingIndex = 0
					m.headingScrollOffset = 0
				}
			}
		}
		return m, nil
	}

	// Normal text view input
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Close window only if sidebar is not visible
			// (If sidebar was open, it was already closed above and we returned early)
			m.state = stateMainMenu
			m.scrollOffset = 0
		case "c", "C":
			// Open sidebar (only if not already open)
			if !m.headingsSidebarOpen {
				m.headingsSidebarOpen = true
				m.headingSearch = ""
				m.headingIndex = 0
				m.headingScrollOffset = 0
				m.filteredHeadings = m.headings
			}
		case "up", "k":
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}
		case "down", "j":
			if m.scrollOffset < maxScroll {
				m.scrollOffset++
			}
		case "pgup":
			m.scrollOffset -= 10
			if m.scrollOffset < 0 {
				m.scrollOffset = 0
			}
		case "pgdown":
			m.scrollOffset += 10
			if m.scrollOffset > maxScroll {
				m.scrollOffset = maxScroll
			}
		}
	}
	return m, nil
}

func (m model) updateViewOPSERawTextHeadings(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "c", "C":
			m.state = stateViewOPSERawText
		case "up", "k":
			if m.headingIndex > 0 {
				m.headingIndex--
			}
		case "down", "j":
			if m.headingIndex < len(m.headings)-1 {
				m.headingIndex++
			}
		case "enter":
			if len(m.headings) > 0 && m.headingIndex < len(m.headings) {
				// Jump to the heading
				m.scrollOffset = m.headings[m.headingIndex].Line
				m.state = stateViewOPSERawText
			}
		}
	}
	return m, nil
}

func (m model) wrapWithBorder(content string) string {
	if m.width == 0 {
		m.width = 80
	}
	if m.height == 0 {
		m.height = 24
	}
	// The interior space must be exactly (width-2) x (height-2)
	// Create a style box that explicitly fills this space
	interiorStyle := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.height - 2).
		Align(lipgloss.Left).
		AlignVertical(lipgloss.Top)

	// Render content in the interior style - this ensures it fills the space
	filledInterior := interiorStyle.Render(content)

	// Apply border - this creates a box of exactly width x height
	return borderStyle.Width(m.width).Height(m.height).Render(filledInterior)
}

func (m model) wrapWithBorderCentered(content string) string {
	if m.width == 0 {
		m.width = 80
	}
	if m.height == 0 {
		m.height = 24
	}
	// Center the content without a border
	centeredStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center)

	return centeredStyle.Render(content)
}

func (m model) overlayPopup(baseContent string, popup string) string {
	// Use lipgloss.Place to center the popup on a full canvas
	centeredPopup := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		popup,
	)

	// Split both into lines
	baseLines := strings.Split(baseContent, "\n")
	popupLines := strings.Split(centeredPopup, "\n")

	// Ensure both have the same number of lines
	for len(baseLines) < m.height {
		baseLines = append(baseLines, strings.Repeat(" ", m.width))
	}
	for len(popupLines) < m.height {
		popupLines = append(popupLines, strings.Repeat(" ", m.width))
	}

	// Trim to exact height
	if len(baseLines) > m.height {
		baseLines = baseLines[:m.height]
	}
	if len(popupLines) > m.height {
		popupLines = popupLines[:m.height]
	}

	// Merge popup over base content character by character
	// For each line, merge where popup has non-space content, preserve base where popup is transparent
	resultLines := make([]string, m.height)
	for lineIdx := 0; lineIdx < m.height; lineIdx++ {
		baseLine := baseLines[lineIdx]
		popupLine := popupLines[lineIdx]

		// Merge character by character, using popup where it has content, base where popup is space
		resultLines[lineIdx] = m.mergeLineCharacterByCharacter(baseLine, popupLine, m.width)
	}

	return strings.Join(resultLines, "\n")
}

// mergeLineCharacterByCharacter merges two lines character by character
// Uses popup character where it's not a space, base character where popup is a space
// Handles ANSI escape sequences and UTF-8 properly
func (m model) mergeLineCharacterByCharacter(baseLine string, popupLine string, maxWidth int) string {
	baseBytes := []byte(baseLine)
	popupBytes := []byte(popupLine)

	var result strings.Builder
	basePos := 0
	popupPos := 0
	resultDisplayWidth := 0

	// Helper to get next character/ANSI sequence from a byte slice
	getNext := func(bytes []byte, pos int) (content []byte, newPos int, isAnsi bool, charWidth int) {
		if pos >= len(bytes) {
			return nil, pos, false, 0
		}

		// Check for ANSI escape sequence
		if pos+1 < len(bytes) && bytes[pos] == '\033' && bytes[pos+1] == '[' {
			// Copy ANSI sequence
			start := pos
			pos += 2
			for pos < len(bytes) {
				if bytes[pos] == 'm' || bytes[pos] == 'H' || bytes[pos] == 'J' || bytes[pos] == 'K' {
					pos++
					return bytes[start:pos], pos, true, 0
				}
				pos++
			}
			return bytes[start:pos], pos, true, 0
		}

		// Decode UTF-8 rune
		r, size := utf8.DecodeRune(bytes[pos:])
		if size == 0 {
			return []byte{bytes[pos]}, pos + 1, false, 1
		}

		charWidth = lipgloss.Width(string(r))
		return bytes[pos : pos+size], pos + size, false, charWidth
	}

	// Merge byte by byte, handling ANSI and UTF-8
	for resultDisplayWidth < maxWidth && (basePos < len(baseBytes) || popupPos < len(popupBytes)) {
		// Check popup first
		if popupPos < len(popupBytes) {
			popupContent, newPopupPos, popupIsAnsi, popupWidth := getNext(popupBytes, popupPos)

			if popupContent != nil {
				// Check if popup has visible content (not just space)
				if popupIsAnsi {
					// Always use ANSI codes from popup
					result.Write(popupContent)
					popupPos = newPopupPos
					// Skip corresponding content in base
					if basePos < len(baseBytes) {
						_, newBasePos, baseIsAnsi, _ := getNext(baseBytes, basePos)
						if baseIsAnsi {
							basePos = newBasePos
						} else {
							basePos = newBasePos
						}
					}
					continue
				} else {
					// Check if it's a space
					r, _ := utf8.DecodeRune(popupContent)
					if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
						// Popup has visible content - use it
						result.Write(popupContent)
						popupPos = newPopupPos
						resultDisplayWidth += popupWidth
						// Skip corresponding content in base
						if basePos < len(baseBytes) {
							_, newBasePos, _, _ := getNext(baseBytes, basePos)
							basePos = newBasePos
						}
						continue
					}
				}
			}
		}

		// Popup is space or empty - use base
		if basePos < len(baseBytes) {
			baseContent, newBasePos, _, baseWidth := getNext(baseBytes, basePos)
			if baseContent != nil {
				result.Write(baseContent)
				basePos = newBasePos
				resultDisplayWidth += baseWidth
			} else {
				break
			}
		} else if popupPos < len(popupBytes) {
			// No more base, but popup has space - use space
			result.WriteByte(' ')
			popupPos++
			resultDisplayWidth++
		} else {
			// Both exhausted, fill with space
			result.WriteByte(' ')
			resultDisplayWidth++
		}

		// Advance popup if it was a space
		if popupPos < len(popupBytes) && popupBytes[popupPos] == ' ' {
			popupPos++
		}
	}

	return result.String()
}

func (m model) View() string {
	var content strings.Builder

	switch m.state {
	case stateMainMenu:
		// Create centered window content
		var menuContent strings.Builder

		// Centered title
		menuContent.WriteString(menuTitleStyle.Width(60).Render("OPSE: One Page Solo Engine"))
		menuContent.WriteString("\n\n")

		// Menu items with [ option ] format
		menuItems := []string{"New Game Log", "Load Game Log", "View License", "View OPSE Raw Text", "Quit"}
		for i, item := range menuItems {
			if i == m.menuIndex {
				menuContent.WriteString("[ ")
				menuContent.WriteString(selectedMenuItemStyle.Render(item))
				menuContent.WriteString(" ]")
			} else {
				menuContent.WriteString("  ")
				menuContent.WriteString(menuItemStyle.Render(item))
				menuContent.WriteString("  ")
			}
			menuContent.WriteString("\n")
		}

		menuContent.WriteString(helpStyle.Render("\nUse ↑/↓ to navigate, Enter to select, Ctrl+C to quit"))

		// Render the menu content with styling
		styledContent := centeredWindowStyle.Render(menuContent.String())

		return m.wrapWithBorderCentered(styledContent)

	case stateNewLog:
		content.WriteString(titleStyle.Render("New Game Log"))
		content.WriteString("\n\n")
		content.WriteString("Enter log title: ")
		content.WriteString(m.textInput)
		content.WriteString("_")
		content.WriteString(helpStyle.Render("\n\nPress Enter to create, Esc to cancel"))

	case stateLoadLog:
		content.WriteString(titleStyle.Render("Load Game Log"))
		content.WriteString("\n\n")
		if len(m.availableLogs) == 0 {
			content.WriteString(menuItemStyle.Render("No log files found."))
		} else {
			for i, log := range m.availableLogs {
				if i == m.logIndex {
					content.WriteString(selectedMenuItemStyle.Render(fmt.Sprintf("> %s", log)))
				} else {
					content.WriteString(menuItemStyle.Render(fmt.Sprintf("  %s", log)))
				}
				content.WriteString("\n")
			}
		}
		content.WriteString(helpStyle.Render("\nUse ↑/↓ to navigate, Enter to load, Esc to cancel"))

	case stateLogView:
		// If menu is open, show it by itself
		if m.logMenuOpen {
			// Build menu content
			var menuContent strings.Builder
			// Center the title within the menu window content width (40 - 4 for border/padding = 36)
			menuTitleCentered := lipgloss.NewStyle().
				Width(36).
				Align(lipgloss.Center).
				Bold(true).
				Foreground(lipgloss.Color("3")).
				Render("Menu")
			menuContent.WriteString(menuTitleCentered)
			menuContent.WriteString("\n\n")

			menuItems := []string{"Save Log", "Load Log", "View License", "View Raw Text", "Exit OPSE"}
			for i, item := range menuItems {
				if i == m.logMenuIndex {
					menuContent.WriteString(selectedMenuItemStyle.Render(fmt.Sprintf("> %s", item)))
				} else {
					menuContent.WriteString(menuItemStyle.Render(fmt.Sprintf("  %s", item)))
				}
				menuContent.WriteString("\n")
			}

			menuContent.WriteString(helpStyle.Render("\nUse ↑/↓ to navigate, Enter to select, Esc to close"))

			// Create menu window
			menuWindow := logMenuPopupStyle.Render(menuContent.String())

			// Center the menu on screen
			centeredMenu := lipgloss.Place(
				m.width,
				m.height,
				lipgloss.Center,
				lipgloss.Center,
				menuWindow,
			)

			return centeredMenu
		}

		// Calculate window dimensions - 4 rows from top margin
		// Sidebar height should equal log window + prompt window (no spacing between them)
		sidebarHeight := m.height - 4 // Full available height
		if sidebarHeight < 15 {
			sidebarHeight = 15 // Minimum height
		}
		promptHeight := 5 // Fixed height for prompt window - never change this
		// Log height = sidebar height - prompt height - 2 (2 rows shorter)
		logHeight := sidebarHeight - promptHeight - 2
		if logHeight < 10 {
			logHeight = 10
			// Don't adjust prompt height - keep it fixed at 5
		}

		// Define sidebar and log window widths
		sidebarWidth := 28 // Reduced by 20% from 35
		spacing := 2
		logWindowWidth := m.width - sidebarWidth - spacing - 4 // Account for margins
		if logWindowWidth < 50 {
			logWindowWidth = 50
		}

		// Content area = log height - border (2) - padding (2) = log height - 4
		// But to allow text within 1 row of bottom, we need to account for the actual usable space
		// If currently cutting off 5 rows, we need to use more of the available height
		// Try using logHeight - 2 (just border) instead of -4, then subtract 1 for margin
		contentHeight := logHeight - 2       // Only account for border (top + bottom)
		maxVisibleLines := contentHeight - 1 // Allow 1 row margin from bottom
		if maxVisibleLines < 1 {
			maxVisibleLines = 1
		}

		// Build sidebar content with commands and oracles
		var sidebarContent strings.Builder
		sidebarContent.WriteString(sidebarTitleStyle.Render("One Page Solo Engine"))
		sidebarContent.WriteString("\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("Esc") + " " + sidebarTextStyle.Render("Menu") + "\n")
		sidebarContent.WriteString("\n")
		sidebarContent.WriteString(sidebarHeadingStyle.Render("Commands"))
		sidebarContent.WriteString("\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("l") + " " + sidebarTextStyle.Render("Oracle Likely") + "\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("e") + " " + sidebarTextStyle.Render("Oracle Even") + "\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("u") + " " + sidebarTextStyle.Render("Oracle Unlikely") + "\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("h") + " " + sidebarTextStyle.Render("Oracle How") + "\n")
		sidebarContent.WriteString("\n")
		sidebarContent.WriteString(sidebarHeadingStyle.Render("Generators"))
		sidebarContent.WriteString("\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("1") + " " + sidebarTextStyle.Render("Scene") + "\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("2") + " " + sidebarTextStyle.Render("Altered Scene") + "\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("3") + " " + sidebarTextStyle.Render("Pacing Move") + "\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("4") + " " + sidebarTextStyle.Render("Failure Move") + "\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("5") + " " + sidebarTextStyle.Render("Random Event") + "\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("6") + " " + sidebarTextStyle.Render("Plot Hook") + "\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("7") + " " + sidebarTextStyle.Render("NPC") + "\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("8") + " " + sidebarTextStyle.Render("Generic") + "\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("9") + " " + sidebarTextStyle.Render("Dungeon") + "\n")
		sidebarContent.WriteString(sidebarShortcutStyle.Render("0") + " " + sidebarTextStyle.Render("Hex") + "\n")

		// Pad sidebar content to match height
		sidebarLines := strings.Split(sidebarContent.String(), "\n")
		sidebarContentHeight := sidebarHeight - 4 // Account for border and padding
		for len(sidebarLines) < sidebarContentHeight {
			sidebarLines = append(sidebarLines, "")
		}
		sidebarText := strings.Join(sidebarLines, "\n")

		// Create sidebar window (full height)
		sidebarWindow := logSidebarStyle.
			Width(sidebarWidth).
			Height(sidebarHeight).
			Render(sidebarText)

		// Build all log content as lines for proper scrolling
		var allLogLines []string
		allLogLines = append(allLogLines, fmt.Sprintf("Game Log: %s", m.currentLog.Title))
		allLogLines = append(allLogLines, "")

		if len(m.currentLog.Entries) == 0 {
			allLogLines = append(allLogLines, menuItemStyle.Render("No entries yet. Type a command or entry below."))
		} else {
			for _, entry := range m.currentLog.Entries {
				entryText := formatLogEntryColored(entry)
				entryLines := strings.Split(entryText, "\n")
				allLogLines = append(allLogLines, entryLines...)
				allLogLines = append(allLogLines, "") // Spacing between entries
			}
		}

		// Calculate scroll bounds
		totalLines := len(allLogLines)
		maxScroll := totalLines - maxVisibleLines
		if maxScroll < 0 {
			maxScroll = 0
		}
		// Clamp scrollOffset to valid bounds (don't modify model in View, just use clamped value)
		scrollOffset := m.scrollOffset
		if scrollOffset > maxScroll {
			scrollOffset = maxScroll
		}
		if scrollOffset < 0 {
			scrollOffset = 0
		}

		// Get visible lines based on scroll offset
		start := scrollOffset
		end := start + maxVisibleLines
		if end > totalLines {
			end = totalLines
		}

		// Build visible content (reserve 1 char for scroll bar)
		contentWidth := logWindowWidth - 6 - 1 // Account for border, padding, and scroll bar
		if contentWidth < 10 {
			contentWidth = 10
		}

		var logContent strings.Builder
		for i := start; i < end; i++ {
			if i < len(allLogLines) {
				line := allLogLines[i]
				// Truncate or wrap line to fit content width
				lineDisplayWidth := lipgloss.Width(line)
				if lineDisplayWidth > contentWidth {
					// Truncate line to fit
					truncated := lipgloss.NewStyle().
						Width(contentWidth).
						Render(line)
					logContent.WriteString(truncated)
				} else {
					logContent.WriteString(line)
				}
				logContent.WriteString("\n")
			}
		}

		// Pad to ensure consistent height (use maxVisibleLines, not +1)
		logText := logContent.String()
		logLines := strings.Split(logText, "\n")
		// Remove trailing newline if present
		if len(logLines) > 0 && logLines[len(logLines)-1] == "" {
			logLines = logLines[:len(logLines)-1]
		}
		for len(logLines) < maxVisibleLines {
			logLines = append(logLines, "")
		}

		// Render scroll bar
		scrollBarHeight := maxVisibleLines
		scrollBarText := renderScrollBar(scrollBarHeight, scrollOffset, totalLines, maxVisibleLines)
		scrollBarLines := strings.Split(scrollBarText, "\n")
		// Ensure scroll bar has same number of lines as log content
		for len(scrollBarLines) < maxVisibleLines {
			scrollBarLines = append(scrollBarLines, " ")
		}
		if len(scrollBarLines) > maxVisibleLines {
			scrollBarLines = scrollBarLines[:maxVisibleLines]
		}

		// Combine log content and scroll bar horizontally, line by line
		var logWithScrollBar strings.Builder
		for i := 0; i < maxVisibleLines; i++ {
			logLine := logLines[i]
			scrollBarChar := " "
			if i < len(scrollBarLines) {
				scrollBarChar = scrollBarLines[i]
			}
			// Pad log line to content width, then add scroll bar
			paddedLogLine := lipgloss.NewStyle().
				Width(contentWidth).
				Render(logLine)
			logWithScrollBar.WriteString(paddedLogLine)
			logWithScrollBar.WriteString(" ")
			logWithScrollBar.WriteString(scrollBarChar)
			if i < maxVisibleLines-1 {
				logWithScrollBar.WriteString("\n")
			}
		}
		logText = logWithScrollBar.String()

		// Create log window with appropriate border color based on selection
		var logStyle lipgloss.Style
		if m.selectedWindow == "log" {
			logStyle = logViewStyleSelected
		} else {
			logStyle = logViewStyle
		}
		logWindow := logStyle.
			Width(logWindowWidth).
			Height(logHeight).
			Render(logText)

		// Build prompt window content with scrolling
		// Build lines - text will wrap naturally within the window width
		var promptContentLines []string
		cursorChar := "_"

		// Build the prompt line with cursor
		var promptLine strings.Builder
		promptLine.WriteString(promptSymbolStyle.Render("> "))
		if m.cursorPos >= len(m.textInput) {
			// Cursor at end
			promptLine.WriteString(m.textInput)
			promptLine.WriteString(cursorChar)
		} else {
			// Cursor in middle
			beforeCursor := m.textInput[:m.cursorPos]
			afterCursor := m.textInput[m.cursorPos:]
			promptLine.WriteString(beforeCursor)
			promptLine.WriteString(cursorChar)
			promptLine.WriteString(afterCursor)
		}

		// The window will handle wrapping, but we need to account for it in scrolling
		// For now, treat the prompt line as a single line (it will wrap in the window)
		promptContentLines = append(promptContentLines, promptLine.String())

		// Add error message lines if any
		if m.errorMsg != "" {
			errorLines := strings.Split(promptHelpStyle.Render(m.errorMsg), "\n")
			promptContentLines = append(promptContentLines, errorLines...)
		}

		// Calculate how many lines the content will take when wrapped
		// Available width = window width - border (2) - padding (4)
		availableWidth := logWindowWidth - 6
		if availableWidth < 10 {
			availableWidth = 10
		}

		// Wrap each line and build all wrapped lines
		var allWrappedLines []string
		for _, line := range promptContentLines {
			// Wrap this line to available width
			wrappedLine := lipgloss.NewStyle().
				Width(availableWidth).
				Render(line)
			// Split the wrapped line (may be multiple lines)
			wrappedLineParts := strings.Split(wrappedLine, "\n")
			// Ensure each wrapped line doesn't exceed available width
			for _, part := range wrappedLineParts {
				// Truncate if necessary to prevent overflow
				partWidth := lipgloss.Width(part)
				if partWidth > availableWidth {
					// Truncate to exact width
					truncated := lipgloss.NewStyle().
						Width(availableWidth).
						Render(part)
					allWrappedLines = append(allWrappedLines, truncated)
				} else {
					allWrappedLines = append(allWrappedLines, part)
				}
			}
		}

		// Calculate prompt scroll bounds
		promptContentHeight := promptHeight - 4 // Account for border and padding
		maxPromptVisibleLines := promptContentHeight - 1
		if maxPromptVisibleLines < 1 {
			maxPromptVisibleLines = 1
		}

		totalPromptLines := len(allWrappedLines)
		maxPromptScroll := totalPromptLines - maxPromptVisibleLines
		if maxPromptScroll < 0 {
			maxPromptScroll = 0
		}

		// Find which wrapped line the cursor is on
		_, cursorWrappedLine, _ := m.calculatePromptWrappedLines()

		// Auto-adjust scroll to keep cursor visible
		promptScrollOffset := m.promptScrollOffset
		if promptScrollOffset == 9999 {
			// Special value set when typing - scroll to show cursor
			promptScrollOffset = cursorWrappedLine - maxPromptVisibleLines + 1
			if promptScrollOffset < 0 {
				promptScrollOffset = 0
			}
			if promptScrollOffset > maxPromptScroll {
				promptScrollOffset = maxPromptScroll
			}
		} else {
			// Ensure cursor is visible
			if cursorWrappedLine < promptScrollOffset {
				// Cursor is above visible area, scroll up
				promptScrollOffset = cursorWrappedLine
			} else if cursorWrappedLine >= promptScrollOffset+maxPromptVisibleLines {
				// Cursor is below visible area, scroll down
				promptScrollOffset = cursorWrappedLine - maxPromptVisibleLines + 1
			}
		}

		// Clamp prompt scroll offset
		if promptScrollOffset > maxPromptScroll {
			promptScrollOffset = maxPromptScroll
		}
		if promptScrollOffset < 0 {
			promptScrollOffset = 0
		}

		// Get visible prompt lines based on scroll offset (vertical scrolling)
		promptStart := promptScrollOffset
		promptEnd := promptStart + maxPromptVisibleLines
		if promptEnd > totalPromptLines {
			promptEnd = totalPromptLines
		}

		var visiblePromptContent strings.Builder
		for i := promptStart; i < promptEnd; i++ {
			if i < len(allWrappedLines) {
				line := allWrappedLines[i]
				// Ensure line doesn't exceed available width
				lineWidth := lipgloss.Width(line)
				if lineWidth > availableWidth {
					// Truncate to exact width
					line = lipgloss.NewStyle().
						Width(availableWidth).
						Render(line)
				} else {
					// Pad to exact width to prevent any overflow
					line = lipgloss.NewStyle().
						Width(availableWidth).
						Render(line)
				}
				visiblePromptContent.WriteString(line)
				visiblePromptContent.WriteString("\n")
			}
		}

		// Pad to ensure consistent height
		promptText := visiblePromptContent.String()
		promptLines := strings.Split(promptText, "\n")
		// Remove trailing empty line if present
		if len(promptLines) > 0 && promptLines[len(promptLines)-1] == "" {
			promptLines = promptLines[:len(promptLines)-1]
		}
		// Ensure all lines are exactly the right width
		for i, line := range promptLines {
			lineWidth := lipgloss.Width(line)
			if lineWidth != availableWidth {
				// Pad or truncate to exact width
				promptLines[i] = lipgloss.NewStyle().
					Width(availableWidth).
					Render(line)
			}
		}
		for len(promptLines) < maxPromptVisibleLines {
			// Add empty lines that are exactly the right width
			emptyLine := lipgloss.NewStyle().
				Width(availableWidth).
				Render("")
			promptLines = append(promptLines, emptyLine)
		}
		promptText = strings.Join(promptLines, "\n")

		// Create prompt window with appropriate border color based on selection
		var promptStyle lipgloss.Style
		if m.selectedWindow == "prompt" {
			promptStyle = promptWindowStyleSelected
		} else {
			promptStyle = promptWindowStyle
		}
		// Final safety check: ensure promptText lines don't exceed available width
		// This prevents any potential overflow from ANSI codes or special characters
		finalPromptLines := strings.Split(promptText, "\n")
		for i, line := range finalPromptLines {
			lineWidth := lipgloss.Width(line)
			if lineWidth > availableWidth {
				// Truncate to exact width
				finalPromptLines[i] = lipgloss.NewStyle().
					Width(availableWidth).
					Render(line)
			}
		}
		promptText = strings.Join(finalPromptLines, "\n")

		promptWindow := promptStyle.
			Width(logWindowWidth).
			Height(promptHeight).
			Render(promptText)

		// Stack log window and prompt window vertically
		logAndPromptStack := lipgloss.JoinVertical(lipgloss.Left, logWindow, promptWindow)

		// Combine sidebar with the stacked log and prompt windows horizontally
		combinedWindows := lipgloss.JoinHorizontal(lipgloss.Top, sidebarWindow, strings.Repeat(" ", spacing), logAndPromptStack)

		// Build the full layout with 4 rows spacing at top
		var fullContent strings.Builder
		fullContent.WriteString(strings.Repeat("\n", 4)) // 4 rows from top
		fullContent.WriteString(combinedWindows)

		return fullContent.String()

	case stateAddEntry:
		content.WriteString(titleStyle.Render("Add Entry"))
		content.WriteString("\n\n")
		content.WriteString("Enter your entry:\n")
		content.WriteString(m.textInput)
		content.WriteString("_")
		content.WriteString(helpStyle.Render("\n\nPress Enter to save, Esc to cancel"))

	case stateOracleYesNo:
		content.WriteString(titleStyle.Render("Oracle (Yes/No)"))
		content.WriteString("\n\n")
		content.WriteString("Select likelihood:\n\n")

		likelihoods := []string{"Likely", "Even", "Unlikely"}
		for i, likelihood := range likelihoods {
			if i == m.menuIndex {
				content.WriteString(selectedMenuItemStyle.Render(fmt.Sprintf("> %s", likelihood)))
			} else {
				content.WriteString(menuItemStyle.Render(fmt.Sprintf("  %s", likelihood)))
			}
			content.WriteString("\n")
		}
		content.WriteString(helpStyle.Render("\nUse ↑/↓ to navigate, Enter to roll, Esc to cancel"))

	case stateViewGeneratorResult:
		content.WriteString(titleStyle.Render("Generator Result"))
		content.WriteString("\n\n")
		content.WriteString(resultStyle.Render(m.generatorResult))
		content.WriteString(helpStyle.Render("\n\nPress Enter or Space to continue"))

	case stateViewLicense:
		licenseText := "CC-BY-SA 4.0\n\nCreative Commons Attribution-ShareAlike 4.0 International\n\nThis work is licensed under the Creative Commons Attribution-ShareAlike 4.0 International License. To view a copy of this license, visit http://creativecommons.org/licenses/by-sa/4.0/ or send a letter to Creative Commons, PO Box 1866, Mountain View, CA 94042, USA.\n\nYou are free to:\n- Share — copy and redistribute the material in any medium or format\n- Adapt — remix, transform, and build upon the material for any purpose, even commercially\n\nUnder the following terms:\n- Attribution — You must give appropriate credit, provide a link to the license, and indicate if changes were made.\n- ShareAlike — If you remix, transform, or build upon the material, you must distribute your contributions under the same license as the original."

		// Calculate popup height: 5 rows from top + 5 rows from bottom
		popupHeight := m.height - 10
		if popupHeight < 10 {
			popupHeight = 10 // Minimum height
		}

		// Content area = popup height - border (2) - padding (2) = popup height - 4
		contentHeight := popupHeight - 4
		maxVisibleLines := contentHeight - 1 // Subtract 1 for help text
		if maxVisibleLines < 1 {
			maxVisibleLines = 1
		}

		// Split license text into lines for scrolling
		lines := strings.Split(licenseText, "\n")

		// Get visible lines
		start := m.scrollOffset
		end := start + maxVisibleLines
		if end > len(lines) {
			end = len(lines)
		}

		visibleLines := lines[start:end]
		// Pad with empty lines to ensure consistent height
		for len(visibleLines) < maxVisibleLines {
			visibleLines = append(visibleLines, "")
		}

		visibleText := strings.Join(visibleLines, "\n")

		// Add help text at the bottom
		helpText := helpStyle.Render("Press Esc to close, ↑/↓ to scroll")

		// Combine license text and help
		popupContent := visibleText + "\n" + helpText

		// Create the popup with blue border - use dynamic height
		popup := licensePopupStyle.
			Width(98).
			Height(popupHeight).
			Render(popupContent)

		// Center the popup on screen using Place method
		centered := lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			popup,
		)

		return centered

	case stateViewOPSERawText:
		// Calculate FIXED popup height: 4 rows from top + 4 rows from bottom
		// This ensures the window NEVER changes size
		popupHeight := m.height - 8
		if popupHeight < 10 {
			popupHeight = 10 // Minimum height
		}

		// Content area = popup height - border (2) - padding (2) = popup height - 4
		contentHeight := popupHeight - 4
		maxVisibleLines := contentHeight - 1 // Subtract 1 for help text
		if maxVisibleLines < 1 {
			maxVisibleLines = 1
		}

		// Split text into lines for scrolling
		lines := strings.Split(opseRawText, "\n")

		// Get visible lines
		start := m.scrollOffset
		end := start + maxVisibleLines
		if end > len(lines) {
			end = len(lines)
		}

		visibleLines := lines[start:end]
		// Pad with empty lines to ensure consistent height
		for len(visibleLines) < maxVisibleLines {
			visibleLines = append(visibleLines, "")
		}

		visibleText := strings.Join(visibleLines, "\n")

		// Add help text at the bottom
		helpText := "Press Esc to close, ↑/↓ to scroll"
		if m.headingsSidebarOpen {
			helpText = "Press Esc to close sidebar, ↑/↓ to navigate, Enter to jump, type to search"
		} else {
			helpText = "Press Esc to close, ↑/↓ to scroll, C to show headings"
		}
		helpTextRendered := helpStyle.Render(helpText)

		// Combine text and help
		popupContent := visibleText + "\n" + helpTextRendered

		// Calculate FIXED widths for split view
		// Main window width is always 98, regardless of sidebar
		mainWidth := 98
		sidebarWidth := 35

		// Create the main popup with blue border - FIXED dimensions
		mainPopup := rawTextPopupStyle.
			Width(mainWidth).
			Height(popupHeight).
			Render(popupContent)

		// If sidebar is open, create sidebar and join
		if m.headingsSidebarOpen {
			// Calculate sidebar content area (same as main window)
			// Content area = popup height - border (2) - padding (2) = popup height - 4
			sidebarContentHeight := popupHeight - 4

			// Build sidebar content
			var sidebarContent strings.Builder
			sidebarContent.WriteString(menuTitleStyle.Render("Headings"))
			sidebarContent.WriteString("\n")

			// Search input
			searchPrompt := "Search: "
			searchLine := searchPrompt + m.headingSearch + "_"
			sidebarContent.WriteString(selectedMenuItemStyle.Render(searchLine))
			sidebarContent.WriteString("\n")

			// Available space for headings list = content height - title (1) - search (1) - spacing (1) = content height - 3
			availableHeadingLines := sidebarContentHeight - 3
			if availableHeadingLines < 1 {
				availableHeadingLines = 1
			}

			// Get visible headings using scroll offset
			headingStart := m.headingScrollOffset
			headingEnd := headingStart + availableHeadingLines
			if headingEnd > len(m.filteredHeadings) {
				headingEnd = len(m.filteredHeadings)
			}
			if headingStart < 0 {
				headingStart = 0
			}
			if headingStart > len(m.filteredHeadings) {
				headingStart = len(m.filteredHeadings)
			}

			var visibleHeadings []Heading
			if headingStart < headingEnd {
				visibleHeadings = m.filteredHeadings[headingStart:headingEnd]
			}

			// Display headings
			headingLines := []string{}
			if len(m.filteredHeadings) == 0 {
				headingLines = append(headingLines, menuItemStyle.Render("No matches"))
			} else {
				for i, heading := range visibleHeadings {
					actualIndex := headingStart + i
					if actualIndex == m.headingIndex {
						headingLines = append(headingLines, selectedMenuItemStyle.Render(fmt.Sprintf("> %s", heading.Title)))
					} else {
						headingLines = append(headingLines, menuItemStyle.Render(fmt.Sprintf("  %s", heading.Title)))
					}
				}
			}

			// Pad with empty lines to ensure consistent height (prevents overflow)
			for len(headingLines) < availableHeadingLines {
				headingLines = append(headingLines, "")
			}

			// Join heading lines
			headingsText := strings.Join(headingLines, "\n")
			sidebarContent.WriteString(headingsText)

			// Create sidebar with border - FIXED dimensions matching main window
			sidebar := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")). // Dark gray
				Padding(1, 2).
				Width(sidebarWidth).
				Height(popupHeight).
				Render(sidebarContent.String())

			// Join main and sidebar
			combined := lipgloss.JoinHorizontal(lipgloss.Top, mainPopup, "  ", sidebar)

			// Center the combined view
			centered := lipgloss.Place(
				m.width,
				m.height,
				lipgloss.Center,
				lipgloss.Center,
				combined,
			)
			return centered
		}

		// No sidebar - center the main popup
		centered := lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			mainPopup,
		)

		return centered

	default:
		content.WriteString("Unknown state")
	}

	return m.wrapWithBorder(content.String())
}
