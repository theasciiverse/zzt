package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"testing"
)

// zzt.zip must be passed on exactly as Epic made it available.
func TestEpicZipUnaltered(t *testing.T) {
	data, err := os.ReadFile("epic/zzt.zip")
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	if got := hex.EncodeToString(sum[:]); got != epicZipSHA256 {
		t.Fatalf("epic/zzt.zip sha256 %s, want %s: it must not be changed", got, epicZipSHA256)
	}
}

func TestEpicWorldsLoad(t *testing.T) {
	z, err := zip.OpenReader("epic/zzt.zip")
	if err != nil {
		t.Fatal(err)
	}
	old := epicZip
	epicZip = z
	defer func() { epicZip.Close(); epicZip = old }()
	dir := t.TempDir()
	wd, _ := os.Getwd()
	os.Chdir(dir) // nothing here: everything comes from the zip
	defer os.Chdir(wd)

	worlds := gameFiles(".ZZT")
	want := map[string]bool{"TOWN.ZZT": true, "CAVES.ZZT": true, "DUNGEONS.ZZT": true, "CITY.ZZT": true, "TOUR.ZZT": true, "DEMO.ZZT": true}
	if len(worlds) != len(want) {
		t.Fatalf("worlds %v", worlds)
	}
	for _, w := range worlds {
		if !want[w] {
			t.Fatalf("unexpected world %s", w)
		}
		f, err := openGameFile(w)
		if err != nil {
			t.Fatalf("%s: %v", w, err)
		}
		head := make([]byte, 4)
		io.ReadFull(f, head)
		f.Close()
		if LoadInt16(head[:2]) != -1 || LoadInt16(head[2:]) <= 0 {
			t.Fatalf("%s does not start like a ZZT world: % x", w, head)
		}
	}
	if f, err := openGameFile("TOWN.HI"); err != nil {
		t.Fatalf("TOWN.HI: %v", err)
	} else {
		f.Close()
	}
}
