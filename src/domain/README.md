## How to create a new /mj subcommand

Create a new file with an object implementing `Subcommand`,
and register it to the container with a name prefixed by `command.mj.`.
It will be auto-loaded when the bot starts.

See the `subcommand_mj_create.go` file in `../commands` for an example of how to do that.

> That's it.


## Buttons

Same goes for buttons.


## Wait, you're using `init()` ?

We register the service to the container in the `init()`.
This should not cause any trouble or unexpected behavior.

> Do strive to not add anything else to the `init()`, though.

