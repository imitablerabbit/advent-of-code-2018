package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	file      = flag.String("file", "input.data", "filepath for the input data")
	duplicate = flag.Bool("duplicate", false, "should we be looking for the first duplicate")
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
	numbers, err := parseLines(lines)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	var output int
	if *duplicate {
		output = findDuplicates(numbers)
	} else {
		output = countFrequency(numbers)
	}
	fmt.Printf("Output: %d\n", output)
}

func parseLines(lines []string) ([]int, error) {
	numbers := make([]int, 0, len(lines))
	for _, line := range lines {
		if line == "" { // file might have ended with a new line
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, n)
	}
	return numbers, nil
}

func countFrequency(deltas []int) int {
	current := 0
	for _, n := range deltas {
		current = current + n
	}
	return current
}

func findDuplicates(deltas []int) int {
	found := make(map[int]bool)
	current := 0
	found[0] = true // include starting 0 in the duplicate list
	for {
		for _, n := range deltas {
			current = current + n
			_, ok := found[current]
			if ok {
				return current
			}
			found[current] = true
		}
	}
}
