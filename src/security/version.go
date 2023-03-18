package security

// Doc: https://github.com/ahmetb/govvv#build-variables

var (
	// GitSummary is populated at compile time by govvv.
	GitSummary string
)

// GetVersion returns the fully described git version of this bot.
func GetVersion() string {
	if GitSummary == "" {
		return "N/A"
	}
	return GitSummary
}
