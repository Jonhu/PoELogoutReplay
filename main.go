package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/getlantern/systray"
	"github.com/google/go-github/github"
	"github.com/google/gopacket/pcap"
	"github.com/inconshreveable/go-update"
)

//go build -ldflags -H=windowsgui
const POE_STEAM string = "PathOfExileSteam.exe"
const POE2_STEAM string = "PathOfExile2Steam.exe"
const POE_STANDALONE string = "PathOfExile.exe"
const POE2_STANDALONE string = "PathOfExile2.exe"
const iconBase64 string = "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"
const version = "v0.2"
const URL = "https://github.com/Jonhu/PoELogoutReplay"
const TagURL = "https://github.com/Jonhu/PoELogoutReplay/releases/tag/"
const EXE = "PoELogoutReplay.exe" //for now windows only

type poeBinding struct {
	destip   string
	srcip    string
	srcport  int
	destport int
	device   string
}

func (p *poeBinding) eq(p2 poeBinding) bool {
	return p.srcport == p2.srcport && p.destport == p2.destport && p.srcip == p2.srcip && p.device == p2.device && p.destip == p2.destip
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Printf("captured %v, stopping packet replay", sig)
			onExit()
		}
	}()

	go getAndRunNewestVersion()
	go systray.Run(onReady, onExit)

	flushDuration := flag.Duration("fl", 69*time.Millisecond, "flush duration of pcap capture")
	repeats := flag.Int("r", 3, "amount of logout repeats")
	instancePollDur := flag.Duration("ip", 1*time.Second, "time waiting between instance data poll")
	logoutSpreadDur := flag.Duration("lp", 200*time.Millisecond, "time waiting between logout packets")
	packetPollDur := flag.Duration("pp", 69*time.Millisecond, "time waiting between polls")
	filterStr := flag.String("fi", "", "BPF filter string")
	log.Printf("PoELogoutReplay version %s, use CTRL+C to exit", version)

	flag.Parse()
	var handle *pcap.Handle
	var binding poeBinding
	var err error
	filter := *filterStr
	ticker := time.NewTicker(*instancePollDur)
	defer ticker.Stop()
	i := make(chan struct{}, 3)
	init := false
	previousBinding := poeBinding{}
	for {
		select {
		case <-ticker.C:
			binding = findPoeBinding()
			if !binding.eq(previousBinding) && binding.srcip != "" && binding.destport != 443 { //we can sometimes capture in client https
				if !init {
					init = true
					i <- struct{}{}
				}
				if *filterStr == "" {
					filter = fmt.Sprintf("tcp[tcpflags] & (tcp-rst|tcp-ack) == (tcp-rst|tcp-ack) and tcp src port %d and tcp dst port %d and dst host %s",
						binding.srcport, binding.destport, binding.destip)
				}
				handle, err = pcap.OpenLive(binding.device, 1024, false, *flushDuration/2)
				if err != nil {
					log.Printf("Error opening pcap handle: %s", err)
					ticker.Reset(*instancePollDur)
					continue
				}
				if err := handle.SetBPFFilter(filter); err != nil {
					log.Printf("Error setting BPF filter: %s", err)
					break //can't recover from malformed filter
				}
				previousBinding = binding
			}
			ticker.Reset(*instancePollDur)
		case <-i:
			data, _, err := handle.ReadPacketData()
			if err != nil {
				time.Sleep(*packetPollDur)
				i <- struct{}{}
			} else {
				handle.SetDirection(pcap.DirectionIn) // let's not capture our own packets
				for repeatSend := 0; repeatSend < *repeats; repeatSend++ {
					if err := handle.WritePacketData(data); err != nil {
						log.Printf("Error writing packet: %v", err)
						handle.SetDirection(pcap.DirectionInOut)
						continue
					}
					time.Sleep(*logoutSpreadDur)
				}
				handle.SetDirection(pcap.DirectionInOut)
				init = false
			}
		}
	}
}

func findPoeBinding() poeBinding {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Printf("error getting devices: %v", err)
		return poeBinding{}
	}

	tabs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Established && (s.Process.Name == POE_STEAM || s.Process.Name == POE_STANDALONE || s.Process.Name == POE2_STEAM || s.Process.Name == POE2_STANDALONE)
	})
	if err != nil {
		log.Printf("error getting PoE process: %v", err)
		return poeBinding{}
	}
	for _, device := range devices {
		for _, address := range device.Addresses {
			for _, poeEntry := range tabs { //take the first if both steam and standalone are running
				if poeEntry.LocalAddr.IP.String() == address.IP.String() {
					return poeBinding{
						srcip:    poeEntry.LocalAddr.IP.String(),
						destip:   poeEntry.RemoteAddr.IP.String(),
						device:   device.Name,
						srcport:  int(poeEntry.LocalAddr.Port),
						destport: int(poeEntry.RemoteAddr.Port),
					}
				}
			}
		}
	}
	return poeBinding{}
}

func onReady() {
	icon, _ := base64.StdEncoding.DecodeString(iconBase64)
	systray.SetIcon(icon)
	systray.SetTitle("PoELogoutReplay")
	systray.SetTooltip("Replays Logout packets of PoE")
	mRelease := systray.AddMenuItem("GitHub", "Go to Github page")
	mAbout := systray.AddMenuItem(version, fmt.Sprintf("Version: %s", version))
	mQuit := systray.AddMenuItem("Quit", "Quit Logout Replay")
	go func() {
		select {
		case <-mQuit.ClickedCh:
			systray.Quit()
			onExit()
		case <-mRelease.ClickedCh:
			OpenBrowser(URL)
		case <-mAbout.ClickedCh:
			OpenBrowser(TagURL + version)
		}

	}()
}

func onExit() {
	os.Exit(0)
}

func OpenBrowser(url string) {
	log.Printf("Opening browser at %s", url)
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Print(err)
	}
}

func updateError(err string) {
	log.Printf("PoELogoutReplay update: %s", err)
}

func getAndRunNewestVersion() {
	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases(context.Background(), "jonhu", "poelogoutreplay", nil)
	if err != nil {
		updateError("internet error:")
		return
	}
	if len(releases) == 0 {
		updateError("repository not available? RIP maintainer?")
		return
	}
	if version < *releases[0].TagName {
		log.Printf("Found new version: %v", *releases[0].TagName)
		url := releases[0].Assets[0].GetBrowserDownloadURL()
		resp, err := http.Get(url)
		if err != nil {
			updateError(fmt.Sprintf("update failed with %s", err.Error()))
			return
		}
		defer resp.Body.Close()
		ex, _ := os.Executable()
		exPath, _ := filepath.Abs(ex)
		fmt.Print("Current Executable path: " + exPath)
		log.Printf("Old checksum: %v", hashExe())
		err = update.Apply(resp.Body, update.Options{TargetPath: exPath})
		if err != nil {
			updateError(fmt.Sprintf("update exe failed with %s", err.Error()))
			return
		} else {
			log.Printf("New checksum: %v", hashExe())
			log.Printf("PoELogoutReplay updated! release notes: %s", *releases[0].HTMLURL)
			time.Sleep(5 * time.Second)
			restartProgram() //restarts for hot reload
		}
	}
}

//creates child and abandons it afterwards
func restartProgram(args ...string) {
	log.Print("PoELogoutReplay restart...")
	cmd := exec.Command(EXE, flag.Args()...)
	err := cmd.Start()
	if err != nil {
		log.Printf("gracefulRestart: Failed to launch, error: %v", err)
	}
	onExit()
}

func hashExe() string {
	hasher := sha256.New()
	f, err := os.Open(EXE)
	if err != nil {
		log.Print(err)
		return ""
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Print(err)
		return ""
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
