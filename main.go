// innopsi project main.go
// using https://godoc.org/github.com/montanaflynn/stats
package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
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
	rejected  bool
}

type confInterval struct {
	min       float64
	max       float64
	diff      float64
	middle    float64
	closeness float64
	t1Min     float64
	t1Max     float64
	diffSd    float64
}

type confInterval2 struct {
	t0min, t0max, t1min, t1max float64
	overlap                    bool
}

const subjects int = 240
const minCriteria = 6
const maxCriteria = 6

const datasets int = 1200

//const datasets int = 4

const datafilename string = "./data/InnoCentive_9933623_Data.csv"

//const datafilename string = "./data/InnoCentive_9933623_Training_Data.csv"

var (
	data               []coreData
	levels             [][]rowCriteria
	levelOne           [][]rowCriteria
	rand_numSets           = 1000
	rand_maxSetMembers int = 9
	maxExperiments     int = 1
	filename           string
	rowThreshhold      int
	scoreCutoff        float64
	zScore             float64 = 2.58
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
				if maxCriteria == 6 {
					if t < 33.3 {
						data[count].xi[i] = 0
					} else if t >= 33.3 && t < 66.6 {
						data[count].xi[i] = 1
					} else {
						data[count].xi[i] = 2
					}
				}

				if maxCriteria == 30 {
					if t < 20.0 {
						data[count].xi[i] = 0
					} else if t >= 20.0 && t < 40.0 {
						data[count].xi[i] = 1
					} else if t >= 40.0 && t < 60.0 {
						data[count].xi[i] = 2
					} else if t >= 60.0 && t < 80.0 {
						data[count].xi[i] = 3
					} else {
						data[count].xi[i] = 4
					}
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

	fmt.Printf("rejected: %t \n", rejected(s))
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

// initialize criteria arrays
func criteriaH(v int, c int) bool {
	switch c {
	case 0:
		return v == 0
	case 1:
		return v == 1
	case 2:
		return v == 2
	case 3:
		return v == 3
	case 4:
		return v == 4
	case 5:
		return v == 0 || v == 1
	case 6:
		return v == 0 || v == 2
	case 7:
		return v == 0 || v == 3
	case 8:
		return v == 0 || v == 4
	case 9:
		return v == 1 || v == 2
	case 10:
		return v == 1 || v == 3
	case 11:
		return v == 2 || v == 3
	case 12:
		return v == 2 || v == 4
	case 13:
		return v == 3 || v == 4
	case 14:
		return v == 1 || v == 4
	case 15:
		return v == 0 || v == 1 || v == 2
	case 16:
		return v == 0 || v == 1 || v == 3
	case 17:
		return v == 0 || v == 1 || v == 4
	case 18:
		return v == 0 || v == 2 || v == 3
	case 19:
		return v == 0 || v == 2 || v == 4
	case 20:
		return v == 0 || v == 3 || v == 4
	case 21:
		return v == 1 || v == 2 || v == 3
	case 22:
		return v == 1 || v == 2 || v == 4
	case 23:
		return v == 1 || v == 3 || v == 4
	case 24:
		return v == 2 || v == 3 || v == 4
	case 25:
		return v == 0 || v == 1 || v == 2 || v == 3
	case 26:
		return v == 0 || v == 1 || v == 2 || v == 4
	case 27:
		return v == 0 || v == 1 || v == 3 || v == 4
	case 28:
		return v == 0 || v == 2 || v == 3 || v == 4
	case 29:
		return v == 1 || v == 2 || v == 3 || v == 4
	}

	return false
}

// Evaluation Criteria
func criteriaL(v int, c int) bool {
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

func criteria(x, v, c int) bool {

	//	if x < 20 {
	//		return criteriaL(v, c)
	//	}

	//	return criteriaH(v, c)
	return criteriaL(v, c)
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
			// k.r is the x index
			if !(criteria(k.r, each.xi[k.r], k.c)) {
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

			pv = resultsArray(s[col].t1, s[col])
			// with the scores for individuals
			for sub := 0; sub < subjects; sub++ {

				dataSet[sub][col+1] = pv[sub]

				// If output is less then expMin
				if rejected(s[col]) {
					dataSet[sub][col+1] = 0
				}
			}
		}
	}

	// Output to file
	t := time.Now()
	filename = fmt.Sprintf("./results/o_%d%02d%02dT%02d%02d%02d.csv",
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
func resultsArray(positive []coreData, s scoreResult) []int {
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

func rejected(s scoreResult) bool {
	return s.score > scoreCutoff
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

	// then calculate the median, also experiment with average
	var mean0, _ = stats.Mean(t0s)
	var mean1, _ = stats.Mean(t1s)
	//var meanAll, _ = stats.Mean(allTs)
	//var sd, _ = stats.StandardDeviationPopulation(allTs)

	// subtract the two t0-t1, we want t1 to be smaller
	// Note: use spooled
	// square root of ((Nt-1)St^2 + (Nc-1)Sc^2)/(Nt+Nc))
	var St, _ = stats.StandardDeviation(t1s)
	var Sc, _ = stats.StandardDeviation(t0s)
	var Nt = float64(len(t1s))
	var Nc = float64(len(t0s))

	//func calculateConfidenceInterval2(nt, nc, mt, mc, sdt, sdc float64) confInterval2
	var ci = calculateConfidenceInterval2(Nt, Nc, mean1, mean0, St, Sc)

	// If the confidence intervals overlap then not valid range
	if ci.overlap {
		return s
	}

	var St2 = St * St
	var Sc2 = Sc * Sc
	//var Ntm1 = float64(Nt - 1)
	//var Ncm1 = float64(Nc - 1)
	//var kt = Ntm1 * St2
	//var kc = Ncm1 * Sc2
	//var ksum = kt + kc
	//var Nsum = Nt + Nc
	//var sPooled = math.Sqrt(ksum / Nsum)

	//http://www.uccs.edu/~lbecker/
	var sPooled = math.Sqrt((St2 + Sc2) / 2)

	s.t0 = t0
	s.t1 = t1
	//sPooled = math.Sqrt((St2 * Sc2) / 2)

	//var _, t1confh = NormalConfidenceInterval(t1s)
	//var _, t0confh = NormalConfidenceInterval(t0s)

	//var meanValue = mean1 - mean0
	//var meanValue = (mean1/St - mean0/Sc) / sPooled

	// Score Type 1
	var meanDifference = mean1 - mean0
	//s.score = meanDifference / meanAll
	//var max, _ = stats.Max(allTs)

	//s.score = meanDifference / sPooled
	var cohensd = meanDifference / sPooled
	var a = ((Nt + Nc) * (Nt + Nc)) / (Nt + Nc)
	s.score = cohensd / math.Sqrt((cohensd*cohensd)+a)

	//	if math.Abs(s.score) >= 1.0 {
	//		s.score = 0
	//	}

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

	for x := 0; x < 20; x++ {
		// for each criteria
		for cr := 0; cr < minCriteria; cr++ {
			var v []rowCriteria
			var k rowCriteria
			k.r = x
			k.c = cr
			v = append(v, k)
			r = append(r, v)
		}
	}

	for x := 20; x < 40; x++ {
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

func outputScoreList(s []scoreResult) {
	for _, each := range s {
		fmt.Printf("%d, %f \n", each.dataSetId, each.score)
	}
}

func evaluateScores(s []scoreResult) scoreResult {
	// look for significance, which includes a high score, and
	// a larger number of members in the set

	var topRange = 5
	//	var scoreWeighted = 0.0
	var scored scoreResult
	var outputResults = false
	var ci float64 = 100.0
	for i := 0; i < topRange; i++ {
		var scoreWeightedN = s[i].score * float64(len(s[i].t0)+len(s[i].t1))
		var ciN = calculateConfidenceInterval(s[i])
		if ciN.diff < ci {
			ci = ciN.diff
			scored = s[i]
		}

		//		if scoreWeightedN < scoreWeighted {
		//			// Find the best weighted score and return that
		//			scoreWeighted = scoreWeightedN
		//			scored = s[i]
		//		}

		if outputResults {
			fmt.Printf("score: %f, weighted score: %f, t=%d, t1=%d, t0=%d, ",
				s[i].score,
				scoreWeightedN,
				len(s[i].t0)+len(s[i].t1),
				len(s[i].t1),
				len(s[i].t0))

			for _, each := range s[i].rc {
				// r+1 to match naming offsets vs slice zero base
				fmt.Printf("x=%d, c=%d, ", each.r+1, each.c)
			}

			fmt.Println()
		}

	}

	if outputResults {
		fmt.Printf("winning: %f\n", scored.score)
	}

	return scored
}

func compareTrainingDataWithResults() {
	// Open training file
	// Open result file
	trainingFile, _ := ioutil.ReadFile(filename)
	testingFile, _ := ioutil.ReadFile("./data/InnoCentive_9933623_Training_Data_truth_subjects.csv")
	fmt.Printf("lengths: %d %d \n", len(trainingFile), len(testingFile))
	// Files are 3066 bytes, (hack)
	var differences = 0
	for i := 0; i < len(trainingFile); i++ {
		// One by one byte comparison
		if trainingFile[i] != testingFile[i] {
			differences += 1
		}
	}

	// output count of differences
	fmt.Printf("differences: %d \n", differences)
}

func calculateConfidenceInterval2(nt, nc, mt, mc, sdt, sdc float64) confInterval2 {

	var ci confInterval2
	//var z = 1.96 // http://www.dummies.com/how-to/content/creating-a-confidence-interval-for-the-difference-.html
	//var z = 2.58 // http://www.dummies.com/how-to/content/creating-a-confidence-interval-for-the-difference-.html

	ci.t1min = mt - zScore*(sdt/math.Sqrt(nt))
	ci.t0min = mc - zScore*(sdt/math.Sqrt(nc))

	ci.t1max = mt + zScore*(sdt/math.Sqrt(nt))
	ci.t0max = mt + zScore*(sdt/math.Sqrt(nc))

	ci.overlap = false

	if ci.t1max <= ci.t0max && ci.t1max >= ci.t0min {
		ci.overlap = true
	}

	if ci.t1min <= ci.t0max && ci.t1min >= ci.t0min {
		ci.overlap = true
	}

	// check for encirclement
	if ci.t0min <= ci.t1max && ci.t0min >= ci.t1min {
		ci.overlap = true
	}

	return ci
}

func calculateConfidenceInterval(s scoreResult) confInterval {
	var t0s []float64
	var t1s []float64

	// Partition the data into treatment 0 and treatment 1
	// and save the score for evaluation
	for _, each := range s.t0 {
		t0s = append(t0s, each.y)
	}

	for _, each := range s.t1 {
		t1s = append(t1s, each.y)
	}

	var ci confInterval
	//var z = 1.96 // http://www.dummies.com/how-to/content/creating-a-confidence-interval-for-the-difference-.html
	//var z = 1.645 // http://www.dummies.com/how-to/content/creating-a-confidence-interval-for-the-difference-.html
	//var z = 2.58

	var m0, _ = stats.Mean(t0s)
	var n0 = float64(len(t0s))
	var sd0, _ = stats.StandardDeviation(t0s)

	var m1, _ = stats.Mean(t1s)
	var n1 = float64(len(t1s))
	var sd1, _ = stats.StandardDeviation(t1s)

	var mDiff = m0 - m1
	var sd0s = sd0 * sd0
	var sd1s = sd1 * sd1

	ci.min = mDiff - zScore*math.Sqrt(sd1s/n1+sd0s/n0)
	ci.max = mDiff + zScore*math.Sqrt(sd1s/n1+sd0s/n0)
	ci.diff = ci.min - ci.max

	ci.t1Max = m1 + ci.max
	ci.t1Min = m1 + ci.min
	ci.diffSd = ci.diff / sd1

	// how close is the score to the middle of the confidence interval
	ci.middle = (ci.min + ci.max) / 2
	//ci.closeness = math.Abs(s.score - ci.middle)
	//ci.closeness = math.Abs(ci.diffSd - s.score)

	// Difference in sample means +- confidence interval

	//fmt.Printf("conf interval: %f to %f,  conf diff: %f, t1: %f, t1max: %f, t1min: %f, diffSd: %f\n", ci.min, ci.max, ci.diff, m1, ci.t1Max, ci.t1Min, ci.diffSd)

	return ci
}

func main() {
	t := time.Now()
	fmt.Println(t.Format(time.RFC3339))

	rand.Seed(1)

	// Read in data
	readData()

	// Set one level with all row criteria,
	// this is used to start the set creation
	levelOne = fullOneLevel()

	//levels = fullTwoLevel()
	outputRowCriteria(levels)

	// experiment variables
	rand_numSets = 100000
	rand_maxSetMembers = 12
	maxExperiments = 1

	var expMin []float64
	var expMax []float64
	scoreCutoff = -0.89
	rowThreshhold = 2
	zScore = 2.58

	for experiment := 1; experiment <= maxExperiments; experiment++ {
		// experiment variables, changes per experiment
		rand_numSets += 0
		rand_maxSetMembers += 0
		scoreCutoff += -0.00

		// Setup experiment variables
		var scores []scoreResult
		var minScore float64 = -100
		var maxScore float64 = 0
		levels = randLevels()
		fmt.Printf("sets count: %d, max set members: %d, level 1 count: %d, rowThreshhold: %d, scoreCutoff: %f\n", len(levels), rand_maxSetMembers+2, len(levelOne), rowThreshhold, scoreCutoff)

		for dataSetId := 1; dataSetId <= datasets; dataSetId++ {
			s := levelEval(dataSetId)
			sort.Sort(scoreResults(s))

			// s contains a list of scores for one dataset, sorted
			// this is were we can get some info on that data
			//outputScoreList(s)

			if len(s) > 0 {
				//var sEval = evaluateScores(s)
				// pick the top score
				var sEval = s[0]
				scores = append(scores, sEval)
				fmt.Printf("%d, %f \n", sEval.dataSetId, sEval.score)

				if minScore < sEval.score {
					minScore = sEval.score
				}
				if maxScore > sEval.score {
					maxScore = sEval.score
				}
			}
		}

		expMin = append(expMin, minScore)
		expMax = append(expMax, maxScore)

		//scoreCutoff = (minScore * (percentRofMin / 100.0)) + minScore
		//fmt.Printf(" scoreCutoff: %f \n", scoreCutoff)

		outputScores(scores)
		// Write output file
		outputResults(scores)

		// Compare to training truth data
		// compareTrainingDataWithResults()
	}

	t = time.Now()
	fmt.Println(t.Format(time.RFC3339))

	// Output min max scores per experiment
	for _, each := range expMin {
		fmt.Printf("min: %f, ", each)
	}
	fmt.Println()
	for _, each := range expMax {
		fmt.Printf("max: %f, ", each)
	}
}
