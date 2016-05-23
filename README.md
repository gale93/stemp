# stemp
This package aims to help the loading, compiling and managing data of Golang's Templates.

### Installation

```
go get "github.com/gale93/stemp"
```

### Code Sample

You will have a base.tmpl file that looks like this:

base.tmpl
```html5
{{define "base"}}
<!doctype html>

<html lang="en">
<head>
  <meta charset="utf-8">

  <title>Stemp</title>
  <meta name="description" content="Stemp">
  <meta name="author" content="Matteo Galeotti">
  
  {{template "includes" .}}
  
</head>

<body>
	<nav></nav>
	
	<div id="content">
		{{template "content" .}}
	</div>
</body>
</html>
{{end}}

{{define "includes" }}{{end}}
{{define "content" }}{{end}}
```

And a index.html who is the specific content while the user land on main page:

index.html
```html5
{{define "content"}}
	<h1> Welcome </h1>
{{end}}
```


So your Golang file will looks like this !

```Go
package main

import (
	"fmt"
	"net/http"
	"os"
	"stemplate/stemp"
)

var st *stemp.Stemplate

func handler(w http.ResponseWriter, r *http.Request) {
	st.Render(&w, "index.html")
}

func main() {
	var err error
	// Just give to Stemplate Plugin the folder of your templates and it will do de rest
	st, err = stemp.NewStemplate("./views/")
	
	// Use this when you are in development and dont want to restart server after modify html files 
	st.LiveReload = true
	

	if err != nil {
		fmt.Println("error initialitating Stemplate")
		os.Exit(1)
	}

	http.HandleFunc("/index", handler)

	// Let's get this party started
	http.ListenAndServe(":9393", nil)
}

```


### Template's folder

In this folder we have 2 different type of files

* .tmpl
* .html

stemp will compiles the templates and make them viable in the following way:
Pseudo code:
```
For Each [Html file]
	Join it all [.tmpl] files...
```