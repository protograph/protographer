all: clean dh u2f

dh:
	go run cmd/protographer/gen.go
	pdflatex -halt-on-error tex/dh.tex
	open dh.pdf
	cd ..

u2f:
	go run cmd/protographer/gen.go
	pdflatex -halt-on-error tex/u2f.tex
	open u2f.pdf
	cd ..

test: dh

clean:
	rm -f tex/dh.* dh.*
	rm -f tex/u2f.* u2f.*
