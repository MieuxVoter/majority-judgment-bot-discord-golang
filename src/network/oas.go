package network

import (
	"fmt"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"log"
	"main/src/container"
	"main/src/database"
	"main/src/services"
	"net/url"
)

// Oas service helps communicating with the MJ OAS API.
// We mostly use this to generate merit profiles.
// See https://github.com/MieuxVoter/mv-api-server-apiplatform
type Oas struct {
	logger *logrus.Logger
	config *services.Config
}

// GetMeritProfileUrl builds the URL to the merit profile image
func (oas *Oas) GetMeritProfileUrl(
	poll *database.Poll,
	proposals []database.Proposal,
	pollTally *judgment.PollTally,
	pollResult *judgment.PollResult,
	extension string,
) (string, error) {

	fileNameNoExt := ""
	query := fmt.Sprintf("?subject=%s", url.QueryEscape(poll.Subject))
	for proposalResultIndex, proposalResult := range pollResult.ProposalsSorted {
		proposal := proposals[proposalResult.Index]

		if proposalResultIndex > 0 {
			fileNameNoExt += "_"
		}
		for gradeLevel := range poll.GetGradingSlice() {
			gradeAmount := pollTally.Proposals[proposalResult.Index].Tally[gradeLevel]

			if gradeLevel > 0 {
				fileNameNoExt += "-"
			}
			fileNameNoExt += fmt.Sprintf("%d", gradeAmount)
		}

		medal := ""
		if proposalResult.Rank == 1 {
			medal = "🥇 "
		} else if proposalResult.Rank == 2 {
			medal = "🥈 "
		} else if proposalResult.Rank == 3 {
			medal = "🥉 "
		}
		query += fmt.Sprintf("&proposals[]=%s", url.QueryEscape(medal+proposal.Name))
	}

	oasDomain := oas.config.Get("OAS_DOMAIN")
	imageUrl := fmt.Sprintf(
		"%s/%s.%s%s",
		oasDomain, fileNameNoExt, extension, query,
	)

	return imageUrl, nil
}

// GetOas returns the Oas service
func GetOas() *Oas {
	return container.Get("oas").(*Oas)
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "oas",
		Build: func(ctn di.Container) (interface{}, error) {
			oas := &Oas{
				logger: ctn.Get("logger").(*logrus.Logger),
				config: ctn.Get("config").(*services.Config),
			}
			return oas, nil
		},
	})
	if err != nil {
		log.Fatalln("service oas failed to build", err)
	}
}
