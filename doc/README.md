## Architectural Decision Records

### Discord Client

This leverages the excellent `disgord` https://github.com/andersfylling/disgord

#### Options

- https://github.com/bwmarrin/discordgo
  More stars.
- https://github.com/andersfylling/disgord
  Cache, higher level.

---

## Architectural Decision Drafts

### Merit Profile Generation using a web service

We use MieuxVoter's OpenApi implementation.

It was faster.  Ideally, perhaps, the bot handles its own SVG and PNG generation.
Careful, though, since image handling would likely introduce significant dependencies, so we should ensure that they stay optional and fallback to the OAS (if available itself).

### Gateways & Webhooks

Right now we're using re-connecting gateways to connect with Discord.

Ideally, we'd use a gateway once to set up some webhooks, close the gateway, listen to webhooks during runtime, and remove the webhooks before exiting, whenever we exit.
This would require a domain configuration and webserver goodies.

### Providers (Discord, Telegram)

We're trying to factor the code to support multiple providers, through adapters.
Might just clone the bot in the end.  Undecided.
We could also perhaps have the same bot handle two databases, since we have _di_.

### Integration Testing

- [ ] Discord mock
- [ ] Telegram mock
- [ ] Golang gherkin runner
