## Screenshots

### With no game running

![Default display](/screenshots/default-display.png)

### With Streets of Rage 4 running

![SoR 4 display](/screenshots/sor4-display.png)

### With Crypt of Necrodancer running

![Crypt of Necrodancer display](/screenshots/cotn-display.png)

# Steam game display

This app is tool for displaying logo of currently running
game in Steam on separate display (banner displays are
recommended), written in Javascript and ported to Go. When game 
has no logo, or no game is running, displays Steam logo.

Now works on every Steam-supported OS. If you run Steam on 
non-supported OS, please cross-compile this app for OS of
your Steam runtime

## How it works

It parses `appinfo.vdf` on startup and serves cached images with
align settings from steam library. When server stops, it displays
error message.

## Dependencies

This project needs [go-winres](https://github.com/tc-hib/go-winres) 
for Windows and `libgtk-3-dev` `libappindicator3-dev` for Debian 
and Ubuntu (Mint users should also install `libxapp-dev`).

## How to use

Run `go build` in project directory, and you are ready to go.
For Windows use `go generate` and `go build -ldflags "-H=windowsgui"`
Start `steamview-go.exe` in project directory, and it will connect to
<http://YOUR_IP:3000> with your default browser (don't forget to 
drag it to your cool monitor, modded into case and put it in 
fullscreen mode). When moving exe file around, don't forget `assets` 
folder, it should be in same folder as binary.

## Google Chrome startup flags

You can start Chrome on selected display with following commandline flags: 
`--window-position=x,y`, where x,y is coordinates of point inside needed
monitor (monitor marked as 1st has 0,0 coordinates at top left) and 
`--kiosk` will start it fullscreen, don't forget `http://localhost:3000/`
as url param.

Steam is &trade; & &reg; of Valve Corporation, I'm not affiliated 
with them.
