//
// Split quoted comma separated list
//

package fields

import "fmt"
import "bufio"
import "strings"
import "unicode"
import "unicode/utf8"

// import "github.com/ondi/go-log"

type Split_t struct {
	Sep map[rune]int
	Ignore map[rune]int
	last_quote rune
	last_split rune
	last_rune rune
	last_size int
}

func (self * Split_t) Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// log.Debug("data = '%s', EOF = %v", data, atEOF)
	
	if len(data) == 0 && atEOF {
		if self.last_quote != 0 {
			self.last_quote = 0
			err = fmt.Errorf("unmatched quote")
			return
		}
		if self.last_split != 0 {
			self.last_split = 0
			token = []byte{}
		}
		return
	}
	
	for {
		if self.last_rune, self.last_size = utf8.DecodeRune(data[advance:]); self.last_size == 0 {
			return
		}
		// log.Debug("rune = '%c', size = %d", self.last_rune, self.last_size)
		advance += self.last_size
		switch {
		case self.last_quote == self.last_rune:
			self.last_quote = 0
		case unicode.In(self.last_rune, unicode.Quotation_Mark):
			self.last_quote = self.last_rune
		case self.last_quote != 0:
			self.last_split = 0
			token = append(token, data[advance - self.last_size:advance]...)
		case self.Sep[self.last_rune] != 0:
			self.last_split = self.last_rune
			if token == nil {
				token = []byte{}
			}
			return
		case self.Ignore[self.last_rune] != 0 && len(token) == 0:
			//
		default:
			self.last_split = 0
			token = append(token, data[advance - self.last_size:advance]...)
		}
	}
	// return 0, data, bufio.ErrFinalToken
	return
}

func Split(in string, s * Split_t) (res []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(in))
	scanner.Split(s.Split)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	err = scanner.Err()
	return
}

func SplitCSV(in string) ([]string, error) {
	s := &Split_t {
		Sep: map[rune]int{',': 1},
		Ignore: map[rune]int{'\v': 1, '\f': 1, '\r': 1, '\n': 1, '\t': 1, ' ': 1},
	}
	return Split(in, s)
}

func SplitTSV(in string) ([]string, error) {
	s := &Split_t {
		Sep: map[rune]int{'\t': 1},
		Ignore: map[rune]int{'\v': 1, '\f': 1, '\r': 1, '\n': 1, '\t': 1, ' ': 1},
	}
	return Split(in, s)
}

type Strings_t []string

func (self * Strings_t) Set(value string) (err error) {
	var temp []string
	if temp, err = SplitCSV(value); err != nil {
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
