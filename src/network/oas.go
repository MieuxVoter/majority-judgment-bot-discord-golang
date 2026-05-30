package network

import (
	"fmt"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"github.com/sarulabs/di/v2"
	"log"
	"log/slog"
	"main/src/container"
	"main/src/database"
	"main/src/security"
	"main/src/services"
	"net/url"
	"regexp"
)

// Oas service helps communicating with the MJ OAS API.
// We mostly use this to generate merit profiles.
// See https://github.com/MieuxVoter/mv-api-server-apiplatform
type Oas struct {
	logger   *slog.Logger
	config   *services.Config
	gradings *services.Gradings
}

// GetMeritProfileUrl builds the URL to the merit profile image
func (oas *Oas) GetMeritProfileUrl(
	poll *database.Poll,
	proposals []database.Proposal,
	pollTally *judgment.PollTally,
	pollResult *judgment.PollResult,
	extension string,
	maxLength int,
) (string, error) {

	fileNameNoExt := ""
	oasDomain := oas.config.Get("OAS_DOMAIN")
	subject := security.TruncateString(poll.Subject, 40)
	query := fmt.Sprintf("?s=%s", url.QueryEscape(subject))

	for proposalResultIndex, proposalResult := range pollResult.ProposalsSorted {
		if proposalResultIndex > 0 {
			fileNameNoExt += "_"
		}
		for gradeLevel := range poll.GetGradingSlice(oas.gradings) {
			gradeAmount := pollTally.Proposals[proposalResult.Index].Tally[gradeLevel]

			if gradeLevel > 0 {
				fileNameNoExt += "-"
			}
			fileNameNoExt += fmt.Sprintf("%d", gradeAmount)
		}
	}

	imageUrlPath := fmt.Sprintf(
		"%s/%s.%s",
		oasDomain, fileNameNoExt, extension,
	)

	amountOfCharactersLeft := maxLength - len(imageUrlPath) - len(query)
	// Conservative approximation through euclidean division:
	amountOfCharactersLeftPerProposal := amountOfCharactersLeft / len(proposals)

	for _, proposalResult := range pollResult.ProposalsSorted {
		proposal := proposals[proposalResult.Index]
		medal := ""
		if proposalResult.Rank == 1 {
			medal = "🥇 "
		} else if proposalResult.Rank == 2 {
			medal = "🥈 "
		} else if proposalResult.Rank == 3 {
			medal = "🥉 "
		}
		queryKey := "p[]"
		queryProposalName := url.QueryEscape(medal + proposal.Name)
		maxProposalNameLength := len(queryProposalName)
		expectedProposalNameLength := len(queryProposalName) + len(queryKey) + 2
		if expectedProposalNameLength > amountOfCharactersLeftPerProposal {
			maxProposalNameLength = amountOfCharactersLeftPerProposal - len(queryKey) - 2
		}
		if maxProposalNameLength > 0 {
			queryProposalName = queryProposalName[:maxProposalNameLength]
			// Since we kind of need to truncate AFTER encoding, we want to remove truncated encoded chars
			trailingEncodedAndTruncated := regexp.MustCompile("[%][a-zA-Z0-9]{0,2}$")
			queryProposalName = trailingEncodedAndTruncated.ReplaceAllString(queryProposalName, "")
			query += fmt.Sprintf("&%s=%s", queryKey, queryProposalName)
		}
	}

	imageUrl := fmt.Sprintf("%s%s", imageUrlPath, query)
	//imageUrl = "https://upload.wikimedia.org/wikipedia/commons/thumb/6/66/SMPTE_Color_Bars.svg/1344px-SMPTE_Color_Bars.svg.png"
	//imageUrl = "https://i.imgur.com/ZGPxFN2.jpg" // why yes Discord dumped support for PNG

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
				logger:   ctn.Get("logger").(*slog.Logger),
				config:   ctn.Get("config").(*services.Config),
				gradings: ctn.Get("gradings").(*services.Gradings),
			}
			return oas, nil
		},
	})
	if err != nil {
		log.Fatalln("service oas failed to build", err)
	}
}
