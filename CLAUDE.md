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
- **Worlds:** zztgo bundles TOWN.ZZT. The Reconstruction ships DEMO.ZZT
  (in its MIT release). The other official worlds (Caves, Dungeons, City,
  Tour) and Town came free on 1997-10-10 in Epic's zzt.zip (Museum of ZZT,
  museumofzzt.com/file/z/zzt.zip), under Epic's "Software License
  Agreement" (license.txt / LICENSE.DOC, 1997-10-09), read 2026-11-03:
  distributed at no cost but "NOT public domain or free software or
  'freeware'" (clause 3: must not be described so); copy and distribute
  "only in its original, unaltered form, with all files included
  unmodified", as a complete package, no added or renamed files except
  installation routines (4); no derivative works (5); no fee beyond a
  distribution cost (6); Epic credited as owner (2); online distribution
  allowed under those terms; no CD-ROM, retail or cover disks. The
  Museum of ZZT holds 3,000+ community worlds, each under its own
  author's terms.

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

1. Flush once per frame; real tick timing: **done** (writes mark the
   screen dirty, a 60 Hz ticker shows it; time in PC timer ticks x 6 as in
   ZZT, 9.1 cycles a second at the default speed; ~1% CPU in play)
2. Shooting in every terminal: **done** (tcell v2 reads Shift on all four
   arrows; Space shoots too, as in ZZT)
3. Sound: **done** (`audio.go`: the 18.2 Hz timer as a goroutine, a square
   wave through oto v3, eased on and off; silent without a device; B and
   `--mute`)
4. The PC character set in Unicode: was already done (`codePage437ToRune`);
   colours now the exact EGA palette
5. ZZT-OOP: **done** (oop.go, elements.go and game.go compared line by
   line with OOP.PAS, ELEMENTS.PAS and GAME.PAS, the element table
   mechanically; fixed, each with a test in `oop_test.go`: #ZAP and
   #RESTORE wrote one byte early, so Town's Jazzman never got his labels
   back; saves lost the flags (never written, and the buffer was not
   cleared); a stat without code kept the last board's, so a later
   object could be saved bound to it and lose its code; #BIND copies
   shared nothing (`Data` is now a shared `*string`, as ZZT's pointer);
   black keys and doors crashed (now ZZT's gem high byte); words and labels
   past 20 characters matched (ZZT's string[20] cuts them); statless
   objects, bombs and pushers crashed; an empty debug prompt crashed)
6. About screen, editor transfer: **done** (help from the Reconstruction's
   DOC/*.HLP embedded, `help/`; About at start; T in the editor)
7. Session logging, data folder: **done** (`ZZT_LOG=1 zzt 2>zzt.log`;
   `datadir.go`: XDG data home on Linux, Town and Demo put there on first
   run, a world given by path plays from its own folder)
8. Worlds: Town and Demo bundled now; the 1997 terms above allow passing
   on zzt.zip only whole and unaltered, so the other worlds would come from
   that zip as is (read from it, not extracted into ours); David to decide

## Notes

- castrec keys `sup sdown sleft sright` send Shift+arrows (added for this).
- A fresh world in the editor asks for the room title first.
- `SoundFreqTable` truncates 32 * 2^3 to 255 Hz, as the float maths gives.
