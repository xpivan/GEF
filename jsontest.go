package main

import "encoding/json"
import "fmt"
import "github.com/EUDAT-GEF/GEF/gefserver/def"
import "flag"
import "log"
import "io/ioutil"
import "os/user"

var configFilePath = "gefserver/config.json"

func getHomeDir() string{
	usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
    return usr.HomeDir
}

var ip = "123.123.123.123"

func configGEFJson(ip string, configFilePath string) {
	flag.StringVar(&configFilePath, "config", configFilePath, "configuration file")
	flag.Parse()

	config, err := def.ReadConfigFile(configFilePath)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}

	homedir := getHomeDir()

	endpoint := "tcp://"+ip+":2376"
	//certpath := homedir

	//var dockerConfig def.DockerConfig
	config.Docker = def.DockerConfig{
		Endpoint: endpoint,
		TLSVerify: true,
		CertPath: homedir+"/cert.pem",
		KeyPath: homedir+"/key.pem",
		CAPath: homedir+"/ca.pem",
	}

	output, err := json.MarshalIndent(&config, "", "\t\t")

	if err != nil {
	  fmt.Println("Error marshalling to JSON:", err)
	  return
	}
	err = ioutil.WriteFile(configFilePath, output, 0644)
	if err != nil {
	  fmt.Println("Error writing JSON to file:", err)
	  return
	}
}

func main() {
	configGEFJson(ip, configFilePath)
}