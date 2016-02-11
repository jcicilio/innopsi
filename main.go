// innopsi project main.go
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type coreData struct {
	dataset   int
	id        int
	treatment int
	y         float64
	xi        [40]int
}

var data []coreData

func readData() {
	csvfile, err := os.Open("./data/InnoCentive_9933623_Training_Data.csv")
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

	// The length of the dataset, minus the column headings
	dataSetLength := len(rawCSVdata) - 1
	data = make([]coreData, dataSetLength)

	// sanity check, display to standard output
	count := -1
	for _, each := range rawCSVdata {
		// skip the first row
		if count != -1 {
			//data[count] = new coreData()
			//fmt.Printf("%s, %s, %d \n", each[0], each[1], len(each))
			data[count].dataset, _ = strconv.Atoi(each[0])
			data[count].id, _ = strconv.Atoi(each[1])
			data[count].treatment, _ = strconv.Atoi(each[2])
			data[count].y, _ = strconv.ParseFloat(each[3], 10)

			// Get the integer values
			for i := 0; i < 20; i++ {
				// Offset of 4 into values
				data[count].xi[i], _ = strconv.Atoi(each[i+4])
			}

			// Get and convert the numeric values to three levels
			for i := 21; i < 40; i++ {
				// Offset of 4 into values
				t, _ := strconv.ParseFloat(each[i+4], 10)
				if t < 33.3 {
					data[count].xi[i] = 0
				} else if t >= 33.3 && t < 66.6 {
					data[count].xi[i] = 1
				} else {
					data[count].xi[i] = 2
				}
			}
		}

		count += 1
	}
}

func outputData() {
	for _, each := range data {
		fmt.Printf("%d, %d, %d, %f ", each.dataset, each.id, each.treatment, each.y)
		for _, x := range each.xi {
			fmt.Printf("%d, ", x)
		}
		fmt.Println()
	}
}

func main() {
	readData()
	outputData()
}
