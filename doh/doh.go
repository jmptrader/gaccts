package main


import "fmt"

import "gaccts"

func main() {
	fmt.Println("gaccts/doh v03")

	gaccts.Init()
	gaccts.ParseComsFile()
	ntrans := gaccts.ParseNtransFile()
	fmt.Println(ntrans)

	gaccts.ParseEtransFile()


	//gaccts.ReadData()
	fmt.Println("TODO Now sometimes DownloadComs, sometimes RetrieveComs")
	gaccts.RetrieveComs()
	
}
