package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/protograph/protographer"
)

var (
	outFile = flag.String("out", "", "output file")
	inFile  = flag.String("in", "", "input file")
)

func main() {
	flag.Parse()

	if *inFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *outFile == "" {
		ext := path.Ext(*inFile)
		*outFile = (*inFile)[0:len(*inFile)-len(ext)] + ".tex"
	}

	i, err := os.Open(*inFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	o, err := os.OpenFile(*outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	p := protographer.NewProtograph(bufio.NewReader(i))
	p.ToLaTeX(o)

	fmt.Printf("transformed %s to %s\n", *inFile, *outFile)
}
