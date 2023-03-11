## Haw to create a new subcommand

Create a new file with an object implementing `Command`,
and register it to the container with a name prefixed by `command.`.
It will be auto-loaded when the bot starts.

See the `command_create.go` file for an example of how to do that.

> That's it.


### Wait, you're using `init()` ?

We register the service to the container in the `init()`.
This should not cause any trouble or unexpected behavior.

Do strive to not add anything else to the `init()`, though.
