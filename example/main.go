package main

import (
	gab "github.com/right-hearted/gab"
)

func main() {
	g, e := gab.New("780002928:AAGwW1BIp2Px1BUJgImde71HAM56m28uPKY")
	if e != nil {
		panic(e)
	}
	g.RunTelBot()
}
