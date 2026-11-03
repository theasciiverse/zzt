#!/bin/sh
# Records asciiverse/zzt.cast, the catalog preview (100x30), with hub's
# castrec (a sibling checkout of hub, ../../hub). Run from the repo after
# go build -o zzt .
#
# The About screen, the title screen with its creatures, then Town of ZZT:
# around the welcome object, three torches, three gems, and into the Bank
# of ZZT, dark inside until a torch is lit.
set -eu
hub=${HUB:-../../hub}
home=$(mktemp -d)
s="w2;kesc;w4.5;kp;w1.2"
s="$s;r5:up:0.15;w0.2;r12:right:0.15;w0.6"
s="$s;kright;w0.3;kup;w0.6;kdown;w0.2;r2:right:0.15;w0.6"
s="$s;r10:right:0.15;w0.4;r2:right:0.15;w0.4;kup;w0.2;kleft;w0.8"
s="$s;r12:left:0.15;w0.2;r3:up:0.15;w0.3;kleft;w1.5;kup;w0.3;kt;w1;r4:up:0.2;w0.5;r5:right:0.2;w2"
CAST_COLS=100 CAST_ROWS=30 HOME="$home" XDG_DATA_HOME="$home/data" \
  python3 "$hub/tools/castrec.py" asciiverse/zzt.cast . "$s" -- ./zzt
rm -rf "$home"
