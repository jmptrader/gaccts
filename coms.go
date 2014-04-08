/*
This is being references from the gaccts/coms.go package. 
*/
package gaccts

//import "container/list"
//import "encoding/csv"
import  (
	"fmt"
	//import "io"
	//import "encoding/json"
	//"math"
	//import "os"
	"path/filepath"
	"sort"
	//"strings"

	//"money"
)

const (
Yafi = iota
Oeic
Indx
)

const ( P = 0.01 ; GBP = 1.0 ; NIL = 0.0)

type Com struct {
	Sym string // 1
	live bool  // 2 - should use qty count instead
	Ctype int  // 3
	scale float64 // 4
	//exchange string
	Gepic string // 5
	Qty float64
	Ucost float64 
	Ubefore float64 // TODO this is never set
	Uvalue float64 // 7
	Download bool
	name string // 8
}

var Coms []Com




func ParseComsFile() {
	coms := ParseFile(dataDir + "coms.txt")
	processComFile(coms)
}

func DownloadComs() (jmaps []Jmap) {
	gepics := make([]string, len(Coms))
	mask := make([]bool, len(Coms))
	for i,c := range Coms {
		gepics[i] = c.Gepic 
		mask[i] = c.Download
	}

	jmaps = GetGepicsMasked(gepics, mask)
	return
}

// Load latest Com values from file, as opposed to downloading them
func RetrieveComs() {
	gofiDir := rootDir + ".accts/gofi/"	
	gofiFiles, err := filepath.Glob (gofiDir+ "2*.txt")
	if err != nil { panic(err) }
	//sortedGofiFiles := gofiFiles.sort()
	//fmt.Println(sortedGofiFiles)
	sort.Sort(sort.StringSlice(gofiFiles))
	//fmt.Println(gofiFiles)
	ngfiles := len(gofiFiles)
	if ngfiles == 0 { panic("Can't find any prices") }	
	gofiFile := gofiFiles[ngfiles-1] // use last one
	prices := ParseFile(gofiFile)
	//fmt.Println(prices)
	for _, rec := range prices {
		//fmt.Println(rec[0])
		sym := rec[3]
		var com *Com
		com = FindCom(sym, false)
		if com != nil {
			uvalue := AsFloat64(rec[4])
			com.Uvalue = uvalue
		}
	}
}


func processComFile (records [][]string) {
	Coms = make([]Com, len(records))
	for i, rec := range records {
		com := &Coms[i]
		com.Sym = rec[1]

		//live:= rec[1]
		switch rec[2] {
		case "LIVE": com.live = true
		case "DEAD": com.live = false
		default: panic("Expecting either LIVE or DEAD")
		}

		switch rec[3] { 
		case "YAFI":  com.Ctype = Yafi
		case "OEIC": com.Ctype = Oeic
		case "INDX": com.Ctype =Indx
		default: panic("Expecting either YAFI, OEIC or INDX")
		}

		switch rec[4] {
		case "P": com.scale = P
		case "GBP": com.scale = GBP
		case "NIL": com.scale = NIL
		default: panic("Expecting unit of P, GBP or NIL")
		}

		//arr := strings.Split(rec[5] , ":")
		//Coms[i].exchange = arr[0]
		//Coms[i].gepic = arr[1]
		com.Gepic = rec[5]

		uvalue := rec[7]
		com.Download = uvalue == "?"
		if ! com.Download { com.Uvalue = AsFloat64(uvalue) }


		com.name = rec[8]
	}
}
/*
// Returns -1 if sym not found	
func FindCom(sym string) (index int) {
	if len(Coms) == 0 {
		fmt.Println("No commodities. Did you load them?")
	}

	index = -1
	//fmt.Println(Coms)
	for index := 0 ; index < len(Coms) ; index++ {
		//fmt.Println(Coms[index])
		if Coms[index].Sym == sym { return index}
	}

	//panic( "Couldn't find sym " + sym)
	return index
}

func FindComOrDie(sym string) (index int) {
	index = FindCom(sym)
	if index == -1 { panic( "Couldn't find sym " + sym)}
	return index
}
*/

func FindComIndex(sym string) (index int) {
	for index = 0 ; index < len(Coms) ; index++ {
		if Coms[index].Sym == sym { return index }
	}
	return -1
}

func FindCom(sym string, die bool) (com *Com) {
	if len(Coms) == 0 {
		fmt.Println("No commodities. Did you load them?")
	}

	//for index := 0 ; index < len(Coms) ; index++ {
	//	if Coms[index].Sym == sym { return &Coms[index]}
	//}
	index := FindComIndex(sym)
	if index > -1 { return &Coms[index] }

	if die {
		panic( "Couldn't find sym " + sym)
	}
	return nil
}

func transactCom (sym string, Qty , value float64) {
	var com *Com // THIS SEMS TO BE CORRECT NOW
	com = FindCom(sym, true)
	if Qty > 0 { // buy
		cost := value + com.Qty * com.Ucost * com.scale
		//com.uvalue = value / Qty // useful in case 
		com.Qty = com.Qty + Qty // need to be calculated after working out today cost
		com.Ucost = cost / com.Qty / com.scale
	} else { // sell
		//fmt.Println(com)
		com.Qty = Qty + com.Qty // rememeber Qty is negative
		// NB don't adjust Ucost when there's a sale!
		if com.Qty == 0 { com.Ucost = 0.0 }
		if com.Qty < 0 { 
			fmt.Println(com)
			fmt.Println(com.Qty)
			//fmt.Println("TODO I should probably panic here")
			panic("Can't have negative quantities") 
		}
	}
	//fmt.Println(com)
}


func QtyAndPriceToMoney(qty float64, price float64) Money {
	v := qty*price
	return Money(v*100.0 +0.5)
}
	
