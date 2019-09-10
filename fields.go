//
// Split quoted comma separated list
//

package fields

import "strings"

// import "github.com/ondi/go-log"

type NextState func() (NextState, bool)

type Lexer_t struct {
	sep map[rune]rune
	trim map[rune]rune
	quote map[rune]rune
	
	reader * strings.Reader
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

func (self * Lexer_t) Set(in string) {
	self.reader = strings.NewReader(in)
	self.state = self.begin
}

func (self * Lexer_t) Next() (res string, ok bool) {
	for self.state != nil && ok == false {
		self.state, ok = self.state()
	}
	if ok {
		res = self.last_token.String()
		self.last_token.Reset()
	}
	return
}

func (self * Lexer_t) Err() error {
	return self.err
}

func (self * Lexer_t) begin() (NextState, bool) {
	last_rune, last_size, _ := self.reader.ReadRune()
	// log.Debug("Begin: rune=`%c`, len=%v, tokens=%#v", last_rune, last_size, self.tokens)
	switch {
	case self.quote[last_rune] > 0:
		self.last_quote = self.quote[last_rune]
		return self.quoted, false
	case self.sep[last_rune] > 0:
		return self.begin, true
	case self.trim[last_rune] > 0:
		return self.begin, false
	case last_size > 0:
		self.last_token.WriteRune(last_rune)
		return self.unquoted, false
	default:
		return nil, true
	}
}

func (self * Lexer_t) unquoted() (NextState, bool) {
	last_rune, last_size, _ := self.reader.ReadRune()
	// log.Debug("Unquoted: rune=`%c`, len=%v, tokens=%#v", last_rune, last_size, self.tokens)
	switch {
	case self.sep[last_rune] > 0:
		self.last_trim.Reset()
		return self.begin, true
	case self.trim[last_rune] > 0:
		self.last_trim.WriteRune(last_rune)
		return self.unquoted, false
	case last_size > 0:
		if self.last_trim.Len() > 0 {
			self.last_token.WriteString(self.last_trim.String())
			self.last_trim.Reset()
		}
		self.last_token.WriteRune(last_rune)
		return self.unquoted, false
	default:
		return nil, true
	}
}

func (self * Lexer_t) quoted() (NextState, bool) {
	last_rune, last_size, _ := self.reader.ReadRune()
	// log.Debug("Quoted : rune=`%c`, len=%v, tokens=%#v", last_rune, last_size, self.tokens)
	switch {
	case self.last_quote == last_rune:
		return self.unquoted, false
	case last_size > 0:
		self.last_token.WriteRune(last_rune)
		return self.quoted, false
	default:
		return nil, true
	}
}

func Split(in string, sep ...rune) (res []string, err error) {
	l := NewLexer(sep,
		[]rune{'\v', '\f', '\r', '\n', '\t', ' '},
		[]Quote_t{Quote_t{'"', '"'}, Quote_t{'\'', '\''}, Quote_t{'«', '»'}},
	)
	l.Set(in)
	for {
		if token, ok := l.Next(); !ok {
			break
		} else {
			res = append(res, token)
		}
	}
	return res, l.Err()
}

type Strings_t []string

func (self * Strings_t) Set(value string) (error) {
	temp, err := Split(value, ',')
	if err != nil {
		for _, v := range temp {
			*self = append(*self, v)
		}
	}
	return err
}

func (self * Strings_t) String() string {
	return strings.Join(*self, ",")
}
