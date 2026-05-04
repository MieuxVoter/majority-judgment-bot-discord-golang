# Abandon of the version 0.x.x of the bot

## Rationale

[Disgord], the foundational library this bot was based on in v0, is abandoned.

[Disgord]: https://github.com/andersfylling/disgord

## Options

### I.a. Move to another Discord wrapper library

#### Pros

- We're all in this together — effort to comply with the Discord API is shared

#### Cons

- We might depend on a library that gets abandoned as well

### I.b. Directly use the Discord REST API with hand-crafted HTTP requests

#### Pros

- No dependency on anything but Discord

#### Cons

- Discord does break their own API from time to time
- So much work !
- We'd end up writing our own Discord API wrapper (wasteful)

### II.a. Keep using the Gateway API

#### Pros

- Maybe faster ?
- Maybe a little bit simpler ?
- Gateway API is easier during dev, as it does not require a webserver and a domain name

#### Cons

- Expensive in Bandwidth
- Carbon footprint is bigger than it could be


### II.b. Move to REST API and Webhooks

#### Pros

- Less expensive, energy/bandwidth/carbon wise

#### Cons

- More work (radically different from what we've been doing)
- Not 100% sure we can do all we want to do

## Decisions

### Choosing I.a.

Because I.b. without a good reason is just dumb and high-maintenance.

### Choosing II.b.

Because carbon emission reduction is paramount in 2026.

Even if it's just a little ; parable of the colibri.
