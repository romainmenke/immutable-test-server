package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router()))

}

func router() http.Handler {
	homeRoute := home()
	immutableRoute := immutable()
	etagRoute := etag()
	etagAndImmutableRoute := etagAndImmutable()
	maxAgeRoute := maxAge()
	maxAgeAndImmutableRoute := maxAgeAndImmutable()
	maxAgeAndImmutableVersionedRoute := maxAgeAndImmutableVersioned()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/immutable.css" {
			immutableRoute.ServeHTTP(w, r)
		} else if r.URL.Path == "/etag.css" {
			etagRoute.ServeHTTP(w, r)
		} else if r.URL.Path == "/etag-and-immutable.css" {
			etagAndImmutableRoute.ServeHTTP(w, r)
		} else if r.URL.Path == "/max-age.css" {
			maxAgeRoute.ServeHTTP(w, r)
		} else if r.URL.Path == "/max-age-and-immutable.css" {
			maxAgeAndImmutableRoute.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/max-age-and-immutable-versioned") {
			maxAgeAndImmutableVersionedRoute.ServeHTTP(w, r)
		} else {
			homeRoute.ServeHTTP(w, r)
		}
	})
}

func home() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "private")
		w.Header().Set("Content-Type", "text/html")

		b := fmt.Sprintf(exampleHTML, fmt.Sprint(time.Now().Round(time.Minute).Unix()), time.Now().Format("04:05"))
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

func etag() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		etag := fmt.Sprintf(`"%s"`, time.Now().Format("04:05")[:3]+"00")

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

		etag := fmt.Sprintf(`"%s"`, time.Now().Format("04:05")[:3]+"00")

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
		w.Header().Set("Cache-Control", "max-age=60")
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
		w.Header().Set("Cache-Control", "max-age=60, immutable")
		w.Header().Set("Age", "0")
		w.Header().Set("Date", time.Now().Format(http.TimeFormat))
		w.Header().Set("Content-Type", "text/css")

		b := fmt.Sprintf(css, "max-age-and-immutable-timestamp", time.Now().Format("04:05"))
		w.Write([]byte(b))
		return

	})
}

func maxAgeAndImmutableVersioned() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=31536000, immutable")
		w.Header().Set("Age", "0")
		w.Header().Set("Date", time.Now().Format(http.TimeFormat))
		w.Header().Set("Content-Type", "text/css")

		stamp := strings.TrimSuffix(r.URL.Path, ".css")
		parts := strings.Split(stamp, "-")
		stamp = parts[len(parts)-1]

		var t time.Time

		intStamp, err := strconv.ParseInt(stamp, 10, 64)
		if err != nil {
			t = time.Now()
		}

		t = time.Unix(intStamp, 0)

		b := fmt.Sprintf(css, "max-age-and-immutable-versioned-timestamp", t.Format("04:05"))
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
	<link rel="stylesheet" type="text/css" href="/immutable.css">
	<link rel="stylesheet" type="text/css" href="/etag.css">
	<link rel="stylesheet" type="text/css" href="/etag-and-immutable.css">
	<link rel="stylesheet" type="text/css" href="/max-age.css">
	<link rel="stylesheet" type="text/css" href="/max-age-and-immutable.css">
	<link rel="stylesheet" type="text/css" href="/max-age-and-immutable-versioned-%s.css">
	<style>
		h1 {
			font-size: 18px;
			text-align: center;
		}

		p {
			font-size: 16px;
			text-align: center;
		}

		li {
			font-size: 15px;
			padding-bottom: 10px;
		}

		.stamps {
			width: 300px;
			margin: 40px auto;
			padding: 0 0 0 40px;
		}

		.stamps > * {
			padding: 3px 0;
		}

		.info {
			width: 600px;
			margin: 40px auto;
		}
	</style>
</head>
<body>
<h1>Immutable</h1>
<p>Best to test over https and a browser that supports "immutable" like Firefox.</p>
<p><a href="https://github.com/romainmenke/immutable-test-server" target="_blank">code</a></p>
<div class="stamps">
	<div id="html-timestamp">%s&nbsp;&nbsp;html</div>
	<div id="immutable-timestamp">&nbsp;&nbsp;immutable</div>
	<div id="etag-timestamp">&nbsp;&nbsp;etag</div>
	<div id="etag-and-immutable-timestamp">&nbsp;&nbsp;etag and immutable</div>
	<div id="max-age-timestamp">&nbsp;&nbsp;max-age</div>
	<div id="max-age-and-immutable-timestamp">&nbsp;&nbsp;max-age and immutable</div>
	<div id="max-age-and-immutable-versioned-timestamp">&nbsp;&nbsp;max-age and immutable versioned</div>
</div>

<div class="info">
	<ul>
	<li>The html has "Cache-Control: private" and the timestamp will update on each request.</li>
	<li>All other timestamps are set through css and will update when the css file is updated.</li>
	<li>"max-age" is set to 60 seconds and Etags change every 60 seconds.</li>
	</ul>
</div>
</body>
`
