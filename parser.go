package gaccts

import (
	"container/list"
	"fmt"
	//"os"
	//"path/filepath"
	//"runtime"
	//"sort"
)
/*
type Command struct { name string ; nargs int }

var Commands =
[...]Commands{
{"etran", 

*/

func ParseFile(filename string) (records [][]string) {
	fmt.Println("Parsing file", filename)
	lexemes := list.New()
	Lexify(filename, lexemes)

	lexListEl := lexemes.Front()
	lexSlice := make([]Lexeme, lexemes.Len())
	//fmt.Println(len(lexSlice))
	for i := 0; i< len(lexSlice); i++ {
		//fmt.Println(i)
		//var v Lexeme
		//e := lexListEl.Next()
		v := lexListEl.Value.(*Lexeme)
		lexSlice[i].ltype  = v.ltype
		lexSlice[i].line = v.line
		lexSlice[i].text = v.text
		//fmt.Println(v)
		lexListEl = lexListEl.Next() 
	}

	const ( scanning = iota )
	pos0 :=0
	for i, lex := range lexSlice {
		if (lex.ltype == ltfeed || lex.ltype == lteof) {
			if i != pos0 { 
				record := make([]string, i - pos0)
				for j := pos0; j< i; j++ {
					text := lexSlice[j].text
					//fmt.Println("field", text)
					record[j-pos0] = text
				}
				//fmt.Println("record", record)
			}
			pos0 = i+1
		}
	}

	type genfn func(a,b int)

	

	gen := func ( f genfn) {
		pos0 :=0
		for i, lex := range lexSlice {
			if (lex.ltype == ltfeed || lex.ltype == lteof) {
				if i != pos0 { f(pos0, i) }
				pos0 = i+1
			}
		}
	}


	numrecords := 0 ; 
	fcount := func(a, b int) {numrecords++}
	gen(fcount)
	//fmt.Println(numrecords)

	records = make([][]string, numrecords)
	recordnum := 0
	faddrecord := func(a, b int) {
		record := make([]string, b-a)
		for i:= a; i<b; i++ { record[i-a] = lexSlice[i].text }
		records[recordnum] = record
		recordnum++
	}
	gen(faddrecord)

	return
}




