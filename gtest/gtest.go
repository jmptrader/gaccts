package main

import "fmt"

import "gaccts"

func main() {
	fmt.Println("gtest version 01")
	//gaccts.ParseFile("etrans.txt")
	gaccts.ReadData()
	gaccts.TestGofi()
}
