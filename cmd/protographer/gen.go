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

//type Protograph struct {
//	Title    string
//	Actor    []map[string]string
//	Sequence []map[string]map[string]interface{} // Sequence is an array of this form (Src: { Dst: Label })
//}
//
//func (p *Protograph) GenerateMermaid() {
//	tmpl, _ := template.New("mermaid").Parse(mermaidTemplate)
//	tmpl.Execute(os.Stdout, p)
//}
//
//func (p *Protograph) GeneratePGFUMLSD(wr io.Writer) {
//	tmpl, _ := template.New("pgfumlsd").Delims("((", "))").
//		Parse(pgfumlsdTemplate)
//	tmpl.Execute(wr, p)
//}

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

		//fmt.Println("===MermaidJS===")
		//protograph.GenerateMermaid()

		fmt.Printf("PGF-UMLSD: %s\n", file.Name())
		out, _ := os.Create(*outDir + "/" + baseName + ".tex")
		defer out.Close()

		p.GeneratePGFUMLSD(out)
	}
}
