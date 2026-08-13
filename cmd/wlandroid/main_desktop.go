//go:build !android

package main

import mobilegame "github.com/wicanr2/wolong_cht/mobile/wolong"

func main() {
	if err := mobilegame.RunDesktop(); err != nil {
		panic(err)
	}
}
