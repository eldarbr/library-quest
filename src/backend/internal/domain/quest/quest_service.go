package quest

import "math/rand"

type QuestShuffler interface {
	DetermineQuestForTeam(teamID int) int
}

type RandomQuestShuffler interface {
	GetRandomQuestID() int
}

type DeterminedTeamToQuestShuffler struct {
	questIDPattern []int
	seed           int64
	rngd           *rand.Rand
}

func NewDeterminedTeamToQuestShuffler(numQuests int, seed int64) DeterminedTeamToQuestShuffler {
	quests := make([]int, 0, numQuests)
	for i := range numQuests {
		quests = append(quests, i)
	}

	rngd := rand.New(rand.NewSource(seed))
	rngd.Shuffle(numQuests, func(i, j int) {
		quests[i], quests[j] = quests[j], quests[i]
	})

	return DeterminedTeamToQuestShuffler{questIDPattern: quests, seed: seed, rngd: rngd}
}

func (dttqs DeterminedTeamToQuestShuffler) DetermineQuestForTeam(teamID int) int {
	return dttqs.questIDPattern[teamID%len(dttqs.questIDPattern)]
}

func (dttqs DeterminedTeamToQuestShuffler) GetRandomQuestID() int {
	return dttqs.rngd.Intn(len(dttqs.questIDPattern))
}

type ConstantKeyworder struct {
	Keyword string
}

func (kw ConstantKeyworder) GetKeyword() string {
	return kw.Keyword
}
