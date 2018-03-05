package sshegi

import (
	"errors"
	"fmt"
	"path/filepath"
	"log"
	"net"
	"os"
	"sync"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const (
	Username    = "egieudat"
	DefaultPort = 22
)

var Ch chan string = make(chan string)

func KeyScanCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	Ch <- fmt.Sprintf("%s %s", hostname[:len(hostname)-3], string(ssh.MarshalAuthorizedKey(key)))
	return nil
}

func dial(server string, config *ssh.ClientConfig, wg *sync.WaitGroup) {
	_ , err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server, DefaultPort), config)
	if err != nil {
		//client.Close()
		log.Fatalln("Failed to dial:", err)
	}
	wg.Done()

}

func out(wg *sync.WaitGroup) {
	for s := range Ch {
		fmt.Printf("%s", s)

		f, err := os.OpenFile(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"),os.O_APPEND|os.O_WRONLY, 0600)

		if err != nil {
		    log.Fatal(err)
		}	

		defer f.Close()

		if _, err = f.WriteString(s); err != nil {
		    panic(err)
		}
		wg.Done()
	}
}

func KeyScanVM(server string) {

	auth_socket := os.Getenv("SSH_AUTH_SOCK")
	if auth_socket == "" {
		log.Fatal(errors.New("no $SSH_AUTH_SOCK defined"))
	}

	conn, err := net.Dial("unix", auth_socket)
	if err != nil {
		log.Fatal(err)
	}
	
	defer conn.Close()
	ag := agent.NewClient(conn)
	auths := []ssh.AuthMethod{ssh.PublicKeysCallback(ag.Signers)}

	config := &ssh.ClientConfig{
		User:            Username,
		Auth:            auths,
		HostKeyCallback: KeyScanCallback,
	}

	var wg sync.WaitGroup

	go out(&wg)
	wg.Add(2)                       
	go dial(server, config, &wg)
	wg.Done()
	
	wg.Wait() 
}