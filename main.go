package main

import (
	"fmt"
	"os"

	"github.com/dantdj/terraform-manager/cmd"
	"github.com/dantdj/terraform-manager/config"
)

func main() {
	config.Load()
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error occurred: %s", err)
		os.Exit(1)
	}
}
