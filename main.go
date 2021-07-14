package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/google/gopacket/pcap"
)

const POE_STEAM string = "PathOfExileSteam.exe"
const POE_STANDALONE string = "PathOfExile.exe"

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
	flushDuration := flag.Duration("fl", 69*time.Millisecond, "flush duration of pcap capture")
	repeats := flag.Int("r", 3, "amount of logout repeats")
	instancePollDur := flag.Duration("ip", 1*time.Second, "time waiting between instance data poll")
	logoutSpreadDur := flag.Duration("lp", 200*time.Millisecond, "time waiting between logout packets")
	packetPollDur := flag.Duration("pp", 69*time.Millisecond, "time waiting between polls")
	filterStr := flag.String("fi", "", "BPF filter string")

	flag.Parse()

	var handle *pcap.Handle
	var binding poeBinding
	var err error
	filter := *filterStr
	ticker := time.NewTicker(*instancePollDur)
	defer ticker.Stop()
	done := make(chan bool)
	i := make(chan struct{}, 3)
	init := false
	previousBinding := poeBinding{}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Printf("captured %v, stopping packet replay", sig)
			done <- true
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			binding = findPoeBinding()
			//log.Printf("%v", findPoeBinding())
			if !binding.eq(previousBinding) && binding.srcip != "" {
				if !init {
					init = true
					i <- struct{}{}
				}
				if *filterStr == "" {
					filter = fmt.Sprintf("tcp[tcpflags] & (tcp-rst|tcp-ack) == (tcp-rst|tcp-ack) and tcp src port %d and tcp dst port %d and dst host %s", binding.srcport, binding.destport, binding.destip)
				}
				handle, err = pcap.OpenLive(binding.device, 1024, false, *flushDuration/2)
				if err != nil {
					log.Printf("Error opening pcap handle: %s", err)
				}
				if err := handle.SetBPFFilter(filter); err != nil {
					log.Printf("Error setting BPF filter: %s", err)
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
				for repeatSend := 0; repeatSend < *repeats; repeatSend++ {
					//log.Printf("received package: %v", data)
					if err := handle.WritePacketData(data); err != nil {
						log.Printf("Error: %v", err)
					}
					time.Sleep(*logoutSpreadDur)
				}
				init = false
			}
		}
	}
}

func findPoeBinding() poeBinding {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	tabs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Established && (s.Process.Name == POE_STEAM || s.Process.Name == POE_STANDALONE)
	})
	if err != nil {
		log.Printf("error getting PoE process: %v", err)
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
