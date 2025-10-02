package http

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Data struct {
	Title   string
	Headers map[string]string
	EnvVars map[string]string
}

func (s *Server) InfoHandler(w http.ResponseWriter, r *http.Request) {
	httpBody := Data{}
	httpBody.Title = "Ciavash"
	httpBody.Headers = make(map[string]string)
	httpBody.EnvVars = make(map[string]string)
	for k, v := range r.Header {
		httpBody.Headers[k] = strings.Join(v, ", ")
	}

	for _, element := range os.Environ() {
		key := strings.Split(element, "=")[0]
		value := strings.Split(element, "=")[1]
		httpBody.EnvVars[key] = value
	}

	tmpl, err := template.New("infoPage").Parse(`
	<!DOCTYPE html>
	<html lang="en">
	  <head>
	  <style>
	    body {
		    background-color: lightblue;
	    }
	    h1   {color: red;}
	    p {
		max-width: 640px;
		border: 2px solid powderblue;
		margin: 50px;
		color: black;
	    }
	    li   {
		padding: 2px;
		margin-left: 50px;
		color: black;
	    }

	  </style>
	    <title>My Page</title>
	  </head>
	  <body>
		  <H1>Welcome to {{.Title}} page</H1>
		<H2> Headers </H2>
		{{range $k, $v := .Headers}}<li><strong>{{$k}}</strong>: {{$v}}</li>{{end}}
		<H2> Environment Variables </H2>
		{{range $k, $v := .EnvVars}}<li><strong>{{$k}}</strong>: {{$v}}</li>{{end}}
	  </body>
	</html>
	`)
	if err != nil {
		log.Print(err)
	}

	err = tmpl.Execute(w, httpBody)
}
