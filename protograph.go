package protographer

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

// protograph is the parsed format of a protograph description.
type protograph struct {
	Title       string
	Description string
	ActorList   [][2]string // list of actor (acronym and full name)
	Sequence    sequences
	Note        string
	Style       style
	Expressions [][2]string
}

type state int

const (
	start state = iota
	title
	config
	actors
	flows
	notes
	expressions
	none
)

type style struct {
	Delay  int
	Color  string
	Format string
	L2R    bool // left to right
	Sep    int
	PadTop int
	Line   string
}

func NewProtograph(r io.Reader) *protograph {
	p := protograph{Style:style{Sep:6}}

	s := bufio.NewScanner(r)
	t := start
	f := invalidFlow
	for s.Scan() {
		l := s.Text()
		switch {
		case strings.HasPrefix(l, "# ") && t == start:
			p.Title = strings.TrimPrefix(l, "# ")
			t = title
		case strings.HasPrefix(l, "## "):
			// this is either style, actors, flows, footnote
			section := normalize(l, "## ")
			switch section {
			case "config":
				t = config
			case "actors":
				t = actors
			case "flows":
				t = flows
			case "notes":
				t = notes
			case "expressions":
				t = expressions
			default:
				t = none
			}
		case strings.HasPrefix(l, "### ") && t == flows:
			var a, b string
			f, a, b = getFlowType(strings.TrimPrefix(l, "### "))
			if f == invalidFlow {
				break
			}
			p.Sequence = p.Sequence.add(a, b)
			p.addFlowStyle(l, f)
			p.addFlowDir()
		default:
			switch t {
			case title:
				p.Description = strings.TrimSpace(p.Description + "\n" + l)
			case config:
				p.addConfig(l)
			case actors:
				p.addActors(l)
			case flows:
				p.addFlows(expandLine(l, p.Expressions), f)
			case notes:
				p.addNotes(expandLine(l, nil))
			case expressions:
				p.addExpressions(l)
			}
		}
	}
	return &p
}

func (p *protograph) addConfig(l string) {
	l = strings.TrimSpace(l)
	if !strings.Contains(l, ":"){
		return
	}
	pcs := strings.SplitN(l, ":", 2)
	a := strings.TrimSpace(pcs[0])
	b := strings.TrimSpace(pcs[1])
	switch a {
	case "separation":
		p.Style.Sep, _ = strconv.Atoi(b)
	}
}

func (p *protograph) addExpressions(l string) {
	l = strings.TrimSpace(l)
	if !strings.Contains(l, ":"){
		return
	}
	pcs := strings.SplitN(l, ":", 2)
	a := strings.TrimSpace(pcs[0])
	b := strings.TrimSpace(pcs[1])
	p.Expressions = append(p.Expressions, [2]string{a, b})
}

func (p *protograph) addActors(l string) {
	l = strings.TrimSpace(l)
	if l == "" {
		return
	}
	pcs := strings.SplitN(l, ":", 2)
	a := strings.TrimSpace(pcs[0])
	b := a
	if len(pcs) >= 2 {
		b = strings.TrimSpace(pcs[1])
	}
	p.ActorList = append(p.ActorList, [2]string{a, b})
}

func (p *protograph) addNotes(l string) {
	p.Note = strings.TrimSpace(p.Note + "\n" + l)
}

func (p *protograph) addFlowStyle(l string, f flowType) {
	// default style
	c := &p.Sequence[len(p.Sequence)-1]
	if c.Style.Color == "" {
		c.Style.Color = "black"
	}
	if c.Style.Format == "" {
		c.Style.Format = "->"

	}

	if f == session {
		c.Style.Format = "<->"
	}

	if f == thinks {
		// remove the arrow by setting the color to white
		c.Style.Color = "white"
	}

	l = strings.ToLower(l)
	// extract styles in the form of (k1=v1, k2=v2) from the input
	re := regexp.MustCompile("[(](.*)[)]")
	matches := re.FindStringSubmatch(l)
	if len(matches) < 2 {
		return
	}
	for _, s := range strings.Split(matches[1], ",") {
		pcs := strings.SplitN(s, "=", 2)
		if len(pcs) < 2 {
			continue
		}
		pcs[0] = strings.TrimSpace(pcs[0])
		pcs[1] = strings.TrimSpace(pcs[1])
		switch pcs[0] {
		case "delay":
			c.Style.Delay, _ = strconv.Atoi(pcs[1])
		case "color":
			c.Style.Color = pcs[1]
		case "padtop":
			c.Style.PadTop, _ = strconv.Atoi(pcs[1])
		case "line":
			c.Style.Line = pcs[1]
		}
	}
}

func (p *protograph) ToLaTeX(w io.Writer) {
	funcMap := template.FuncMap{
		"l2r": l2r,
		"expand": expand,
		"padtop": padtop,
	}
	err := template.Must(template.New("pgfumlsd").
		Funcs(funcMap).Delims("<<", ">>").
		Parse(theTemplate)).Execute(w, p)
	if err != nil {
		fmt.Println("err converting protograph to latex: ", err.Error())
	}
}
