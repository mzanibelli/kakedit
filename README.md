# KakEdit

Make calls to `$EDITOR` act upon an existing [Kakoune](https://github.com/mawww/kakoune) instance.

## Disclaimer

You should be using [connect.kak](https://github.com/alexherbo2/connect.kak).

## Why

Mainly for fun and getting back at writing Go code.
Also, leveraging client-server architecture for a text editor is still underrated.

## Build & test

See `Makefile`.

## Usage

- Make sure `kak_session` and `kak_client` environment variables are defined.
- Run `kakedit PROGRAM`. This project was mainly written with [broot(1)](https://github.com/Canop/broot) in mind.

This will run the given program and replace `$EDITOR` (and `$VISUAL`) with
something that will send an edit command to an existing Kakoune client.

## How it works

### `kakedit`

Start listening to a socket and run the program given as
argument. `$EDITOR` is replaced by `kakpipe <socket>` and this will
make subsequent `$EDITOR` invocations just write the filename to the
socket. The running server will catch that filename and send an edit
command to an existing Kakoune instance.

### `kakpipe`

The first argument is the path to a socket and all the remaining arguments form a string that is written to the socket.

### `kakwrap`

Start a new Kakoune client inside an existing session that is unique for each directory.
If the session does not exist yet, start it too.

## Recommended tweaks

- [Make Broot open all files in `$EDITOR`.](https://dystroy.org/broot/tricks/#change-standard-file-opening)
- Tweak Kakoune runtime to forward `kak_client` and `kak_session` to the windowing system. This makes `:terminal kakedit` become the fun bit.

## See also

The infamous [connect.kak](https://github.com/alexherbo2/connect.kak) is much better in (almost) every aspect.
It's too big for my use cases and I prefer a single binary over a collection of runtime scripts.
