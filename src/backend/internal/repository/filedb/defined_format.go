package filedb

import (
	"errors"
)

// Columns separated by tabs:
// - quest id        number
// - quest word idx  number
// - value           string
// - book            string
// - page            string
// - line            string
// - word in line    string

const (
	tokenSeparator = "\t"
	tokensPerLine  = 7
)
const (
	questIDToken = iota
	questWordIDXToken
	valueToken
	bookToken
	pageToken
	lineToken
	wordInLineToken
)

var (
	ErrTokensPerLine = errors.New("wrong amount of tokens")
)
