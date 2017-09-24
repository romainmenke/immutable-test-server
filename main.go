package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router()))

}

func router() http.Handler {
	homeRoute := home()
	immutableRoute := immutable()
	etagRoute := etag()
	immutableAndEtagRoute := immutableAndEtag()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/immutable-and-etag" || r.URL.Path == "/immutable-and-etag/stylesheet.css" {
			immutableAndEtagRoute.ServeHTTP(w, r)
		} else if r.URL.Path == "/etag" || r.URL.Path == "/etag/stylesheet.css" {
			etagRoute.ServeHTTP(w, r)
		} else if r.URL.Path == "/immutable" || r.URL.Path == "/immutable/stylesheet.css" {
			immutableRoute.ServeHTTP(w, r)
		} else {
			homeRoute.ServeHTTP(w, r)
		}
	})
}

func home() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "private")
		w.Header().Set("Content-Type", "text/html")

		b := fmt.Sprintf(homeHTML)
		w.Write([]byte(b))
		return

	})
}

func immutableAndEtag() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/immutable-and-etag/stylesheet.css" {
			etag := fmt.Sprintf(`"%s"`, time.Now().Format("04:05")[:4]+"0")

			if r.Header.Get("If-None-Match") == etag {
				w.WriteHeader(304)
				w.Write([]byte{})
				return
			}

			w.Header().Set("Cache-Control", "max-age=10, must-revalidate, immutable")
			w.Header().Set("Age", "0")
			w.Header().Set("Date", time.Now().Format(http.TimeFormat))
			w.Header().Set("Content-Type", "text/css")
			w.Header().Set("Etag", etag)

			b := fmt.Sprintf(css, time.Now().Format("04:05"))
			w.Write([]byte(b))
			return
		}

		w.Header().Set("Cache-Control", "private")
		w.Header().Set("Content-Type", "text/html")

		b := fmt.Sprintf(exampleHTML, "immutable-and-etag", "immutable-and-etag", time.Now().Format("04:05"))
		w.Write([]byte(b))
		return

	})
}

func etag() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/etag/stylesheet.css" {
			etag := fmt.Sprintf(`"%s"`, time.Now().Format("04:05")[:4]+"0")

			if r.Header.Get("If-None-Match") == etag {
				w.WriteHeader(304)
				w.Write([]byte{})
				return
			}

			w.Header().Set("Cache-Control", "max-age=10, must-revalidate")
			w.Header().Set("Age", "0")
			w.Header().Set("Date", time.Now().Format(http.TimeFormat))
			w.Header().Set("Content-Type", "text/css")
			w.Header().Set("Etag", etag)

			b := fmt.Sprintf(css, time.Now().Format("04:05"))
			w.Write([]byte(b))
			return
		}

		w.Header().Set("Cache-Control", "private")
		w.Header().Set("Content-Type", "text/html")

		b := fmt.Sprintf(exampleHTML, "etag", "etag", time.Now().Format("04:05"))
		w.Write([]byte(b))
		return

	})
}

func immutable() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/immutable/stylesheet.css" {
			w.Header().Set("Cache-Control", "max-age=10, immutable")
			w.Header().Set("Age", "0")
			w.Header().Set("Date", time.Now().Format(http.TimeFormat))
			w.Header().Set("Content-Type", "text/css")

			b := fmt.Sprintf(css, time.Now().Format("04:05"))
			w.Write([]byte(b))
			return
		}

		w.Header().Set("Cache-Control", "private")
		w.Header().Set("Content-Type", "text/html")

		b := fmt.Sprintf(exampleHTML, "immutable", "immutable", time.Now().Format("04:05"))
		w.Write([]byte(b))
		return

	})
}

const css = `
#timestamp {
	position: absolute;
	left: 20px;
	top: 100px;
}

#timestamp::after {
	content: '%s';
	position: absolute;
	left: 0px;
	top: 20px;
}
`

const homeHTML = `
<!DOCTYPE html>
<html>
<head>
</head>
<body>
<a href="/etag">etag</a>
<a href="/immutable">immutable</a>
<a href="/immutable-and-etag">immutable-and-etag</a>
</body>
`

const exampleHTML = `
<!DOCTYPE html>
<html>
<head>
  <link rel="stylesheet" type="text/css" href="/%s/stylesheet.css">
</head>
<body>
<p>%s</p>
<a href="/">back</a>
<div id="timestamp">%s</div>
</body>
`
