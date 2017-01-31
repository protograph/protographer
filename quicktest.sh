rm tex/dh.*
go run cmd/protographer/gen.go
mkdir -p tex && cd tex
pdflatex -halt-on-error dh.tex
open dh.pdf
cd ..
