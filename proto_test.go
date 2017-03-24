package protographer

import (
	"bytes"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

const backtickvar = "`aaa`"
const content = `
title: Hello World

actor:
- A: Alice
- B: Bob
- C: Carol

sequence:
- A: {B: hello}
- B: {C: world}
- C: {A: "!"}
- C: {A: "ha $ ` + backtickvar + ` $ haha $ bbb $ hahaha"}
`

func TestGeneratePGFUMLSD(t *testing.T) {
	data := ProtographYAML{}
	var b bytes.Buffer
	expected := `\end{document}`
	err := yaml.Unmarshal([]byte(content), &data)
	if err != nil {
		t.Error(err)
	}
	p := New(&data)
	p.GeneratePGFUMLSD(&b)
	t.Log(b.String())
	if !strings.HasSuffix(strings.TrimSpace(b.String()), expected) {
		t.Errorf("Generated output does not contain %s: Got %s", expected, b.String())
	}
}
