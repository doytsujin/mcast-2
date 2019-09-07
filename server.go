package Chat

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var IP = [4]byte{127, 0, 0, 1}
const (
	Fd = "myown" // AF_UNIX
	Port = 3000 // AF_INET | AF_INET6
)


func main() {
	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Socket:", s)

	if err = syscall.SetsockoptInt(s, syscall.IPPROTO_IP, syscall.SO_REUSEADDR, 1); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Set opt: SO_REUSEADDR")

	if err = syscall.SetsockoptInt(s, syscall.IPPROTO_IP, syscall.IP_ADD_MEMBERSHIP, 1); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Set opt: IP_MULTICAST_LOOP")

	linger := syscall.Linger{
		Onoff: 1,
		Linger: 1,
	}
	if err = syscall.SetsockoptLinger(s, syscall.SOL_SOCKET, syscall.SO_LINGER, &linger); err != nil {
		log.Fatalln(err)
	}

	//if err = syscall.SetNonblock(s, true); err != nil {
	//	log.Fatalln(err)
	//}

	if err = syscall.Bind(s, &syscall.SockaddrInet4{Addr: IP, Port: Port}); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Bind socket:", s)

	backlog := 0
	if err = syscall.Listen(s, backlog); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Listen to socket:", s)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGABRT, syscall.SIGHUP)
	go func() {
		sigerr := <-sig
		fmt.Println("Caught signal:", sigerr)
		if err = syscall.Close(s); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Closed:", s)

		//if err = syscall.Unlink(Fd); err != nil {
		//	log.Fatalln(err)
		//}
		//fmt.Println("Unlink:", Fd)
		os.Exit(1)
	}()

	var wg sync.WaitGroup
	for {
		client_sock, _, err := syscall.Accept(s)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Accepted:", client_sock)
		//go handle(client_sock, &wg)
		handle(client_sock, &wg)
	}
}

func handle(fd int, wg *sync.WaitGroup) {
	//defer wg.Done()
	defer func() {
		if err := syscall.Close(fd); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Closed:", fd)
	}()

	buf := make([]byte, 1024)
	for {
		// read from socket
		n, err := syscall.Read(fd, buf[:])
		if n <= 0 || err != nil {
			break
		}

		syscall.Write(fd, buf[:n])
		//fmt.Println("Socket:", fd)
		//fmt.Println("Accepted bytes:", n)
		fmt.Printf("[%d] >>> %s", fd, buf)
	}

	return
}