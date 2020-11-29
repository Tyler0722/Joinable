package main

// Options store configuration for a session or session store
// Fields are a subset of http.Cookie fields
type Options struct {
	Path string
	Domain string
	MaxAge int
	Secure bool
	HttpOnly bool
}