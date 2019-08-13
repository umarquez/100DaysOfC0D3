package main

type CookieType string
type internalCookieStatus int

const (
	mixed internalCookieStatus = iota
	served
	baked
	packaged
)
