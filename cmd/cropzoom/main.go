package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
)

func main() {
	in := flag.String("in", "", "")
	out := flag.String("out", "", "")
	x := flag.Int("x", 0, "")
	y := flag.Int("y", 0, "")
	w := flag.Int("w", 64, "")
	h := flag.Int("h", 64, "")
	scale := flag.Int("scale", 6, "")
	flag.Parse()
	f, _ := os.Open(*in)
	src, err := png.Decode(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	f.Close()
	fmt.Println("size", src.Bounds())
	dst := image.NewRGBA(image.Rect(0, 0, *w**scale, *h**scale))
	for j := 0; j < *h**scale; j++ {
		for i := 0; i < *w**scale; i++ {
			dst.Set(i, j, src.At(*x+i/(*scale), *y+j/(*scale)))
		}
	}
	g, _ := os.Create(*out)
	png.Encode(g, dst)
	g.Close()
}
