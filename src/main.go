package main

// A Bot for Discord to create polls, using majority judgment.
// Usage:   /mj create <subject> <proposalA> <proposalB> …

import (
	"fmt"
	"main/src/container"
	"main/src/provider/discord"
	"main/src/security"
)

func main() {
	fmt.Printf("=== ⚖  MAJORITY JUDGMENT BOT 🤖 v%s ===\n", security.GetVersion())

	discord.Run()
}

func init() {
	// Each service registers into the container in their own init.
	// init() of main is always last, so let's build the container.
	container.Build()
}
