//
// Split quoted comma separated list
//

package fields

import "io"
import "strings"

// import "github.com/ondi/go-log"

type Token_t int

const (
	TOKEN_NONE Token_t = 0
	TOKEN_OPEN_QUOTE Token_t = 1
	TOKEN_CLOSE_QUOTE Token_t = 2
	TOKEN_SEPARATOR Token_t = 3
	TOKEN_TRIM Token_t = 4
	TOKEN_STRING Token_t = 5
	TOKEN_EOF Token_t = 6
)

type NextState func() (NextState, Token_t)

type Lexer_t struct {
	sep map[rune]rune
	trim map[rune]rune
	quote map[rune]rune
	
	reader io.RuneReader
	state NextState
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

func (self * Lexer_t) begin() (NextState, Token_t) {
	last_rune, last_size, _ := self.reader.ReadRune()
	// log.Debug("Begin: rune=`%c`, len=%v, tokens=%#v", last_rune, last_size, self.tokens)
	switch {
	case self.quote[last_rune] > 0:
		self.last_quote = self.quote[last_rune]
		return self.quoted, TOKEN_OPEN_QUOTE	// false
	case self.sep[last_rune] > 0:
		return self.begin, TOKEN_SEPARATOR		// true
	case self.trim[last_rune] > 0:
		return self.begin, TOKEN_TRIM			// false
	case last_size > 0:
		self.last_token.WriteRune(last_rune)
		return self.unquoted, TOKEN_STRING		// false
	default:
		return nil, TOKEN_EOF					// true
	}
}

func (self * Lexer_t) unquoted() (NextState, Token_t) {
	last_rune, last_size, _ := self.reader.ReadRune()
	// log.Debug("Unquoted: rune=`%c`, len=%v, tokens=%#v", last_rune, last_size, self.tokens)
	switch {
	case self.sep[last_rune] > 0:
		self.last_trim.Reset()
		return self.begin, TOKEN_SEPARATOR		// true
	case self.trim[last_rune] > 0:
		self.last_trim.WriteRune(last_rune)
		return self.unquoted, TOKEN_TRIM		// false
	case last_size > 0:
		if self.last_trim.Len() > 0 {
			self.last_token.WriteString(self.last_trim.String())
			self.last_trim.Reset()
		}
		self.last_token.WriteRune(last_rune)
		return self.unquoted, TOKEN_STRING		// false
	default:
		return nil, TOKEN_EOF					// true
	}
}

func (self * Lexer_t) quoted() (NextState, Token_t) {
	last_rune, last_size, _ := self.reader.ReadRune()
	// log.Debug("Quoted : rune=`%c`, len=%v, tokens=%#v", last_rune, last_size, self.tokens)
	switch {
	case self.last_quote == last_rune:
		return self.unquoted, TOKEN_CLOSE_QUOTE	// false
	case last_size > 0:
		self.last_token.WriteRune(last_rune)
		return self.quoted, TOKEN_STRING		// false
	default:
		return nil, TOKEN_EOF					// true
	}
}

func (self * Lexer_t) Set(in io.RuneReader) {
	self.reader = in
	self.state = self.begin
}

func (self * Lexer_t) Next() (res string, status Token_t) {
	for self.state != nil {
		self.state, status = self.state()
		switch status {
		case TOKEN_SEPARATOR, TOKEN_EOF:
			res = self.last_token.String()
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
		if token, ok := l.Next(); ok == 0 {
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
