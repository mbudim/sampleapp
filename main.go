package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	consul "github.com/hashicorp/consul/api"
	vault "github.com/mch1307/vaultlib"
)

var VAULT_ADDR string
var VAULT_TOKEN string
var CONSUL_ADDR string

func handler(w http.ResponseWriter, r *http.Request) {

	//// GET CONFIG FROM CONSUL ////
	//var consulConfig *consul.Config
	//consulConfig.Address = CONSUL_ADDR
	conn, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		log.Fatalln(err)
	}

	consulData, _, err := conn.KV().Get("config/sampleapp", nil)
	if err != nil {
		log.Panicln(err)
	}

	config := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(consulData.Value)))

	fmt.Fprintf(w, "%s\n", "-------- FROM CONSUL --------")
	for scanner.Scan() {
		stringarray := strings.Split(scanner.Text(), "=")
		config[stringarray[0]] = stringarray[1]
		fmt.Fprintf(w, "key: %s value: %s\n", stringarray[0], config[stringarray[0]])
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
		log.Println("unset VAULT_ADDR")
	}

	VAULT_TOKEN = os.Getenv("VAULT_TOKEN")
	if VAULT_TOKEN == "" {
		panic("please set environment variable VAULT_TOKEN")
	}

	CONSUL_ADDR = os.Getenv("CONSUL_ADDR")
	if CONSUL_ADDR == "" {
		CONSUL_ADDR = "http://localhost:8500"
		log.Println("unset CONSUL_ADDR")
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
