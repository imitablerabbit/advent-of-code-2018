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
	print = flag.Bool("print", false, "display the claimant map")

	width  = flag.Int("width", 1000, "width of the fabric")
	height = flag.Int("height", 1000, "height of the fabric")
)

func init() {
	flag.Parse()
}

type Claim struct {
	ID     string
	X      int
	Y      int
	Width  int
	Height int
}

func main() {
	data, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	claims := parseClaims(lines)
	claimMap := mapClaims(claims, *width, *height)

	// print out the map
	if *print {
		printMap(claimMap)
	}

	n := calculateOverlaps(claimMap)
	goodClaims := filterOverlapping(claims, claimMap)

	fmt.Printf("Overlaps: %d\n", n)
	fmt.Println("Non Overlapping Claims:")
	for _, c := range goodClaims {
		fmt.Println(c.ID)
	}
}

func parseClaims(lines []string) []*Claim {
	claims := make([]*Claim, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		split := strings.Split(line, " ")
		id := split[0]
		pos := split[2]
		dim := split[3]
		splitPos := strings.Split(pos[:len(pos)-1], ",")
		x, _ := strconv.Atoi(splitPos[0])
		y, _ := strconv.Atoi(splitPos[1])
		splitDim := strings.Split(dim, "x")
		w, _ := strconv.Atoi(splitDim[0])
		h, _ := strconv.Atoi(splitDim[1])
		c := &Claim{
			ID:     id,
			X:      x,
			Y:      y,
			Width:  w,
			Height: h,
		}
		claims = append(claims, c)
	}
	return claims
}

func mapClaims(claims []*Claim, width, height int) [][]*Claim {
	fabric := make([][]*Claim, width*height)
	for _, c := range claims {
		start := c.X + (c.Y * width)

		// Loop through the claimaint area
		for x := 0; x < c.Width; x++ {
			for y := 0; y < c.Height; y++ {
				i := start + x + (y * width)
				sqi := fabric[i]
				sqi = append(sqi, c)
				fabric[i] = sqi
			}
		}
	}
	return fabric
}

func printMap(m [][]*Claim) {
	for w := 0; w < *width; w++ {
		for h := 0; h < *height; h++ {
			i := w + (h * *width)
			c := m[i]
			s := "."
			if len(c) == 1 {
				s = c[0].ID[1:]
			} else if len(c) > 1 {
				s = "x"
			}
			fmt.Print(s)
		}
		fmt.Println()
	}
}

func calculateOverlaps(claimMap [][]*Claim) int {
	count := 0
	for _, claims := range claimMap {
		if len(claims) > 1 {
			count++
		}
	}
	return count
}

func filterOverlapping(claims []*Claim, claimMap [][]*Claim) []*Claim {
	// Convert to map for easy deletion
	remaining := make(map[string]*Claim)
	for _, c := range claims {
		remaining[c.ID] = c
	}

	// Delete claims from map if they have overlapped
	for _, cs := range claimMap {
		if len(cs) <= 1 {
			continue
		}
		for _, c := range cs {
			delete(remaining, c.ID)
		}
	}

	// Convert back to slice as it makes more sense
	r := make([]*Claim, len(remaining))
	i := 0
	for _, v := range remaining {
		r[i] = v
		i++
	}
	return r
}
