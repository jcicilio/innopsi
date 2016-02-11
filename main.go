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

type rowCriteria struct {
	r int
	c int
}

// Evaluation Criteria

func criteria(v int, c int) bool {
	if c == 0 {
		return v == 0
	}

	if c == 1 {
		return v == 1
	}

	if c == 2 {
		return v == 2
	}

	if c == 3 {
		return v == 0 || v == 1
	}

	if c == 4 {
		return v == 0 || v == 2
	}

	if c == 5 {
		return v == 1 || v == 2
	}

	return false
}

var data []coreData

// Read in the raw data, and then classify numerica data into three values also
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

// Output a slice of data
func outputData(d []coreData) {
	for _, each := range d {
		fmt.Printf("%d, %d, %d, %f ", each.dataset, each.id, each.treatment, each.y)
		for _, x := range each.xi {
			fmt.Printf("%d, ", x)
		}
		fmt.Println()
	}
}

// Get a partition of the dataset
func partitionByDataset(dataSetId int) []coreData {
	var r []coreData

	for i := 0; i < len(data); i++ {
		if data[i].dataset == dataSetId {
			r = append(r, data[i])
		}
	}

	return r
}

// each row has an associated criteria
func partitionByRowCriteria(d []coreData, rc []rowCriteria) []coreData {
	var r []coreData

	// for each data row see that it matches each row criteria
	// and if it does append to output
	for _, each := range d {
		var b = true
		// check that each row criteria value matches
		for _, k := range rc {
			if !(criteria(each.xi[k.r], k.c)) {
				b = false
			}
		}

		if b {
			r = append(r, each)
		}
	}

	return r
}

func main() {
	readData()
	d2 := partitionByDataset(2)
	//outputData(d2)

	var rc, rc1 rowCriteria
	rc.c = 0
	rc.r = 0
	rc1.r = 39
	rc1.c = 1

	var rcData []rowCriteria
	rcData = append(rcData, rc)
	rcData = append(rcData, rc1)

	t2 := partitionByRowCriteria(d2, rcData)
	outputData(t2)
}
