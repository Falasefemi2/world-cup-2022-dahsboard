package ui

import (
	"github.com/Falasefemi2/worldcupdashboard/data"
	"github.com/Falasefemi2/worldcupdashboard/ui/playerstats"
	tea "charm.land/bubbletea/v2"
)

type dataFetcher interface {
	GroupTables() ([]data.GroupTable, error)
	SortedMatches() ([]data.Match, error)
	Name() string
}

type dataFetchMsg struct {
	groupTablesByLetter map[string]data.GroupTable
	sortedMatches       []data.Match
	playerStatsByTeam   map[string]playerstats.PlayerStats
}

type dataFetchErrMsg struct{ err error }

func dataFetchCmd(fetcher dataFetcher) tea.Cmd {
	return func() tea.Msg {
		groupTables, err := fetcher.GroupTables()
		if err != nil {
			return dataFetchErrMsg{err: err}
		}
		groupTablesByLetter := make(map[string]data.GroupTable, len(groupTables))
		for _, g := range groupTables {
			groupTablesByLetter[g.Letter] = g
		}

		sortedMatches, err := fetcher.SortedMatches()
		if err != nil {
			return dataFetchErrMsg{err: err}
		}

		playerStatsByTeam := playerstats.PlayerStatsByTeam(sortedMatches)

		return dataFetchMsg{groupTablesByLetter, sortedMatches, playerStatsByTeam}
	}
}
