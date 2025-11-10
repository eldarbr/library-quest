package filedb

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/eldarbr/library-quest/backend/internal/domain/quest"
	"github.com/eldarbr/library-quest/backend/internal/domain/word"
)

type FileToMemoDB struct {
	quests map[int]quest.Quest
}

func NewFileToMemoDB(fp string) (*FileToMemoDB, error) {
	data, err := readData(fp)
	if err != nil {
		return nil, fmt.Errorf("read data: %w", err)
	}

	return &FileToMemoDB{quests: data}, nil
}

func readData(fp string) (map[int]quest.Quest, error) {
	basePath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get pwd: %w", err)
	}

	fp = filepath.Join(basePath, filepath.Clean(fp))
	if !strings.HasPrefix(fp, basePath) {
		return nil, fmt.Errorf("invalid path")
	}

	f, err := os.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("open fp: %w", err)
	}

	defer func() {
		fErr := f.Close()
		if fErr != nil {
			slog.Error("close file read database", slog.Any("err", fErr))
		}
	}()

	var (
		scan = bufio.NewScanner(f)

		wordObj  = word.Word{}
		questID  int
		dataRead = make(map[int][]word.Word)
	)

	for linesRead := 0; ; linesRead++ {
		if !scan.Scan() {
			break
		}

		lineTokens := strings.Split(scan.Text(), tokenSeparator)

		questID, wordObj, err = lineTokensToWord(lineTokens)
		if err != nil {
			return nil, fmt.Errorf("line %d to quest: %w", linesRead, err)
		}

		dataRead[questID] = append(dataRead[questID], wordObj)
	}

	if len(dataRead) < quest.GetMinQuestCnt() {
		return nil, ErrNotEnoughQuests
	}

	return readDataToMap(dataRead), nil
}

func readDataToMap(data map[int][]word.Word) map[int]quest.Quest {
	result := make(map[int]quest.Quest, len(data))

	for questID := range data {
		slices.SortFunc(data[questID], func(a, b word.Word) int {
			return int(a.QuestWordIDx - b.QuestWordIDx)
		})

		que := quest.NewQuest(questID, data[questID])
		result[questID] = que
	}

	return result
}

func lineTokensToWord(lineTokens []string) (int, word.Word, error) {
	if len(lineTokens) != tokensPerLine {
		return 0, word.Word{}, ErrTokensPerLine
	}

	questID, err := strconv.Atoi(lineTokens[questIDToken])
	if err != nil {
		return 0, word.Word{}, fmt.Errorf("parse quest id: %w", err)
	}

	questWordIDX, err := strconv.Atoi(lineTokens[questWordIDXToken])
	if err != nil {
		return 0, word.Word{}, fmt.Errorf("parse quest word idx: %w", err)
	}

	return questID,
		word.NewWord(
			lineTokens[valueToken],
			word.Position{
				Book: lineTokens[bookToken],
				Page: lineTokens[pageToken],
				Line: lineTokens[lineToken],
				Word: lineTokens[wordInLineToken],
			},
			int64(questWordIDX),
		),
		nil
}
