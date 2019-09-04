//
//
//

package fields

import "testing"

type Data_t struct {
	Input string
	Expect []string
	Err bool
}

var data1 = []Data_t {
	{
		Input: "",
		Expect: []string{""},
	},
	{
		Input: "''",
		Expect: []string{""},
	},
	{
		Input: "1",
		Expect: []string{"1"},
	},
	{
		Input: "'1'",
		Expect: []string{"1"},
	},
	{
		Input: "''1''",
		Expect: []string{"1''"},
	},
	{
		Input: "(1,2,3)",
		Expect: []string{"1,2,3"},
	},
	{
		Input: "((1,2,3))",
		Expect: []string{"(1,2,3)"},
	},
	{
		Input: "(((1,2,3)))",
		Expect: []string{"((1,2,3))"},
	},
	{
		Input: "(((1,2,3",
		Expect: []string{"((1,2,3"},
	},
	{
		Input: ",",
		Expect: []string{"", ""},
	},
	{
		Input: "'', it's test",
		Expect: []string{"", "it's test"},
	},
	{
		Input: "'', 'it's test'",
		Expect: []string{"", "its test'"},
	},
	{
		Input: "'',''",
		Expect: []string{"", ""},
	},
	{
		Input: ",,",
		Expect: []string{"", "", ""},
	},
	{
		Input: "«», »«, « »",
		Expect: []string{"«»", "»«", "« »"},
	},
	{
		Input: "aaa'b'b'b'ccc,",
		Expect: []string{"aaa'b'b'b'ccc", ""},
	},
	{
		Input: ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4' , ,aaa'5'bbb , it's test   ",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "it's test"},
	},
	{
		Input: ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4' , ,aaa'5'bbb , it's test   ,",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "it's test", ""},
	},
	{
		Input: ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4' , ,aaa'5'bbb , it's test   ,   ",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "it's test", ""},
	},
	{
		Input: ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4' , ,aaa'5'bbb , it's test   ,''",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "it's test", ""},
	},
	{
		Input: ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4' , ,aaa'5'bbb , it's test   ,'",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "it's test", ""},
	},
}

func Test001(t * testing.T) {
	for _, v := range data1 {
		s := NewLexer([]rune{','}, []rune{'\v', '\f', '\r', '\n', '\t', ' '}, []Quote_t{Quote_t{'"', '"'}, Quote_t{'\'', '\''}, Quote_t{'(', ')'}})
		res, err := s.Split(v.Input)
		t.Logf("Input  = %v", v.Input)
		t.Logf("Expect = %#v, Error=%v", v.Expect, v.Err)
		t.Logf("Result = %#v, Error=%v", res, err)
		if v.Err {
			if err == nil {
				t.Fatalf("NO ERROR FOUND")
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
