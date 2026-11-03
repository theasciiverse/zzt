package main

import "embed"

// The documentation ZZT 3.2 shipped in ZZT.DAT, in the Reconstruction of
// ZZT's source form (DOC/*.HLP).
//
//go:embed help/*.HLP
var helpFiles embed.FS
