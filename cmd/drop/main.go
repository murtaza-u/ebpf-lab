package main

import (
	"log"
	"os"

	"github.com/murtaza-u/ebpf-lab/internal/cli"
)

func main() {
	err := cli.Cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
