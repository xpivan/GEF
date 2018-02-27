package runcontainer

import (
	"fmt"
	"os/user"
	"os"
	"log"
	"bytes"
	"context"
	"time"
	"github.com/fsouza/go-dockerclient"
	"github.com/EUDAT-GEF/GEF/gefserver/def"
)

var configFilePath = "config.json"
var EGIconfig = "configEGI.json"

func CheckContState(c *docker.Client, contId string) string {
	cont, err := c.InspectContainer(contId)
	if err != nil {
		panic(err)
	}
	return cont.State.Status
}

func getWorkingDir() string {
	currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
	}
	return currentDir
}

func getHomeDir() string{
	usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
    return usr.HomeDir
}

func RunVmContainer() (docker.Client, string, string) {

	config, err := def.ReadConfigFile(configFilePath)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}

	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}

	var stdout bytes.Buffer
	var ContConf docker.Config
	var CreateContOpts docker.CreateContainerOptions
	var hostConf docker.HostConfig
	var optsAttach docker.AttachToContainerOptions

	// get the current working directory
	currentDir := getWorkingDir()

	// get the Home directory
	usrHome := getHomeDir()


	//Set up the path for the output
	TargetList := make([]string,3)

	TargetList[0] = "/home/configEGI.json"
	TargetList[1] = "/root/.ssh"
	TargetList[2] = "/root/.globus"

	// Set up the path of the folder or file to be mounted within the container
	SourceList := make([]string,3)

	SourceList[0] = currentDir+"/"+EGIconfig
	SourceList[1] = usrHome+"/.ssh"
	SourceList[2] = usrHome+"/.globus"

	readOnly := make([]bool,3)

	readOnly[0] = false
	readOnly[1] = true
	readOnly[2] = true

	// Fill the docker config with the mounted file
	VolsMount := make([]docker.HostMount,3)

	vtm := new(docker.HostMount)

	for i := range VolsMount {
		vtm.Target = TargetList[i]
		vtm.Source = SourceList[i]
		vtm.Type = "bind"
		vtm.ReadOnly = readOnly[i]
		VolsMount[i] = *vtm
	}

	ContConf.Image = "egi"
	ContConf.AttachStdout = true
	ContConf.AttachStderr = true
	ContConf.OpenStdin = true

	hostConf.Mounts = VolsMount

	createContainerContext, cancel := context.WithTimeout(context.Background(), time.Duration(config.Timeouts.Preparation)*time.Second)
	defer cancel()

	CreateContOpts.Config = &ContConf
	CreateContOpts.HostConfig = &hostConf
	CreateContOpts.Context = createContainerContext

	var contId *docker.Container 

	contId , err = client.CreateContainer(CreateContOpts)
	if err != nil {
		panic(err)
	}

	optsAttach.Container = contId.ID

	attached := make(chan struct{})
	go func() {
		client.AttachToContainer(docker.AttachToContainerOptions{
			Container:    contId.ID,
			OutputStream: &stdout,
			ErrorStream:  &stdout,
			Logs:         true,
			Stdin:        true,
			Stdout:       true,
			Stderr:       true,
			Stream:       true,
			RawTerminal:  true,
			Success:      attached,
		})
	}()

	<-attached
	attached <- struct{}{}

	jobExecutionContext, cancel := context.WithTimeout(context.Background(), time.Duration(config.Timeouts.JobExecution)*time.Second)
	defer cancel()

	err = client.StartContainerWithContext(contId.ID, &hostConf, jobExecutionContext)
	if err != nil {
		fmt.Println("err: ",err)
	}
	return *client, contId.ID, SourceList[0]
}





