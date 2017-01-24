rm tex/dh.*
go run cmd/protographer/gen.go
cd tex
pdflatex -halt-on-error dh.tex
open dh.pdf
cd ..
