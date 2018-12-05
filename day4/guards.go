package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	file  = flag.String("file", "input.data", "filepath for the input data")
	print = flag.Bool("print", false, "display the sleep times")
)

func init() {
	flag.Parse()
}

type Guard struct {
	ID         string
	SleepTimes []Sleep
}

type Sleep struct {
	Start time.Time
	End   time.Time
}

func main() {
	data, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	sort.Slice(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})

	items := parseLines(lines)
	sm := sleepMap(items)
	if *print {
		printSleepMap(sm)
	}
	gsts := guardSleepTimes(sm)

	// Strategy 1
	gst := sleepiestGuard(gsts)
	fmt.Println("Sleepiest Guard:", gst.ID)
	fmt.Println("Sleep time:", gst.TotalSleepTime)
	fmt.Println("Most common minute:", gst.HighestFrequencyIndex)
	idInt, _ := strconv.Atoi(gst.ID[1:])
	fmt.Println("Strategy 1 Output: ", idInt*gst.HighestFrequencyIndex)
	fmt.Println()

	// Strategy 2
	gst2 := mostCommonSleepGuard(gsts)
	fmt.Println("Guard with most common minutes slept:", gst2.ID)
	fmt.Println("Number of times slept during same minute:", gst2.HighestFrequency)
	fmt.Println("Most common minute:", gst.HighestFrequencyIndex)
	id2Int, _ := strconv.Atoi(gst2.ID[1:])
	fmt.Println("Strategy 2 Output: ", id2Int*gst2.HighestFrequencyIndex)
}

func parseLines(lines []string) []Item {
	items := make([]Item, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		item := parseLine(line)
		items = append(items, item)
	}
	return items
}

func sleepMap(items []Item) map[string][]Sleep {
	sleepMap := make(map[string][]Sleep)
	currentID := ""
	for i := 0; i < len(items); i++ {
		item := items[i]
		switch item.ItemType {
		case TypeNewShift:
			currentID = idFromNewShift(item)
		case TypeSleep:
			i++
			wakeUp := items[i]
			s := Sleep{
				Start: item.Timestamp,
				End:   wakeUp.Timestamp,
			}
			ss := sleepMap[currentID]
			ss = append(ss, s)
			sleepMap[currentID] = ss
		default:
			// Should not hit this as wake up is peeled after sleep
		}
	}
	return sleepMap
}

func idFromNewShift(item Item) string {
	if item.ItemType != TypeNewShift {
		return ""
	}
	split := strings.Split(item.Text, " ")
	id := split[1]
	return id
}

type Item struct {
	Timestamp time.Time
	Text      string
	ItemType  Type
}

const (
	TypeNewShift = iota
	TypeSleep
	TypeWakeUp
)

type Type int

func parseLine(line string) Item {
	split := strings.SplitN(line, "]", 2)
	timeString := split[0]
	time, err := time.Parse("2006-01-02 15:04", timeString[1:len(timeString)])
	if err != nil {
		fmt.Println(err)
	}
	var t Type
	text := split[1][1:]
	switch text {
	case "falls asleep":
		t = TypeSleep
	case "wakes up":
		t = TypeWakeUp
	default:
		t = TypeNewShift
	}
	return Item{
		Timestamp: time,
		Text:      text,
		ItemType:  t,
	}
}

func printSleepMap(sm map[string][]Sleep) {
	fmt.Printf("Date\tID\tMinute\n")
	fmt.Printf("\t\t")
	for t := 0; t < 6; t++ {
		fmt.Printf("%d%d%d%d%d%d%d%d%d%d", t, t, t, t, t, t, t, t, t, t)
	}
	fmt.Println()
	fmt.Printf("\t\t")
	for t := 0; t < 6; t++ {
		fmt.Printf("0123456789")
	}
	fmt.Println()
	for k, v := range sm {
		for _, s := range v {
			fmt.Printf("%s\t%s\t", s.Start.Format("01-02"), k)
			for i := 0; i < 60; i++ {
				if i >= s.Start.Minute() && i < s.End.Minute() {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Println()
		}
	}
}

type GuardSleepTime struct {
	ID                    string
	Mins                  []int
	TotalSleepTime        int
	HighestFrequencyIndex int
	HighestFrequency      int
}

func guardSleepTimes(sleepMap map[string][]Sleep) []GuardSleepTime {
	gsts := make([]GuardSleepTime, len(sleepMap))
	i := 0
	for id, sleepTimes := range sleepMap {
		mins := minutesSlept(sleepTimes)
		t := totalSleep(mins)
		fIndex := highestFrequency(mins)
		gst := GuardSleepTime{
			Mins:                  mins,
			ID:                    id,
			TotalSleepTime:        t,
			HighestFrequencyIndex: fIndex,
			HighestFrequency:      mins[fIndex],
		}
		gsts[i] = gst
		i++
	}
	return gsts
}

func sleepiestGuard(gsts []GuardSleepTime) GuardSleepTime {
	longestSlept := GuardSleepTime{}
	for _, gst := range gsts {
		if gst.TotalSleepTime > longestSlept.TotalSleepTime {
			longestSlept = gst
		}
	}
	return longestSlept
}

func mostCommonSleepGuard(gsts []GuardSleepTime) GuardSleepTime {
	g := GuardSleepTime{}
	for _, gst := range gsts {
		if gst.HighestFrequency > g.HighestFrequency {
			g = gst
		}
	}
	return g
}

func minutesSlept(sleepTimes []Sleep) []int {
	mins := make([]int, 60)
	for _, s := range sleepTimes {
		for c := s.Start; c.Before(s.End); c = c.Add(time.Minute) {
			m := c.Minute()
			mins[m]++
		}
	}
	return mins
}

func totalSleep(mins []int) int {
	total := 0
	for _, min := range mins {
		total = total + min
	}
	return total
}

func highestFrequency(mins []int) int {
	highest := -1
	highestIndex := -1
	for i, min := range mins {
		if min > highest {
			highest = min
			highestIndex = i
		}
	}
	return highestIndex
}
