package protographer

const theTemplate = `
<<- template "header" >>
<<- template "body" . >>
<<- template "footer" >>

<<- define "header" >>
\documentclass[tikz,border=3mm]{standalone}
\usepackage[underline=false,roundedcorners=true]{pgf-umlsd}
\usepackage{underscore}
\usepackage{syntax}
\usepackage{hyperref}
\usepackage{amsmath}
\usepackage{amsfonts}
\usepackage{mathtools}
\usepackage{textcomp}
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

%% make text appear as op in math mode
\newcommand{\textop}[1]{\relax\ifmmode\mathop{\text{#1}}\else\text{#1}\fi}
\newcommand{\stextop}[1]{\relax\ifmmode\mathop{\text{ #1 }}\else\text{#1}\fi}

\begin{document}
\sffamily
\small

\begin{sequencediagram}
\tikzstyle{inststyle}+=[drop shadow={opacity=0.9,fill=lightgray}]
\def\unitfactor{1.2}

<<- end >>

<<- define "footer" >>

\end{sequencediagram}
\end{document}
<<- end >>

<<- define "body" >>
<<- template "actors" . >>
<<- template "sequences" . >>
<<- template "notes" . >>
<<- end >>

<<- define "actors" >>

%
% actors
%
<<- $actorList := .ActorList >>
<<- $w := .Style.Sep >>
<<- range $idx, $actor := $actorList >>
<<- if eq $idx 0 >>
\newinst[0]{<< index $actor 0 >>}{<< index $actor 1 >>}
<<- else >>
\newinst[<< $w >>]{<< index $actor 0 >>}{<< index $actor 1 >>}
<<- end >>
<<- end >>
<<- end >>

<<- define "sequences" >>

%
% sequences
%
<<- range .Sequence >>
<< padtop .Style.PadTop >>
\mess[<< .Style.Delay >>]{<< .Sender >>}{<< expand .Description.Message >>}{<< .Receiver >>}{<< .Style.Color >>,<< .Style.Format >><< with .Style.Line >>,<<.>><< end >>}
\node [anchor=<< l2r .Style.L2R 0 >>] at (mess from) {<< expand .Description.From >>};
\node [anchor=<< l2r .Style.L2R 1 >>] at (mess to) {<< expand .Description.To >>};
<<- end >>
<<- end >>

<<- define "notes" >>
<<- with .Note >>

%
% notes
%
\node [anchor=north west] (rect) at ([yshift=-20] current bounding box.south west) [draw,thick] {<< expand . >>};
<<- end >>
<<- end >>
`