//
//
//

package fields

import "testing"

func Test01(t * testing.T) {
	test := ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4', ,aaa'5'bbb , this is test   "
	res, err := SplitCSV(test)
	if err != nil {
		t.Logf("ERROR: %v", err)
		return
	}
	t.Logf("test = %v", test)
	t.Logf("res  = %#v", res)
	if len(res) != 12 {
		t.Fail()
	}
}

func Test02(t * testing.T) {
	test := ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4', ,aaa'5'bbb , this is test   ,"
	res, err := SplitCSV(test)
	if err != nil {
		t.Logf("ERROR: %v", err)
		return
	}
	t.Logf("test = %v", test)
	t.Logf("res  = %#v", res)
	if len(res) != 13 {
		t.Fail()
	}
}

func Test03(t * testing.T) {
	test := ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4', ,aaa'5'bbb , this is test   ,''"
	res, err := SplitCSV(test)
	if err != nil {
		t.Logf("ERROR: %v", err)
		return
	}
	t.Logf("test = %v", test)
	t.Logf("res  = %#v", res)
	if len(res) != 13 {
		t.Fail()
	}
}

func Test04(t * testing.T) {
	test := ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4', ,aaa'5'bbb , this is test   ,'"
	_, err := SplitCSV(test)
	if err == nil {
		t.Fail()
	}
}

func Test05(t * testing.T) {
	test := ", , ',1,',,' 2 ', , ' ,3 ,', , aaa'4', ,aaa'5'bbb , this is test   ,"
	res, err := SplitCSV(test)
	if err != nil {
		t.Logf("ERROR: %v", err)
		return
	}
	t.Logf("test = %v", test)
	t.Logf("res  = %#v", res)
	if len(res) != 13 {
		t.Fail()
	}
}
