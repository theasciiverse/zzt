package main

import "testing"

// one object at 10,10 on an empty board running code
func setupObject(code ...string) {
	BoardCreate()
	Board.StatCount = 0
	for i, c := range code {
		id := int16(i + 1)
		Board.StatCount = id
		Board.Stats[id] = StatTemplateDefault
		Board.Stats[id].X, Board.Stats[id].Y = byte(10+i), 10
		Board.Tiles[10+i][10].Element = E_OBJECT
		if c != "" {
			data := c
			Board.Stats[id].Data = &data
			Board.Stats[id].DataLen = int16(len(c))
		}
	}
}

func run(id int16) {
	pos := int16(0)
	OopExecute(id, &pos, "Test")
}

func TestZapAndRestoreLabels(t *testing.T) {
	setupObject("#zap touch\r#end\r:touch\r#end\r")
	run(1)
	if got, want := *Board.Stats[1].Data, "#zap touch\r#end\r'touch\r#end\r"; got != want {
		t.Fatalf("after #zap: %q, want %q", got, want)
	}
	// Town's Jazzman: three zaps, then a restore brings all back
	code := "#zap touch\r#zap touch\r#zap touch\r#restore touch\r#end\r:touch\r#end\r:touch\r#end\r:touch\r#end\r"
	setupObject(code)
	run(1)
	if got := *Board.Stats[1].Data; got != code {
		t.Fatalf("after zaps and restore: %q, want %q", got, code)
	}
}

func TestBoundObjectsShareCode(t *testing.T) {
	setupObject("@guard\r#end\r:touch\r#end\r", "#bind guard\r#end\r")
	run(2)
	if Board.Stats[2].Data != Board.Stats[1].Data {
		t.Fatal("#bind did not share the code")
	}
	setupObject("@guard\r#zap touch\r#end\r:touch\r#end\r", "#bind guard\r#end\r")
	run(2)
	pos := int16(0)
	OopExecute(1, &pos, "Test")
	if *Board.Stats[2].Data != *Board.Stats[1].Data || *Board.Stats[2].Data == "@guard\r#zap touch\r#end\r:touch\r#end\r" {
		t.Fatalf("a #zap in one bound object should show in the other: %q", *Board.Stats[2].Data)
	}
}

func TestLongWordsCutAt20(t *testing.T) {
	setupObject("#set abcdefghijklmnopqrstuvwxyz\r#end\r")
	for i := range World.Info.Flags {
		World.Info.Flags[i] = ""
	}
	run(1)
	if World.Info.Flags[0] != "ABCDEFGHIJKLMNOPQRST" {
		t.Fatalf("flag %q, want the first 20 letters", World.Info.Flags[0])
	}
}

func TestFlagsSurviveSaving(t *testing.T) {
	var src, dst TWorldInfo
	src.Flags[0] = "GOTKEY"
	src.Flags[9] = "LAST"
	buf := make([]byte, SizeOfWorldInfo)
	for i := range buf {
		buf[i] = 0xAA
	}
	StoreWorldInfo(buf, &src)
	LoadWorldInfo(buf, &dst)
	if dst.Flags != src.Flags {
		t.Fatalf("flags after saving %q, want %q", dst.Flags, src.Flags)
	}
}

// A stat without code keeps none after boards are switched, and an object
// after it keeps its own code.
func TestNoStaleCodeAcrossBoards(t *testing.T) {
	WorldCreate()
	World.BoardCount = 1
	setupObject("@guard\r#end\r")
	World.Info.CurrentBoard = 0
	BoardClose()
	setupObject("", "@guard\r#end\r")
	World.Info.CurrentBoard = 1
	BoardClose()
	BoardOpen(0)
	BoardClose()
	BoardOpen(1)
	if Board.Stats[1].Data != nil {
		t.Fatalf("stat without code got %q", *Board.Stats[1].Data)
	}
	BoardClose()
	BoardOpen(1)
	if Board.Stats[2].Data == nil || *Board.Stats[2].Data != "@guard\r#end\r" {
		t.Fatalf("object lost its code: %v", Board.Stats[2].Data)
	}
}

func TestBlackKeyUsesGemsHighByte(t *testing.T) {
	World.Info.Gems = 5
	keySet(0, true)
	if !keyHeld(0) || World.Info.Gems != 256+5 {
		t.Fatalf("black key: held %v gems %d", keyHeld(0), World.Info.Gems)
	}
	keySet(0, false)
	if keyHeld(0) || World.Info.Gems != 5 {
		t.Fatalf("black door: held %v gems %d", keyHeld(0), World.Info.Gems)
	}
}

func TestStatlessTileDoesNotCrash(t *testing.T) {
	BoardCreate()
	if s := statRef(-1); s == nil {
		t.Fatal("no blank stat")
	}
	MoveStat(-1, 1, 1)
	RemoveStat(-1)
}

func TestKeyNames(t *testing.T) {
	for key, want := range map[int16]string{0: "Black", 1: "Blue", 5: "Purple", 7: "White"} {
		if got := keyName(key); got != want {
			t.Fatalf("key %d named %q, want %q", key, got, want)
		}
	}
}
