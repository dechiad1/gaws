package util_test

import (
	"github.com/dechiad1/gaws/util"
	"testing"
)

func TestSetHeaders(t *testing.T) {
	du := util.SetHeaders("carb", "arnone")

	if len(du.Header) != 2 {
		t.Error("expected 2, got: ", du.Header)
	}

	if len(du.Body) != 0 {
		t.Error("expected 0 initial rows, got:", len(du.Body))
	}

}

func TestAddRow(t *testing.T) {
	du := util.SetHeaders("don", "miguel", "the", "conquistador")
	result := du.AddRow("1", "2", "3", "4")

	if !result {
		t.Error("expected true response from valid add")
	}

	if len(du.Body[0]) != 4 {
		t.Error("expected first row of length 4, got:", du.Body[0])
	}

	if len(du.Body) != 1 {
		t.Error("expect 1 row is appended", du.Body)
	}

	if du.Row_indicator != 1 {
		t.Error("expected indicator to increase to mark addition of row")
	}
}

func TestAddRow_false(t *testing.T) {
	du := util.SetHeaders("don", "miguel", "the", "conquistador")
	result := du.AddRow("1", "2", "3", "4", "5")

	if result {
		t.Error("expected false response from invalid add")
	}
}

/*
Testing the print display

func TestPrintDisplay(t * testing.T) {
	du := util.SetHeaders("don", "miguel", "the", "conquistador")
	du.AddRow("1","20","3","muuust")
	du.AddRow("asdfasfdasdf","asdfafdasdfaswerqwerdf","09080809","muuust")
	du.PrintDisplay()
}
*/
