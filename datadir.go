package main

// Where ZZT keeps its files. ZZT worked in the current directory: worlds,
// saved games (.SAV), high scores (.HI) and exported boards (.BRD) all
// side by side, and W listed the worlds there. This keeps that, in one
// folder per player: the data folder, with the demo world put there on
// first run (Epic's worlds are read from zzt.zip, see epicworlds.go). A world named on the command line by its path
// (zzt ~/worlds/FOO.ZZT) is played from its own folder instead.

import (
	"embed"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed DEMO.ZZT
var bundledWorlds embed.FS

func dataDir() string {
	if runtime.GOOS == "linux" || runtime.GOOS == "freebsd" || runtime.GOOS == "openbsd" || runtime.GOOS == "netbsd" {
		if d := os.Getenv("XDG_DATA_HOME"); d != "" {
			return filepath.Join(d, "zzt")
		}
		if h, err := os.UserHomeDir(); err == nil {
			return filepath.Join(h, ".local", "share", "zzt")
		}
	}
	if d, err := os.UserConfigDir(); err == nil {
		return filepath.Join(d, "zzt")
	}
	return "."
}

// worldFile finds arg as a world file, with or without its .ZZT.
func worldFile(arg string) string {
	for _, p := range []string{arg, arg + ".ZZT", arg + ".zzt"} {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	return ""
}

// DataDirSetup moves into the folder ZZT works in and returns the world to
// start with ("" for the default).
func DataDirSetup(arg string) string {
	if arg != "" && (strings.ContainsAny(arg, `/\`) || worldFile(arg) != "") {
		if p := worldFile(arg); p != "" {
			abs, _ := filepath.Abs(p)
			if err := os.Chdir(filepath.Dir(abs)); err == nil {
				logf("playing from %s", filepath.Dir(abs))
			}
			name := filepath.Base(abs)
			return strings.TrimSuffix(strings.TrimSuffix(name, ".ZZT"), ".zzt")
		}
	}
	dir := dataDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		logf("no data folder %s: %v", dir, err)
		return arg
	}
	entries, _ := bundledWorlds.ReadDir(".")
	for _, e := range entries {
		dst := filepath.Join(dir, e.Name())
		if _, err := os.Stat(dst); err == nil {
			continue // the player's copy, maybe edited, stays
		}
		data, _ := bundledWorlds.ReadFile(e.Name())
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			logf("could not put %s in %s: %v", e.Name(), dir, err)
		}
	}
	if err := os.Chdir(dir); err != nil {
		logf("could not enter %s: %v", dir, err)
	} else {
		logf("data folder %s", dir)
	}
	return strings.TrimSuffix(strings.TrimSuffix(arg, ".ZZT"), ".zzt")
}
