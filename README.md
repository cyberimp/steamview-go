# Steam game display

This app is tool for displaying logo of currently running
game in Steam on separate display (banner displays are
recommended), written in Javascript and ported to Go. When game 
has no logo, or no game is running, displays Steam logo.

Now works only on Windows, if you know, how to detect running
Steam game on Linux/macOS/whatever, message me or make a PR 
(Steam handling logic is in `./steam/steam.go`)


## How to use

Run `go build` in project directory, and you are ready to go.
Start `steamview-go.exe` in project directory, and connect to
<http://YOUR_IP:3000> with your favourite browser (don't forget to 
drag it to your cool monitor, modded into case and put it in 
fullscreen mode). There are settings of aligning logo on handle
<http://YOUR_IP:3000/align> that you can change by pressing 
buttons. When moving exe file around, don't forget `assets` 
folder and `database.db` file, they should be in same folder
as binary.

Steam is &trade; & &reg; of Valve Corporation, I'm not affiliated 
with them.
