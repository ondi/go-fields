//
//
//

package fields

import (
	"strings"
	"testing"

	"gotest.tools/assert"
)

type Data_t struct {
	Input  string
	Expect []string
	Err    bool
}

var data1 = []Data_t{
	{
		Input:  "",
		Expect: []string{""},
	},
	{
		Input:  "''",
		Expect: []string{""},
	},
	{
		Input:  "1",
		Expect: []string{"1"},
	},
	{
		Input:  "'1'",
		Expect: []string{"1"},
	},
	{
		Input:  "'`1`'",
		Expect: []string{"`1`"},
	},
	{
		Input:  "« 1,2,3 »",
		Expect: []string{" 1,2,3 "},
	},
	{
		Input:  "«« 1,2,3 »»",
		Expect: []string{" 1,2,3 "},
	},
	{
		Input:  "««« 1,2,3 »»»",
		Expect: []string{" 1,2,3 "},
	},
	{
		Input:  "««« 1,2,3, ",
		Expect: []string{" 1,2,3, "},
		Err:    true,
	},
	{
		Input:  ",",
		Expect: []string{"", ""},
	},
	{
		Input:  "'', it's test",
		Expect: []string{"", "it's test"},
	},
	{
		Input:  "'', 'it\"s test'",
		Expect: []string{"", "it\"s test"},
	},
	{
		Input:  "'',''",
		Expect: []string{"", ""},
	},
	{
		Input:  ",,",
		Expect: []string{"", "", ""},
	},
	{
		Input:  " ( ), ) (, ()",
		Expect: []string{"( )", ") (", "()"},
	},
	{
		Input:  "aaa'b'b'b'ccc,",
		Expect: []string{"aaa'b'b'b'ccc", ""},
	},
	{
		Input:  ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4' , ,aaa'5'bbb , it's test   ",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "it's test"},
	},
	{
		Input:  ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4' , ,aaa'5'bbb , it's test   ,",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "it's test", ""},
	},
	{
		Input:  ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4' , ,aaa'5'bbb , it's test   ,   ",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "it's test", ""},
	},
	{
		Input:  ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4' , ,aaa'5'bbb , it's test   ,''",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "it's test", ""},
	},
	{
		Input:  ", \n ',1,',\n' 2 ', , ' ,3 ,'\n , aaa'4' , \naaa'5'bbb , it's test   ,''",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "it's test", ""},
	},
	{
		Input:  ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4' , ,aaa'5'bbb , it's test   ,'",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "it's test", ""},
		Err:    true,
	},
}

func Test001(t *testing.T) {
	for _, v := range data1 {
		res, err := Split(v.Input, ',')
		t.Logf("Input  = %#v", v.Input)
		t.Logf("Expect = %#v, Error=%v", v.Expect, v.Err)
		t.Logf("Result = %#v, Error=%v", res, err)
		if v.Err {
			if err == nil {
				t.Fatalf("ERROR EXPECTED")
				return
			}
		} else {
			if err != nil {
				t.Fatalf("ERROR: %v", err)
				return
			}
		}
		if len(res) != len(v.Expect) {
			t.Fatalf("LENGTH")
			return
		}
		for i, j := range res {
			if j != v.Expect[i] {
				t.Fatalf("ELEMENT: '%v' != '%v'", j, v.Expect[i])
				return
			}
		}
	}
}

func Test002(t *testing.T) {
	var temp Strings_t
	temp.Set("default")

	t.Logf("TEMP: %v", temp)
	if len(temp) == 0 || temp[0] != "default" {
		t.Fatalf("DEFAULT")
	}
}

func Test003(t *testing.T) {
	l := NewLexer([]rune{','},
		[]rune{'\n'},
		[]rune{'\v', '\f', '\r', '\t', ' '},
		[]Quote_t{
			{'"', '"'},
			{'\'', '\''},
			{'«', '»'},
			{'[', ']'},
		},
	)
	reader := strings.NewReader("[\"test\"]")
	var res []string
	for {
		token, state := l.Next(reader)
		if state == STATE_NONE {
			break
		}
		res = append(res, token)
	}
	assert.Assert(t, len(res) > 0)
	assert.Assert(t, res[0] == "test", res[0])
}
