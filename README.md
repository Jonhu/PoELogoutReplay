# PoELogoutReplay
Protects PoE TCP RST logout traffic from packet loss by replaying logout packets.
The program runs in the background, captures a logout packet to the PoE server and replays it (default 3) times in a time interval.

Please be polite with the amount of packets that are sent to the servers.

This tool can be used in conjunction with [Lutbot](http://lutbot.com/#/) to increase the probability to log out on a patchy connection.
You can add it to the macro.ahk from Lutbot by adding: 

```run, %A_MyDocuments%\AutoHotKey\LutTools\PoELogoutReplay_v0.11.exe``` 
after ```RunWait, verify.ahk```

and moving the ```PoELogoutReplay_v0.11.exe``` to ```C:\Users\%USER%\Documents\AutoHotKey\LutTools```

# Requirements
[npcap](https://nmap.org/npcap/)

A Windows installation

# Download
Here: https://github.com/Jonhu/PoELogoutReplay/releases

# Building

In the root directory:

```go build -ldflags -H=windowsgui```

```PoELogoutReplay.exe```
This will run the app in the system tray where it can be closed.

For more options please run ```PoELogoutReplay.exe --help```

# TODO
* Test on different systems for device/interface/process grabbing.
* Make system tray prettier
* implement update?
