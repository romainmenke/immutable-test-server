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
	etagRoute := etag()
	etagAndImmutableRoute := etagAndImmutable()
	maxAgeRoute := maxAge()
	maxAgeAndImmutableRoute := maxAgeAndImmutable()
	immutableRoute := immutable()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/etag.css" {
			etagRoute.ServeHTTP(w, r)
		} else if r.URL.Path == "/etag-and-immutable.css" {
			etagAndImmutableRoute.ServeHTTP(w, r)
		} else if r.URL.Path == "/max-age.css" {
			maxAgeRoute.ServeHTTP(w, r)
		} else if r.URL.Path == "/max-age-and-immutable.css" {
			maxAgeAndImmutableRoute.ServeHTTP(w, r)
		} else if r.URL.Path == "/immutable.css" {
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

		b := fmt.Sprintf(exampleHTML, time.Now().Format("04:05"))
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

			w.Header().Set("Cache-Control", "max-age=30, must-revalidate, immutable")
			w.Header().Set("Age", "0")
			w.Header().Set("Date", time.Now().Format(http.TimeFormat))
			w.Header().Set("Content-Type", "text/css")
			w.Header().Set("Etag", etag)

			b := fmt.Sprintf(css, "", time.Now().Format("04:05"))
			w.Write([]byte(b))
			return
		}

		w.Header().Set("Cache-Control", "private")
		w.Header().Set("Content-Type", "text/html")

		// b := fmt.Sprintf(exampleHTML, "immutable-and-etag", "immutable-and-etag", time.Now().Format("04:05"))
		// w.Write([]byte(b))
		return

	})
}

func etag() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		etag := fmt.Sprintf(`"%s"`, time.Now().Format("04:05")[:4]+"0")

		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(304)
			w.Write([]byte{})
			return
		}

		w.Header().Set("Cache-Control", "must-revalidate")
		w.Header().Set("Age", "0")
		w.Header().Set("Date", time.Now().Format(http.TimeFormat))
		w.Header().Set("Content-Type", "text/css")
		w.Header().Set("Etag", etag)

		b := fmt.Sprintf(css, "etag-timestamp", time.Now().Format("04:05"))
		w.Write([]byte(b))
		return

	})
}

func etagAndImmutable() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		etag := fmt.Sprintf(`"%s"`, time.Now().Format("04:05")[:4]+"0")

		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(304)
			w.Write([]byte{})
			return
		}

		w.Header().Set("Cache-Control", "must-revalidate, immutable")
		w.Header().Set("Age", "0")
		w.Header().Set("Date", time.Now().Format(http.TimeFormat))
		w.Header().Set("Content-Type", "text/css")
		w.Header().Set("Etag", etag)

		b := fmt.Sprintf(css, "etag-and-immutable-timestamp", time.Now().Format("04:05"))
		w.Write([]byte(b))
		return

	})
}

func maxAge() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=30")
		w.Header().Set("Age", "0")
		w.Header().Set("Date", time.Now().Format(http.TimeFormat))
		w.Header().Set("Content-Type", "text/css")

		b := fmt.Sprintf(css, "max-age-timestamp", time.Now().Format("04:05"))
		w.Write([]byte(b))
		return

	})
}

func maxAgeAndImmutable() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=30, immutable")
		w.Header().Set("Age", "0")
		w.Header().Set("Date", time.Now().Format(http.TimeFormat))
		w.Header().Set("Content-Type", "text/css")

		b := fmt.Sprintf(css, "max-age-and-immutable-timestamp", time.Now().Format("04:05"))
		w.Write([]byte(b))
		return

	})
}

func immutable() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Cache-Control", "immutable")
		w.Header().Set("Age", "0")
		w.Header().Set("Date", time.Now().Format(http.TimeFormat))
		w.Header().Set("Content-Type", "text/css")

		b := fmt.Sprintf(css, "immutable-timestamp", time.Now().Format("04:05"))
		w.Write([]byte(b))
		return
	})
}

const css = `
#%s::before {
	content: '%s';
}
`

const exampleHTML = `
<!DOCTYPE html>
<html>
<head>
	<link rel="stylesheet" type="text/css" href="/etag.css">
	<link rel="stylesheet" type="text/css" href="/etag-and-immutable.css">
	<link rel="stylesheet" type="text/css" href="/max-age.css">
	<link rel="stylesheet" type="text/css" href="/max-age-and-immutable.css">
	<link rel="stylesheet" type="text/css" href="/immutable.css">
</head>
<body>
<a href="/">back</a>
<div id="html-timestamp">%s&nbsp;&nbsp;html</div>
<div id="etag-timestamp">&nbsp;&nbsp;etag</div>
<div id="etag-and-immutable-timestamp">&nbsp;&nbsp;etag and immutable</div>
<div id="max-age-timestamp">&nbsp;&nbsp;max-age</div>
<div id="max-age-and-immutable-timestamp">&nbsp;&nbsp;max-age and immutable</div>
<div id="immutable-timestamp">&nbsp;&nbsp;immutable</div>
</body>
`
