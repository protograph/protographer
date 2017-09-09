package protographer

import (
	"regexp"
	"strings"
)

// sequence describes the flow
type sequence struct {
	Sender      string
	Receiver    string
	Style       style
	Description struct {
		From    string
		To      string
		Message string
	}
}

type sequences []sequence

type flowType int

const (
	invalidFlow flowType = iota
	from
	to
	message
	session
	thinks
)

func (p *protograph) addFlows(l string, f flowType) {
	if l == "" || len(p.Sequence) == 0 {
		return // don't know the actors, skip this.
	}
	seq := &p.Sequence[len(p.Sequence)-1]
	switch f {
	case thinks:
		seq.Description.From = strings.TrimSpace(seq.Description.From + "\n" + l)
	case session:
		fallthrough
	case message:
		seq.Description.Message = strings.TrimSpace(seq.Description.Message + "\n" + l)
	case from:
		seq.Description.From = strings.TrimSpace(seq.Description.From + "\n" + l)
	case to:
		seq.Description.To = strings.TrimSpace(seq.Description.To + "\n" + l)
	}
}

func (p *protograph) addFlowDir() {
	c := &p.Sequence[len(p.Sequence)-1]
	c.Style.L2R = dir(c.Sender, c.Receiver, p.ActorList)
}

func (s sequences) add(a, b string) sequences {
	if len(s) == 0 {
		s = append(s, sequence{Sender: a, Receiver: b})
		return s
	}
	c := &s[len(s)-1]

	switch {
	case b == "":
		// case of "A computes"
		s = append(s, sequence{Sender: a, Receiver: b})
	case c.Receiver == "" && a == c.Sender && b != "":
		// case of "A computes follow with A sends B"
		c.Receiver = b
	case c.Receiver != "" && c.Sender != "" && a != "" && b != "":
		// case of "A sends B" follows by "A sends B"
		s = append(s, sequence{Sender: a, Receiver: b})
	case a == "" && c.Receiver == b:
		// case of "Then B computes"
		fallthrough
	case a == "" && c.Receiver == "":
		// case of "Then B computes"
		c.Receiver = b
	case c.Sender != a || c.Receiver != b:
		// case that this is a new message
		s = append(s, sequence{Sender: a, Receiver: b})
	}
	return s
}

func getFlowType(l string) (flowType, string, string) {
	l = strings.TrimSpace(l)

	// match "A sends B"
	re := regexp.MustCompile("^([A-Za-z0-9_]+) sends ([A-Za-z0-9_]+)")
	matches := re.FindStringSubmatch(l)
	if len(matches) > 2 {
		return message, matches[1], matches[2]
	}

	// match "A and B exchange messages about"
	re = regexp.MustCompile("^([A-Za-z0-9_]+) and ([A-Za-z0-9_]+) exchange messages about")
	matches = re.FindStringSubmatch(l)
	if len(matches) > 2 {
		return session, matches[1], matches[2]
	}

	// match "A thinks"
	re = regexp.MustCompile("^([A-Za-z0-9_]+) thinks")
	matches = re.FindStringSubmatch(l)
	if len(matches) > 1 {
		return thinks, matches[1], matches[1]
	}

	// match "A computes"
	re = regexp.MustCompile("^([A-Za-z0-9_]+) computes")
	matches = re.FindStringSubmatch(l)
	if len(matches) > 1 {
		return from, matches[1], ""
	}

	// match "Then B computes"
	re = regexp.MustCompile("^Then ([A-Za-z0-9_]+) computes")
	matches = re.FindStringSubmatch(l)
	if len(matches) > 1 {
		return to, "", matches[1]
	}

	return invalidFlow, "", ""
}

func dir(a, b string, l [][2]string) bool {

	// case without message
	if a == b {
		if l[len(l)-1][0] != a {
			return true
		}
		// a is last element, so we put the dialog box on the right (right to left)
		return false
	}

	for _, s := range l {
		if a == s[0] {
			return true
		}
		if b == s[0] {
			return false
		}
	}
	return true
}