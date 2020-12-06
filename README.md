# KakEdit

Make calls to `$EDITOR` act upon an existing [Kakoune](https://github.com/mawww/kakoune) instance.

## Disclaimer

You should be using [connect.kak](https://github.com/alexherbo2/connect.kak).

## Why

Mainly for fun and getting back at writing Go code.
Also, leveraging client-server architecture for a text editor is still underrated.

## Build & test

Like any standard Go program.

## Usage

- Make sure `kak_session` and `kak_client` environment variables are defined.
- Run `kakedit PROGRAM`. This project was mainly written with [broot(1)](https://github.com/Canop/broot) in mind.

Without any specific option, this will run the given program and replace
`$EDITOR` (and `$VISUAL`) with something that will send an edit command
to an existing Kakoune client.

With `-mode local`, instead of sending the edit request to an existing
client, any call to `$EDITOR` will open a new client connected to the
same session. This is a blocking behavior, preferred for use-cases like
Git commit message edition where the parent process waits for its child
to exit.

## How it works

### Mode `server`

This is the default mode.

Using this mode, KakEdit will start listening to a socket and run the
program given as argument. `$EDITOR` is replaced by `kakedit -mode client <socket>` and this will make subsenquent `$EDITOR` invocations
run KakEdit in a specific mode that will just write the filename to the
socket. The running server will catch that filename and execute `echo 'evaluate-commands -client <client> edit <filename>' | kak -p <session>`
to forward the edit request to an existing Kakoune instance.

### Mode `local`

`$EDITOR` is simply replaced by a call to `kak -c <session>` where the session is read from the environment.

### Mode `client`

This is not intended for humans, see `server` mode.

## Recommended tweaks

- [Make Broot open all files in `$EDITOR`.](https://dystroy.org/broot/tricks/#change-standard-file-opening)
- Tweak Kakoune runtime to forward `kak_client` and `kak_session` to the windowing system. This makes `:terminal kakedit` become the fun bit.

## See also

The infamous [connect.kak](https://github.com/alexherbo2/connect.kak) is much better in (almost) every aspect.
If I had a reason not to choose it, it would be because it's too big for my use cases.
I also prefer a single binary over a collection of runtime scripts.
