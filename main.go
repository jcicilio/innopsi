// innopsi project main.go
// using https://godoc.org/github.com/montanaflynn/stats
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/montanaflynn/stats"
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

type scoreResult struct {
	dataSetId int
	score     float64
	d         []coreData
	rc        []rowCriteria
	t0        []coreData
	t1        []coreData
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

const subjects int = 240
const datasets int = 4
const rowThreshhold int = 10

var (
	data   []coreData
	levels [][]rowCriteria
)

// Sorting interface implementation for scoreResults
type scoreResults []scoreResult

func (s scoreResults) Len() int {
	return len(s)
}

func (s scoreResults) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s scoreResults) Less(i, j int) bool {
	return s[i].score < s[j].score
}

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

func outputScore(s scoreResult) {

	if s.score >= 0.0 {
		return
	}

	fmt.Printf("%d, %f, ", s.dataSetId, s.score)
	for _, each := range s.rc {
		fmt.Printf("x=%d, c=%d, ", each.r+1, each.c)
	}
	fmt.Println()

}

func outputScores(s []scoreResult) {
	for _, each := range s {
		outputScore(each)
	}
}

func negativeScoreCount(s []scoreResult) int {
	var count = 0
	for _, each := range s {
		if each.score < 0.0 {
			count += 1
		}
	}

	return count
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

func outputResults(s []scoreResult) {
	// Add one column for the id
	var dataSet [subjects][datasets + 1]int

	// Write headings, id, dataset_1 through n
	// Write ids
	for row := 0; row < subjects; row++ {
		dataSet[row][0] = row + 1
	}

	// With the top score per dataset
	// build and array of subject rows and dataset + 1 columns (1200 data, 1 id)
	// output the array
	for row := 0; row < subjects; row++ {
		for col := 0; col < datasets; col++ {
			// data is organized one dataset per column
			pv := resultsArray(s[col].t1)
			// with the scores for individuals
			for sub := 0; sub < subjects; sub++ {
				dataSet[sub][col+1] = pv[sub]
			}
		}
	}

	// Output to file
	f, _ := os.Create("/temp/output.csv")
	defer f.Close()

	// Write headings, id, dataset_1 through n
	for h := 0; h <= datasets; h++ {
		if h == 0 {
			f.WriteString("\"id\",")
		} else {
			var name = fmt.Sprintf("\"dataset_%d\"", h)
			f.WriteString(name)
			if h < datasets {
				f.WriteString(",")
			}
		}
	}

	f.WriteString("\r\n")

	// Write data
	for row := 0; row < subjects; row++ {
		for col := 0; col <= datasets; col++ {

			f.WriteString(strconv.Itoa(dataSet[row][col]))
			if col != datasets {
				f.WriteString(",")
			}
		}
		f.WriteString("\r\n")
	}

}

// Convert a sparse list of subject Id in a t1 array into a binary vector
// of 0 / 1 for each subject that matches
func resultsArray(positive []coreData) []int {
	var r []int

	// first add all zero values
	for row := 0; row < subjects; row++ {
		r = append(r, 0)
	}

	// create row vector
	for _, each := range positive {
		r[each.id-1] = 1
	}

	return r
}

// for a partition in the set of data, calculate the effective treatement
// score using (mean t1 - mean t0) / population standard deviation
func evalScore(d []coreData, rc []rowCriteria, dataSetId int) scoreResult {
	var s scoreResult
	s.rc = rc
	s.d = d
	s.score = 0
	s.dataSetId = dataSetId

	// check for minimum row threshhold
	if len(d) <= rowThreshhold {
		return s
	}

	var t0 []coreData
	var t1 []coreData
	var t0s []float64
	var t1s []float64
	var allTs []float64

	// Partition the data into treatment 0 and treatment 1
	// and save the score for evaluation
	for _, each := range d {
		// Save all responses for later SD calculation
		allTs = append(allTs, each.y)

		if each.treatment == 0 {
			t0 = append(t0, each)
			t0s = append(t0s, each.y)
		} else {
			t1 = append(t1, each)
			t1s = append(t1s, each.y)
		}
	}

	if len(t0) <= rowThreshhold/2 || len(t1) <= rowThreshhold/2 {
		return s
	}

	s.t0 = t0
	s.t1 = t1

	// then calculate the median, also experiment with average
	var mean0, _ = stats.Mean(t0s)
	var mean1, _ = stats.Mean(t1s)
	var sd, _ = stats.StandardDeviationPopulation(allTs)

	// subtract the two t0-t1, we want t1 to be smaller
	var meanValue = mean1 - mean0

	s.score = meanValue / sd

	return s
}

// sample criteria selection distributions
func fullOneLevel() [][]rowCriteria {
	var r [][]rowCriteria

	for x := 0; x < 40; x++ {
		// for each criteria
		for cr := 0; cr < 6; cr++ {
			var v []rowCriteria
			var k rowCriteria
			k.r = x
			k.c = cr
			v = append(v, k)
			r = append(r, v)
		}
	}

	return r
}

// sample criteria selection distributions
func fullTwoLevel() [][]rowCriteria {
	var r [][]rowCriteria

	for x := 0; x < 40; x++ {
		for cr := 0; cr < 6; cr++ {
			var v0 []rowCriteria
			var k0 rowCriteria
			k0.r = x
			k0.c = cr
			v0 = append(v0, k0)
			r = append(r, v0)

			for x1 := 0; x1 < 40; x1++ {
				// for each criteria
				for cr1 := 0; cr1 < 6; cr1++ {
					if x1 > x && cr1 > cr {
						var v1 []rowCriteria
						var k1 rowCriteria
						k1.r = x1
						k1.c = cr1
						v1 = append(v1, k1)
						v1 = append(v1, k0)
						r = append(r, v1)
					}
				}
			}
		}
	}

	return r
}

// run the testing
func levelEval(dataSetId int) []scoreResult {
	var (
		r []scoreResult
	)

	// Get the partition to work on
	d := partitionByDataset(dataSetId)

	// globally set
	src := levels

	for _, src1 := range src {
		t := partitionByRowCriteria(d, src1)
		s := evalScore(t, src1, dataSetId)
		r = append(r, s)
	}

	return r
}

func main() {
	readData()

	var scores []scoreResult
	levels = fullTwoLevel()

	for dataSetId := 1; dataSetId <= datasets; dataSetId++ {
		s := levelEval(dataSetId)
		sort.Sort(scoreResults(s))
		if len(s) > 0 {
			// pick the top score
			scores = append(scores, s[0])
			fmt.Printf("%d, %f \n", s[0].dataSetId, s[0].score)
		}
	}

	outputScores(scores)
	outputResults(scores)
}
