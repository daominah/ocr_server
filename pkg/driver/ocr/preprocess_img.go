package ocr

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

// minHeight is the minimum image height in pixels. Images smaller than this
// are upscaled (nearest-neighbor) before erosion so that character strokes
// are thick enough for erosion to be effective.
const minHeight = 100

// ErodeImage applies horizontal morphological erosion to separate overlapping
// characters. A black pixel stays black only if all horizontal neighbors
// within the given radius are also black. Vertical neighbors are not checked,
// preserving vertical strokes while breaking horizontal connections between
// adjacent glyphs.
//
// Small images are upscaled first so erosion has enough pixel resolution.
func ErodeImage(imgBytes []byte, radius int) ([]byte, error) {
	src, err := png.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, err
	}

	// Upscale small images so strokes are wide enough for erosion.
	bounds := src.Bounds()
	h := bounds.Max.Y - bounds.Min.Y
	if h < minHeight {
		scale := (minHeight + h - 1) / h // ceil division
		src = scaleNN(src, scale)
		bounds = src.Bounds()
	}

	gray := toGray(src)
	dst := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if isForegroundH(gray, x, y, radius) {
				dst.SetGray(x, y, color.Gray{Y: 0})
			} else {
				dst.SetGray(x, y, color.Gray{Y: 255})
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// scaleNN upscales an image by an integer factor using nearest-neighbor
// interpolation. This preserves hard edges in captcha text.
func scaleNN(src image.Image, factor int) image.Image {
	bounds := src.Bounds()
	w := (bounds.Max.X - bounds.Min.X) * factor
	h := (bounds.Max.Y - bounds.Min.Y) * factor
	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst.Set(x, y, src.At(bounds.Min.X+x/factor, bounds.Min.Y+y/factor))
		}
	}
	return dst
}

// toGray converts an image to grayscale, compositing transparent pixels
// onto a white background (so transparent = white = background).
func toGray(src image.Image) *image.Gray {
	bounds := src.Bounds()
	gray := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			if a == 0 {
				gray.SetGray(x, y, color.Gray{Y: 255})
				continue
			}
			// Un-premultiply alpha, then composite onto white.
			r = r * 0xffff / a
			g = g * 0xffff / a
			b = b * 0xffff / a
			// Blend with white: out = fg*alpha + white*(1-alpha)
			fa := float64(a) / 0xffff
			rf := float64(r)/0xffff*fa + (1 - fa)
			gf := float64(g)/0xffff*fa + (1 - fa)
			bf := float64(b)/0xffff*fa + (1 - fa)
			lum := 0.299*rf + 0.587*gf + 0.114*bf
			gray.SetGray(x, y, color.Gray{Y: uint8(lum * 255)})
		}
	}
	return gray
}

// isForegroundH returns true only if the pixel at (x,y) and all horizontal
// neighbors within radius are dark (below threshold).
func isForegroundH(gray *image.Gray, x, y, radius int) bool {
	bounds := gray.Bounds()
	const threshold = 128
	if y < bounds.Min.Y || y >= bounds.Max.Y {
		return false
	}
	for dx := -radius; dx <= radius; dx++ {
		nx := x + dx
		if nx < bounds.Min.X || nx >= bounds.Max.X {
			return false
		}
		if gray.GrayAt(nx, y).Y >= threshold {
			return false
		}
	}
	return true
}
