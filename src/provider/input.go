package provider

// Input holds data coming from userland through the vendor.
// Trying to be generic so as to add other platforms/vendors than Discord at some point.
// This might not work, or might become troublesome, but let's strive for vendor abstraction anyway.
type Input interface {
	GetOptionString(subcommand string, name string, defaultValue string) (string, error)
	GetActorVendorId() (string, error)
	GetActorName() (string, error)
	GetActorLanguage() string
	GetGuildVendorId() (string, error)
	IsDirectMessage() bool
}

// ButtonInput holds data coming from userland through the vendor when a button was pressed.
type ButtonInput interface {
	Input
	GetButtonName() (string, error)
}
