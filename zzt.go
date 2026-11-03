// ZZT ported to Go

/*
TODO:
- ensure len>0 for things like: switch state.Lines[lpos-1][0]
- timer handling, timeout of message display
- editor: EditorTransferBoard
- performance
  + optimize when to call VideoShow
  + stat.Data should probably be []byte instead of string
*/

package main

import (
	"fmt"
	"os"
	"strings"
)

const version = "3.2.0"

const help = `zzt [WORLD] [--mute]

ZZT, Tim Sweeney's 1991 text-mode adventure and world editor. Starts on
the title screen with WORLD (default TOWN, The Town of ZZT); W picks another
world, P plays, E opens the editor.

Worlds, saved games and high scores live in the data folder
(~/.local/share/zzt on Linux); put .ZZT worlds there and W lists them.
WORLD given as a path (zzt ~/worlds/FOO.ZZT) plays from that folder.

In play: arrows move, Shift+arrows (or Space) shoot, T torch, B sound off
and on, S save, P pause, H help, Esc quits.

  --mute     start with the sound off

Logging: ZZT_LOG=1 zzt 2>zzt.log
`

func main() {
	muted := false
	startWorld := ""
	for _, a := range os.Args[1:] {
		switch {
		case a == "--help" || a == "-h":
			fmt.Print(help)
			return
		case a == "--version" || a == "-V":
			fmt.Println("zzt", version)
			return
		case a == "--mute":
			muted = true
		case !strings.HasPrefix(a, "-"):
			startWorld = a
		}
	}

	WorldFileDescCount = 7
	WorldFileDescKeys[0] = "TOWN"
	WorldFileDescValues[0] = "TOWN       The Town of ZZT"
	WorldFileDescKeys[1] = "DEMO"
	WorldFileDescValues[1] = "DEMO       Demo of the ZZT World Editor"
	WorldFileDescKeys[2] = "CAVES"
	WorldFileDescValues[2] = "CAVES      The Caves of ZZT"
	WorldFileDescKeys[3] = "DUNGEONS"
	WorldFileDescValues[3] = "DUNGEONS   The Dungeons of ZZT"
	WorldFileDescKeys[4] = "CITY"
	WorldFileDescValues[4] = "CITY       Underground City of ZZT"
	WorldFileDescKeys[5] = "BEST"
	WorldFileDescValues[5] = "BEST       The Best of ZZT"
	WorldFileDescKeys[6] = "TOUR"
	WorldFileDescValues[6] = "TOUR       Guided Tour ZZT's Other Worlds"

	startWorld = DataDirSetup(startWorld)
	StartupWorldFileName = "TOWN"
	if startWorld != "" {
		StartupWorldFileName = startWorld
	}
	ResourceDataFileName = "ZZT.DAT"
	GameTitleExitRequested = false
	EditorEnabled = true

	logf("start zzt %s world %s muted %v", version, StartupWorldFileName, muted)
	AudioInstall(muted)
	SoundEnabled = !muted
	VideoInstall()
	TextWindowInit(5, 3, 50, 18)
	VideoHideCursor()
	VideoClrScr()
	TickSpeed = 4
	DebugEnabled = false
	SavedGameFileName = "SAVED"
	SavedBoardFileName = "TEMP"
	GenerateTransitionTable()
	WorldCreate()

	GameTitleLoop()

	SoundUninstall()
	SoundClearQueue()
	VideoUninstall()
}
