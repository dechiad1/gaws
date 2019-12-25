package util

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type DisplayUnit struct {
	Header        []string
	Row_len       int //how long is the header col - for add row indicator
	Body          [][]string
	Row_indicator int //current row that is populated
	Column_border int //longest field in column - will act as col length in display
}

func SetHeaders(cols ...string) *DisplayUnit {
	du := DisplayUnit{
		Header:  cols,
		Row_len: len(cols),
		//start with 5 rows
		Body:          make([][]string, 0),
		Row_indicator: 0,
		Column_border: 0, //make([]int, len(cols)),
	}
	return &du
}

func (du *DisplayUnit) AddRow(cols ...string) bool {
	if len(cols) != du.Row_len {
		return false
	} else {
		du.Body = append(du.Body, cols)
		du.Row_indicator += 1
		return true
	}
}

func (du *DisplayUnit) setColumnBorder() {
	//find the longest string in each column - except for the last column as there is nothing after it
	for row := 0; row < (du.Row_len - 1); row++ {
		longest := 0
		for col := 0; col < du.Row_indicator; col++ {
			temp := len(du.Body[col][row])
			if temp > longest {
				longest = temp
			}
		}
		du.Column_border = longest
	}

	//check back to compare against the length of the header items - except for the last column as there is nothing after it
	for row := 0; row < (du.Row_len - 1); row++ {
		if du.Column_border < len(du.Header[row]) {
			du.Column_border = len(du.Header[row])
		}
	}
}

func buildFmtString(f int) string {
	var fmtString strings.Builder
	for i := 0; i < f; i++ {
		fmtString.WriteString("%s")
		if i < (f - 1) {
			fmtString.WriteString("\t")
		}
	}
	fmtString.WriteString("\n")
	return fmtString.String()
}

func (du *DisplayUnit) PrintDisplay() {
	du.setColumnBorder()
	//param 2: minimum col width is the longest word in the table
	//param 3: tabstop of 8
	//param 4: 1 character buffer after col
	w := tabwriter.NewWriter(os.Stdout, du.Column_border, 8, 1, '\t', 0)
	fmtString := buildFmtString(du.Row_len)

	//format the table head
	border := createHeaderBorders(du.Header)
	inter_head := createInters(du.Header)
	inter_border := createInters(border)
	fmt.Fprintf(w, fmtString, inter_head...)
	fmt.Fprintf(w, fmtString, inter_border...)

	//format the table body
	for i, row := range du.Body {
		if i < du.Row_indicator {
			inter := createInters(row)
			fmt.Fprintf(w, fmtString, inter...)
		}
	}
	w.Flush()
}

func createInters(cols []string) []interface{} {
	inter := make([]interface{}, len(cols))
	for i, v := range cols {
		inter[i] = v
	}
	return inter
}

func createHeaderBorders(cols []string) []string {
	result := make([]string, len(cols))
	r := rune('-')
	for i, col := range cols {
		l := len(col)
		var s strings.Builder
		for i := 0; i < l; i++ {
			//write hyphen as rune, int32 representation of ascii char
			s.WriteRune(r)
		}
		result[i] = s.String()
	}
	return result
}
