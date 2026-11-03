# ZZT

Tim Sweeney's ZZT (Potomac Computer Systems / Epic MegaGames, 1991), the
real-time text-mode action adventure and world editor that started Epic,
for the Asciiverse catalog. Started 2026-11-03 (David: "go, start from
zztgo"). Catalog: `kind = "game"`, `genre = "adventure"` (it is played with
the arrow keys, not a parser text adventure).

## Where it comes from

- **The Reconstruction of ZZT** (Adrian Siekierka "asie", 2020,
  github.com/asiekierka/reconstruction-of-zzt): ZZT 3.2's source rebuilt in
  Turbo Pascal, compiling byte for byte to the original. MIT, with Tim
  Sweeney's permission (Museum of ZZT, Retronauts); the only exception is
  the LZEXE tool in its build kit.
- **zztgo** (Ben Hoyt, 2020, github.com/benhoyt/zztgo, MIT): the
  Reconstruction converted to Go with his own Pascal-to-Go converter, on the
  tcell terminal library. 37 commits, last 2020-05-29. This repo is it,
  `upstream` its remote, our work on the `asciiverse` branch.
- **Worlds:** zztgo bundles TOWN.ZZT, Town of ZZT, the original shareware
  world (always free to share). The other official worlds (Caves, Dungeons,
  City) were made freeware on 1997-10-10 (Wikipedia, Museum of ZZT); the
  exact terms of that release not yet found, confirm before bundling. The
  Reconstruction ships DEMO.ZZT. The Museum of ZZT holds 3,000+ community
  worlds, each under its own author's terms.

## State found (2026-11-03)

Builds as is (Go 1.26). Runs well: title screen and Town of ZZT drawn right
in colour, sidebar, moving about. Gaps, from Ben Hoyt's write-up
(benhoyt.com/writings/zzt-in-go) and the code:

- **No sound**: `sounds.go`'s output commented out (ZZT's PC-speaker
  music and effects are square-wave tones: synthesise with a Go audio
  library, e.g. oto, as other titles do with rodio)
- **Timing**: a fixed `time.Sleep` instead of waiting for the next tick
  (game.go:1489)
- **Shooting**: Shift+arrow (ZZT's shoot) does not arrive through tcell
  (input.go:119); terminals send it differently, kitty keyboard protocol
  or an alternative key needed
- **Output**: `VideoShow()` after every text write (video.go:47, "TODO: is
  this inefficient?"): 8.8 MB of output in 9 s on the title screen; flush
  once per frame instead
- **ZZT-OOP bugs**: "definitely some bugs" (a world's title screen drawn
  wrong); test against the Reconstruction's behaviour
- **Missing**: the About screen (game.go:1404), EditorTransferBoard
  (editor.go:791)
- **Characters**: expects an IBM EGA font and zero line spacing; map the
  PC character set to Unicode so any terminal font works

## Plan

1. Flush once per frame; real tick timing
2. Shooting that works in every terminal
3. Sound
4. The PC character set in Unicode, no special font
5. ZZT-OOP: check bugs against the Reconstruction (DEMO.ZZT, Town)
6. About screen, editor transfer
7. Session logging, saves and high scores in the data folder, worlds found
   beside the binary
8. Worlds: Town and Demo; the other three once the 1997 terms are confirmed
