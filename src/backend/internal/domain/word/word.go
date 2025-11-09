package word

type Word struct {
	value        string
	Position     Position
	QuestWordIDx int64
}

type Position struct {
	Book string
	Page string
	Line string
	Word string
}

func NewWord(
	value string,
	position Position,
	questWordIDx int64,
) Word {
	return Word{
		value:        value,
		Position:     position,
		QuestWordIDx: questWordIDx,
	}
}

func (word Word) CompareValue(value string) bool {
	return word.value == value
}
