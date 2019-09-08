//
// Split quoted comma separated list
//

package fields

import "strings"

// import "github.com/ondi/go-log"

type NextState func() NextState

type Lexer_t struct {
	sep map[rune]rune
	trim map[rune]rune
	quote map[rune]rune
	
	reader * strings.Reader
	last_quote rune
	last_token strings.Builder
	last_trim strings.Builder
	tokens []string
	err error
}

type Quote_t struct {
	Open rune
	Close rune
}

func NewLexer(sep []rune, trim []rune, quote []Quote_t) (self * Lexer_t) {
	self = &Lexer_t{sep: map[rune]rune{}, trim: map[rune]rune{}, quote: map[rune]rune{}, tokens: []string{}}
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

func (self * Lexer_t) Split(in string) ([]string, error) {
	self.reader = strings.NewReader(in)
	state := self.begin
	for state != nil {
		state = state()
	}
	return self.tokens, self.err
}

func (self * Lexer_t) begin() NextState {
	last_rune, last_size, _ := self.reader.ReadRune()
	// log.Debug("Begin: rune=`%c`, len=%v, tokens=%#v", last_rune, last_size, self.tokens)
	switch {
	case self.quote[last_rune] > 0:
		self.last_quote = self.quote[last_rune]
		return self.quoted
	case self.trim[last_rune] > 0:
		return self.begin
	case self.sep[last_rune] > 0:
		self.tokens = append(self.tokens, self.last_token.String())
		return self.begin
	case last_size > 0:
		self.last_token.WriteRune(last_rune)
		return self.unquoted
	default:
		self.tokens = append(self.tokens, self.last_token.String())
		return nil
	}
}

func (self * Lexer_t) unquoted() NextState {
	last_rune, last_size, _ := self.reader.ReadRune()
	// log.Debug("Unquoted: rune=`%c`, len=%v, tokens=%#v", last_rune, last_size, self.tokens)
	switch {
	case self.trim[last_rune] > 0:
		self.last_trim.WriteRune(last_rune)
		return self.unquoted
	case self.sep[last_rune] > 0:
		self.tokens = append(self.tokens, self.last_token.String())
		self.last_token.Reset()
		self.last_trim.Reset()
		return self.begin
	case last_size > 0:
		if self.last_trim.Len() > 0 {
			self.last_token.WriteString(self.last_trim.String())
			self.last_trim.Reset()
		}
		self.last_token.WriteRune(last_rune)
		return self.unquoted
	default:
		self.tokens = append(self.tokens, self.last_token.String())
		return nil
	}
}

func (self * Lexer_t) quoted() NextState {
	last_rune, last_size, _ := self.reader.ReadRune()
	// log.Debug("Quoted : rune=`%c`, len=%v, tokens=%#v", last_rune, last_size, self.tokens)
	switch {
	case self.last_quote == last_rune:
		return self.unquoted
	case last_size > 0:
		self.last_token.WriteRune(last_rune)
		return self.quoted
	default:
		self.tokens = append(self.tokens, self.last_token.String())
		return nil
	}
}

func Split(in string, sep ...rune) ([]string, error) {
	return NewLexer(sep,
		[]rune{'\v', '\f', '\r', '\n', '\t', ' '},
		[]Quote_t{Quote_t{'"', '"'}, Quote_t{'\'', '\''}, Quote_t{'«', '»'}},
	).Split(in)
}

type Strings_t []string

func (self * Strings_t) Set(value string) (err error) {
	var temp []string
	if temp, err = Split(value, ','); err != nil {
		return
	}
	for _, v := range temp {
		*self = append(*self, v)
	}
	return
}

func (self * Strings_t) String() string {
	return strings.Join(*self, ",")
}
