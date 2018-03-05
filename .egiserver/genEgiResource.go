package main

import (
    "fmt"
    "time"

    "github.com/EUDAT-GEF/GEF/egiserver/sshegi"
    "github.com/EUDAT-GEF/GEF/egiserver/egidef"
    "github.com/EUDAT-GEF/GEF/egiserver/runcontainer"

)

func main() {

    fmt.Println("Start container to create VM")

    client, contID, sourceList := runcontainer.RunVmContainer()
    state := runcontainer.CheckContState(&client, contID)

    for state != "exited" {
        state = runcontainer.CheckContState(&client, contID)
        fmt.Println("state: ",state)
        time.Sleep(10 * time.Second)
    }

    fmt.Println("The Virtual Machine is now active on EGI")

    configVM, err := egidef.ReadConfigFile(sourceList)
    if err != nil {
        panic(err)
    }

    gefConfigFilePath := "../gefserver/config.json"

    egidef.ConfigGEFJson(configVM.PublicIP, gefConfigFilePath)

    EGIserver := configVM.PublicIP
    sshegi.KeyScanVM(EGIserver)

}

