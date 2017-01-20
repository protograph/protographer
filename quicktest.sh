go run cmd/protographer/gen.go
cd tex
pdflatex -halt-on-error oauth2.authcode.tex
pdflatex -halt-on-error dh.tex
open oauth2.authcode.pdf
open dh.pdf
cd ..
