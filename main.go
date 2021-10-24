package main

import (
	"fmt"
	"os"

	"github.com/dantdj/terraform-manager/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error occurred: %s", err)
		os.Exit(1)
	}
}
