package server

import "net/http"

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func FooHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Foo"))
}

func BarHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bar"))
}

func BazHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Baz"))
}
