# jstruct
## A CLI written in Go for the translation of JSON into Go struct types

Designed as a shortcut for handcoding structs from JSON responses tests and example code. 
The program should take the JSON and output modularised structs for dropping straight into your code or as a baseline for more advanced use.

### #Prerequisites
- Goland
- Git
- Light knowledge of terminal and Git

---
### #Installation 
Download the files using 
>`git clone https://github.com/Karlsburg87/jstruct`

Then `cd` into the directory and run 
>`go install .`

(notice the period)

You could be able to run the program using 
>`jstruct <flags>` 

if you have $GOPATH set

---
### #How to use:
There are only a few options:
| Flag | Purpose |	
| --- | --- |
| -data | "The raw json data. if '-file' or 'f' flag is given then file contents will overwrite this." |
| -file | "The path to the file containing json data" |
| -f | "Shorthand for -file. The path to the file containing JSON data" |
| -out | "Location to create the go file containing struct" |
| -name | "The name of the struct in the output file" |
  
  
---
### #Examples

__*JSON input*__
```json
{"widget": {
    "debug": "on",
    "window": {
        "title": "Sample Konfabulator Widget",
        "name": "main_window",
        "width": 500,
        "height": 500
    },
    "image": { 
        "src": "Images/Sun.png",
        "name": "sun1",
        "hOffset": 250,
        "vOffset": 250,
        "alignment": "center"
    },
    "text": {
        "data": "Click Here",
        "size": 36,
        "style": "bold",
        "name": "text1",
        "hOffset": 250,
        "vOffset": 100,
        "alignment": "center",
        "onMouseUp": "sun1.opacity = (sun1.opacity / 100) * 90;"
    }
}} 
 
```
Command line input example:
>jstruct -file testData.json -out new.go -name test

__*Golang Strut output*__
```golang
package test 

type Test struct {
	Widget TestWidget `json:"widget,omitempty"`
}
 

type TestWidget struct {
	Window TestWidgetWindow `json:"window,omitempty"`
	Image TestWidgetImage `json:"image,omitempty"`
	Text TestWidgetText `json:"text,omitempty"`
	Debug string `json:"debug,omitempty"`
}
 

type TestWidgetWindow struct {
	Title string `json:"title,omitempty"`
	Name string `json:"name,omitempty"`
	Width float64 `json:"width,omitempty"`
	Height float64 `json:"height,omitempty"`
}
 

type TestWidgetImage struct {
	Src string `json:"src,omitempty"`
	Name string `json:"name,omitempty"`
	HOffset float64 `json:"hOffset,omitempty"`
	VOffset float64 `json:"vOffset,omitempty"`
	Alignment string `json:"alignment,omitempty"`
}
 

type TestWidgetText struct {
	Name string `json:"name,omitempty"`
	HOffset float64 `json:"hOffset,omitempty"`
	VOffset float64 `json:"vOffset,omitempty"`
	Alignment string `json:"alignment,omitempty"`
	OnMouseUp string `json:"onMouseUp,omitempty"`
	Data string `json:"data,omitempty"`
	Size float64 `json:"size,omitempty"`
	Style string `json:"style,omitempty"`
}


/***********
A Rhythmic Sound Project
***********/

```