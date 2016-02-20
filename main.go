// innopsi project main.go
// using https://godoc.org/github.com/montanaflynn/stats
package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

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
	switch c {
	case 0:
		return v == 0
	case 1:
		return v == 1
	case 2:
		return v == 2
	case 3:
		return v == 0 || v == 1
	case 4:
		return v == 0 || v == 2
	case 5:
		return v == 1 || v == 2
	case 6:
		return v == 0 || v == 1 || v == 2
	}

	return false
}

const subjects int = 240
const rowThreshhold int = 10
const maxCriteria = 6

// Actual set max will be two more than this, see rand function
const rand_maxSetMembers = 11
const rand_numSets = 100000

const validCriteriaThreshhold = -0.03

const datasets int = 1200

//const datasets int = 4

const datafilename string = "./data/InnoCentive_9933623_Data.csv"

//const datafilename string = "./data/InnoCentive_9933623_Training_Data.csv"

var (
	data     []coreData
	levels   [][]rowCriteria
	levelOne [][]rowCriteria
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

	csvfile, err := os.Open(datafilename)
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

			// Get and convert the numeric values to five levels
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

	//if s.score >= 0.0 {
	//	return
	//}

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

func outputRowCriteria(c [][]rowCriteria) {
	for _, each := range c {
		for _, v := range each {
			fmt.Printf("%d, %d, ", v.r, v.c)
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
				break
			}

			// temp, just the ordinal data
			//if k.r >= 20 {
			//	b = false
			//}
		}

		if b {
			r = append(r, each)
		}
	}

	return r
}

func outputResults(s []scoreResult) {
	var zeroVector, pv []int
	// first add all zero values
	for row := 0; row < subjects; row++ {
		zeroVector = append(zeroVector, 0)
	}
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
			//if s[col].score > validCriteriaThreshhold {
			//	pv = zeroVector
			//} else {
			pv = resultsArray(s[col].t1)
			//}
			// with the scores for individuals
			for sub := 0; sub < subjects; sub++ {
				dataSet[sub][col+1] = pv[sub]
			}
		}
	}

	// Output to file
	t := time.Now()
	filename := fmt.Sprintf("./results/o_%d%02d%02dT%02d%02d%02d.csv",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	f, _ := os.Create(filename)

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

	//fmt.Printf("positive: %d", len(positive))

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
	//var meanAll, _ = stats.Mean(allTs)
	var sd, _ = stats.StandardDeviationPopulation(allTs)

	// subtract the two t0-t1, we want t1 to be smaller
	// Note: use spooled
	// square root of ((Nt-1)St^2 + (Nc-1)Sc^2)/(Nt+Nc))
	//var St, _ = stats.StandardDeviation(t1s)
	//var Sc, _ = stats.StandardDeviation(t0s)
	//	var Nt = float64(len(t1s))
	//	var Nc = float64(len(t0s))
	//	var sPooled = math.Sqrt((Nt-1)*square(St) + (Nc-1)*square(Sc)/(Nt+Nc))

	//var _, t1confh = NormalConfidenceInterval(t1s)
	//var _, t0confh = NormalConfidenceInterval(t0s)

	//var meanValue = mean1 - mean0
	//var meanValue = (mean1/St - mean0/Sc) / sPooled

	// Score Type 1
	var meanDifference = mean1 - mean0
	//s.score = meanDifference / meanAll
	//var max, _ = stats.Max(allTs)

	s.score = meanDifference / sd

	//s.score = (mean1/St - mean0/Sc) / St

	return s
}

func square(f float64) float64 {
	return f * f
}

// https://github.com/hermanschaaf/stats/blob/master/stats.go
func NormalConfidenceInterval(nums []float64) (lower float64, upper float64) {
	conf := 1.95996 // 95% confidence for the mean, http://bit.ly/Mm05eZ
	mean, _ := stats.Mean(nums)
	dev, _ := stats.StandardDeviation(nums)
	dev = dev / math.Sqrt(float64(len(nums)))
	return mean - dev*conf, mean + dev*conf
}

// sample criteria selection distributions
func fullOneLevel() [][]rowCriteria {
	var r [][]rowCriteria

	for x := 0; x < 40; x++ {
		// for each criteria
		for cr := 0; cr < maxCriteria; cr++ {
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
	var f, r [][]rowCriteria

	f = levelOne

	// Append single criteria
	for i := 0; i < len(f); i++ {
		var v0 []rowCriteria
		var vi rowCriteria
		vi.c = f[i][0].c
		vi.r = f[i][0].r
		v0 = append(v0, vi)
		r = append(r, v0)
	}

	// Append two level criteria
	for i := 0; i < len(f); i++ {
		for j := i + 1; j < len(f); j++ {
			var v0 []rowCriteria
			var vi, vj rowCriteria
			vi.c = f[i][0].c
			vi.r = f[i][0].r
			vj.c = f[j][0].c
			vj.r = f[j][0].r
			v0 = append(v0, vi)
			v0 = append(v0, vj)
			r = append(r, v0)
		}
	}

	return r
}

func randLevels() [][]rowCriteria {
	var f, r [][]rowCriteria
	var flen int

	f = levelOne
	flen = len(f)

	// Append single criteria
	for i := 0; i < len(f); i++ {
		var v0 []rowCriteria
		var vi rowCriteria
		vi.c = f[i][0].c
		vi.r = f[i][0].r
		v0 = append(v0, vi)
		r = append(r, v0)
	}

	for i := 0; i < rand_numSets; i++ {
		var s []rowCriteria
		var sets = rand.Intn(rand_maxSetMembers) + 2
		for j := 0; j < sets; j++ {
			var random = rand.Intn(flen)

			var vi rowCriteria
			vi.c = f[random][0].c
			vi.r = f[random][0].r
			s = append(s, vi)
		}

		r = append(r, s)
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
		//fmt.Printf("number of rows: %d \n", len(t))
		s := evalScore(t, src1, dataSetId)
		r = append(r, s)
	}

	return r
}

func main() {
	t := time.Now()
	fmt.Println(t.Format(time.RFC3339))

	//rand.Seed(time.Now().UTC().UnixNano())

	// Read in data
	readData()

	// Set one level with all row criteria
	levelOne = fullOneLevel()

	var scores []scoreResult

	//var level1 = fullOneLevel()
	//levels = fullTwoLevel()

	//levels = level1
	levels = randLevels()
	//outputRowCriteria(levels)

	fmt.Printf("levels count: %d \n", len(levels))

	for dataSetId := 1; dataSetId <= datasets; dataSetId++ {
		s := levelEval(dataSetId)
		sort.Sort(scoreResults(s))
		if len(s) > 0 {
			// pick the top score
			scores = append(scores, s[0])
			fmt.Printf("%d, %f \n", s[0].dataSetId, s[0].score)
			//outputScore(s[0])
			//outputScore(s[1])
			//outputScore(s[2])
		}
	}

	outputScores(scores)
	outputResults(scores)

	t = time.Now()
	fmt.Println(t.Format(time.RFC3339))

}
