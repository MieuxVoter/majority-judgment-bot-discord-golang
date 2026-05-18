package services

import (
	"errors"
	"fmt"
	"github.com/canhlinh/svg2png"
	"github.com/mieuxvoter/merit-profile-library-go/merit"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"image/color"
	"log"
	"main/src/container"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"
)

var svgDimensionsRegex = regexp.MustCompile(`<svg width="(?P<width>[^"]+)" height="(?P<height>[^"]+)"`)

type Analysis struct {
	logger *logrus.Logger
}

func (service *Analysis) GenerateMeritProfileSVG(
	proposals []merit.Proposal,
) (svg string, err error) {

	svg, err = merit.RenderLinearProfileSVG(
		proposals,
		merit.WithWidth(600),
		merit.WithBgColor(color.NRGBA{R: 32, G: 32, B: 32, A: 255}),
		merit.WithFontFamily("Noto Sans, sans-serif"),
		merit.WithTallyFontSize("1em"),
	)
	if err != nil {
		return
	}

	return
}

func (service *Analysis) GenerateMeritProfilePNG(
	proposals []merit.Proposal,
) (png []byte, err error) {

	svg, err := service.GenerateMeritProfileSVG(proposals)
	if err != nil {
		return
	}

	matches := svgDimensionsRegex.FindStringSubmatch(svg)
	if matches == nil {
		err = errors.New("cannot find dimensions in svg")
		return
	}

	widthString := matches[1]
	heightString := matches[2]
	widthFloat, err := strconv.ParseFloat(widthString, 64)
	if err != nil {
		return
	}
	heightFloat, err := strconv.ParseFloat(heightString, 64)
	if err != nil {
		return
	}

	now := time.Now().UnixNano()
	tmpDir := "/tmp"

	svgPath := fmt.Sprintf("%s/profile%d.svg", tmpDir, now)
	pngPath := fmt.Sprintf("%s/profile%d.png", tmpDir, now)

	err = os.WriteFile(
		svgPath,
		[]byte(svg),
		0664,
	)
	if err != nil {
		return
	}

	chromiumConverter := svg2png.NewChrome().
		SetWith(int(math.Round(widthFloat))).
		SetHeight(int(math.Round(heightFloat)))
	err = chromiumConverter.Screenshoot(
		svgPath,
		pngPath,
	)

	png, err = os.ReadFile(pngPath)

	_ = os.Remove(svgPath)
	_ = os.Remove(pngPath)

	return
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "analysis",
		Build: func(ctn di.Container) (interface{}, error) {
			service := &Analysis{
				logger: ctn.Get("logger").(*logrus.Logger),
			}
			return service, nil
		},
	})
	if err != nil {
		log.Fatalln("analysis failed to build", err)
	}
}
