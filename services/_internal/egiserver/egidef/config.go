package egidef

import (
	"encoding/json"
	"os"
	"fmt"
	"io/ioutil"
	"flag"
	"log"
	"os/user"

	"github.com/EUDAT-GEF/GEF/gefserver/def"
)

func getHomeDir() string{
	usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
    return usr.HomeDir
}

type Configuration struct {
	Endpoint 				string `json:"endpoint"`
	ResourceTpl 			string `json:"resourceTpl"`
	OsTpl 					string `json:"osTpl"`
	PublicKey 				string `json:"publicKey"`
	Contextualisation 		string `json:"contextualisation"`
	ProxyPath 				string `json:"proxyPath"`
	Auth 					string `json:"auth"`
	Vo 						string `json:"vo"`
	VomsDir 				string `json:"vomsDir"`
	TrustedCertificatesPath string `json:"trustedCertificatesPath"`
	PublicIP 				string `json:"publicIP"`
	State 					string `json:"vmState"`
}

// ReadConfigFile reads a configuration file
func ReadConfigFile(configFilepath string) (Configuration, error) {
	var config Configuration

	file, err := os.Open(configFilepath)
	if err != nil {
		return config, Err(err, "Cannot open config file %s", configFilepath)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	
	err = decoder.Decode(&config)
	if err != nil {
		return config, Err(err, "Cannot read config file %s", configFilepath)
	}

	return config, nil
}

func ConfigGEFJson(ip string, configFilePath string) {

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
		CertPath: homedir+"/.docker/cert.pem",
		KeyPath: homedir+"/.docker/key.pem",
		CAPath: homedir+"/.docker/ca.pem",
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
