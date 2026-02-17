package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type helpPage struct {
	title   string
	content string
}

type HelpModel struct {
	page     int
	pages    []helpPage
	viewport viewport.Model
	ready    bool
}

func NewHelp() HelpModel {
	return HelpModel{
		pages: buildHelpPages(),
	}
}

func (h *HelpModel) Reset() {
	h.page = 0
	if h.ready {
		h.viewport.SetContent(h.pages[0].content)
		h.viewport.GotoTop()
	}
}

func (h *HelpModel) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if h.page > 0 {
				h.page--
				h.viewport.SetContent(h.pages[h.page].content)
				h.viewport.GotoTop()
			}
		case "right", "l":
			if h.page < len(h.pages)-1 {
				h.page++
				h.viewport.SetContent(h.pages[h.page].content)
				h.viewport.GotoTop()
			}
		default:
			h.viewport, _ = h.viewport.Update(msg)
		}
	}
}

func (h *HelpModel) View(width, height int) string {
	boxW := width - 4
	if boxW < 50 {
		boxW = 50
	}
	boxH := height - 2
	if boxH < 12 {
		boxH = 12
	}

	contentW := boxW - 4 // padding(2×2) — border is outside Width
	contentH := boxH - 4 // padding(2) + header(1) + footer(1)

	if !h.ready {
		h.viewport = viewport.New(contentW, contentH)
		h.viewport.SetContent(h.pages[h.page].content)
		h.ready = true
	} else {
		h.viewport.Width = contentW
		h.viewport.Height = contentH
	}

	header := ResultLabelStyle.Render(h.pages[h.page].title)

	var dots []string
	for i := range h.pages {
		if i == h.page {
			dots = append(dots, "●")
		} else {
			dots = append(dots, "○")
		}
	}
	nav := fmt.Sprintf("← %s → | j/k scroll | Esc close", strings.Join(dots, " "))
	footer := DimStyle.Render(nav)

	content := lipgloss.JoinVertical(lipgloss.Left, header, h.viewport.View(), footer)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(boxW).
		Height(boxH).
		Render(content)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
}

func buildHelpPages() []helpPage {
	return []helpPage{
		{title: "Quick Reference", content: pageQuickRef},
		{title: "How to Play", content: pageHowToPlay},
		{title: "The Oracle", content: pageOracle},
		{title: "Suit Domains", content: pageSuitDomains},
		{title: "GM Moves", content: pageGMMoves},
		{title: "Scene Management", content: pageScenes},
		{title: "Generators", content: pageGenerators},
		{title: "Commands", content: pageCommands},
		{title: "Tips & About", content: pageTips},
	}
}

var pageQuickRef = `KEYBOARD SHORTCUTS

  ORACLE                     SCENE
  [1] Yes/No (Likely)        [8] Set the Scene
  [2] Yes/No (Even)          [9] Random Event
  [3] Yes/No (Unlikely)
  [4] How                    GM MOVES
                             [0] Pacing Move
  FOCUS                      [-] Failure Move
  [5] Action Focus
  [6] Detail Focus           GENERATORS
  [7] Topic Focus            [=] Generic

  NAVIGATION                 COMMANDS
  Tab    Switch panel         /roll NdS   Roll dice
  Esc    Back to input        /flip [N]   Flip coins
  j/k    Scroll up/down       /draw [N]   Draw cards
  ?      Toggle this help     /weather    Weather
  Ctrl+S Save journal         /color      Random color
  Ctrl+R Saved rolls          /sound      Random sound
  Ctrl+P Portraits            /dir        Direction
  Ctrl+Q Quit                 /shuffle    Reshuffle deck
                              /portrait   Portrait browser
  Number shortcuts work       /scene      Set the Scene
  from the sidebar or
  log view (not while
  typing in the input).`

var pageHowToPlay = `HOW TO PLAY

One Page Solo Engine lets you play tabletop RPGs without a
Game Master. You write the story, and the engine provides the
unpredictability and prompts that a GM normally supplies.

GETTING STARTED
1. Create characters using your chosen game system.
2. Run a Plot Hook generator for your starting quest.
3. Run a Random Event [9] for an opening twist.
4. Set the Scene [8] to establish your first situation.
5. Start asking the Oracle [1-3] yes/no questions.

BASIC GAMEPLAY LOOP
• Describe what your character wants to do in the input area.
  Press Enter to add narrative text to your adventure log.
• When you need the "GM" to decide something, ask the Oracle
  a yes/no question [1-3] or use How [4] for degree.
• When you need inspiration, use Focus tables [5-7] to
  generate ideas. Interpret results using the suit domain.
• Use GM Moves [0/-] when the action stalls or PCs fail.
• Set the Scene [8] when moving to a new situation.

THE INTERFACE
• Input Area: Type narrative or /commands here. Press Enter.
• Log View: Scroll through your adventure. Tab to focus.
• Sidebar: Browse all generators. Tab to focus, Enter to run.

SAVING YOUR WORK
Press Ctrl+S to save. Your adventure is stored as a standard
Markdown file that looks great in any viewer or browser.`

var pageOracle = `THE ORACLE

The Oracle replaces the GM's judgment with dice rolls.

YES/NO ORACLE [1-3]
Choose a likelihood based on how probable "yes" is:

  Likely   [1] — Yes on 3+ (d6)
    Use when it would probably happen.
  Even     [2] — Yes on 4+ (d6)
    Use when it could go either way.
  Unlikely [3] — Yes on 5+ (d6)
    Use when it probably wouldn't happen.

A second d6 roll adds a modifier:
  1       → "but..."   Complicates the answer.
  2-5     → (none)     Straight answer.
  6       → "and..."   Amplifies the answer.

Example: "Is the door locked?" (Even)
  Answer: 5 → Yes.  Modifier: 6 → "and..."
  Result: "Yes, and the lock is enchanted."

HOW ORACLE [4]
Use when you need to know how much, how strong, etc.
  1: Surprisingly lacking    4: About average
  2: Less than expected      5: More than expected
  3: About average           6: Extraordinary

FOCUS TABLES [5-7]
For open-ended questions, draw a card from a Focus table:
  Action Focus [5] — What does it do?
  Detail Focus [6] — What kind of thing is it?
  Topic Focus  [7] — What is this about?

The card rank gives the table entry and the suit adds a
domain (see Suit Domains page) for interpretation.`

var pageSuitDomains = `SUIT DOMAINS

Every card drawn carries a suit that adds context. This is
the key to interpreting all card-based results.

  ♣ Clubs    — PHYSICAL
    Appearance, existence, the body, the tangible world.
    Examples: strength, damage, objects, terrain, looks.

  ♦ Diamonds — TECHNICAL
    Mental, operation, skill, knowledge, mechanisms.
    Examples: plans, devices, intelligence, craft, systems.

  ♠ Spades   — MYSTICAL
    Meaning, capability, magic, fate, the supernatural.
    Examples: destiny, spells, omens, power, the unknown.

  ♥ Hearts   — SOCIAL
    Personal, connection, emotion, relationships.
    Examples: loyalty, love, betrayal, reputation, bonds.

INTERPRETING RESULTS
Combine the table entry with the suit domain:

  "Harm" + ♥ Hearts = Social harm
    → Betrayal, insult, broken trust, emotional damage.

  "Harm" + ♣ Clubs = Physical harm
    → Violence, destruction, bodily injury.

  "Create" + ♦ Diamonds = Technical creation
    → Building a device, writing a plan, crafting a tool.

  "Create" + ♠ Spades = Mystical creation
    → Casting a spell, summoning, a prophetic vision.

JOKER
When you draw a Joker, shuffle the deck and add a Random
Event to the scene, then draw again for your original result.`

var pageGMMoves = `GM MOVES

GM Moves replace the Game Master's ability to advance the
story. Use them to keep the action moving and create drama.

PACING MOVES [0]
Use when there's a lull in the action, or "what now?"
  1: Foreshadow Trouble — Hint at something bad coming.
  2: Reveal a New Detail — Add information to the scene.
  3: An NPC Takes Action — An NPC does something notable.
  4: Advance a Threat    — Make a danger worse or closer.
  5: Advance a Plot      — Move a story arc forward.
  6: Random Event        — Something unexpected happens.

FAILURE MOVES [-]
Use when PCs fail a check and you want consequences:
  1: Cause Harm           — Someone gets hurt.
  2: Put Someone in a Spot — Create a difficult situation.
  3: Offer a Choice       — Present a tough decision.
  4: Advance a Threat     — Make the danger worse.
  5: Reveal Unwelcome Truth — Bad news surfaces.
  6: Foreshadow Trouble   — Hint at more problems.

WHEN TO USE THEM
Not every failed roll needs a GM Move. Use them when failure
has real consequences and the action needs to keep moving.

  Failed Spot check? Probably not.
  Failed to pick a lock while guards approach? Yes.
  Wondering what happens next? Pacing Move.
  Character failed a critical climb? Failure Move.

GM Moves should drive the story forward, not punish players.`

var pageScenes = `SCENE MANAGEMENT

SET THE SCENE [8]
Begin each new situation by Setting the Scene. Describe where
your character is and what they're trying to accomplish.

Every scene starts with a COMPLICATION (d6):
  1: Hostile forces oppose you
  2: An obstacle blocks your way
  3: Wouldn't it suck if...
  4: An NPC acts suddenly
  5: All is not as it seems
  6: Things actually go as planned

ALTERED SCENE
After the complication, a d6 is rolled. On 5+, the scene is
altered with an additional twist:
  1: A major detail is enhanced or somehow worse
  2: The environment is different
  3: Unexpected NPCs are present
  4: Add another Scene Complication
  5: Add a Pacing Move
  6: Add a Random Event

When an Altered Scene triggers another generator (options 4-6),
the result cascades automatically.

RANDOM EVENTS [9]
Generated by drawing two cards:
  What happens: from the Action Focus table
  Involving:    from the Topic Focus table

Combine both draws with their suit domains to create a
surprising event that makes sense in your story's context.

SCENE FLOW
1. Describe the situation → Set the Scene [8]
2. Play through the scene → Oracle, GM Moves as needed
3. Scene resolves → Set the Scene for next situation`

var pageGenerators = `GENERATORS

GENERIC [=]
Use for towns, factions, items, monsters, or anything.
  What it does:    Action Focus (card draw)
  How it looks:    Detail Focus (card draw)
  How significant: How Oracle (d6)

PLOT HOOK
Generates quests or missions for the PCs.
  Objective:  What needs to be done (d6)
  Adversary:  Who opposes you (d6)
  Reward:     What you'll gain (d6)

NPC
Generates a non-player character.
  Identity: Who they are (card)    Goal: What they want (card)
  Feature:  Notable trait (d6 + Detail Focus)
  Attitude: How they feel about you (How Oracle)
  Topic:    Conversation starter (Topic Focus)

DUNGEON THEME
Sets the tone for a dangerous location.
  How it looks: Detail Focus (card)
  How it's used: Action Focus (card)

DUNGEON ROOM
Generates a new area in a dungeon.
  Location (d6), Encounter (d6), Object (d6), Exits (d6)

HEX
Generates a hex for wilderness exploration.
  Terrain (d6), Contents (d6), Feature (d6), Event (d6)
  Events may trigger a Random Event and scene.`

var pageCommands = `SLASH COMMANDS

Type these in the input area and press Enter.

DICE ROLLER
  /roll NdS       Roll N dice with S sides.
  /roll NdS+M     Add a modifier (+ or -).
  /roll NdSkH     Roll N, keep highest H.
  /r NdS          Shorthand for /roll.

  Examples: /roll 2d6, /roll 4d6k3, /r 1d20+5

COINS
  /flip            Flip one coin.
  /flip N          Flip N coins.
  /f               Shorthand for /flip.

CARDS
  /draw            Draw one card from the utility deck.
  /draw N          Draw N cards.
  /card            Same as /draw.
  /shuffle         Reshuffle the utility deck.

DIRECTION
  /dir             Random 8-point compass direction.
  /dir N           N-point compass (4, 8, or 16).

ATMOSPHERE
  /weather         Random weather (condition, temp, wind).
  /color           Random color.
  /sound           Random sound effect.
  /sound CATEGORY  Sound from a category.
                   Categories: nature, urban, combat,
                   social, mechanical, animal, weather,
                   supernatural, domestic, musical

SCENE
  /scene           Same as Set the Scene [8].

CHARACTER VOICE & PORTRAITS
  /char NAME TEXT   Add a log entry attributed to NAME.
                   Example: /char Elara I search the room.
                   Shows "HH:MM  Elara" instead of "User".
  /portrait        Open the portrait browser (also Ctrl+P).
                   Generate random portraits, save favorites
                   with a name. Saved portraits display next
                   to /char dialogue matching that name.

SAVED ROLLS
  Ctrl+R           Open the saved rolls manager.
                   Create, organize, and execute saved
                   dice expressions from a modal dialog.`

var pageTips = `TIPS FOR BEST RESULTS

• Ask mostly yes/no questions to the Oracle.
• Loose interpretations are okay — the results are meant
  to inspire, not dictate.
• Always go with what's cool. If a result suggests something
  exciting, lean into it.
• If it doesn't make sense, try again. Don't force nonsensical
  results into your story.
• Use GM Moves to drive the action forward when things stall.
• Try group play with no GM — it works great with friends.

THE POWER OF INTERPRETATION
Results are meant to inspire ideas that make sense in the
context of your game. The answer should have meaning, not
just be a random detail. The result may be surprising, but
it should always be logical.

Give all results meaning. Embrace the unexpected.
Reject the nonsensical.

ABOUT
One Page Solo Engine v1.6
Created by Karl Hendricks — Inflatable Studios
License: CC-BY-SA 4.0

Inspired by:
  Mythic (Tana Pigeon)
  World vs Hero (John Fiore)
  Conjecture, UNE (Zach Best)
  Dungeon World (Koebel, LaTorra)
  The Black Hack (David Black)
  Maze Rats (Ben Milton)
  The Lone Wolf Solo RPG community`
