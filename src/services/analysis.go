package services

import (
	"context"
	"github.com/mieuxvoter/merit-profile-library-go/merit"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"image/color"
	"log"
	"main/src/container"
)

//var svgDimensionsRegex = regexp.MustCompile(`<svg width="(?P<width>[^"]+)" height="(?P<height>[^"]+)"`)

type Analysis struct {
	logger     *logrus.Logger
	rasterizer *Rasterizer
}

func (service *Analysis) GenerateMeritProfileSVG(
	proposals []merit.Proposal,
	gradesOutlines [][]uint8,
) (svg string, err error) {

	// Type conversion [][]uint8 → [][]int
	intGradesOutlines := make([][]int, len(gradesOutlines))
	for i, uint8Grades := range gradesOutlines {
		intGradesOutlines[i] = make([]int, len(uint8Grades))
		for j, uintGrade := range uint8Grades {
			intGradesOutlines[i][j] = int(uintGrade)
		}
	}

	svg, err = merit.RenderLinearProfileSVG(
		proposals,
		merit.WithWidth(600),
		merit.WithGradeHeight(48),
		merit.WithBgColor(color.NRGBA{R: 32, G: 32, B: 32, A: 255}),
		merit.WithFontFamily("Noto Sans, sans-serif"),
		merit.WithProposalFontSize("28"),
		merit.WithTallyFontSize("20"),
		merit.WithBestGradeOnLeft(true),
		merit.WithGradesOutlines(intGradesOutlines),
		merit.WithGradesOutlinesWidth(3.0),
		merit.WithGradesOutlinesColor(color.White),
	)
	if err != nil {
		return
	}

	return
}

func (service *Analysis) GenerateMeritProfilePNG(
	ctx context.Context,
	proposals []merit.Proposal,
	gradesOutlines [][]uint8,
) (png []byte, err error) {

	svg, err := service.GenerateMeritProfileSVG(proposals, gradesOutlines)
	if err != nil {
		return
	}

	return service.rasterizer.ConvertSVGToPNG(ctx, []byte(svg))
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "analysis",
		Build: func(ctn di.Container) (interface{}, error) {
			service := &Analysis{
				logger:     ctn.Get("logger").(*logrus.Logger),
				rasterizer: ctn.Get("rasterizer").(*Rasterizer),
			}
			return service, nil
		},
	})
	if err != nil {
		log.Fatalln("analysis failed to build", err)
	}
}
