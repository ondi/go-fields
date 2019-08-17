//
// Split quoted comma separated list
//

package fields

import "bufio"
import "strings"
import "unicode/utf8"

// import "github.com/ondi/go-log"

type Quote_t struct {
	Open rune
	Close rune
}

type Split_t struct {
	Sep map[rune]int
	Trim map[rune]int
	Quote map[rune]rune		// map[open_quote]close_quote
	last_quote rune
	produce_token bool
}

func NewSplit(sep []rune, trim []rune, quote []Quote_t) (self * Split_t) {
	self = &Split_t{Sep: map[rune]int{}, Trim: map[rune]int{}, Quote: map[rune]rune{}}
	for _, v := range sep {
		self.Sep[v] = 1
	}
	for _, v := range trim {
		self.Trim[v] = 1
	}
	for _, v := range quote {
		self.Quote[v.Open] = v.Close
	}
	return
}

func (self * Split_t) Token(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var last_rune rune
	var last_size int
	var trim_len int
	
	for {
		last_rune, last_size = utf8.DecodeRune(data[advance:])
		advance += last_size
		// log.Debug("rune = '%c', size = %d", last_rune, last_size)
		switch {
		case last_size == 0:
			if self.produce_token {
				self.produce_token = false
				token = []byte{}
				return
			}
			if len(token) > 0 {
				token = token[:trim_len]
			}
			return
		case self.last_quote == last_rune:
			self.last_quote = 0
		case self.Quote[last_rune] > 0:
			if len(token) > 0 {
				token = append(token, data[advance - last_size:advance]...)
				self.produce_token = false
				trim_len = len(token)
			} else {
				self.last_quote = self.Quote[last_rune]
				self.produce_token = true
			}
		case self.last_quote > 0:
			token = append(token, data[advance - last_size:advance]...)
			self.produce_token = false
			trim_len = len(token)
		case self.Sep[last_rune] > 0:
			self.produce_token = true
			if token == nil {
				token = []byte{}
			} else {
				token = token[:trim_len]
			}
			return
		case self.Trim[last_rune] > 0:
			if len(token) > 0 {
				token = append(token, data[advance - last_size:advance]...)
				self.produce_token = false
			}
		default:
			token = append(token, data[advance - last_size:advance]...)
			self.produce_token = false
			trim_len = len(token)
		}
	}
	return
}

func (self * Split_t) Split(in string) (res []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(in))
	scanner.Split(self.Token)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	err = scanner.Err()
	return
}

func Split(in string, sep ...rune) ([]string, error) {
	return NewSplit(sep, []rune{'\v', '\f', '\r', '\n', '\t', ' '}, []Quote_t{Quote_t{'"', '"'}, Quote_t{'\'', '\''}}).Split(in)
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
