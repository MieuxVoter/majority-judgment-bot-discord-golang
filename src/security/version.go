package security

// Doc: https://github.com/ahmetb/govvv#build-variables

var (
	// GitSummary is populated at compile time by govvv.
	GitSummary string
)

func GetVersion() string {
	if GitSummary == "" {
		return "N/A"
	}
	return GitSummary
}
