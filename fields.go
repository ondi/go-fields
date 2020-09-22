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
	STATE_NEW_LINE           State_t = 4
	STATE_TRIM               State_t = 5
	STATE_DATA               State_t = 6
	STATE_EOF                State_t = 7
	STATE_ERROR_NO_QUOTE     State_t = 8
	STATE_ERROR_NO_SEPARATOR State_t = 9
)

type NextState func(rune, int, State_t) (NextState, State_t)

type Lexer_t struct {
	sep        map[rune]rune
	new_line   map[rune]rune
	trim       map[rune]rune
	open_quote map[rune]rune

	state       NextState
	close_quote []rune
	last_token  strings.Builder
	last_trim   strings.Builder
}

type Quote_t struct {
	Open  rune
	Close rune
}

func NewLexer(sep []rune, new_line []rune, trim []rune, quote []Quote_t) (self *Lexer_t) {
	self = &Lexer_t{
		sep:        map[rune]rune{},
		new_line:   map[rune]rune{},
		trim:       map[rune]rune{},
		open_quote: map[rune]rune{},
	}
	self.state = self.begin
	for _, v := range sep {
		self.sep[v] = v
	}
	for _, v := range new_line {
		self.new_line[v] = v
	}
	for _, v := range trim {
		self.trim[v] = v
	}
	for _, v := range quote {
		self.open_quote[v.Open] = v.Close
	}
	return
}

func (self *Lexer_t) begin(last_rune rune, last_size int, last_state State_t) (NextState, State_t) {
	switch {
	case self.open_quote[last_rune] > 0:
		self.close_quote = append(self.close_quote, self.open_quote[last_rune])
		return self.quoted, STATE_OPEN_QUOTE
	case self.sep[last_rune] > 0:
		return self.begin, STATE_SEPARATOR
	case self.new_line[last_rune] > 0:
		return self.begin, STATE_NEW_LINE
	case self.trim[last_rune] > 0:
		return self.begin, STATE_TRIM
	case last_size > 0:
		self.last_token.WriteRune(last_rune)
		return self.not_quoted, STATE_DATA
	default:
		return nil, STATE_EOF
	}
}

func (self *Lexer_t) not_quoted(last_rune rune, last_size int, last_state State_t) (NextState, State_t) {
	switch {
	case self.sep[last_rune] > 0:
		self.last_trim.Reset()
		return self.begin, STATE_SEPARATOR
	case self.new_line[last_rune] > 0:
		self.last_trim.Reset()
		return self.begin, STATE_NEW_LINE
	case self.trim[last_rune] > 0:
		self.last_trim.WriteRune(last_rune)
		return self.not_quoted, STATE_TRIM
	case last_size > 0:
		if self.last_trim.Len() > 0 {
			self.last_token.WriteString(self.last_trim.String())
			self.last_trim.Reset()
		}
		self.last_token.WriteRune(last_rune)
		return self.not_quoted, STATE_DATA
	default:
		return nil, STATE_EOF
	}
}

func (self *Lexer_t) quoted(last_rune rune, last_size int, last_state State_t) (NextState, State_t) {
	q_len := len(self.close_quote)
	switch {
	case self.close_quote[q_len-1] == last_rune:
		self.close_quote = self.close_quote[:q_len-1]
		if q_len == 1 {
			return self.separator, STATE_CLOSE_QUOTE
		}
		return self.quoted, STATE_CLOSE_QUOTE
	case last_state == STATE_OPEN_QUOTE && self.open_quote[last_rune] > 0:
		self.close_quote = append(self.close_quote, self.open_quote[last_rune])
		return self.quoted, STATE_OPEN_QUOTE
	case last_state != STATE_CLOSE_QUOTE && last_size > 0:
		self.last_token.WriteRune(last_rune)
		return self.quoted, STATE_DATA
	default:
		return nil, STATE_ERROR_NO_QUOTE
	}
}

func (self *Lexer_t) separator(last_rune rune, last_size int, last_state State_t) (NextState, State_t) {
	switch {
	case self.sep[last_rune] > 0:
		return self.begin, STATE_SEPARATOR
	case self.new_line[last_rune] > 0:
		return self.begin, STATE_NEW_LINE
	case self.trim[last_rune] > 0:
		return self.separator, STATE_TRIM
	case last_size == 0:
		return nil, STATE_EOF
	default:
		return nil, STATE_ERROR_NO_SEPARATOR
	}
}

func (self *Lexer_t) Next(in io.RuneReader) (token string, state State_t) {
	var last_rune rune
	var last_size int
	for self.state != nil {
		last_rune, last_size, _ = in.ReadRune()
		self.state, state = self.state(last_rune, last_size, state)
		if state == STATE_SEPARATOR || state == STATE_NEW_LINE || state >= STATE_EOF {
			token = self.last_token.String()
			self.last_token.Reset()
			return
		}
	}
	return
}

func Split(in string, sep ...rune) (res []string, err error) {
	l := NewLexer(
		sep,
		[]rune{'\n'},
		[]rune{'\v', '\f', '\r', '\t', ' '},
		[]Quote_t{
			{'"', '"'},
			{'\'', '\''},
			{'«', '»'},
		},
	)
	reader := strings.NewReader(in)
	for {
		token, state := l.Next(reader)
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
