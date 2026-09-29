package data

import "time"

type EventType string

const (
	EventTypeGoal                = "Goal"
	EventTypePenaltyKickGoal     = "Goal (P)"
	EventTypePenaltyShootoutGoal = "Goal (S)"
	EventTypeOwnGoal             = "Own Goal"
	EventTypeYellowCard          = "Yellow Card"
	EventTypeSeconYellowCard     = "Second Yellow Card"
	EventTypeRedCard             = "Red Card"
	EventTypeSubIn               = "Substitution In"
	EventTypeSubOut              = "Substitution Out"
)

type Stage string

const (
	StageGroup   Stage = "Group"
	StageLast16  Stage = "1/8"
	StageQuarter Stage = "1/4"
	StageSemi    Stage = "1/2"
	StageThird   Stage = "3rd"
	StageFinal   Stage = "Final"
)

type Status string

const (
	StatusScheduled = "Scheduled"
	StatusLive      = "Live"
	StatusFinished  = "Finished"
)

type GroupTable struct {
	Letter string
	Table  []GroupTableTeam
}

type GroupTableTeam struct {
	Code              string
	Points            int
	Wins              int
	Draws             int
	Losses            int
	MatchesPlayed     int
	GoalsFor          int
	GoalsAgainst      int
	GoalsDifferential int
}

type Match struct {
	ID             int
	HomeTeamCode   string
	AwayTeamCode   string
	Date           time.Time
	Venue          string
	HomeTeamScore  uint64
	AwayTeamScore  uint64
	WinnerTeamCode string
	Minute         string
	HomeTeamEvents []Event
	AwayTeamEvents []Event
	Status         Status
	HomeTeamLineup []Player
	AwayTeamLineup []Player

	// Stage is a Stage when the type is known - otherwise a string
	Stage string
}

type Event struct {
	// Type is a EventType when the type is known - otherwise a string
	Type string

	Minute   string
	Player   string
	Canceled bool
}

type TeamInfo struct {
	Name        string
	Group       string
	FirstColor  string
	SecondColor string
}

type Player struct {
	Name        string
	ShirtNumber int
}
