package protographer

import (
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/protograph/protographer/templates/pgfumlsd"
)

// ProtographYAML is the YAML format of a protograph description.
type ProtographYAML struct {
	Title    string                              //
	Actor    []map[string]string                 //
	Sequence []map[string]map[string]interface{} // Sequence is an array of this form (Src: { Dst: Label })
}

// Protograph is the parsed format of a protograph description.
type Protograph struct {
	Title     string              //
	Actor     []map[string]string //
	ActorList []string            // abbreviated names, ordered.
	Sequence  []Sequence
}

// Sequence denotes a source to destination and the annotations.
type Sequence struct {
	Source         string
	Destination    string
	Label          string
	AnnotationFrom string
	AnnotationTo   string
	Delay          int
}

// New creates a protogrph
func New(y *ProtographYAML) *Protograph {

	p := Protograph{}

	p.Title = y.Title
	p.Actor = y.Actor
	p.ActorList = make([]string, len(y.Actor))
	for i, a := range y.Actor {
		for k := range a {
			p.ActorList[i] = k
			break
		}
	}

	p.Sequence = make([]Sequence, len(y.Sequence))

	for i, s := range y.Sequence {

		for k, v := range s {
			p.Sequence[i].Source = k
			for k2, v2 := range v {
				p.Sequence[i].Destination = k2
				switch v3 := v2.(type) {
				case map[interface{}]interface{}:
					fmt.Println(v3["textFrom"])
					p.Sequence[i].Label = v3["label"].(string)
					p.Sequence[i].AnnotationFrom = v3["textFrom"].(string)
					p.Sequence[i].AnnotationTo = v3["textTo"].(string)
				case string:
					p.Sequence[i].Label = v3
				case int:
					p.Sequence[i].Label = strconv.Itoa(v3)
				default:
					fmt.Println("unknown type")
					fmt.Println(reflect.TypeOf(v3))
				}

				break
			}
			break
		}

	}

	return &p
}

// GeneratePGFUMLSD generates the pgf umlsd output
func (p *Protograph) GeneratePGFUMLSD(wr io.Writer) {

	pgfumlsd.GetTemplate().Execute(wr, p)

}
