package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	vault "github.com/mch1307/vaultlib"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var VAULT_HOST string
var VAULT_PORT string
var VAULT_TOKEN string
var VAULT_PATH string
var CONSUL_HOST string
var CONSUL_PORT string
var CONSUL_PATH string

func handler(w http.ResponseWriter, r *http.Request) {

	//// GET CONFIG FROM CONSUL ////
	viper.AddRemoteProvider("consul", CONSUL_HOST+":"+CONSUL_PORT, CONSUL_PATH)
	viper.SetConfigType("properties")
	err := viper.ReadRemoteConfig()
	if err != nil {
		fmt.Println(err)
	}
	config := viper.AllSettings()
	fmt.Fprintf(w, "%s\n", "-------- FROM CONSUL --------")
	for key, val := range config {
		fmt.Fprintf(w, "key: %s value: %s\n", key, val)
	}
	fmt.Fprintf(w, "\n")

	//// GET CONFIG FROM VAULT ////
	vaultConfig := vault.NewConfig()
	vaultConfig.Address = "http://" + VAULT_HOST + ":" + VAULT_PORT
	vaultConfig.Token = VAULT_TOKEN

	vaultConn, err := vault.NewClient(vaultConfig)
	if err != nil {
		panic(err)
	}

	secret, err := vaultConn.GetSecret(VAULT_PATH)
	if err != nil {
		panic(err)
	}

	mapString := secret.KV
	fmt.Fprintf(w, "%s\n", "-------- FROM VAULT --------")
	for key, val := range mapString {
		fmt.Fprintf(w, "key: %s value: %s\n", key, val)
	}
	fmt.Fprintf(w, "\n")

	//// OVERRIDE CONSUL CONFIG BY VAULT ////
	for key, val := range mapString {
		if config[key] != "" {
			config[key] = val
		} else {
			config[key] = mapString[key]
		}
	}

	fmt.Fprintf(w, "%s\n", "-------- FINAL CONFIG --------")
	for key, val := range config {
		fmt.Fprintf(w, "key: %s value: %s\n", key, val)
	}

}

func init() {
	log.SetPrefix("sampleapp: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	//log.Println("init started")
}

func main() {

	VAULT_HOST = os.Getenv("VAULT_HOST")
	if VAULT_HOST == "" {
		VAULT_HOST = "localhost"
		log.Println("VAULT_HOST is not set. It will be set to localhost")
	} else {
		log.Println("VAULT_HOST is set to ", VAULT_HOST)
	}

	VAULT_PORT = os.Getenv("VAULT_PORT")
	if VAULT_PORT == "" {
		VAULT_PORT = "8200"
		log.Println("VAULT_PORT is not set. It will be set to 8200")
	} else {
		log.Println("VAULT_PORT is set to ", VAULT_PORT)
	}

	VAULT_TOKEN = os.Getenv("VAULT_TOKEN")
	if VAULT_TOKEN == "" {
		panic("please set environment variable VAULT_TOKEN")
	}

	VAULT_PATH = os.Getenv("VAULT_PATH")
	if VAULT_PATH == "" {
		VAULT_PATH = "secret/sampleapp"
		log.Println("VAULT_PATH is not set. It will be set to ", VAULT_PATH)
	} else {
		log.Println("VAULT_PATH is set to", VAULT_PATH)
	}

	CONSUL_HOST = os.Getenv("CONSUL_HOST")
	if CONSUL_HOST == "" {
		CONSUL_HOST = "localhost"
		log.Println("CONSUL_HOST is not set. It will be set to localhost")
	} else {
		log.Println("CONSUL_HOST is set to ", CONSUL_HOST)
	}

	CONSUL_PORT = os.Getenv("CONSUL_PORT")
	if CONSUL_PORT == "" {
		log.Println("CONSUL_PORT is not set. It will be set to 8500")
		CONSUL_PORT = "8500"
	} else {
		log.Println("CONSUL_PORT is set to ", CONSUL_PORT)
	}

	CONSUL_PATH = os.Getenv("CONSUL_PATH")
	if CONSUL_PATH == "" {
		CONSUL_PATH = "config/sampleapp"
		log.Println("CONSUL_PATH is not set. It will be set to http://localhost:8500")
	} else {
		log.Println("CONSUL_PATH is set to ", CONSUL_PATH)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
