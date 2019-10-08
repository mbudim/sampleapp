package main

import (
	"fmt"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	// viper.AddRemoteProvider("consul", "localhost:8500", "config/sampleapp")
	// viper.SetConfigType("yaml")
	// err := viper.ReadRemoteConfig()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(viper.Get("name"))
	// fmt.Println(viper.Get("address"))

	viper.SetConfigName("config")                                                                        // name of config file (without extension)
	viper.AddConfigPath("/Users/moka/Documents/budionduty/golang/src/github.com/mbudim/sampleapp/viper") // path to look for the config file in
	//viper.AddConfigPath("/Users/moka/Documents/budionduty/golang/src/github.com/mbudim/sampleapp/.viper") // call multiple times to add many search paths
	//viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// fmt.Println(viper.Get("name"))
	// fmt.Println(viper.Get("address.province"))
	data := viper.Get("")
	fmt.Println(data)
	fmt.Printf("%s", data)
}
