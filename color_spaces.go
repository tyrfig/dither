package dither

import (
	"image/color"
	"math"
)

// linearize1 linearizes an R, G, or B channel value from an sRGB color.
// Must be in the range [0, 1].
func linearize1(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func linearize65535(i uint16) uint16 {
	v := float64(i) / 65535.0
	return uint16(math.RoundToEven(linearize1(v) * 65535.0))
}

func linearize255to65535(i uint8) uint16 {
	v := float64(i) / 255.0
	return uint16(math.RoundToEven(linearize1(v) * 65535.0))
}

// toLinearRGB converts a non-linear sRGB color to a linear RGB color space.
// RGB values are taken directly and alpha value is ignored, so this will not
// handle non-opaque colors properly.
func toLinearRGB(c color.Color) (uint16, uint16, uint16) {
	// Optimize for different color types
	switch v := c.(type) {
	case color.Gray:
		g := linearize255to65535(v.Y)
		return g, g, g
	case color.Gray16:
		g := linearize65535(v.Y)
		return g, g, g
	case color.NRGBA:
		return linearize255to65535(v.R), linearize255to65535(v.G), linearize255to65535(v.B)
	case color.NRGBA64:
		return linearize65535(v.R), linearize65535(v.G), linearize65535(v.B)
	case color.RGBA:
		return linearize255to65535(v.R), linearize255to65535(v.G), linearize255to65535(v.B)
	case color.RGBA64:
		return linearize65535(v.R), linearize65535(v.G), linearize65535(v.B)
	}

	r, g, b, _ := c.RGBA()
	return linearize65535(uint16(r)), linearize65535(uint16(g)), linearize65535(uint16(b))
}


func linearRGBtoXYZ(r uint16, g uint16, b uint16) (float64, float64, float64) {
	 x := (0.412453*float64(r) + 0.357580*float64(g) + 0.180423*float64(b)) / 65535.0
	 y := (0.212671*float64(r) + 0.715160*float64(g) + 0.072169*float64(b)) / 65535.0
         z := (0.019334*float64(r) + 0.119193*float64(g) + 0.950227*float64(b)) / 65535.0
    return x, y, z
}

func xyz2lab(x float64, y float64, z float64) (float64, float64, float64) {    
    x_scaled := x / 0.95047
    y_scaled := y
    z_scaled := z / 1.08883
    
    x_int := 0.0
    if x_scaled > .008856 {
	x_int = math.Pow(x_scaled, 1.0/3.0)
    } else {
	x_int = 7.787*x_scaled + 16./116.
    }

    y_int := 0.0
    if y_scaled > .008856 {
	y_int = math.Pow(y_scaled, 1.0/3.0)
    } else {
	y_int = 7.787*y_scaled + 16./116.
    }
	
    z_int := 0.0
    if z_scaled > .008856 {
	z_int = math.Pow(z_scaled, 1.0/3.0)
    } else {
	z_int = 7.787*z_scaled + 16./116.
    }
	
    L := 116. * y_int -16.
    a := 500.*(x_int - y_int)
    b := 200.*(y_int - z_int)

    return L, a, b
}
