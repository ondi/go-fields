//
//
//

package fields

import "testing"

type Data_t struct {
	In string
	Out []string
	Err bool
}

var data = []Data_t {
	{
		In: ",",
		Out: []string{"", ""},
	},
	{
		In: ",,",
		Out: []string{"", "", ""},
	},
	{
		In: " , , ",
		Out: []string{"", "", ""},
	},
	{
		In: "   1   ,",
		Out: []string{"1", ""},
	},
	{
		In: ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4', ,aaa'5'bbb , this is test   ",
		Out: []string{"", "", "',1,'", "", "' 2 '", "", "' ,3 ,'", "", "aaa4", "", "aaa5bbb", "this is test"},
	},
	{
		In: ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4', ,aaa'5'bbb , this is test   ,",
		Out: []string{"", "", "',1,'", "", "' 2 '", "", "' ,3 ,'", "", "aaa4", "", "aaa5bbb", "this is test", ""},
	},
	{
		In: ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4', ,aaa'5'bbb , this is test   ,''",
		Out: []string{"", "", "',1,'", "", "' 2 '", "", "' ,3 ,'", "", "aaa4", "", "aaa5bbb", "this is test", ""},
	},
	{
		In: ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4', ,aaa'5'bbb , this is test   ,'",
		Out: []string{"", "", "',1,'", "", "' 2 '", "", "' ,3 ,'", "", "aaa4", "", "aaa5bbb", "this is test"},
		Err: true,
	},
}

func Test001(t * testing.T) {
	for _, v := range data {
		res, err := SplitCSV(v.In)
		t.Logf("In   = %v", v.In)
		t.Logf("Out1 = %#v", v.Out)
		t.Logf("Out2 = %#v", res)
		if v.Err {
			if err == nil {
				t.Fatalf("NO ERROR FOUND")
			}
		} else {
			if err != nil {
				t.Fatalf("ERROR: %v", err)
			}
		}
		if len(res) != len(v.Out) {
			t.Fatalf("LENGTH")
		}
		for i, j := range res {
			if j != v.Out[i] {
				t.Fatalf("ELEMENT: '%v' != '%v'", j, v.Out[i])
			}
		}
	}
}
