//
// Split quoted comma separated list
//

package fields

import "io"
import "strings"

// import "github.com/ondi/go-log"

type State_t int

const (
	STATE_NONE State_t = 0
	STATE_OPEN_QUOTE State_t = 1
	STATE_CLOSE_QUOTE State_t = 2
	STATE_SEPARATOR State_t = 3
	STATE_TRIM State_t = 4
	STATE_STRING State_t = 5
	STATE_EOF State_t = 6
)

type NextState func() (NextState, State_t)

type Lexer_t struct {
	sep map[rune]rune
	trim map[rune]rune
	quote map[rune]rune
	
	reader io.RuneReader
	state NextState
	last_rune rune
	last_size int
	last_quote rune
	last_token strings.Builder
	last_trim strings.Builder
	err error
}

type Quote_t struct {
	Open rune
	Close rune
}

func NewLexer(sep []rune, trim []rune, quote []Quote_t) (self * Lexer_t) {
	self = &Lexer_t{sep: map[rune]rune{}, trim: map[rune]rune{}, quote: map[rune]rune{}}
	for _, v := range sep {
		self.sep[v] = 1
	}
	for _, v := range trim {
		self.trim[v] = 1
	}
	for _, v := range quote {
		self.quote[v.Open] = v.Close
	}
	return
}

func (self * Lexer_t) begin() (NextState, State_t) {
	self.last_rune, self.last_size, _ = self.reader.ReadRune()
	// log.Debug("begin: rune=`%c`, len=%v, token=%#v", self.last_rune, self.last_size, self.last_token.String())
	switch {
	case self.quote[self.last_rune] > 0:
		self.last_quote = self.quote[self.last_rune]
		return self.quoted, STATE_OPEN_QUOTE	// false
	case self.sep[self.last_rune] > 0:
		return self.begin, STATE_SEPARATOR		// true
	case self.trim[self.last_rune] > 0:
		return self.begin, STATE_TRIM			// false
	case self.last_size > 0:
		self.last_token.WriteRune(self.last_rune)
		return self.unquoted, STATE_STRING		// false
	default:
		return nil, STATE_EOF					// true
	}
}

func (self * Lexer_t) unquoted() (NextState, State_t) {
	self.last_rune, self.last_size, _ = self.reader.ReadRune()
	// log.Debug("unquoted: rune=`%c`, len=%v, token=%#v", self.last_rune, self.last_size, self.last_token.String())
	switch {
	case self.sep[self.last_rune] > 0:
		self.last_trim.Reset()
		return self.begin, STATE_SEPARATOR		// true
	case self.trim[self.last_rune] > 0:
		self.last_trim.WriteRune(self.last_rune)
		return self.unquoted, STATE_TRIM		// false
	case self.last_size > 0:
		if self.last_trim.Len() > 0 {
			self.last_token.WriteString(self.last_trim.String())
			self.last_trim.Reset()
		}
		self.last_token.WriteRune(self.last_rune)
		return self.unquoted, STATE_STRING		// false
	default:
		return nil, STATE_EOF					// true
	}
}

func (self * Lexer_t) quoted() (NextState, State_t) {
	self.last_rune, self.last_size, _ = self.reader.ReadRune()
	// log.Debug("quoted: rune=`%c`, len=%v, token=%#v", self.last_rune, self.last_size, self.last_token.String())
	switch {
	case self.last_quote == self.last_rune:
		return self.unquoted, STATE_CLOSE_QUOTE	// false
	case self.last_size > 0:
		self.last_token.WriteRune(self.last_rune)
		return self.quoted, STATE_STRING		// false
	default:
		return nil, STATE_EOF					// true
	}
}

func (self * Lexer_t) Set(in io.RuneReader) {
	self.reader = in
	self.state = self.begin
}

func (self * Lexer_t) Next() (token string, last_rune rune, state State_t) {
	for self.state != nil {
		self.state, state = self.state()
		switch state {
		case STATE_SEPARATOR, STATE_EOF:
			token = self.last_token.String()
			last_rune = self.last_rune
			self.last_token.Reset()
			return
		}
	}
	return
}

func (self * Lexer_t) Err() error {
	return self.err
}

func Split(in string, sep ...rune) (res []string, err error) {
	l := NewLexer(sep,
		[]rune{'\v', '\f', '\r', '\n', '\t', ' '},
		[]Quote_t{Quote_t{'"', '"'}, Quote_t{'\'', '\''}, Quote_t{'«', '»'}},
	)
	l.Set(strings.NewReader(in))
	for {
		if token, _, state := l.Next(); state == STATE_NONE {
			break
		} else {
			res = append(res, token)
		}
	}
	return res, l.Err()
}

type Strings_t []string

func (self * Strings_t) Set(value string) (err error) {
	var temp []string
	if temp, err = Split(value, ','); err == nil {
		for _, v := range temp {
			*self = append(*self, v)
		}
	}
	return
}

func (self * Strings_t) String() string {
	return strings.Join(*self, ",")
}
