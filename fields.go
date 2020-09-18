//
// Split quoted comma separated list
//

package fields

import (
	"fmt"
	"io"
	"strings"
)

type State_t int

const (
	STATE_NONE               State_t = 0
	STATE_OPEN_QUOTE         State_t = 1
	STATE_CLOSE_QUOTE        State_t = 2
	STATE_SEPARATOR          State_t = 3
	STATE_TRIM               State_t = 4
	STATE_STRING             State_t = 5
	STATE_EOF                State_t = 6
	STATE_ERROR_NO_QUOTE     State_t = 7
	STATE_ERROR_NO_SEPARATOR State_t = 8
)

type NextState func(State_t) (NextState, State_t)

type Lexer_t struct {
	sep        map[rune]rune
	trim       map[rune]rune
	open_quote map[rune]rune

	reader      io.RuneReader
	state       NextState
	last_rune   rune
	last_size   int
	close_quote []rune
	last_token  strings.Builder
	last_trim   strings.Builder
}

type Quote_t struct {
	Open  rune
	Close rune
}

func NewLexer(sep []rune, trim []rune, quote []Quote_t) (self *Lexer_t) {
	self = &Lexer_t{
		sep:        map[rune]rune{},
		trim:       map[rune]rune{},
		open_quote: map[rune]rune{},
	}
	for _, v := range sep {
		self.sep[v] = v
	}
	for _, v := range trim {
		self.trim[v] = v
	}
	for _, v := range quote {
		self.open_quote[v.Open] = v.Close
	}
	return
}

func (self *Lexer_t) begin(prev State_t) (NextState, State_t) {
	self.last_rune, self.last_size, _ = self.reader.ReadRune()
	switch {
	case self.open_quote[self.last_rune] > 0:
		self.close_quote = append(self.close_quote, self.open_quote[self.last_rune])
		return self.quoted, STATE_OPEN_QUOTE
	case self.sep[self.last_rune] > 0:
		return self.begin, STATE_SEPARATOR
	case self.trim[self.last_rune] > 0:
		return self.begin, STATE_TRIM
	case self.last_size > 0:
		self.last_token.WriteRune(self.last_rune)
		return self.not_quoted, STATE_STRING
	default:
		return nil, STATE_EOF
	}
}

func (self *Lexer_t) not_quoted(prev State_t) (NextState, State_t) {
	self.last_rune, self.last_size, _ = self.reader.ReadRune()
	switch {
	case self.sep[self.last_rune] > 0:
		self.last_trim.Reset()
		return self.begin, STATE_SEPARATOR
	case self.trim[self.last_rune] > 0:
		self.last_trim.WriteRune(self.last_rune)
		return self.not_quoted, STATE_TRIM
	case self.last_size > 0:
		if self.last_trim.Len() > 0 {
			self.last_token.WriteString(self.last_trim.String())
			self.last_trim.Reset()
		}
		self.last_token.WriteRune(self.last_rune)
		return self.not_quoted, STATE_STRING
	default:
		return nil, STATE_EOF
	}
}

func (self *Lexer_t) quoted(prev State_t) (NextState, State_t) {
	self.last_rune, self.last_size, _ = self.reader.ReadRune()
	q_len := len(self.close_quote)
	switch {
	case self.close_quote[q_len-1] == self.last_rune:
		self.close_quote = self.close_quote[:q_len-1]
		if q_len == 1 {
			return self.separator, STATE_CLOSE_QUOTE
		}
		return self.quoted, STATE_CLOSE_QUOTE
	case prev == STATE_OPEN_QUOTE && self.open_quote[self.last_rune] > 0:
		self.close_quote = append(self.close_quote, self.open_quote[self.last_rune])
		return self.quoted, STATE_OPEN_QUOTE
	case prev != STATE_CLOSE_QUOTE && self.last_size > 0:
		self.last_token.WriteRune(self.last_rune)
		return self.quoted, STATE_STRING
	default:
		return nil, STATE_ERROR_NO_QUOTE
	}
}

func (self *Lexer_t) separator(prev State_t) (NextState, State_t) {
	self.last_rune, self.last_size, _ = self.reader.ReadRune()
	switch {
	case self.sep[self.last_rune] > 0:
		self.last_trim.Reset()
		return self.begin, STATE_SEPARATOR
	case self.trim[self.last_rune] > 0:
		self.last_trim.WriteRune(self.last_rune)
		return self.not_quoted, STATE_TRIM
	case self.last_size == 0:
		return nil, STATE_EOF
	default:
		return nil, STATE_ERROR_NO_SEPARATOR
	}
}

func (self *Lexer_t) Set(in io.RuneReader) {
	self.reader = in
	self.state = self.begin
}

func (self *Lexer_t) Next() (token string, state State_t) {
	for self.state != nil {
		self.state, state = self.state(state)
		if state == STATE_SEPARATOR || state >= STATE_EOF {
			token = self.last_token.String()
			self.last_token.Reset()
			return
		}
	}
	return
}

func Split(in string, sep ...rune) (res []string, err error) {
	l := NewLexer(sep,
		[]rune{'\v', '\f', '\r', '\n', '\t', ' '},
		[]Quote_t{
			{'"', '"'},
			{'\'', '\''},
			{'«', '»'},
		},
	)
	l.Set(strings.NewReader(in))
	for {
		token, state := l.Next()
		if state == STATE_NONE {
			return
		}
		res = append(res, token)
		if state > STATE_EOF {
			err = fmt.Errorf("ERROR: %v", state)
			return
		}
	}
}

type Strings_t []string

func (self *Strings_t) Set(value string) (err error) {
	var temp []string
	if temp, err = Split(value, ','); err == nil {
		for _, v := range temp {
			*self = append(*self, v)
		}
	}
	return
}

func (self *Strings_t) String() string {
	return strings.Join(*self, ",")
}
