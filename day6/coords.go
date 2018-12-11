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
	file  = flag.String("file", "input.data", "filepath for the input data")
	print = flag.Bool("print", false, "print out the areas")
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
	coords := parseCoords(lines)
	w, h := getArraySizes(coords)
	// w = w + 1 // example shows 1 larger than largest coord
	distances := createManhattenArray(coords, w, h)
	closestDistances := findClosestDistances(distances)
	if *print {
		printAreas(closestDistances, w, h)
	}
	la := largestArea(closestDistances, w, h)
	fmt.Println("Largest Area:", la)
}

type Coord struct {
	X  int
	Y  int
	ID string
}

func (c *Coord) String() string {
	return fmt.Sprintf("%dx%d", c.X, c.Y)
}

type Distance struct {
	Coord    *Coord
	Distance int
}

func (d *Distance) String() string {
	return fmt.Sprintf("%d", d.Distance)
}

func parseCoords(lines []string) []*Coord {
	letterIndex := 0
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	coords := make([]*Coord, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		split := strings.Split(line, ", ")
		x, _ := strconv.Atoi(split[0])
		y, _ := strconv.Atoi(split[1])
		coord := &Coord{
			X:  x,
			Y:  y,
			ID: string(alphabet[letterIndex]),
		}
		coords = append(coords, coord)
		letterIndex++
		if letterIndex >= len(alphabet) {
			letterIndex = 0
		}
	}
	return coords
}

func getArraySizes(coords []*Coord) (int, int) {
	var maxX, maxY int
	for _, coord := range coords {
		if coord.X > maxX {
			maxX = coord.X
		}
		if coord.Y > maxY {
			maxY = coord.Y
		}
	}
	return maxX + 1, maxY + 1
}

func createManhattenArray(coords []*Coord, w, h int) [][]*Distance {
	size := w * h
	array := make([][]*Distance, size)
	for _, coord := range coords {
		for i := 0; i < size; i++ {
			x := i % w
			y := i / w
			d := manhattenDistance(x, y, coord.X, coord.Y)
			distance := &Distance{
				Coord:    coord,
				Distance: d,
			}
			array[i] = append(array[i], distance)
		}
	}
	return array
}

func manhattenDistance(x1, y1, x2, y2 int) int {
	dx := x1 - x2
	dy := y1 - y2
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

func findClosestDistances(distanceArray [][]*Distance) []*Distance {
	closestDistances := make([]*Distance, len(distanceArray))
	for i, distances := range distanceArray {
		closest := findSmallest(distances)
		closestDistances[i] = closest
	}
	return closestDistances
}

func findSmallest(distances []*Distance) *Distance {
	smallest := distances[0]
	multiple := false
	for i := 1; i < len(distances); i++ {
		d := distances[i]
		if d.Distance < smallest.Distance {
			smallest = d
			multiple = false
		} else if d.Distance == smallest.Distance {
			multiple = true
		}
	}
	if multiple {
		return nil
	}
	return smallest
}

func printAreas(distances []*Distance, w, h int) {
	for i, distance := range distances {
		x := i % w
		y := i / w
		if distance == nil { // no outright closest
			fmt.Print(".")
		} else if distance.Coord.X == x && distance.Coord.Y == y {
			fmt.Print(strings.ToUpper(distance.Coord.ID))
		} else {
			fmt.Print(distance.Coord.ID)
		}
		if x == w-1 {
			fmt.Println()
		}
	}
}

func largestArea(distances []*Distance, w, h int) int {
	count := make(map[*Coord]int)
	edges := make(map[*Coord]bool)
	for i, distance := range distances {
		x := i % w
		y := i / w
		if distance == nil {
			continue
		}
		count[distance.Coord]++
		if x >= w-1 || x == 0 || y >= h-1 || y == 0 {
			edges[distance.Coord] = true
		}
	}
	largest := -1
	for k, v := range count {
		if !edges[k] && v > largest {
			largest = v
		}
	}
	return largest
}
