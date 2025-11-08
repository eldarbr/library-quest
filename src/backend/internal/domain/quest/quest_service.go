package quest

var (
	questIDPattern = []int{2, 3, 5, 9, 0, 7, 6, 1, 4, 8}
	minQuestCnt    = len(questIDPattern)
)

func DetermineQuestForTeam(teamID int) int {
	return questIDPattern[teamID%len(questIDPattern)]
}

func GetMinQuestCnt() int {
	return minQuestCnt
}
