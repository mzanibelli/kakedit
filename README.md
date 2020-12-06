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

## Tweaks

- Make Broot open all files in `$EDITOR`.
- Tweak Kakoune runtime to forward `kak_client` and `kak_session` to the windowing system. This makes `:terminal kakedit` become the fun bit.

## See also

The infamous [connect.kak](https://github.com/alexherbo2/connect.kak) is much better in (almost) every aspect.
If I had a reason not to choose it, it would be because it's too big for my use cases and heavily relies on Unix-like underlying OS.
I also prefer a single binary over a collection of runtime scripts.
