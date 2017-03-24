package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/protograph/protographer"
	"gopkg.in/yaml.v2"
	"path/filepath"
)

var (
	outDir = flag.String("outdir", "tex", "output directory")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [OPTIONS] YAML-FILENAME:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	filePath := flag.Arg(0)
	if filePath == "" {
		flag.Usage()
		os.Exit(1)
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("error opening %s: %s", filePath, err))
	}
	err = genTex(filePath, content)
	if err != nil {
		panic(fmt.Errorf("error generating tex %s: %s", filePath, err))
	}
}

func genTex(fileName string, content []byte) error {
	fName := filepath.Base(fileName)
	base := fName[:len(fName)-len(filepath.Ext(fName))]
	outFile := *outDir + "/" + base + ".tex"
	data := protographer.ProtographYAML{}
	err := yaml.Unmarshal([]byte(content), &data)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	p := protographer.New(&data)
	fmt.Printf("Converting %s to %s with PGF-UMLSD.\n", fileName, outFile)
	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()
	p.GeneratePGFUMLSD(out)
	return nil
}
