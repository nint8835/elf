# elf
advent of code discord intergration

## Running
1. Ensure you have [Go 1.19 or newer](https://golang.org/doc/install) installed
2. Run `go mod download`
3. Copy `.env.dist` to `.env` and populate it's contents
   - `ELF_DISCORD_TOKEN`: Discord token for the bot user to run as
   - `ELF_DISCORD_APP_ID`: ID of the Discord app the bot belongs to
   - `ELF_ADVENT_OF_CODE_SESSION`: Advent of Code session cookie for the bot user
4. Run `go run ./cmd/elf`

## Registering Commands
In order to use the bot, you need to register it's slash commands either globally to your app, or to the guild that you will be testing on. You can do this with the `./cmd/registercommands` tool. You can use it like the following:

```sh
go run ./cmd/registercommands --guild-id 514110851016556567
```

Or, to register globally:

```sh
go run ./cmd/registercommands --global
```

Note that registering commands globally can take up to an hour to fully apply, so for development it is recommended you register at a guild level.

## Unregistering Commands
If you need to clean up command registrations (for example, to remove a deleted command or to migrate from guild to global commands), you can re-use the same `./cmd/registercommands` tool. Just add the `--unregister` flag to any invocation of the tool and it will clean up all registered commands in a given scope, rather than registering the new commands for that scope.

## Manually Registering Guilds

If you are developing (or if the onboarding process still doesn't exist as of you reading this), it is helpful to manually insert a testing guild into the database. This can be done with the `testdata` tool included.

Example:

```sh
go run ./cmd/testdata --channel-id 909857762064871444 --guild-id 514110851016556567 --leaderboard-code 1111111-11111111 --leaderboard-id 0000001
```