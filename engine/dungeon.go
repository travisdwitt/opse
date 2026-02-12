package engine

var dungeonLocations = [7]string{
	"",
	"Typical area",
	"Transitional area",
	"Living area or meeting place",
	"Working or utility area",
	"Area with a special feature",
	"Location for a specialized purpose",
}

func dungeonEncounter(roll int) string {
	switch {
	case roll <= 2:
		return "None"
	case roll <= 4:
		return "Hostile enemies"
	case roll == 5:
		return "An obstacle blocks the way"
	default:
		return "Unique NPC or adversary"
	}
}

func dungeonObject(roll int) string {
	switch {
	case roll <= 2:
		return "Nothing, or mundane objects"
	case roll == 3:
		return "An interesting item or clue"
	case roll == 4:
		return "A useful tool, key, or device"
	case roll == 5:
		return "Something valuable"
	default:
		return "Rare or special item"
	}
}

func dungeonExits(roll int) string {
	switch {
	case roll <= 2:
		return "Dead end"
	case roll <= 4:
		return "1 additional exit"
	default:
		return "2 additional exits"
	}
}

func DungeonTheme(deck *Deck) DungeonThemeResult {
	return DungeonThemeResult{
		Looks: DetailFocus(deck),
		Used:  ActionFocus(deck),
	}
}

func DungeonRoom(rng *Randomizer) DungeonRoomResult {
	lr := rng.RollD6()
	er := rng.RollD6()
	or_ := rng.RollD6()
	xr := rng.RollD6()
	return DungeonRoomResult{
		LocationRoll: lr, Location: dungeonLocations[lr],
		EncounterRoll: er, Encounter: dungeonEncounter(er),
		ObjectRoll: or_, Object: dungeonObject(or_),
		ExitsRoll: xr, Exits: dungeonExits(xr),
	}
}
