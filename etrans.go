package gaccts

import "fmt"

type Etran struct {
	Dstamp string
	Folio string
	Vbefore Money
	Vflow Money
	Vprofit Money
	Vto Money
}

func processEtransFile(records [][]string) (etrans []Etran) {

	fmt.Println("Processing etrans")

	// reset number of coms to 0
	//for _, com:= range Coms { com.qty = 0.0 ; fmt.Println(com.sym, com.qty)}
	// should already be 0
	//fmt.Println("Mojo?")
	etrans = make([]Etran, len(records))
	for i, rec:= range records {
		// assume they are all etrans
		f64 := func (s string) float64 {
			v, err := Enfloat64(s)
			if err != nil {panic(err) } // this pattern needs refactoring
			return v
		}

		dstamp := rec[1]
		way := rec[2]
		var sgn float64 
		if way == "B" {sgn = 1.0 
		} else if way == "S" { sgn = -1.0
		} else { panic("way should be either B or S") }

		folio := rec[3]
		sym := rec[4]
				
		qty := rec[5]
		qty64 := sgn * f64(qty)
		
		value := f64(rec[6])

//fmt.Println("AAA")
		transactCom(sym, qty64, value)
		
		etran := &etrans[i]
		var com *Com
		com = FindCom(sym, true)
		if dstamp <= before() {
			etran.Dstamp = before()
			etran.Vbefore = QtyAndPriceToMoney(qty64, com.Ubefore)
			// etran.vflow
		} else {
			etran.Dstamp = dstamp
			// etran.vbefore - unneeded - golang autosets to 0
			etran.Vflow = Float64ToMoney(value)
		}
		etran.Folio = folio
		etran.Vto = QtyAndPriceToMoney(qty64, com.Uvalue)

		
		
	}

	//fmt.Println("Computed quantities")
	//for _, com:= range Coms1 { fmt.Println(com.sym, com.qty) }
	return

}


func ParseEtransFile() {
	recs := ParseFile(dataDir + "etrans.txt")
	processEtransFile(recs)
}
