package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mieuxvoter/merit-profile-library-go/merit"
	"image/color"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
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
	fmt.Printf("height=%d\n", int(math.Round(heightFloat)))

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

	err = ConvertSvgToPng(
		"./lib/resvg",
		svgPath,
		pngPath,
	)
	if err != nil {
		panic(err)
	}
}

func ConvertSvgToPng(
	rsvgPath string,
	svgPath string,
	pngPath string,
) error {
	args := []string{
		svgPath,
		pngPath,
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	err := exec.CommandContext(ctx, rsvgPath, args...).Run()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return errors.New("timeout when running rsvg")
		}
		return err
	}

	_, err = os.Stat(pngPath)
	if err != nil {
		return err
	}

	return nil
}
