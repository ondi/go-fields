//
// Split quoted comma separated list
//

package fields

import "bufio"
import "strings"
import "unicode"
import "unicode/utf8"

type Fields_t struct {
	Sep map[rune]int
	last_quote rune
}

func (self * Fields_t) test(c rune) bool {
	switch {
	case self.last_quote == c:
		self.last_quote = 0
		return true		// false to keep quotes
	case self.last_quote != 0:
		return false
	case unicode.In(c, unicode.Quotation_Mark):
		self.last_quote = c
		return true		// false to keep quotes
	default:
		_, ok := self.Sep[c]
		return ok
	}
}

func (self * Fields_t) Fields(in string) []string {
	return strings.FieldsFunc(in, self.test)
}

func FieldsTSV(in string) []string {
	return (&Fields_t {Sep: map[rune]int{' ': 0, '\t': 0, '\v': 0, '\r': 0, '\n': 0, '\f': 0}}).Fields(in)
}

func FieldsCSV(in string) []string {
	return (&Fields_t{Sep: map[rune]int{',': 0, ' ': 0, '\t': 0, '\v': 0, '\r': 0, '\n': 0, '\f': 0}}).Fields(in)
}

type Split_t struct {
	Sep map[rune]int
	Ignore map[rune]int
	last_quote rune
	last_rune rune
	last_size int
}

func (self * Split_t) Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for {
		if self.last_rune, self.last_size = utf8.DecodeRune(data[advance:]); self.last_size == 0 {
			return
		}
		advance += self.last_size
		switch {
		case self.last_quote == self.last_rune:
			self.last_quote = 0
		case self.last_quote != 0:
			token = append(token, data[advance - self.last_size:advance]...)
		case unicode.In(self.last_rune, unicode.Quotation_Mark):
			self.last_quote = self.last_rune
		case self.Sep[self.last_rune] != 0:
			if token == nil {
				token = []byte{}
			}
			return
		case self.Ignore[self.last_rune] != 0 && len(token) == 0:
			//
		default:
			token = append(token, data[advance - self.last_size:advance]...)
		}
	}
	return
}

func Split(in string, s * Split_t) (res []string) {
	scanner := bufio.NewScanner(strings.NewReader(in))
	scanner.Split(s.Split)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	return
}

func SplitCSV(in string) []string {
	s := &Split_t {
		Sep: map[rune]int{',': 1},
		Ignore: map[rune]int{'\v': 1, '\f': 1, '\r': 1, '\n': 1, '\t': 1, ' ': 1},
	}
	return Split(in, s)
}

func SplitTSV(in string) []string {
	s := &Split_t {
		Sep: map[rune]int{'\t': 1},
		Ignore: map[rune]int{'\v': 1, '\f': 1, '\r': 1, '\n': 1, '\t': 1, ' ': 1},
	}
	return Split(in, s)
}

type Strings_t []string

func (self * Strings_t) Set(value string) (err error) {
	for _, v := range SplitCSV(value) {
		*self = append(*self, v)
	}
	return
}

func (self * Strings_t) String() string {
	return strings.Join(*self, ",")
}
