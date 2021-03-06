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
	Title      string                              //
	FootNote   string                              //
	Actor      []map[string]string                 //
	Sequence   []map[string]map[string]interface{} // Sequence is an array of this form (Src: { Dst: Label })
	Separation int                                 //
}

// Protograph is the parsed format of a protograph description.
type Protograph struct {
	Title      string              //
	FootNote   string              //
	Actor      []map[string]string //
	ActorList  []string            // abbreviated names, ordered.
	Sequence   []Sequence          //
	Separation int                 //
}

// Sequence denotes a source to destination and the annotations.
type Sequence struct {
	Source         string
	Destination    string
	Label          string
	AnnotationFrom string
	AnnotationTo   string
	Delay          int
	Color          string
	Style          string
}

// New creates a protogrph
func New(y *ProtographYAML) *Protograph {
	p := Protograph{}
	p.Title = y.Title
	p.FootNote = y.FootNote
	p.Actor = y.Actor
	p.ActorList = make([]string, len(y.Actor))
	p.Separation = y.Separation
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
					if s, ok := v3["arrow"].(string); ok {
						p.Sequence[i].Label = s
					}
					if s, ok := v3["from"].(string); ok {
						p.Sequence[i].AnnotationFrom = s
					}
					if s, ok := v3["to"].(string); ok {
						p.Sequence[i].AnnotationTo = s
					}
					if s, ok := v3["time"].(int); ok {
						p.Sequence[i].Delay = s
					}
					if s, ok := v3["color"].(string); ok {
						p.Sequence[i].Color = s
					}
					if s, ok := v3["style"].(string); ok {
						p.Sequence[i].Style = s
					}
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
