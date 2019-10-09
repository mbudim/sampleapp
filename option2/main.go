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

var VAULT_ADDR string
var VAULT_TOKEN string
var CONSUL_ADDR string

func handler(w http.ResponseWriter, r *http.Request) {

	//// GET CONFIG FROM CONSUL ////
	viper.AddRemoteProvider("consul", CONSUL_ADDR, "config/sampleapp")
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
	vaultConfig.Address = VAULT_ADDR
	vaultConfig.Token = VAULT_TOKEN

	vaultConn, err := vault.NewClient(vaultConfig)
	if err != nil {
		panic(err)
	}

	secret, err := vaultConn.GetSecret("secret/sampleapp")
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

	VAULT_ADDR = os.Getenv("VAULT_ADDR")
	if VAULT_ADDR == "" {
		VAULT_ADDR = "http://localhost:8200"
		log.Println("unset VAULT_ADDR. It will be set to http://localhost:8200")
	} else {
		log.Println("VAULT_ADDR is set to ", VAULT_ADDR)
	}

	VAULT_TOKEN = os.Getenv("VAULT_TOKEN")
	if VAULT_TOKEN == "" {
		panic("please set environment variable VAULT_TOKEN")
	}

	CONSUL_ADDR = os.Getenv("CONSUL_ADDR")
	if CONSUL_ADDR == "" {
		CONSUL_ADDR = "http://localhost:8500"
		log.Println("unset CONSUL_ADDR. It will be set to http://localhost:8500")
	} else {
		log.Println("CONSUL_ADDR is set to ", CONSUL_ADDR)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
