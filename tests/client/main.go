package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/gopacket/pcap"
)

type testBinding struct {
	destip   string
	srcip    string
	srcport  int
	destport int
	device   string
}

func (p *testBinding) eq(p2 testBinding) bool {
	return p.srcport == p2.srcport && p.destport == p2.destport && p.srcip == p2.srcip && p.device == p2.device && p.destip == p2.destip
}

func findTestBinding() testBinding {
	return testBinding{}
}

func main() {
	strEcho := "Halo"
	servHost := "localhost"
	servPort := 9003
	servAddr := servHost + ":" + strconv.Itoa(servPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte(strEcho))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	println("write to server = ", strEcho)

	reply := make([]byte, 1024)

	_, err = conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}
	add := strings.Split(conn.LocalAddr().String(), ":")
	localAddr := add[0]
	localPort, _ := strconv.Atoi(add[1])
	binding := testBinding{
		destip:   servHost,
		destport: servPort,
		srcip:    localAddr,
		srcport:  localPort,
		device:   "enps0s3",
	}
	log.Printf("Binding: %v", binding)
	devices, err := pcap.FindAllDevs()
	log.Printf("Devices: %v, err %v", devices, err)
	time.Sleep(500 * time.Second)

	println("reply from server=", string(reply))
	flushDuration := flag.Duration("fl", 69*time.Millisecond, "flush duration of pcap capture")
	repeats := flag.Int("r", 30, "amount of logout repeats")
	instancePollDur := flag.Duration("ip", 1*time.Second, "time waiting between instance data poll")
	logoutSpreadDur := flag.Duration("lp", 20*time.Millisecond, "time waiting between logout packets")
	packetPollDur := flag.Duration("pp", 69*time.Millisecond, "time waiting between polls")
	filterStr := flag.String("fi", "", "BPF filter string")

	flag.Parse()
	var handle *pcap.Handle
	filter := *filterStr
	ticker := time.NewTicker(*instancePollDur)
	defer ticker.Stop()
	i := make(chan struct{}, 3)
	init := false

	for {
		select {
		case <-ticker.C:
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
			conn.Close() //sends RST packet
		case <-i:
			data, _, err := handle.ReadPacketData()
			if err != nil {
				time.Sleep(*packetPollDur)
				i <- struct{}{}
			} else {
				handle.SetDirection(pcap.DirectionIn) // let's not capture our own packets
				for repeatSend := 0; repeatSend < *repeats; repeatSend++ {
					log.Printf("Sending paket %d", repeatSend)
					if err := handle.WritePacketData(data); err != nil {
						log.Printf("Error writing packet: %v", err)
						handle.SetDirection(pcap.DirectionInOut)
						continue
					}
					time.Sleep(*logoutSpreadDur)
				}
				handle.SetDirection(pcap.DirectionInOut)
			}
		}
	}
}
