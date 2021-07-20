# PoELogoutReplay
Protects PoE TCP RST logout traffic from packet loss by replaying logout packets.
The program runs in the background, captures a logout packet to the PoE server and replays (default 3) times in a time interval.

Please be polite with the amount of packets that are sent to the servers.

This tool can be used in conjunction with [Lutbot](http://lutbot.com/#/) to increase the probability to log out on a patchy connection.
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
