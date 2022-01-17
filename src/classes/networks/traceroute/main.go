package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

var (
	traceroutePort = 33434
	empty          = make([]byte, 24)

	startTTL = 1
	maxTTL   = 64

	probes = 3

	emptyIP net.IP
)

func main() {
	domain := os.Args[len(os.Args)-1]
	addr, err := net.LookupHost(domain)
	if err != nil {
		log.Fatalf("failed to look up %s", domain)
	}

	ip := net.ParseIP(addr[0])
	fmt.Printf("pinging %q (%s)\n", domain, ip)

	sender, err := NewSocket(ip, traceroutePort)
	if err != nil {
		log.Fatalf("unable to create socket: %s", err)
	}
	defer syscall.Close(sender.fd)

	data := make([]byte, 65535)
	receiver, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		log.Fatalf("unable to create socket: %s", err)
	}
	defer syscall.Close(receiver)

	err = syscall.SetsockoptTimeval(receiver, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &syscall.Timeval{
		Sec: int64(3),
	})
	if err != nil {
		log.Fatalf("failed to set timeout on icmp replies: %s", err)
	}

	pingThrottler := time.Tick(1 * time.Nanosecond)
	for ttl := startTTL; ttl <= maxTTL; ttl++ {
		err = syscall.SetsockoptInt(sender.fd, 0, syscall.IP_TTL, ttl)
		if err != nil {
			log.Fatalf("set ttl failure: %s", err)
		}

		fmt.Printf("%d ", ttl)

		var ipForStep net.IP
		for p := 0; p < probes; p++ {
			<-pingThrottler

			start := time.Now()
			req := EchoICMPRequest()
			err := Encode(req, sender)
			if err != nil {
				log.Fatalf("sending ping: %s", err)
			}

			_, from, err := syscall.Recvfrom(receiver, data, 0)
			if err != nil {
				if strings.Contains(err.Error(), "resource temporarily unavailable") {
					fmt.Printf("* ")
					continue
				}
				log.Fatalf("receving icmp data: %s", os.NewSyscallError("recvfrom", err))
			}

			var ip net.IP
			{
				addr, ok := from.(*syscall.SockaddrInet4)
				if !ok {
					addr, ok := from.(*syscall.SockaddrInet6)
					if !ok {
						continue
					}
					ip = addr.Addr[:]
				}
				ip = addr.Addr[:]
			}

			hosts, err := net.LookupAddr(ip.String())
			if err != nil {
				hosts = []string{"<unknown>"}
			}
			host := hosts[0]

			typ := GetICMPType(data)
			if typ == EchoReply || typ == TimeExceeded {

				if ipForStep.Equal(emptyIP) {
					ipForStep = ip
					fmt.Printf("%q (%s) ", host, ip)
				}

				fmt.Printf("%s ", time.Since(start))
			} else {
				log.Printf("received unknown message (type %d) from %q (%s) time=%s\n", typ, host, ip, time.Since(start))
			}

			if typ == EchoReply {
				os.Exit(0)
			}
		}
		fmt.Println()
	}
}

func GetICMPType(data []byte) ICMPType {
	dataOffset := (data[0] & 0x0f) << 2

	return ICMPType(data[dataOffset])
}

type ICMPType uint8

const (
	EchoReply              ICMPType = 0
	DestinationUnreachable ICMPType = 3
	EchoRequest            ICMPType = 8
	TimeExceeded           ICMPType = 11
)

type ICMPEcho struct {
	Typ            ICMPType
	Code           uint8
	Checksum       uint16
	Identifier     uint16
	SequenceNumber uint16
}

type ICMPTTLExpired struct {
	Typ      ICMPType
	Code     uint8
	Checksum uint16
	Unused   uint32
}

func EchoICMPRequest() ICMPEcho {
	id := uint16(rand.Uint32() >> 16)
	e := ICMPEcho{
		Typ:            EchoRequest,
		Identifier:     id,
		SequenceNumber: 1,
	}

	e.Checksum = e.CalculateChecksum()
	return e
}

func (e ICMPEcho) CalculateChecksum() uint16 {
	return ^e.binarySum()
}

func (e ICMPEcho) ValidateChecksum() bool {
	return e.Checksum+e.binarySum() == uint16(1)
}

func (e ICMPEcho) binarySum() uint16 {
	sum, carry := Add16((uint16(e.Typ)<<8)+uint16(e.Code), uint16(0), 0)
	sum, carry = Add16(sum, e.Identifier, carry)
	sum, carry = Add16(sum, e.SequenceNumber, carry)

	return sum + carry
}

func Add16(a, b uint16, carry uint16) (sum, carryOut uint16) {
	sum32 := uint32(a) + uint32(b) + uint32(carry)
	sum = uint16(sum32)
	carryOut = uint16(sum32 >> 16)

	return sum, carryOut
}

func Encode(e ICMPEcho, w io.Writer) error {
	return binary.Write(w, binary.BigEndian, e)
}

func Decode(r io.Reader) (ICMPEcho, error) {
	var e ICMPEcho
	err := binary.Read(r, binary.BigEndian, &e)
	if err != nil {
		return ICMPEcho{}, err
	}

	return e, nil
}

type socket struct {
	fd      int
	address syscall.Sockaddr
}

func NewSocket(ip net.IP, port int) (socket, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		return socket{}, os.NewSyscallError("socket", err)
	}

	socket := socket{fd: fd}

	if ip.To4() != nil {
		a := syscall.SockaddrInet4{
			Port: port,
		}
		copy(a.Addr[:], ip.To4())

		socket.address = &a
		return socket, nil
	}

	a := syscall.SockaddrInet6{
		Port: port,
	}
	copy(a.Addr[:], ip.To16())

	socket.address = &a
	return socket, nil

}

func (s socket) Write(data []byte) (n int, err error) {
	if err := syscall.Sendto(s.fd, data, 0, s.address); err != nil {
		return 0, os.NewSyscallError("sendto", err)
	}

	return len(data), nil
}

func (s socket) Read(data []byte) (n int, err error) {
	n, _, err = syscall.Recvfrom(s.fd, data, 0)
	if err != nil {
		return -1, os.NewSyscallError("recvfrom", err)
	}

	return n, nil
}
