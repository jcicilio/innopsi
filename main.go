// innopsi project main.go
package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type coreData struct {
	dataset   int
	id        int
	treatment int
	y         float64
	xi        [20]int
	xd        [20]float64
}

func readData() {
	csvfile, err := os.Open("./InnoCentive_9933623_Training_Data.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer csvfile.Close()
	// Create a new reader.
	reader := csv.NewReader(csvfile)
	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// sanity check, display to standard output
	for _, each := range rawCSVdata {
		fmt.Printf("%s, %s, %d \n", each[0], each[1], len(each))
	}
}

func main() {
	readData()
}
