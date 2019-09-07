package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var IP = [4]byte{127, 0, 0, 1}

const (
	//Fd = "myown" // AF_UNIX
	Port = 3000 // AF_INET | AF_INET6
)

type Server struct {
	serverSocket  int
	clientSockets map[int]int
}

func (s *Server) setup(addr *syscall.SockaddrInet4) error {
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return err
	}

	if err = syscall.SetsockoptInt(sock, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return err
	}
	log.Println("Set opt: SO_REUSEADDR")

	linger := syscall.Linger{
		Onoff:  1,
		Linger: 1,
	}
	if err = syscall.SetsockoptLinger(sock, syscall.SOL_SOCKET, syscall.SO_LINGER, &linger); err != nil {
		log.Fatalln(err)
	}

	//if err = syscall.SetNonblock(s, true); err != nil {
	//	log.Fatalln(err)
	//}

	if err = syscall.Bind(sock, addr); err != nil {
		return err
	}
	log.Printf("Binding to %d\n", sock)

	s.serverSocket = sock
	s.clientSockets = make(map[int]int)
	return nil
}

func (s *Server) listen(backlog int) error {
	if err := syscall.Listen(s.serverSocket, backlog); err != nil {
		return err
	}
	log.Println("Listen to socket:", s)
	return nil
}

func (s *Server) addMemberToGroup(clientSock int) error {
	if _, ok := s.clientSockets[clientSock]; ok {
		return fmt.Errorf("client already exist in group")
	}

	s.clientSockets[clientSock] = 1
	return nil
}

func (s *Server) dropMemberFromGroup(clientSock int) error {
	if _, ok := s.clientSockets[clientSock]; !ok {
		return fmt.Errorf("client is not in group")
	}

	delete(s.clientSockets, clientSock)
	return nil
}

func (s *Server) handle(clientSock int) {
	defer func() {
		if err := syscall.Close(clientSock); err != nil {
			log.Fatalln(err)
		}
		log.Println("Closed:", clientSock)
	}()

	buf := make([]byte, 1024)
	for {
		// append sender info
		buf = append(buf, byte(clientSock))

		// read from socket
		n, err := syscall.Read(clientSock, buf)
		if n <= 0 || err != nil {
			break
		}

		s.sendAll(buf[:n], clientSock)

		//log.Println("Socket:", fd)
		log.Println("Accepted bytes:", n)
		log.Printf("[%d] >>> %s", clientSock, buf)

		// set zeros only on non-zero bytes
		copy(buf, make([]byte, n))
	}

	return
}

func (s *Server) sendAll(msg []byte, sender int) {
	for sock := range s.clientSockets {
		// skip client who sent this message
		if sock == sender {
			continue
		}

		if _, err := syscall.Write(sock, msg); err != nil {
			log.Println("[WARNING]", "failed to send message to", sock)
			continue
		}
	}
}

func (s *Server) teardown() error {
	if err := syscall.Close(s.serverSocket); err != nil {
		return err
	}
	log.Println("Closed:", s)

	//if err = syscall.Unlink(Fd); err != nil {
	//	return err
	//}
	//log.Println("Unlink:", Fd)
	return nil
}

func main() {
	s := Server{}
	if err := s.setup(&syscall.SockaddrInet4{Addr: IP, Port: Port}); err != nil {
		log.Fatalln("failed to setup", err)
	}
	log.Printf("Server initiated %+v\n", s)

	if err := s.listen(5); err != nil {
		log.Fatalln("failed to listen", err)
	}
	log.Printf("Listening to port %d", Port)

	// trap signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGABRT, syscall.SIGHUP)
	go func() {
		sigcall := <-sig
		log.Println("Caught signal:", sigcall)
		if err := s.teardown(); err != nil {
			log.Println(err)
		}

		os.Exit(1)
	}()

	// accept new clients
	for {
		clientSocket, _, err := syscall.Accept(s.serverSocket)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("New client:", clientSocket)

		if err = s.addMemberToGroup(clientSocket); err != nil {
			log.Println("[WARNING]", "failed to add", clientSocket, "to group")
		}
		go func() {
			s.handle(clientSocket)
			if err = s.dropMemberFromGroup(clientSocket); err != nil {
				log.Println("[WARNING]", "failed to drop", clientSocket, "from group")
			}
		}()
	}
}
