//
// Split quoted comma separated list
//

package fields

import (
	"io"
	"strings"
)

type State_t int

const (
	STATE_NONE        State_t = 0
	STATE_OPEN_QUOTE  State_t = 1
	STATE_CLOSE_QUOTE State_t = 2
	STATE_SEPARATOR   State_t = 3
	STATE_TRIM        State_t = 4
	STATE_STRING      State_t = 5
	STATE_ERROR       State_t = 6
	STATE_EOF         State_t = 7
)

type NextState func(State_t) (NextState, State_t)

type Lexer_t struct {
	sep   map[rune]rune
	trim  map[rune]rune
	quote map[rune]rune

	reader      io.RuneReader
	state       NextState
	last_rune   rune
	last_size   int
	last_quotes []rune
	last_token  strings.Builder
	last_trim   strings.Builder
}

type Quote_t struct {
	Open  rune
	Close rune
}

func NewLexer(sep []rune, trim []rune, quote []Quote_t) (self *Lexer_t) {
	self = &Lexer_t{sep: map[rune]rune{}, trim: map[rune]rune{}, quote: map[rune]rune{}}
	for _, v := range sep {
		self.sep[v] = v
	}
	for _, v := range trim {
		self.trim[v] = v
	}
	for _, v := range quote {
		self.quote[v.Open] = v.Close
	}
	return
}

func (self *Lexer_t) begin(prev State_t) (NextState, State_t) {
	self.last_rune, self.last_size, _ = self.reader.ReadRune()
	// log.Debug("begin: rune=`%c`, len=%v, token=%#v", self.last_rune, self.last_size, self.last_token.String())
	switch {
	case self.quote[self.last_rune] > 0:
		self.last_quotes = append(self.last_quotes, self.quote[self.last_rune])
		return self.first_quote, STATE_OPEN_QUOTE // false
	case self.sep[self.last_rune] > 0:
		return self.begin, STATE_SEPARATOR // true
	case self.trim[self.last_rune] > 0:
		return self.begin, STATE_TRIM // false
	case self.last_size > 0:
		self.last_token.WriteRune(self.last_rune)
		return self.no_quote, STATE_STRING // false
	default:
		return nil, STATE_EOF // true
	}
}

func (self *Lexer_t) no_quote(prev State_t) (NextState, State_t) {
	self.last_rune, self.last_size, _ = self.reader.ReadRune()
	// log.Debug("no_quote: rune=`%c`, len=%v, token=%#v", self.last_rune, self.last_size, self.last_token.String())
	switch {
	case self.sep[self.last_rune] > 0:
		self.last_trim.Reset()
		return self.begin, STATE_SEPARATOR // true
	case self.trim[self.last_rune] > 0:
		self.last_trim.WriteRune(self.last_rune)
		return self.no_quote, STATE_TRIM // false
	case self.last_size > 0:
		if self.last_trim.Len() > 0 {
			self.last_token.WriteString(self.last_trim.String())
			self.last_trim.Reset()
		}
		self.last_token.WriteRune(self.last_rune)
		return self.no_quote, STATE_STRING // false
	default:
		return nil, STATE_EOF // true
	}
}

func (self *Lexer_t) first_quote(prev State_t) (NextState, State_t) {
	self.last_rune, self.last_size, _ = self.reader.ReadRune()
	// log.Debug("first_quote: rune=`%c`, len=%v, token=%#v", self.last_rune, self.last_size, self.last_token.String())
	switch {
	case self.last_quotes[0] == self.last_rune:
		self.last_quotes = []rune{}
		return self.no_quote, STATE_CLOSE_QUOTE // false
	case self.quote[self.last_rune] > 0:
		if prev != STATE_OPEN_QUOTE {
			return nil, STATE_ERROR
		}
		self.last_quotes = append(self.last_quotes, self.quote[self.last_rune])
		return self.more_quote, STATE_OPEN_QUOTE
	case self.last_size > 0:
		self.last_token.WriteRune(self.last_rune)
		return self.first_quote, STATE_STRING // false
	default:
		return nil, STATE_EOF // true
	}
}

func (self *Lexer_t) more_quote(prev State_t) (NextState, State_t) {
	self.last_rune, self.last_size, _ = self.reader.ReadRune()
	q_len := len(self.last_quotes)
	switch {
	case self.last_quotes[q_len-1] == self.last_rune:
		self.last_quotes = self.last_quotes[:q_len-1]
		if q_len == 1 {
			return self.no_quote, STATE_CLOSE_QUOTE
		}
		return self.more_unquote, STATE_CLOSE_QUOTE
	case self.quote[self.last_rune] > 0:
		self.last_quotes = append(self.last_quotes, self.quote[self.last_rune])
		return self.more_quote, STATE_OPEN_QUOTE
	case self.last_size > 0:
		self.last_token.WriteRune(self.last_rune)
		return self.more_quote, STATE_STRING // false
	default:
		return nil, STATE_ERROR
	}
}

func (self *Lexer_t) more_unquote(prev State_t) (NextState, State_t) {
	self.last_rune, self.last_size, _ = self.reader.ReadRune()
	q_len := len(self.last_quotes)
	switch {
	case self.last_quotes[q_len-1] == self.last_rune:
		self.last_quotes = self.last_quotes[:q_len-1]
		if q_len == 1 {
			return self.no_quote, STATE_CLOSE_QUOTE
		}
		return self.more_unquote, STATE_CLOSE_QUOTE
	default:
		return nil, STATE_ERROR
	}
}

func (self *Lexer_t) Set(in io.RuneReader) {
	self.reader = in
	self.state = self.begin
}

func (self *Lexer_t) Next() (token string, last_rune rune, state State_t) {
	for self.state != nil {
		self.state, state = self.state(state)
		switch state {
		case STATE_SEPARATOR, STATE_ERROR, STATE_EOF:
			token = self.last_token.String()
			last_rune = self.last_rune
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
		if token, _, state := l.Next(); state == STATE_NONE {
			break
		} else {
			res = append(res, token)
		}
	}
	return res, nil
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
