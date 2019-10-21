package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	vault "github.com/mch1307/vaultlib"
	"github.com/spf13/viper"
)

var VAULT_HOST string
var VAULT_PORT string
var VAULT_TOKEN string
var VAULT_PATH string

func handler(w http.ResponseWriter, r *http.Request) {

	//// GET CONFIG FROM FILE ////
	viper.SetConfigName("config")       // name of config file (without extension)
	viper.AddConfigPath("/app/config/") // path to look for the config file in
	viper.AddConfigPath("config/")      // path to look for the config file in
	viper.AddConfigPath(".")            // path to look for the config file in
	err := viper.ReadInConfig()         // Find and read the config file
	if err != nil {                     // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	config := viper.AllSettings()

	fmt.Fprintf(w, "%s\n", "-------- FROM FILE --------")
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

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
