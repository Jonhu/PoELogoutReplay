# PoELogoutReplay
Protects PoE TCP RST logout traffic from packet loss by replaying logout packets.
The program runs in the background, captures a logout packet to the PoE server and replays (default 3) times in a time interval.

Please be polite with the amount of packets that are sent to the servers.

This tool can be used in conjunction with [Lutbot](http://lutbot.com/#/) to increase the probability to log out on a patchy connection.
# Requirements
[npcap](https://nmap.org/npcap/)

A Windows installation

# Usage

In the root directory:

```go build```

```PoELogoutReplay.exe```

For more options please run ```PoELogoutReplay.exe --help```

# TODO
* Error handling
* Measure performance impact
* Test on different systems for device/interface/process grabbing.
