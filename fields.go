//
// Split quoted comma separated list
//

package fields

import "strings"

import "github.com/ondi/go-log"

type Quote_t struct {
	Open rune
	Close rune
}

type Lexer_t struct {
	sep map[rune]rune
	trim map[rune]rune
	quote map[rune]rune
	
	reader * strings.Reader
	last_quote []rune
	last_token strings.Builder
	last_trim strings.Builder
	tokens []string
	quoted_prev bool
	quoted_curr bool
	err error
}

type StateFunc func(l * Lexer_t) StateFunc

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

func Unquoted(lexer * Lexer_t) StateFunc {
	last_rune, last_size, _ := lexer.reader.ReadRune()
	log.Debug("Unquoted: rune=`%c`, len=%v, tokens=%#v", last_rune, last_size, lexer.tokens)
	lexer.quoted_prev = lexer.quoted_curr
	lexer.quoted_curr = false
	switch {
	case lexer.quote[last_rune] > 0:
		if lexer.last_token.Len() > 0 {
			lexer.last_token.WriteRune(last_rune)
			return Unquoted
		}
		lexer.last_quote = append(lexer.last_quote, lexer.quote[last_rune])
		return Quoted
	case lexer.trim[last_rune] > 0:
		if lexer.last_token.Len() > 0 {
			lexer.last_trim.WriteRune(last_rune)
		}
		return Unquoted
	case lexer.sep[last_rune] > 0:
		lexer.tokens = append(lexer.tokens, lexer.last_token.String())
		lexer.last_token.Reset()
		lexer.last_trim.Reset()
		if _, last_size, _ = lexer.reader.ReadRune(); last_size == 0 {
			lexer.tokens = append(lexer.tokens, lexer.last_token.String())
			return nil
		} else {
			lexer.reader.UnreadRune()
		}
		return Unquoted
	case last_size > 0:
		if lexer.last_trim.Len() > 0 {
			lexer.last_token.WriteString(lexer.last_trim.String())
			lexer.last_trim.Reset()
		}
		lexer.last_token.WriteRune(last_rune)
		return Unquoted
	case last_size == 0:
		if lexer.last_token.Len() > 0 || lexer.quoted_prev {
			lexer.tokens = append(lexer.tokens, lexer.last_token.String())
			lexer.last_token.Reset()
		}
	}
	return nil
}

func Quoted(lexer * Lexer_t) StateFunc {
	last_rune, last_size, _ := lexer.reader.ReadRune()
	log.Debug("Quoted : rune=`%c`, len=%v, tokens=%#v", last_rune, last_size, lexer.tokens)
	lexer.quoted_prev = lexer.quoted_curr
	lexer.quoted_curr = true
	switch {
	case lexer.last_quote[len(lexer.last_quote) - 1] == last_rune:
		lexer.last_quote = lexer.last_quote[:len(lexer.last_quote) - 1]
		if len(lexer.last_quote) > 0 {
			return Quoted
		}
		return Unquoted
	case lexer.quote[last_rune] > 0:
		lexer.last_quote = append(lexer.last_quote, lexer.quote[last_rune])
		return Quoted
	case last_size > 0:
		lexer.last_token.WriteRune(last_rune)
		return Quoted
	case last_size == 0:
		lexer.tokens = append(lexer.tokens, lexer.last_token.String())
		lexer.last_token.Reset()
	}
	return nil
}

func (self * Lexer_t) Split(in string) (res []string, err error) {
	self.reader = strings.NewReader(in)
	for state := Unquoted(self); state != nil; {
		state = state(self)
	}
	return self.tokens, self.err
}

func Split(in string, sep ...rune) ([]string, error) {
	return NewLexer(sep, []rune{'\v', '\f', '\r', '\n', '\t', ' '}, []Quote_t{Quote_t{'"', '"'}, Quote_t{'\'', '\''}}).Split(in)
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
