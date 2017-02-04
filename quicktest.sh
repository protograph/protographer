rm tex/dh.*
go run cmd/protographer/gen.go
mkdir -p tex && cd tex
pdflatex -halt-on-error helloworld.tex
open helloworld.pdf
cd ..
