package main

import (
	"fmt"

	"github.com/volcente/gator/internal/config"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Error! %v\n", err)
		return
	}
	fmt.Printf("%v\n", cfg)

	err = cfg.SetUser("bartosz")
	if err != nil {
		fmt.Printf("Error! %v\n", err)
		return
	}

	cfg, err = config.GetConfig()
	if err != nil {
		fmt.Printf("Error! %v\n", err)
		return
	}
	fmt.Printf("%v\n", cfg)

}
