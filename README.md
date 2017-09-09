# Protographer

[![Build Status](https://travis-ci.org/protograph/protographer.svg?branch=master)](https://travis-ci.org/protograph/protographer) [![GoDoc](https://godoc.org/github.com/protograph/protographer?status.svg)](https://godoc.org/github.com/protograph/protographer)

Protographer is a tool that generates protocol diagrams. It reads the protocol definition (protograph) from markdown file and converts that to a LaTeX file with PGF-UMLSD syntax. 

## There are many sequence diagram generators. Why don't you use them?

None of those support diagram generation with math syntax.

## Getting Started

A protograph has 5 sections: Actors, Flows, Notes, Configs, Expressions.
  
### Actors
 
The `Actors` section defines the acronyms for the actors. Only acronyms are allowed in the flow definition (see next section). Here is an example.

```
## Actors

A: Alice
B: Bob
C: Carol
```

The definition marks A, B, C as the acronyms of Alice, Bob, Carol respectively.

### Flows

The `Flows` section defines the interactions between the actors. Here is an example.

```
## Flows

### A sends B
Hello

### B sends A
World
```

It would be transformed in a sequence diagram in LaTex form as:

```
 A               B
 | --- Hello --> |
 |               |
 | <-- World --- |
```

If the actors need to perform certain operations before sending or after receiving messages, here is another example.

```

## Flows

### A computes
I will say hello

### A sends B
Hello

### B then computes
I got hello
```

It would be transformed in a sequence diagram in LaTex form as:

```
                 A               B
I will say hello | --- Hello --> | I got hello      
```


Sometimes an actor would like to perform an operation without sending a message. Here is an example
 
```

## Flows

### A thinks
I am fine
```

It would be transformed in a sequence diagram in LaTex form as:
```
          A               B
I am fine |               |       
```

Or the actors may exchange messages without worrying too much on who initial it.

```
## Flows

### A and B exchange messages about 
some messages
```

It would be transformed in a sequence diagram in LaTex form as:
```
 A                       B
 | <-- some messages --> |
```


Messages can be customized with different styles. That can be done with a `(...)` block, like:

```
## Flows

### A sends B (delay=1, color=red, padtop=1, line=dashed)
```

The accepted styles are:
* delay: for making the message appeared as delayed
* color: for changing the color of the message
* padtop: for appending space above the message
* line: for changing the message to dashed or dotted form

### Notes

The `Notes` section defines the footnote to be appeared in the generated diagram. 

### Configs

The `Configs` section defines the configuration for the diagram. Currently, it supports only 1 configuration:

* separation: for changing the separation between the actors. Default is 6.

### Expressions

The `Expressions` section defines the regular expressions to replace the text within math mode. Here is an example

```
## Expressions
A: `A`

```

Then the string `A` would be matched (with latex syntax aware word boundary) and replaced with `` `A` ``

## Math Mode 

There are two ways to render text in math mode. Obviously, the first way is to use `$...$` so the output would be converted as is in the latex form, and be rendered as math.
 
The other way is to use the block quote syntax in markdown, and protographer will convert the block to `$...$`. Here is an example.

```
### A sends B
    g^x
```

### Text Mode in Math Mode

In LaTeX, it is a common way to use `\text{}` to denote text inside math. Protographer overrides the `` ` `` character as a deliminator to denote text mode rendering inside math mode. In simpler words, `` `foo` `` would be converted to `\text{foo}`.
 
### Default Expressions in Math Mode

Protographer has bundled the following default expressions:

* `([A-Za-z]{3,})`: `` `$1` `` - this renders any words of 3 or more letters to text mode 
* `'`: `\textnormal{\textquotesingle}` - this formats the single quote correctly in math mode
* `<-`: `\get` - this formats the arrow sign correctly
* `->`: `\to` - ditto
* `||`: `\Vert` - this formats the vertical bar symbol correctly

## Examples 

TBD


