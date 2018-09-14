//
// strings.Fields
//

package fields

import "strings"
import "unicode"

type Fields_t struct {
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
		return c == ',' || unicode.IsSpace(c)
	}
}

func (self * Fields_t) Fields(in string) []string {
	return strings.FieldsFunc(in, self.test)
}

func Fields(in string) []string {
	return (&Fields_t{}).Fields(in)
}

type Strings_t []string

func (self * Strings_t) Set(value string) (err error) {
	for _, v := range Fields(value) {
		*self = append(*self, v)
	}
	return
}

func (self * Strings_t) String() string {
	return strings.Join(*self, ",")
}
