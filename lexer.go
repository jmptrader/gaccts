// Parsing functions for mu standard file format
// You can test it using gtest

/*

This is being referenced from the gaccts/lexer.go file.

*/

package gaccts

import "container/list"
//import "fmt"
import "io/ioutil"

const (
	ltfeed = iota // implies line feed
	lttext
	lteof
)

// This is a comment for the Lexeme Struct
type Lexeme struct {
	ltype int
	line int // line number
	text string // textual representation of the thing being passed
}

// some Lexify comment belongs here
func Lexify(filename string, lexemes *list.List) {
	bytes, err := ioutil.ReadFile(filename)
	if err!= nil { panic(err) }

	linenum := 1

	//lexemes = list.New() // make([]Lexeme, 0)
	addLexeme := func (altype int, atext string) {
		//var l Lexeme
		l  := new(Lexeme)
		l.ltype = altype
		l.line = linenum
		l.text = atext
		lexemes.PushBack(l)
	}


	pos := 0
	//var start int 
	//start = make(int)
	
	blen := len(bytes)
	b  := byte(0)
	pos0 := 0
	set0 := func () {pos0 = pos }
	//foo := func () { fmt.Println("foo here") }
	//foo()

	pushText := func () {
		word := string(bytes[pos0:pos])
		addLexeme(lttext, word)
	}

	pushString := func (pos1 int) {
		text := string(bytes[pos0+1:pos1])
		addLexeme(lttext, text)
	}

/*
	set1 := func(b byte) bool {
		return b == '\r' || b == '\n'
	}
*/

	set2 := func(b byte) bool { return b == '\r' || b == '\n' || b == '#' || b == ' ' || b == '\t'}
	

sniffing:
	set0()
	if pos >=blen {goto finished }
	b = bytes[pos]
	// fmt.Println(b)
	switch b {
	case '\t', ' ': goto eat_whiting
	case '#': goto eat_hash 
	case '\r', '\n': addLexeme(ltfeed, "") ; linenum++
	case '"': goto stringing
	default: goto wording
	}

	pos++
	goto sniffing

eat_whiting:
	pos++
	if pos >= blen{goto finished}
	b = bytes[pos]
	if b == ' ' ||  b == '\t' {goto eat_whiting} else {goto sniffing}

eat_hash:
	pos++
	if pos >= blen{goto finished}
	b = bytes[pos]
	if b == '\r' || b == '\n' {goto sniffing} else {goto eat_hash}

wording:
	pos++
	if pos >=blen {pushText() ; goto finished}
	b = bytes[pos] 
	if set2(b) {pushText() ; goto sniffing} else {goto wording}
	
stringing:
	pos++
	if pos >=blen {pushString(pos) ; goto finished}
	b = bytes[pos]
	pos++
	if (b == '"'|| b == '\r' || b == '\n') {pushString(pos-1) ; goto sniffing}
	pos--
	goto stringing
	

finished:
	//fmt.Println("Job's a good'un")
	addLexeme(lteof, "EOF")
	//b = bytes[pos]
	return

}

