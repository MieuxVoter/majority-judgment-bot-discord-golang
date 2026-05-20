# Lib

These are vendors we use.

It's not great to have those here ; if you manage to remove them, please make a MR.

> Keep in mind that we need a setup that works both natively and in Docker.

## resvg

> Converts SVG to PNG, quickly and properly.

From https://github.com/linebender/resvg/releases

Current vendored version: `0.47.0`, _(linux x86_64)_:
https://github.com/linebender/resvg/releases/download/v0.47.0/resvg-linux-x86_64.tar.gz

We also tried:
- _inkscape_: does not support the filters we use, the end result is terrible
- _chromium_: works, but the docker image is more than 1Gio
