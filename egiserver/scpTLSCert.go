package main

import (
	"fmt"
	"log"
	"path/filepath"
	"bufio"
	"strings"
	"os/exec"
	"os"
	"io/ioutil"
	"golang.org/x/crypto/ssh"
	"github.com/EUDAT-GEF/GEF/egiserver/egidef"
	"os/user"
)

type ScpFile struct {
	FileToCopy 		[]string
	RemoteFolder 	[]string
}


func PublicKeyFile(file string) ssh.AuthMethod {

	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

// Check if the EGI server/VM fingerprint is knownn is ~/.ssh/known_hosts
func CheckKnownHosts(server string) ssh.PublicKey {

	// Open and scan known_hosts
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey

	// Start to scan the file
	for scanner.Scan() {

		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}

		if fields[0]==server {

			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		} 
	}

	if hostKey == nil {
		log.Fatalf("no hostkey for %s", server)
	}
	return hostKey
}

// Connection to the EGI host through ssh
func connectToHost(user string, host string, PublicKeyPath string) (*ssh.Client, *ssh.Session, error) {	

	hostKey := CheckKnownHosts(host)

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(PublicKeyPath)},

		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	port := ":22"
	host = host + port

	client, err := ssh.Dial("tcp", host, sshConfig)
	
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}

func getUserEnv() string {
	usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
    return usr.HomeDir
}

func fromHostToVM(scpFile ScpFile) {

	for i:=0; i<len(scpFile.FileToCopy); i++ {
		fmt.Println("i",scpFile.FileToCopy[i])
	
		fileToCopy, err := ioutil.ReadFile(scpFile.FileToCopy[i])
		if err != nil {
				log.Fatalln(err)
		}
		fileToCopyString := string(fileToCopy)
		fmt.Println("fileToCopy",fileToCopyString)
	}
}

func getWorkingDir() string {
	currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
	}
	return currentDir
}

func main() {

	
	currentDir := getWorkingDir()
    homeDir := getUserEnv()

	configVM, err := egidef.ReadConfigFile(currentDir+"/"+"configEGI.json")
    if err != nil {
        panic(err)
    }

    server := configVM.PublicIP
    PublicKeyPath := homeDir+"/.ssh/fedcloudNoPass"
	name:="egieudat"

	cmd, err := exec.Command("/bin/sh", "gen-docker-certs.sh",server).Output()
  	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cmd)

	client, session, err := connectToHost(name, server, PublicKeyPath)
	if err != nil {
		panic(err)
	}

	//scpFile.FileToCopy = []string{"server.pem","server-key.pem","ca.pem"}
	//scpFile.RemoteFolder = []string{"docker","docker","docker"}

	defer session.Close()
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		fileToCopy, err := ioutil.ReadFile("dockerServerConfig.sh")
		if err != nil {
			log.Fatalln(err)
		}

		fileToCopyString := string(fileToCopy)
		fmt.Fprintln(w, "C0644", len(fileToCopyString), "dockerServerConfig.sh")
		fmt.Fprint(w, fileToCopyString)
		fmt.Fprint(w, "\x00")

		fileToCopy, err = ioutil.ReadFile(homeDir+"/.docker/server.pem")
		if err != nil {
				log.Fatalln(err)
		}

		fileToCopyString = string(fileToCopy)
		fmt.Fprintln(w, "D0755",0, ".docker") 
		fmt.Fprintln(w, "C0644", len(fileToCopyString), "server.pem")
		fmt.Fprint(w, fileToCopyString)
		fmt.Fprint(w, "\x00")

		fileToCopy, err = ioutil.ReadFile(homeDir+"/.docker/server-key.pem")
		if err != nil {
				log.Fatalln(err)
		}

		fileToCopyString = string(fileToCopy)
		fmt.Fprintln(w, "C0644", len(fileToCopyString), "server-key.pem")
		fmt.Fprint(w, fileToCopyString)
		fmt.Fprint(w, "\x00")
	}()
	
	if err := session.Run("/usr/bin/scp -tr ./"); err != nil {
		panic("Failed to run: " + err.Error())
	}

	client.Close()
}