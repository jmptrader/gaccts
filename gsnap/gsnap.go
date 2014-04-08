package main


import (
	"errors"
	"fmt"
	"math"

	"gaccts"
)

func main() {
	fmt.Println("gaccts/gsnap v 04")
	gaccts.Init()
	gaccts.ParseComsFile()
	gaccts.ParseEtransFile()
	jmaps := gaccts.DownloadComs()
	//gaccts.Snapshot()
	tabulate(jmaps)
}


/*
func SnapshotXXX() {
	
	messages := make(chan QuoteResponse)

	//ReadData()

	for i , com := range Coms {
		go getQuote(messages, i, com.gepic, com.exchange)
		//fmt.Println(com.name)
	}

	// assemble the responses in the right order
	//responses := make([]QuoteResponse, len(Coms))
	responses := make([][]byte, len(Coms))
	for i:=0; i< len(Coms); i++ {
		qr := <- messages
		responses[qr.index] = qr.response
	}

	// now extract prices, which is what we really wanted
	jsons := make([]Jmap, len(Coms))
	//for i, jn := range jsons  {
	for i:= 0 ; i< len(Coms); i++ {
		jsons[i] = mapifyResponseOrDie(Coms[i].gepic, responses[i])
	}


	///tabulate(jsons)

}
*/

func tabulate(jmaps []gaccts.Jmap) {
	tvalue := 0.0
	tprofit := 0.0
	//   " 0  ALL   3812   108.8678   116.0000    4150.04     -76.24 -1.69"
	hdr:=" #  SYM    QTY      UCOST     UVALUE      VALUE     PROFIT  CHG%"
	fmt.Println(hdr)

	index := 0
	// TODO gofi.go should do most of this decoding
	for i, com := range gaccts.Coms {
		if com.Qty >0 && com.Ctype == gaccts.Yafi {
			index++

			bail := func (e error) { 
				if e != nil { 
					fmt.Println(com) 
					panic(e) }}
			jn := jmaps[i]
			if jn == nil { bail(errors.New("jmap is nil")) }

			val64 := func (key string) float64 {			
				s := jn[key].(string)
				v, err := gaccts.Enfloat64(s)
				bail(err)
				return v
			}
			com.Uvalue = val64("l")


			value := com.Qty * com.Uvalue / 100.0
			tvalue += value
			chgpc := val64("cp")
			chg := val64("c")
			profit := chg * com.Qty /100.0
			tprofit += profit
			f := "%2d %4.4s %6.0f %10.4f %10.4f %10.2f %10.2f %5.2f"
			s1 := fmt.Sprintf(f, index, com.Sym, com.Qty, com.Ucost, 
				com.Uvalue, value, profit, chgpc)
			fmt.Println(s1)
			if math.Mod(float64(index), 3) == 0.0 { fmt.Println("") }
		}

	}
	

	gainpc := 100.0 * tprofit / (tvalue - tprofit)
	fmt.Println(fmt.Sprintf("%36.36s %10.2f %10.2f %5.2f", "", 
		tvalue, tprofit, gainpc))

	// TODO surely we can use coms.go for most of this
	/*
	messages := make(chan QuoteResponse)	
	go getQuote(messages, 0, "ASX", "INDEXFTSE")
	qr := <- messages
	jn := mapifyResponseOrDie("ASX", qr.response)
*/
	i := gaccts.FindComIndex("FTAS")
	jn := jmaps[i]
	c, err := gaccts.Enfloat64(jn["c"].(string))
	if err != nil { panic(err) }
	l, err := gaccts.Enfloat64(jn["l"].(string))
	if err != nil { panic(err) }
	prev := l-c
	f := "FTAS: Prev: %.2f Curr: %.2f Chg%%: %.2f"
	s := fmt.Sprintf(f, prev, l,  c/prev * 100.0)
	fmt.Println(s)

}

