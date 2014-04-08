package main

import "os"
//import "flag"

//import "fmt"

import "gaccts"

func main () {
	gepics := os.Args[1:]
	//fmt.Println( os.Args[1:])
	gaccts.CreateGofiReport(gepics)
}
