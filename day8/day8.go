package main

import "github.com/dotMaro/AoC2019/utils"

import "strings"

func main() {
	input := utils.ReadFile("day8/input.txt")
	image := newImage(input, 25, 6)

	utils.Print("Part 1: %d", image.findFewestZeroesLayer())
	utils.Print("Part 2:\n%v", image.decode())
}

func newImage(input string, width, height int) image {
	var cur int
	var layers []layer
	for cur < len(input) {
		var layer layer
		layer.pixels = make([][]rune, height)
		for y := 0; y < height; y++ {
			layer.pixels[y] = []rune(input[cur : cur+width])
			cur += width
		}
		layers = append(layers, layer)
	}
	return image{
		layers: layers,
		width:  width,
		height: height,
	}
}

type image struct {
	layers        []layer
	width, height int
}

type layer struct {
	pixels [][]rune // y - x
}

func (l layer) String() string {
	var b strings.Builder
	for _, row := range l.pixels {
		for _, p := range row {
			var r rune
			if p == '1' {
				r = '█'
			} else {
				r = ' '
			}
			b.WriteRune(r)
		}
		b.WriteRune('\n')
	}
	return b.String()
}

// findFewestZeroesLayer and return the one digit and two digit count
// multiplied in that layer.
func (i image) findFewestZeroesLayer() int {
	var oneCount, twoCount int
	fewestZeroes := 99999
	for _, layer := range i.layers {
		var layerZeroCount, layerOneCount, layerTwoCount int
		for _, row := range layer.pixels {
			for _, p := range row {
				switch p {
				case '0':
					layerZeroCount++
				case '1':
					layerOneCount++
				case '2':
					layerTwoCount++
				}
			}
		}
		if layerZeroCount < fewestZeroes {
			fewestZeroes = layerZeroCount
			oneCount = layerOneCount
			twoCount = layerTwoCount
		}
	}
	return oneCount * twoCount
}

// decode into a merged layer.
func (i image) decode() layer {
	var decoded layer
	// init decoded layer with transparent pixels
	decoded.pixels = make([][]rune, i.height)
	for y, row := range decoded.pixels {
		row = make([]rune, i.width)
		for p := 0; p < i.width; p++ {
			row[p] = '2'
		}
		decoded.pixels[y] = row
	}

	for _, layer := range i.layers {
		for y, row := range layer.pixels {
			for x, p := range row {
				if decoded.pixels[y][x] == '2' {
					decoded.pixels[y][x] = p
				}
			}
		}
	}

	return decoded
}
