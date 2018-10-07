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
		Expect: []string{},
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
		Input: ",",
		Expect: []string{"", ""},
	},
	{
		Input: "'',",
		Expect: []string{"", ""},
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
		Input: "«»,»«, « »",
		Expect: []string{"«»", "»«", "« »"},
	},
	{
		Input: "aaa'bbb'ccc,",
		Expect: []string{"aaa'bbb'ccc", ""},
	},
	{
		Input: ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4' , ,aaa'5'bbb , this is test   ",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "this is test"},
	},
	{
		Input: ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4' , ,aaa'5'bbb , this is test   ,",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "this is test", ""},
	},
	{
		Input: ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4' , ,aaa'5'bbb , this is test   ,''",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "this is test", ""},
	},
	{
		Input: ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4' , ,aaa'5'bbb , this is test   ,'",
		Expect: []string{"", "", ",1,", "", " 2 ", "", " ,3 ,", "", "aaa'4'", "", "aaa'5'bbb", "this is test"},
		Err: true,
	},
}

func Test001(t * testing.T) {
	for _, v := range data1 {
		res, err := SplitCSV(v.Input)
		t.Logf("Input  = %v", v.Input)
		t.Logf("Expect = %#v", v.Expect)
		t.Logf("Result = %#v", res)
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
