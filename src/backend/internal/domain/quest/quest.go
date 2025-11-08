package quest

import (
	"errors"
	"fmt"

	"github.com/eldarbr/library-quest/backend/internal/domain/word"
)

type Quest struct {
	id    int
	words []word.Word
}

func NewQuest(id int, words []word.Word) Quest {
	return Quest{
		id:    id,
		words: words,
	}
}

func (quest Quest) GetID() int {
	return quest.id
}

func (quest Quest) GetWords() []word.Word {
	return quest.words
}

func (quest Quest) ValidateAnswer(answer []string) error {
	if len(quest.words) != len(answer) {
		return ErrWordsCntMismatch
	}

	wae := WrongAnswerError{
		Mistakes: []int64{},
	}

	for i, word := range quest.words {
		if !word.CompareValue(answer[i]) {
			wae.Mistakes = append(wae.Mistakes, int64(i))
		}
	}

	if len(wae.Mistakes) != 0 {
		return wae
	}

	return nil
}

type WrongAnswerError struct {
	Mistakes []int64
}

func (wae WrongAnswerError) Error() string {
	return fmt.Sprintf("wrong answer, mistakes: %v", wae.Mistakes)
}

var (
	ErrWordsCntMismatch = errors.New("number of words did not match")
)
