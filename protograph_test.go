package protographer

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestTransform(t *testing.T) {

	type testcase struct {
		input  string
		output string
	}

	tcs := []testcase{
		{"testdata/helloworld.md", "testdata/helloworld.tex"},
		{"testdata/dh.md", "testdata/dh.tex"},
	}

	for _, tc := range tcs {
		dmp := diffmatchpatch.New()

		i, err := os.Open(tc.input)
		if err != nil {
			t.Fatalf("could not read %s: %v", tc.input, err)
		}
		defer i.Close()

		o, err := os.Open(tc.output)
		if err != nil {
			t.Fatalf("could not read %s: %v", tc.output, err)
		}
		defer o.Close()

		outBuf, err := ioutil.ReadAll(o)
		if err != nil {
			t.Fatalf("could not read %s: %v", tc.output, err)
		}

		p := NewProtograph(bufio.NewReader(i))
		transformed := bytes.NewBuffer([]byte{})
		p.ToLaTeX(transformed)

		t.Logf("%#v\n", p)

		diffs := dmp.DiffMain(string(outBuf), transformed.String(), false)

		if dmp.DiffLevenshtein(diffs) != 0 {
			t.Errorf("transform from %s to %s failed", tc.input, tc.output)
			t.Error("differences are: ", dmp.DiffPrettyText(diffs))
			t.Log("expected: ", string(outBuf))
			t.Log("got:      ", transformed.String())
		}
	}
}
