package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
)

var serverIP = net.ParseIP("0.0.0.0")
var serverPort = 9999

var targetIP = net.ParseIP("127.0.0.1")
var targetPort = 6666

func main() {
	fmt.Println("hi")

	s, err := listen(serverIP, serverPort)

	if err != nil {
		log.Fatalf("failed to listen on %s:%d: %s", serverIP, serverPort, err)
	}

	log.Printf("listening on %s:%d...", serverIP, serverPort)

	for {
		from, err := accept(s)
		if err != nil {
			log.Fatalln(os.NewSyscallError("accept", err))
		}

		go func() {
			to, err := connect(targetIP, targetPort)
			if err != nil {
				log.Printf("failed to connect to %s:%d :%s", targetIP, targetPort, err)
			}

			err = proxyOnce(from, to)
			if err != nil {
				log.Printf("while proxying request: %s", err)
			}

			log.Printf("closing sockets...")

			if err = from.close(); err != nil {
				log.Printf("failed to close 'from' socket: %s", err)
			}

			if err = to.close(); err != nil {
				log.Printf("failed to close 'to' socket: %s", err)
			}

			log.Printf("closed sockets...")
		}()
	}
}

func echo(s socket) error {
	return proxyOnce(s, s)
}

func proxyOnce(from, to socket) error {
	var errGroup errgroup.Group

	var requestData = make([]byte, 65535)

	for {
		data := make([]byte, 65535)
		n, err := from.read(data)
		if err != nil {
			return fmt.Errorf("while receiving request: %s", err)
		}

		requestData = append(requestData, data[:n]...)

		if strings.HasSuffix(string(requestData), "\r\n\r\n") {
			break
		}
	}

	err := to.write(requestData)
	if err != nil {
		return fmt.Errorf("failed to write: %s", err)
	}

	errGroup.Go(func() error {
		return forwardOnce(from, to)
	})

	errGroup.Go(func() error {
		return forwardOnce(to, from)
	})

	return errGroup.Wait()
}

func forwardOnce(from, to socket) error {
	data := make([]byte, 65535)
	n, err := from.read(data)
	if err != nil {
		return fmt.Errorf("while receiving request: %s", err)
	}

	// I don't know how else to detect if the other side has hung up
	if n == 0 {
		log.Printf("%s hung up!", from)
		return io.EOF
	}

	log.Printf("got data (%d bytes): %s", n, data[:n])

	err = to.write(data[:n])
	if err != nil {
		return fmt.Errorf("while sending request: %s", err)
	}

	return nil
}

type socket struct {
	fd      int
	address syscall.Sockaddr
}

func (s socket) read(data []byte) (int, error) {
	n, _, err := syscall.Recvfrom(s.fd, data, syscall.MSG_WAITALL)
	if err != nil {
		return -1, os.NewSyscallError("recvfrom", err)
	}

	return n, nil
}

func (s socket) write(data []byte) error {
	if err := syscall.Sendto(s.fd, data, 0, s.address); err != nil {
		return os.NewSyscallError("sendto", err)
	}

	return nil
}

func (s socket) close() error {
	if err := syscall.Close(s.fd); err != nil {
		return os.NewSyscallError("close", err)
	}

	return nil
}

func accept(s socket) (socket, error) {
	nfd, sa, err := syscall.Accept(s.fd)
	if err != nil {
		return socket{}, os.NewSyscallError("accept", err)
	}

	return socket{nfd, sa}, nil
}

func connect(ip net.IP, port int) (socket, error) {
	s, err := NewSocket(ip, port)
	if err != nil {
		return socket{}, fmt.Errorf("creating socket: %w", err)
	}

	if err := syscall.Connect(s.fd, s.address); err != nil {
		return socket{}, fmt.Errorf("while establishing connection: %w", os.NewSyscallError("connect", err))
	}

	return s, err
}

func listen(ip net.IP, port int) (socket, error) {
	s, err := NewSocket(ip, port)
	if err != nil {
		return socket{}, fmt.Errorf("creating socket: %w", err)
	}

	if err := syscall.Bind(s.fd, s.address); err != nil {
		return socket{}, os.NewSyscallError("bind", err)
	}

	if err := syscall.Listen(s.fd, syscall.SOMAXCONN); err != nil {
		return socket{}, os.NewSyscallError("listen", err)
	}

	return s, nil
}

func NewSocket(ip net.IP, port int) (socket, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
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
