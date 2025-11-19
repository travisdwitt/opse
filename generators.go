package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Card represents a playing card
type Card struct {
	Rank int    // 2-14 (2-10, J=11, Q=12, K=13, A=14)
	Suit string // Clubs, Diamonds, Spades, Hearts
}

// SuitDomain maps suits to their domains
var SuitDomain = map[string]string{
	"Clubs":    "Physical (appearance, existence)",
	"Diamonds": "Technical (mental, operation)",
	"Spades":   "Mystical (meaning, capability)",
	"Hearts":   "Social (personal, connection)",
}

// DrawCard draws a random card
func DrawCard() Card {
	rand.Seed(time.Now().UnixNano())
	rank := rand.Intn(13) + 2 // 2-14
	suits := []string{"Clubs", "Diamonds", "Spades", "Hearts"}
	suit := suits[rand.Intn(4)]
	return Card{Rank: rank, Suit: suit}
}

// RollD6 rolls a d6
func RollD6() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(6) + 1
}

// Roll2D6 rolls 2d6
func Roll2D6() (int, int) {
	return RollD6(), RollD6()
}

// RollD4 rolls a d4
func RollD4() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(4) + 1
}

// RollD12 rolls a d12
func RollD12() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(12) + 1
}

// RankToString converts rank number to string
func RankToString(rank int) string {
	switch rank {
	case 11:
		return "J"
	case 12:
		return "Q"
	case 13:
		return "K"
	case 14:
		return "A"
	default:
		return fmt.Sprintf("%d", rank)
	}
}

// OracleYesNo provides a yes/no answer based on likelihood
func OracleYesNo(likelihood string) (bool, string) {
	d1, d2 := Roll2D6()
	total := d1 + d2
	
	var threshold int
	switch likelihood {
	case "Likely":
		threshold = 3
	case "Even":
		threshold = 4
	case "Unlikely":
		threshold = 5
	default:
		threshold = 4
	}
	
	yes := total >= threshold
	mod := RollD6()
	
	var modText string
	if mod == 1 {
		modText = "but…"
	} else if mod == 6 {
		modText = "and…"
	}
	
	result := fmt.Sprintf("Rolled %d+%d = %d (threshold: %d+)", d1, d2, total, threshold)
	if yes {
		result += " → Yes"
	} else {
		result += " → No"
	}
	if modText != "" {
		result += " " + modText
	}
	
	return yes, result
}

// OracleHow provides a "how" answer
func OracleHow() string {
	roll := RollD6()
	switch roll {
	case 1:
		return "Surprisingly lacking"
	case 2:
		return "Less than expected"
	case 3, 4:
		return "About average"
	case 5:
		return "More than expected"
	case 6:
		return "Extraordinary"
	}
	return ""
}

// SceneComplication generates a scene complication
func SceneComplication() string {
	roll := RollD6()
	complications := map[int]string{
		1: "Hostile forces oppose you",
		2: "An obstacle blocks your way",
		3: "Wouldn't it suck if…",
		4: "An NPC acts suddenly",
		5: "All is not as it seems",
		6: "Things actually go as planned",
	}
	return complications[roll]
}

// AlteredScene generates an altered scene result
func AlteredScene() string {
	roll := RollD6()
	altered := map[int]string{
		1: "A major detail of the scene is enhanced or somehow worse",
		2: "The environment is different",
		3: "Unexpected NPCs are present",
		4: "Add a SCENE COMPLICATION",
		5: "Add a PACING MOVE",
		6: "Add a RANDOM EVENT",
	}
	return altered[roll]
}

// PacingMove generates a pacing move
func PacingMove() string {
	roll := RollD6()
	moves := map[int]string{
		1: "Foreshadow Trouble",
		2: "Reveal a New Detail",
		3: "An NPC Takes Action",
		4: "Advance a Threat",
		5: "Advance a Plot",
		6: "Add a RANDOM EVENT to the scene",
	}
	return moves[roll]
}

// FailureMove generates a failure move
func FailureMove() string {
	roll := RollD6()
	moves := map[int]string{
		1: "Cause Harm",
		2: "Put Someone in a Spot",
		3: "Offer a Choice",
		4: "Advance a Threat",
		5: "Reveal an Unwelcome Truth",
		6: "Foreshadow Trouble",
	}
	return moves[roll]
}

// ActionFocus generates an action focus
func ActionFocus() (Card, string) {
	card := DrawCard()
	actions := map[int]string{
		2:  "Seek",
		3:  "Oppose",
		4:  "Communicate",
		5:  "Move",
		6:  "Harm",
		7:  "Create",
		8:  "Reveal",
		9:  "Command",
		10: "Take",
		11: "Protect",
		12: "Assist",
		13: "Transform",
		14: "Deceive",
	}
	action := actions[card.Rank]
	return card, action
}

// DetailFocus generates a detail focus
func DetailFocus() (Card, string) {
	card := DrawCard()
	details := map[int]string{
		2:  "Small",
		3:  "Large",
		4:  "Old",
		5:  "New",
		6:  "Mundane",
		7:  "Simple",
		8:  "Complex",
		9:  "Unsavory",
		10: "Specialized",
		11: "Unexpected",
		12: "Exotic",
		13: "Dignified",
		14: "Unique",
	}
	detail := details[card.Rank]
	return card, detail
}

// TopicFocus generates a topic focus
func TopicFocus() (Card, string) {
	card := DrawCard()
	topics := map[int]string{
		2:  "Current Need",
		3:  "Allies",
		4:  "Community",
		5:  "History",
		6:  "Future Plans",
		7:  "Enemies",
		8:  "Knowledge",
		9:  "Rumors",
		10: "A Plot Arc",
		11: "Recent Events",
		12: "Equipment",
		13: "A Faction",
		14: "The PCs",
	}
	topic := topics[card.Rank]
	return card, topic
}

// RandomEvent generates a random event
func RandomEvent() string {
	actionCard, action := ActionFocus()
	topicCard, topic := TopicFocus()
	
	actionDomain := SuitDomain[actionCard.Suit]
	topicDomain := SuitDomain[topicCard.Suit]
	
	return fmt.Sprintf("What happens: %s (%s - %s)\nInvolving: %s (%s - %s)",
		action, RankToString(actionCard.Rank)+actionCard.Suit[0:1], actionDomain,
		topic, RankToString(topicCard.Rank)+topicCard.Suit[0:1], topicDomain)
}

// PlotHook generates a plot hook
func PlotHook() string {
	objectiveRoll := RollD6()
	adversaryRoll := RollD6()
	rewardRoll := RollD6()
	
	objectives := map[int]string{
		1: "Eliminate a threat",
		2: "Learn the truth",
		3: "Recover something valuable",
		4: "Escort or deliver to safety",
		5: "Restore something broken",
		6: "Save an ally in peril",
	}
	
	adversaries := map[int]string{
		1: "A powerful organization",
		2: "Outlaws",
		3: "Guardians",
		4: "Local inhabitants",
		5: "Enemy horde or force",
		6: "A new or recurring villain",
	}
	
	rewards := map[int]string{
		1: "Money or valuables",
		2: "Money or valuables",
		3: "Knowledge and secrets",
		4: "Support of an ally",
		5: "Advance a plot arc",
		6: "A unique item of power",
	}
	
	return fmt.Sprintf("OBJECTIVE: %s\nADVERSARIES: %s\nREWARDS: %s",
		objectives[objectiveRoll], adversaries[adversaryRoll], rewards[rewardRoll])
}

// NPCGenerator generates an NPC
func NPCGenerator() string {
	identityCard, _ := identityFocus()
	goalCard, _ := goalFocus()
	notableRoll := RollD6()
	_, detail := DetailFocus()
	attitudeRoll := RollD6()
	_, topic := TopicFocus()
	
	identities := map[int]string{
		2:  "Outlaw",
		3:  "Drifter",
		4:  "Tradesman",
		5:  "Commoner",
		6:  "Soldier",
		7:  "Merchant",
		8:  "Specialist",
		9:  "Entertainer",
		10: "Adherent",
		11: "Leader",
		12: "Mystic",
		13: "Adventurer",
		14: "Lord",
	}
	
	goals := map[int]string{
		2:  "Obtain",
		3:  "Learn",
		4:  "Harm",
		5:  "Restore",
		6:  "Find",
		7:  "Travel",
		8:  "Protect",
		9:  "Enrich Self",
		10: "Avenge",
		11: "Fulfill Duty",
		12: "Escape",
		13: "Create",
		14: "Serve",
	}
	
	notableFeatures := map[int]string{
		1: "Unremarkable",
		2: "Notable nature",
		3: "Obvious physical trait",
		4: "Quirk or mannerism",
		5: "Unusual equipment",
		6: "Unexpected age or origin",
	}
	
	attitudes := map[int]string{
		1: "Surprisingly lacking",
		2: "Less than expected",
		3: "About average",
		4: "More than expected",
		5: "Extraordinary",
		6: "Extraordinary",
	}
	
	identityName := identities[identityCard.Rank]
	goalName := goals[goalCard.Rank]
	
	return fmt.Sprintf("IDENTITY: %s\nGOAL: %s\nNOTABLE FEATURE: %s (%s)\nATTITUDE TO PCs: %s\nCONVERSATION: %s",
		identityName, goalName, notableFeatures[notableRoll], detail, attitudes[attitudeRoll], topic)
}

func identityFocus() (Card, string) {
	card := DrawCard()
	return card, ""
}

func goalFocus() (Card, string) {
	card := DrawCard()
	return card, ""
}

// GenericGenerator generates a generic thing
func GenericGenerator() string {
	actionCard, action := ActionFocus()
	detailCard, detail := DetailFocus()
	how := OracleHow()
	
	return fmt.Sprintf("What it does: %s (%s - %s)\nHow it looks: %s (%s - %s)\nHow significant: %s",
		action, RankToString(actionCard.Rank)+actionCard.Suit[0:1], SuitDomain[actionCard.Suit],
		detail, RankToString(detailCard.Rank)+detailCard.Suit[0:1], SuitDomain[detailCard.Suit],
		how)
}

// DungeonCrawler generates dungeon content
func DungeonCrawler() string {
	themeDetailCard, themeDetail := DetailFocus()
	themeActionCard, themeAction := ActionFocus()
	
	locationRoll := RollD6()
	encounterRoll := RollD6()
	objectRoll := RollD6()
	exitsRoll := RollD6()
	
	locations := map[int]string{
		1: "Typical area",
		2: "Transitional area",
		3: "Living area or meeting place",
		4: "Working or utility area",
		5: "Area with a special feature",
		6: "Location for a specialized purpose",
	}
	
	encounters := map[int]string{
		1: "None",
		2: "None",
		3: "Hostile enemies",
		4: "Hostile enemies",
		5: "An obstacle blocks the way",
		6: "Unique NPC or adversary",
	}
	
	objects := map[int]string{
		1: "Nothing, or mundane objects",
		2: "Nothing, or mundane objects",
		3: "An interesting item or clue",
		4: "A useful tool, key, or device",
		5: "Something valuable",
		6: "Rare or special item",
	}
	
	exits := map[int]string{
		1: "Dead end",
		2: "Dead end",
		3: "1 additional exit",
		4: "1 additional exit",
		5: "2 additional exits",
		6: "2 additional exits",
	}
	
	return fmt.Sprintf("DUNGEON THEME:\n  How it looks: %s (%s)\n  How it is used: %s (%s)\n\nLOCATION: %s\nENCOUNTER: %s\nOBJECT: %s\nTOTAL EXITS: %s",
		themeDetail, RankToString(themeDetailCard.Rank)+themeDetailCard.Suit[0:1],
		themeAction, RankToString(themeActionCard.Rank)+themeActionCard.Suit[0:1],
		locations[locationRoll], encounters[encounterRoll], objects[objectRoll], exits[exitsRoll])
}

// HexCrawler generates hex crawl content
func HexCrawler() string {
	terrainRoll := RollD6()
	contentsRoll := RollD6()
	eventRoll := RollD6()
	
	terrains := map[int]string{
		1: "Same as current hex",
		2: "Same as current hex",
		3: "Common terrain",
		4: "Common terrain",
		5: "Uncommon terrain",
		6: "Rare terrain",
	}
	
	contents := map[int]string{
		1: "Nothing notable",
		2: "Nothing notable",
		3: "Nothing notable",
		4: "Nothing notable",
		5: "Nothing notable",
		6: "Roll a FEATURE",
	}
	
	events := map[int]string{
		1: "None",
		2: "None",
		3: "None",
		4: "None",
		5: "RANDOM EVENT then SET THE SCENE",
		6: "RANDOM EVENT then SET THE SCENE",
	}
	
	result := fmt.Sprintf("TERRAIN: %s\nCONTENTS: %s\nEVENT: %s",
		terrains[terrainRoll], contents[contentsRoll], events[eventRoll])
	
	if contentsRoll == 6 {
		featureRoll := RollD6()
		features := map[int]string{
			1: "Notable structure",
			2: "Dangerous hazard",
			3: "A settlement",
			4: "Strange natural feature",
			5: "New region (set new terrain types)",
			6: "DUNGEON CRAWLER entrance",
		}
		result += fmt.Sprintf("\nFEATURE: %s", features[featureRoll])
	}
	
	return result
}

