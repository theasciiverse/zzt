# ZZT

**A smiley face, a town full of lions, five purple keys, and a world editor that started a thirty-year community.**

In January 1991 Tim Sweeney, a mechanical engineering student, released ZZT from his bedroom, with the orders mailed to his parents' address: an adventure drawn entirely in the PC's text characters, where you walk the board as a ☻, shoot at ruffians and tigers, push boulders, light torches in dark caves and talk to objects that have their own little programs. Then he handed every player the same editor he used to build it. Thousands of worlds followed, and the company he founded to sell it became Epic.

This is ZZT 3.2 in your terminal, rebuilt from Adrian Siekierka's Reconstruction of ZZT, the original source recreated byte for byte, by way of Ben Hoyt's Go translation.

- **Town of ZZT**, the original adventure: find the five purple keys and reach the palace
- **The demo world**, a guided tour of every creature, item and terrain
- **The world editor** with ZZT-OOP, the little language every object speaks; export and import boards
- **The PC speaker**, synthesised: every bleep, drum and tune, at ZZT's own 18.2 Hz
- **ZZT's own clock**: 9.1 cycles a second at the default speed, just as on a 1991 PC
- **The EGA palette**, exact, whatever your terminal's theme
- **Your own worlds**: drop any `.ZZT` file into the data folder and press W

## Controls

- **Arrows**: move; **Shift+arrows**: shoot (**Space** shoots the way you last moved)
- **T**: light a torch; **B**: sound off or on; **H**: help
- **S**: save; **P**: pause; **Esc** or **Q**: quit to the title
- On the title screen: **P** plays, **W** picks a world, **R** restores a saved game, **E** opens the editor, **S** sets the speed, **A** about, **H** high scores
- In the editor: **H** shows its help, **T** transfers a board

## Options

- `zzt WORLD`: start with another world (`zzt DEMO`), or a world file by its path (`zzt ~/worlds/FOO.ZZT`), which then plays from its own folder
- `--mute`: start with the sound off

Worlds, saved games and high scores live in the data folder: `~/.local/share/zzt` on Linux, `~/Library/Application Support/zzt` on macOS, `%AppData%\zzt` on Windows.

Session log: `ZZT_LOG=1 zzt 2>zzt.log`

## Credits

- ZZT: Tim Sweeney, Epic MegaGames, 1991
- The Reconstruction of ZZT: Adrian Siekierka, 2020 (MIT, with Tim Sweeney's permission; `LICENSE-reconstruction`); its documentation and demo world are built in
- zztgo, the Go translation: Ben Hoyt, 2020 (MIT; `LICENSE.txt`)
- Sound, timing, keys, help, the data folder and the ZZT-OOP fixes: Asciiverse
