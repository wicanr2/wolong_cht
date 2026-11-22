//go:build ignore

// Export source frames for the local trailer without altering game resources.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/wicanr2/wolong_cht/internal/assets/cutscene"
	"github.com/wicanr2/wolong_cht/internal/assets/palette"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func read(dir, name string) []byte {
	data, err := os.ReadFile(filepath.Join(dir, name))
	must(err)
	return data
}

func write(dir, name string, img image.Image) {
	f, err := os.Create(filepath.Join(dir, name))
	must(err)
	must(png.Encode(f, img))
	must(f.Close())
}

func animationFrame(frame []byte, colors [16]color.RGBA) image.Image {
	const width, height = 320, 200
	const rowBytes, planeBytes = width / 8, width / 8 * height
	dst := image.NewPaletted(image.Rect(0, 0, width, height), make(color.Palette, len(colors)))
	for i, c := range colors {
		dst.Palette[i] = c
	}
	for y := range height {
		for x := range width {
			index := 0
			byteOffset := y*rowBytes + x/8
			mask := byte(0x80) >> (x % 8)
			for plane := range 4 {
				if frame[plane*planeBytes+byteOffset]&mask != 0 {
					index |= 1 << plane
				}
			}
			dst.SetColorIndex(x, y, uint8(index))
		}
	}
	return dst
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: promo_opening_frames ORIGINAL OUTPUT")
	}
	orig, out := os.Args[1], os.Args[2]
	must(os.MkdirAll(out, 0o755))
	pal, err := palette.Parse(read(orig, "OPENPAL.BRG"))
	must(err)
	const frameBytes = 320 * 200 / 8 * 4
	for scene := 2; scene <= 4; scene++ {
		data, err := cutscene.Decode(read(orig, fmt.Sprintf("OPEN_S%d.DAT", scene)))
		must(err)
		if len(data) != frameBytes*12 {
			log.Fatalf("scene %d: unexpected decoded length %d", scene, len(data))
		}
		colors, err := pal.Bank(scene - 1)
		must(err)
		for number := range 30 {
			frame := number % 4
			if number >= 18 {
				frame = number - 18
			}
			start := frame * frameBytes
			write(out, fmt.Sprintf("scene%d-%02d.png", scene, number), animationFrame(data[start:start+frameBytes], colors))
		}
	}
	data, err := cutscene.Decode(read(orig, "OPEN_S6.DAT"))
	must(err)
	colors, err := pal.Bank(5)
	must(err)
	img, err := cutscene.Image(data, colors[:])
	must(err)
	write(out, "scene6.png", img)
}
