package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

var (
	file         = flag.String("file", "input.data", "filepath for the input data")
	checksumFlag = flag.Bool("checksum", false, "verify that the ids are correct")
)

func init() {
	flag.Parse()
}

func main() {
	data, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	if *checksumFlag {
		c := checksum(lines)
		fmt.Printf("Checksum: %d\n", c)
		return
	}
	id, id2 := findSimilarBoxes(lines, 1)
	letters := findCommonLetters(id, id2)
	fmt.Printf("Common Letters: %s\n", string(letters))
}

func checksum(ids []string) int {
	count2 := 0
	count3 := 0
	for _, id := range ids {
		freq := make(map[rune]int)
		for _, c := range id {
			freq[c]++
		}

		var has2, has3 bool
		for _, count := range freq {
			if count == 2 {
				has2 = true
			}
			if count == 3 {
				has3 = true
			}
			if has2 && has3 {
				break
			}
		}

		if has2 {
			count2++
		}
		if has3 {
			count3++
		}
	}
	return count2 * count3
}

func findSimilarBoxes(ids []string, diff int) (string, string) {
	halfway := math.Ceil(float64(len(ids)) / 2)
	for i := 0; i < int(halfway); i++ {
		id := ids[i]
		if id == "" {
			continue
		}
		for j, id2 := range ids {
			if i == j {
				continue
			}
			if id2 == "" {
				continue
			}
			n := letterDifferences(id, id2)

			// Check if the number of characters that are the same is
			// less than or equal to diff
			if n <= diff {
				return id, id2
			}
		}
	}
	return "", ""
}

func letterDifferences(id, id2 string) int {
	count := 0
	for i, c := range id {
		if c != rune(id2[i]) {
			count++
		}
	}
	return count
}

func findCommonLetters(id, id2 string) []rune {
	common := make([]rune, 0, len(id))
	for i, r := range id {
		if r == rune(id2[i]) {
			common = append(common, r)
		}
	}
	return common
}
