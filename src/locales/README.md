
The `BCP 47` language tag in the filename of the `toml` language files must be recognized by Discord.

That is it must be one of the keys of the map `discord.Locales`.

Otherwise, Discord will yell.  So we ignore languages that are not in that list.

Sadly, the list is not well-defined, sometimes using subtags and sometimes not.

It may also be subject to changes as Discord improves their l10n support.

You can find the up-to-date list here : https://docs.discord.com/developers/reference#locales
