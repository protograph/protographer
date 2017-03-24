package pgfumlsd

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

func expand(a string) string {
	a = strings.TrimSpace(a)
	// expand EOL to \n and shortstack
	if strings.Contains(a, "\n") {
		a = fmt.Sprintf("\\shortstack[l]{%s}", strings.Replace(a, "\n", "\\\\", -1))
	}

	// expand $ `...` $ to $ \textsf{...}
	re := regexp.MustCompile("[$][^$]*[$]")
	a = re.ReplaceAllStringFunc(a, func(s string) string {
		re2 := regexp.MustCompile("[`]([^`]*)[`]")
		s = re2.ReplaceAllString(s, "\\mathsf{$1}")
		s = strings.Replace(s, "<-", `\gets`, -1)
		s = strings.Replace(s, "->", `\to`, -1)
		s = strings.Replace(s, "||", `\Vert `, -1)
		return s
	})
	return a
}

func anchor(list []string, src, dst, fromOrTo string) string {
	for _, v := range list {
		if v == src {
			switch fromOrTo {
			case "from":
				return "east"
			default:
				return "west"
			}
		}
		if v == dst {
			switch fromOrTo {
			case "from":
				return "west"
			default:
				return "east"
			}
		}
	}
	return "east"
}

func arrowStyle(color, style string) string {
	result := "black"
	if color != "" {
		result = color
	}
	if style != "" {
		result += "," + style
	} else {
		result += ",->"
	}
	return result
}

func instSize(list []string, abbr string) int {
	if list[0] == abbr {
		return 0
	}
	return 6
}

// GetTemplate return a parsed template
func GetTemplate() *template.Template {
	funcMap := template.FuncMap{
		"expand":     expand,
		"anchor":     anchor,
		"instSize":   instSize,
		"arrowStyle": arrowStyle,
	}
	return template.Must(template.New("pgfumlsd").Funcs(funcMap).Delims("##", "##").Parse(theTemplate))

}

const theTemplate = `
##- template "header" ##
##- template "document" . ##

##- define "header" ##
\documentclass[tikz,border=3mm]{standalone}
\usepackage[underline=false,roundedcorners=true]{pgf-umlsd}
\usepackage{underscore}
\usepackage{syntax}
\usepackage{hyperref}
\usetikzlibrary{shadows,positioning}
\tikzset{every shadow/.style={fill=none,shadow xshift=0pt,shadow yshift=0pt}}

%% Redefine mess for node with non-trival label, add support arrow shape
\renewcommand{\mess}[5][0]{
  \stepcounter{seqlevel}
  \path
  (#2)+(0,-\theseqlevel*\unitfactor-0.7*\unitfactor) node (mess from) {};
  \addtocounter{seqlevel}{#1}
  \path
  (#4)+(0,-\theseqlevel*\unitfactor-0.7*\unitfactor) node (mess to) {};
  \draw[#5,>=angle 60] (mess from) -- (mess to) node[midway, above]
  {#3};
}
##- end ##

##- define "document" ##
\begin{document}
\sffamily
\small
\begin{sequencediagram}
\tikzstyle{inststyle}+=[drop shadow={opacity=0.9,fill=lightgray}]
\def\unitfactor{1.2}
##- template "actors" . ##
##- template "sequences" . ##
\end{sequencediagram}
\end{document}
##- end ##

##- define "actors" ##
    ##- $actorList := .ActorList ##
    ##- range $idx, $actor := .Actor ##
        ##- range $abbr, $name := $actor ##
\newinst[## instSize $actorList $abbr ##]{##$abbr##}{##$name##}
        ##- end ##
    ##- end ##
##- end ##

##- define "sequences" ##

    ##- $actorList := .ActorList ##
    ##- range .Sequence ##
      ##- if eq .Source "EMPTYLINE" ##
\postlevel
      ##- else ##
\mess[## .Delay ##]{## .Source ##}{## expand .Label ##}{## .Destination ##}{## arrowStyle .Color .Style ##}
\node [anchor=## anchor $actorList .Source .Destination "from"  ##] at (mess from) {## expand .AnnotationFrom ##};
\node [anchor=## anchor $actorList .Source .Destination "to"  ##] at (mess to) {## expand .AnnotationTo ##};
      ##- end ##
    ##- end##
    ##- if .FootNote ##
\node [anchor=north west] at ([yshift=-20] current bounding box.south west){## expand .FootNote ##};
    ##- end##
##- end ##
`
