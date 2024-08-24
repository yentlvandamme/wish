package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func main () {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if _, err := os.Stat("./config.yaml"); os.IsNotExist(err) {
		_, err := os.OpenFile("./config.yaml", os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("API key: ")
		apiKey, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid input", err)
		}
		
		viper.Set("api_key", strings.TrimSpace(apiKey))
		if err := viper.WriteConfig(); err != nil {
			fmt.Println(err)
		}
	}
}