package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"updatedb/download"
)

func main() {

	var token string
	var name string
	var dstPath string

	pflag.StringVar(&dstPath, "path", "", "--path=/path/dir/")
	pflag.StringVar(&name, "name", "", "--name")
	pflag.StringVar(&token, "token", "", "--token=")
	pflag.Parse()

	e := download.Custom(token, name, dstPath)
	if e != nil {
		fmt.Println()
		fmt.Println(e)
		fmt.Println()
	} else {
		fmt.Println()
		fmt.Println("ok")
		fmt.Println()
	}
}
