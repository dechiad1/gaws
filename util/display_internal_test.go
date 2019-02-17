package util

import (
	"testing"
)

func TestSetColumnBorder(t *testing.T) {
	du := SetHeaders("don", "miguel", "the", "conquistador")
	du.AddRow("1","20","3","muuust")
	du.AddRow("crust","2","caaarb","4")
	du.setColumnBorder()

	if du.Column_border != 6 {
		t.Error("expect 6 longest characters in first column, but got", du.Column_border)
	}
}

func TestBuildFmtString(t *testing.T) {
	fmtString := buildFmtString(3)
	expected := "%s\t%s\t%s\n"
	if fmtString != expected {
		t.Error("expected [1] but got [2]: ", expected, fmtString)
	}
}

func TestCreateInters(t *testing.T) {
	s := []string {"afasdf","asfasd"}
	i := createInters(s)

	if len(i) != len(s) {
		t.Error("expected interface slice of [1] but got len [2]", len(i), len(s))
	}

	//test that each item in the []interface{} contains a concrete value of type string
	for k,v := range i {
		_, ok := v.(string)
		if !ok {
			t.Error("expecting all elements converted can be asserted to strings. element of index [1] cant: ", k, v)
		}
	}
}

func TestCreateHeaderBorders(t *testing.T) {
	s := []string {"11111", "222222"}
	b := createHeaderBorders(s)

	for k:=0; k<len(s); k ++ {
		if len(s[k]) != len(b[k]) {
			t.Error("expecting equal length strings for input and output but didnt get that for index",k)
		}
	}
}

