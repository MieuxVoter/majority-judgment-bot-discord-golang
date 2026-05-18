
The `var/` directory need to be writable by non-root since our docker image uses the `nonroot` user with id `65532`.

It will hold the sqlite database.

> This is the only way I managed to make it work with a hardened image without being root inside the container.
> I guess it's kind of okay.

Make sure you backup this directory !
