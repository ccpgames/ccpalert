package ccpalertql

import (
	"strconv"
	"strings"
)

const (
	//ALERT keyword
	ALERT TokenType = iota
	//IDENTIFIER is a string token for rule IDs and keys
	IDENTIFIER
	//STRING is unstructured text enclosed in quotes
	STRING
	//OP is a token representing one of <,> or ==
	OP
	//IF keyworkd
	IF
	//TEXT keyword
	TEXT
	//ILLEGAL token
	ILLEGAL
	//VALUE is a float value
	VALUE
)

type (
	// TokenType denote the type of a lexical token.
	TokenType int

	//Token denotes a lexical token
	Token struct {
		literal   string
		tokenType TokenType
	}

	//Scanner is a scanner
	Scanner struct {
		pos      int
		runes    []string
		hasRunes bool
	}
)

//NewScanner returns an instance of scanner
func NewScanner(rawString string) *Scanner {
	return &Scanner{pos: 0, runes: strings.Fields(rawString), hasRunes: len(rawString) > 0}
}
func (s *Scanner) next() {
	s.pos++
	if s.pos >= len(s.runes) {
		s.hasRunes = false
	}
}

func (s *Scanner) previous() {
	s.pos--
}

func (s *Scanner) readRune() string {
	return s.runes[s.pos]
}

func (s *Scanner) scan() []Token {
	var tokens []Token

	for s.hasRunes {
		rune := s.readRune()
		if rune == "<" || rune == ">" || rune == "==" {
			tokens = append(tokens, scanOP(rune))
		} else {
			switch strings.ToUpper(rune) {
			case "ALERT":
				alertToken := new(Token)
				alertToken.literal = rune
				alertToken.tokenType = ALERT
				tokens = append(tokens, *alertToken)
			case "IF":
				ifToken := new(Token)
				ifToken.literal = rune
				ifToken.tokenType = IF
				tokens = append(tokens, *ifToken)
			case "TEXT":
				textToken := new(Token)
				textToken.literal = rune
				textToken.tokenType = TEXT
				tokens = append(tokens, *textToken)
			default:

				if rune[0] == '"' {
					tokens = append(tokens, s.scanString())
				} else if _, err := strconv.ParseFloat(rune, 64); err == nil {
					valueToken := new(Token)
					valueToken.literal = rune
					valueToken.tokenType = VALUE
					tokens = append(tokens, *valueToken)

				} else {
					idToken := new(Token)
					idToken.literal = rune
					idToken.tokenType = IDENTIFIER
					tokens = append(tokens, *idToken)
				}
			}
		}
		s.next()
	}
	return tokens
}

func (s *Scanner) scanString() Token {
	var fullString []string
	for s.hasRunes {
		rune := s.readRune()
		fullString = append(fullString, rune)
		if rune[len(rune)-1] == '"' {
			stringToken := new(Token)
			stringToken.literal = strings.Join(fullString, "")
			stringToken.tokenType = STRING
			return *stringToken
		}
		s.next()
	}
	illegalToken := new(Token)
	illegalToken.literal = strings.Join(fullString, "")
	illegalToken.tokenType = ILLEGAL
	return *illegalToken
}

func scanOP(s string) Token {
	opToken := new(Token)
	opToken.literal = s
	opToken.tokenType = OP
	return *opToken
}
