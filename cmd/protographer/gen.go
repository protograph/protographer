package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/protograph/protographer"
	"gopkg.in/yaml.v2"
)

var (
	inDir  = flag.String("in", "yaml", "input directory")
	outDir = flag.String("out", "tex", "output directory")
)

func main() {

	flag.Parse()

	files, _ := ioutil.ReadDir(*inDir)

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		pcs := strings.Split(file.Name(), ".")
		baseName := strings.Join(pcs[:len(pcs)-1], ".")

		content, _ := ioutil.ReadFile(*inDir + "/" + file.Name())

		data := protographer.ProtographYAML{}

		err := yaml.Unmarshal([]byte(content), &data)
		if err != nil {
			log.Printf("error: %v\n", err)
		}

		p := protographer.New(&data)

		fmt.Printf("Converting %s to TeX with PGF-UMLSD.\n", file.Name())
		out, _ := os.Create(*outDir + "/" + baseName + ".tex")
		defer out.Close()

		p.GeneratePGFUMLSD(out)
	}
}
