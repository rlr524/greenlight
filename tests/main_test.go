package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const body = `{"title":"Moana","year":2016,"runtime":107, "genres":["animation","adventure"]}`

func BenchmarkUnmarshal(b *testing.B) {
	w := httptest.NewRecorder()

	for n := 0; n < b.N; n++ {
		r, err := http.NewRequest(http.MethodPost, "", strings.NewReader(body))
		if err != nil {
			b.Fatal(err)
		}

		createMovieHandlerUnmarshal(w, r)
	}
}

func BenchmarkDecoder(b *testing.B) {
	w := httptest.NewRecorder()

	for n := 0; n < b.N; n++ {
		r, err := http.NewRequest(http.MethodPost, "", strings.NewReader(body))
		if err != nil {
			b.Fatal(err)
		}

		createMovieHandlerDecoder(w, r)
	}
}

/*
(base)  rob@yuki:~/Developer/GoLandProjects/github.com/rlr524/greenlight/tests/ [main+*] ./run_tests.sh
goos: darwin
goarch: arm64
pkg: github.com/rlr524/greenlight/tests
BenchmarkUnmarshal-8     4475164              1336 ns/op            1464 B/op         20 allocs/op
BenchmarkUnmarshal-8     4456177              1373 ns/op            1464 B/op         20 allocs/op
BenchmarkUnmarshal-8     4527135              1331 ns/op            1464 B/op         20 allocs/op
BenchmarkDecoder-8       4162125              1399 ns/op            1672 B/op         22 allocs/op
BenchmarkDecoder-8       4303933              1399 ns/op            1672 B/op         22 allocs/op
BenchmarkDecoder-8       4305110              1411 ns/op            1672 B/op         22 allocs/op
PASS
ok      github.com/rlr524/greenlight/tests      44.547s

// ns/op: speed
// B/op: memory use
// allocs/op: heap allocations
*/
