package engine

type OracleYesNoResult struct {
	Likelihood string
	AnswerRoll int
	Answer     bool
	ModRoll    int
	Modifier   string
}

type OracleHowResult struct {
	Roll   int
	Result string
}

type CardTableResult struct {
	Draw      DrawResult
	TableName string
	Entry     string
}

type SceneComplicationResult struct {
	Roll   int
	Result string
}

type AlteredSceneResult struct {
	Roll    int
	Result  string
	Cascade interface{}
}

type SetTheSceneResult struct {
	Complication SceneComplicationResult
	AlteredRoll  int
	Altered      bool
	AlteredScene *AlteredSceneResult
}

type PacingMoveResult struct {
	Roll        int
	Result      string
	RandomEvent *RandomEventResult
}

type FailureMoveResult struct {
	Roll   int
	Result string
}

type RandomEventResult struct {
	Action CardTableResult
	Topic  CardTableResult
}

type GenericGeneratorResult struct {
	Action       CardTableResult
	Detail       CardTableResult
	Significance OracleHowResult
}

type PlotHookResult struct {
	ObjectiveRoll int
	Objective     string
	AdversaryRoll int
	Adversary     string
	RewardRoll    int
	Reward        string
}

type NPCResult struct {
	Identity      CardTableResult
	Goal          CardTableResult
	FeatureRoll   int
	Feature       string
	FeatureDetail CardTableResult
	Attitude      OracleHowResult
	Topic         CardTableResult
}

type DungeonThemeResult struct {
	Looks CardTableResult
	Used  CardTableResult
}

type DungeonRoomResult struct {
	LocationRoll  int
	Location      string
	EncounterRoll int
	Encounter     string
	ObjectRoll    int
	Object        string
	ExitsRoll     int
	Exits         string
}

type HexResult struct {
	TerrainRoll  int
	Terrain      string
	ContentsRoll int
	Contents     string
	FeatureRoll  int
	Feature      string
	EventRoll    int
	Event        string
	RandomEvent  *RandomEventResult
}

type DiceExpression struct {
	Count    int
	Sides    int
	Explode  bool
	KeepMode string
	KeepN    int
	Modifier int
	Raw      string
}

type DiceRollResult struct {
	Expression DiceExpression
	Rolls      []int
	Kept       []bool
	Subtotal   int
	Total      int
}

type CoinFlipResult struct {
	Flips []bool
	Heads int
	Tails int
}

type CardDrawResult struct {
	Cards     []Card
	Remaining int
}

type DirectionResult struct {
	Direction string
	Abbrev    string
	Arrow     string
}

type WeatherResult struct {
	Condition   string
	Temperature string
	Wind        string
}

type ColorResult struct {
	Color string
}

type SoundResult struct {
	Sound    string
	Category string
}
