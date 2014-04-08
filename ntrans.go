package gaccts

type Ntran struct {
	Dstamp string
	Dr string
	Cr string
	Amount int64
	Desc string
}

func ParseNtransFile() (ntrans []Ntran) {
	recs := ParseFile(dataDir + "ntrans.txt")
	ntrans = make([]Ntran, len(recs))
	for i, rec := range recs {
		ntran := &ntrans[i]
		ntran.Dstamp = rec[1]
		ntran.Dr = rec[2]
		ntran.Cr = rec[3]
		
		amount := int64(AsFloat64(rec[4]) * 100.00 + 0.5)
		ntran.Amount = amount
		
		ntran.Desc = rec[5]
	}
	return ntrans
}
