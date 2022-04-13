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

	pflag.StringVar(&dstPath, "path", "d:/tmp", "--path=/path/dir/")
	pflag.StringVar(&name, "name", "", "--name")
	pflag.StringVar(&token, "token", "4b53024029758e53e3ab119f956ef41ae66acc6a", "--token=")
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
