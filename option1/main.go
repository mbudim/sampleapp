package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	vault "github.com/mch1307/vaultlib"
	"github.com/spf13/viper"
)

var VAULT_ADDR string
var VAULT_TOKEN string

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
	}

	VAULT_TOKEN = os.Getenv("VAULT_TOKEN")
	if VAULT_TOKEN == "" {
		panic("please set environment variable VAULT_TOKEN")
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
