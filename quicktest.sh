rm tex/.*
go run cmd/protographer/gen.go yaml/helloworld.yaml
go run cmd/protographer/gen.go yaml/dh.yaml
mkdir -p tex && cd tex
pdflatex -halt-on-error helloworld.tex
pdflatex -halt-on-error dh.tex
open *.pdf
cd ..
