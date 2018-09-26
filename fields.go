//
// strings.Fields
//

package fields

import "bufio"
import "strings"
import "unicode"
import "unicode/utf8"

// import "github.com/ondi/go-log"

type fields_t struct {
	Sep map[rune]int
	last_quote rune
}

func (self * fields_t) test(c rune) bool {
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

func (self * fields_t) Fields(in string) []string {
	return strings.FieldsFunc(in, self.test)
}

func FieldsTSV(in string) []string {
	return (&fields_t{Sep: map[rune]int{' ': 0, '\t': 0, '\v': 0, '\r': 0, '\n': 0, '\f': 0}}).Fields(in)
}

func FieldsCSV(in string) []string {
	return (&fields_t{Sep: map[rune]int{',': 0, ' ': 0, '\t': 0, '\v': 0, '\r': 0, '\n': 0, '\f': 0}}).Fields(in)
}

type Split_t struct {
	Sep map[rune]int
	Ignore map[rune]int
}

func (self * Split_t) Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var last_rune rune
	var last_size int
	var last_quote rune
	
	for {
		if last_rune, last_size = utf8.DecodeRune(data[advance:]); last_size == 0 {
			break
		}
		// log.Debug("rune: '%c', last_size: %v", last_rune, last_size)
		advance += last_size
		
		if last_quote == 0 && unicode.In(last_rune, unicode.Quotation_Mark) {
			last_quote = last_rune
		} else if last_quote == last_rune {
			last_quote = 0
		} else if _, ok := self.Sep[last_rune]; ok && last_quote == 0 {
			if token == nil {
				token = []byte{}
			}
			break
		} else if _, ok := self.Ignore[last_rune]; ok && last_quote == 0 {
			continue
		} else {
			token = append(token, data[advance - last_size:advance]...)
		}
	}
	// log.Debug("RETURN: advance=%v, token='%s' (%v)", advance, token, token == nil)
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
		Sep: map[rune]int{',': 0},
		Ignore: map[rune]int{' ': 0, '\t': 0, '\v': 0, '\r': 0, '\n': 0, '\f': 0},
	}
	return Split(in, s)
}

func SplitTSV(in string) []string {
	s := &Split_t {
		Sep: map[rune]int{'\t': 0},
		Ignore: map[rune]int{' ': 0, '\t': 0, '\v': 0, '\r': 0, '\n': 0, '\f': 0},
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
