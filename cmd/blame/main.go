package main

import (
	"flag"
	"fmt"
	"os"

	blame "github.com/skanehira/github-blame"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	url, err := blame.GetPRURL(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(url)
}
