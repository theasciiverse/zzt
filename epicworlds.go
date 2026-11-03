package main

// Epic's official worlds: Town, Caves, Dungeons, City, the Tour and the
// demo, as Epic MegaGames made ZZT available at no cost on 1997-10-10. Its
// licence lets the package be passed on only whole and unaltered, so it is
// shipped as it came, zzt.zip beside the game, and the worlds are read
// from inside it, never unpacked. A world or high score file of the same
// name in the folder ZZT works in comes first.

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// sha256 of Epic's zzt.zip (Museum of ZZT, museumofzzt.com/file/z/zzt.zip)
const epicZipSHA256 = "fdf5a0f9f45c32447199aaad3d3caa5dbabcb6ec6a5f55cd954d060515488d5e"

var epicZip *zip.ReadCloser

// epicZipPath finds zzt.zip beside the binary (an install, the release
// archives) or in the source tree's epic/ folder (a development build).
func epicZipPath() string {
	var dirs []string
	if exe, err := os.Executable(); err == nil {
		if real, err := filepath.EvalSymlinks(exe); err == nil {
			exe = real
		}
		d := filepath.Dir(exe)
		dirs = append(dirs, d, filepath.Join(d, "epic"))
	}
	for _, d := range dirs {
		p := filepath.Join(d, "zzt.zip")
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	return ""
}

// EpicWorldsOpen opens zzt.zip; call it before moving into the data folder.
func EpicWorldsOpen() {
	p := epicZipPath()
	if p == "" {
		logf("no zzt.zip found beside the game")
		return
	}
	z, err := zip.OpenReader(p)
	if err != nil {
		logf("zzt.zip at %s: %v", p, err)
		return
	}
	epicZip = z
	logf("Epic's worlds from %s", p)
}

func epicEntry(name string) *zip.File {
	if epicZip == nil {
		return nil
	}
	for _, f := range epicZip.File {
		if strings.EqualFold(f.Name, name) {
			return f
		}
	}
	return nil
}

type readCloser struct{ *bytes.Reader }

func (readCloser) Close() error { return nil }

// openGameFile opens name in the folder ZZT works in, or else from zzt.zip.
func openGameFile(name string) (io.ReadCloser, error) {
	f, err := os.Open(name)
	if err == nil {
		return f, nil
	}
	if e := epicEntry(name); e != nil {
		r, zerr := e.Open()
		if zerr != nil {
			return nil, zerr
		}
		defer r.Close()
		data, zerr := io.ReadAll(r)
		if zerr != nil {
			return nil, zerr
		}
		return readCloser{bytes.NewReader(data)}, nil
	}
	return nil, err
}

// gameFiles lists files with extension (".ZZT", ".SAV") in the folder ZZT
// works in and, for worlds, in zzt.zip, each name once, sorted.
func gameFiles(extension string) []string {
	seen := map[string]bool{}
	var out []string
	add := func(name string) {
		key := strings.ToUpper(name)
		if !seen[key] {
			seen[key] = true
			out = append(out, name)
		}
	}
	if matches, err := filepath.Glob("*" + extension); err == nil {
		for _, m := range matches {
			add(m)
		}
	}
	if epicZip != nil && extension == ".ZZT" {
		for _, f := range epicZip.File {
			if strings.HasSuffix(strings.ToUpper(f.Name), extension) {
				add(f.Name)
			}
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return strings.ToUpper(out[i]) < strings.ToUpper(out[j]) })
	return out
}
