package tracer

import (
	"image"
	"log"
	"os"

	// Support for .hdr files
	_ "github.com/mdouchement/hdr/codec/rgbe"

	"github.com/mdouchement/hdr"
	"github.com/mdouchement/hdr/tmo"
)

// HDRToImage return an hdr image as an image.Image interface
func HDRToImage(filename string) image.Image {
	fi, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Decoding hdr image...")
	m, format, err := image.Decode(fi)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("decoded image format: %v", format)

	if hdrm, ok := m.(hdr.Image); ok {
		log.Println("Applying tonemap to hdr image...")
		// t := tmo.NewLinear(hdrm)
		// t := tmo.NewLogarithmic(hdrm)
		// t := tmo.NewDefaultDrago03(hdrm)
		t := tmo.NewDefaultDurand(hdrm)
		// t := tmo.NewDefaultCustomReinhard05(hdrm)
		// t := tmo.NewDefaultReinhard05(hdrm)
		// t := tmo.NewDefaultICam06(hdrm)
		m = t.Perform()
	}

	return m
}
