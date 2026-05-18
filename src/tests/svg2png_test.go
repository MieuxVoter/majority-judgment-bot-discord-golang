package main

import (
	"github.com/canhlinh/svg2png"
	"github.com/mieuxvoter/merit-profile-library-go/merit"
	"image/color"
	"math"
	"os"
	"regexp"
	"strconv"
)

func main() {
	proposals := []merit.Proposal{
		{
			Name:  "Pizza 4 Dimensions",
			Tally: []uint64{5, 4, 11}, // 3 grades, 20 judgments
		},
		{
			Name:  "Lasagnes Assange",
			Tally: []uint64{9, 5, 6}, // same
		},
		{
			Name:  "Jurassique Pâtes",
			Tally: []uint64{14, 0, 6}, // same
		},
	}

	svg, err := merit.RenderLinearProfileSVG(
		proposals,
		merit.WithBgColor(color.Black),
		merit.WithWidth(800),
	)
	if err != nil {
		panic(err)
	}

	heightRegex := regexp.MustCompile(`<svg width="(?P<width>[^"]+)" height="(?P<height>[^"]+)"`)

	matches := heightRegex.FindStringSubmatch(svg)
	if matches == nil {
		panic("cannot find height in svg")
	}

	heightString := matches[2]
	heightFloat, err := strconv.ParseFloat(heightString, 64)
	if err != nil {
		panic(err)
	}

	//fmt.Print(svg)
	svgPath := "/tmp/test.svg"
	pngPath := "/tmp/test.png"

	err = os.WriteFile(
		svgPath,
		[]byte(svg),
		0666,
	)
	if err != nil {
		panic(err)
	}

	chromiumConverter := svg2png.NewChrome()
	chromiumConverter.SetWith(800).SetHeight(int(math.Round(heightFloat)))
	err = chromiumConverter.Screenshoot(
		svgPath,
		pngPath,
	)

}
