//
// Split quoted comma separated list
//

package fields

import "bufio"
import "strings"
import "unicode/utf8"

// import "github.com/ondi/go-log"

type Split_t struct {
	Split map[rune]int
	Quote map[rune]int
	Trim map[rune]int
	last_quote rune
	produce_token bool
}

func (self * Split_t) Token(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var last_rune rune
	var last_size int
	var last_len int
	
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
				token = token[:last_len]
			}
			return
		case self.last_quote == last_rune:
			self.last_quote = 0
		case self.Quote[last_rune] != 0:
			if len(token) > 0 {
				token = append(token, data[advance - last_size:advance]...)
				self.produce_token = false
				last_len = len(token)
			} else {
				self.last_quote = last_rune
				self.produce_token = true
			}
		case self.last_quote != 0:
			token = append(token, data[advance - last_size:advance]...)
			self.produce_token = false
			last_len = len(token)
		case self.Split[last_rune] != 0:
			self.produce_token = true
			if token == nil {
				token = []byte{}
			} else {
				token = token[:last_len]
			}
			return
		case self.Trim[last_rune] != 0:
			if len(token) > 0 {
				token = append(token, data[advance - last_size:advance]...)
				self.produce_token = false
			}
		default:
			token = append(token, data[advance - last_size:advance]...)
			self.produce_token = false
			last_len = len(token)
		}
	}
	return
}

func split(in string, s * Split_t) (res []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(in))
	scanner.Split(s.Token)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	err = scanner.Err()
	return
}

func Split(in string, sep ...rune) ([]string, error) {
	s := &Split_t {
		Split: map[rune]int{},
		Quote: map[rune]int{'"': 1, '\'': 1},
		Trim: map[rune]int{'\v': 1, '\f': 1, '\r': 1, '\n': 1, '\t': 1, ' ': 1},
	}
	for _, r := range sep {
		s.Split[r] = 1
	}
	return split(in, s)
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
