package pixelding

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

const MaxX = 400
const MaxY = 200

const OutOfBoundsError = "out of bounds"
const AlreadySetError = "already set"
const DimensionError = "dimension error"

//todo Screen Trim Nur ausgeben was wirklich da ist
//todo Return Screen als String Array
//todo Clear Screen
//todo Stamp Bitmuster
//todo Font?

type PixelDing struct {
	init    bool
	matrix  [][]bool
	x, y    int
	aspectX int
	aspectY int
	scalef  float64
	debug   bool
	invert  bool
	toggle  bool
	render  int
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) Toggle(b bool) {
	p.toggle = b
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) Invert(b bool) {
	p.invert = b
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) Debug(b bool) {
	p.debug = b
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) Scale(s float64) {
	p.scalef = s
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) Aspect(x, y int) {
	p.aspectX = 0
	p.aspectY = 0
	if x > 0 {
		p.aspectX = 1
	}
	if y > 0 {
		p.aspectY = 1
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) Display() {
	cox := []string{
		string(32),     // 0
		string(0x2597), // 1
		string(0x2596), // 2
		string(0x2584), // 3
		string(0x259D), // 4
		string(0x2590), // 5
		string(0x259E), // 6
		string(0x259F), // 7
		string(0x2598), // 8
		string(0x259A), // 9
		string(0x258C), // 10
		string(0x2599), // 11
		string(0x2580), // 12
		string(0x259C), // 13
		string(0x259B), // 14
		string(0x2588), // 15
	}
	cmp := true
	if p.invert {
		cmp = !cmp
	}
	for y := 0; y < p.y-1; y = y + 2 - p.aspectY {
		if p.debug {
			fmt.Print(fmt.Sprintf("%03d", y))
		}
		for x := 0; x < p.x-1; x = x + 2 - p.aspectX { // = x + 2 {
			bit := 0
			if p.getPixel(x, y) == cmp {
				bit += 8
			}
			if p.getPixel(x+1, y) == cmp {
				bit += 4
			}
			if p.getPixel(x, y+1) == cmp {
				bit += 2
			}
			if p.getPixel(x+1, y+1) == cmp {
				bit += 1
			}
			fmt.Print(cox[bit])
		}
		fmt.Println()
	}
}

func (p *PixelDing) Clear() {

	p.matrix = make([][]bool, p.y)

	for i := range p.matrix {
		p.matrix[i] = make([]bool, p.x)
	}

}


//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) check(x, y int) bool {
	if x < 0 || x > p.x-1 || y < 0 || y > p.y-1 {
		return false
	}
	return true
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) scale(x, y int) (int, int) {
	if p.scalef != 0.0 {
		x = int(float64(x) * p.scalef)
		y = int(float64(y) * p.scalef)
	}
	return x, y
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) sscale(x int) int {
	if p.scalef != 0.0 {
		x = int(float64(x) * p.scalef)
	}
	return x
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) Dimensions(x, y int) error {
	if p.init {
		return errors.New(AlreadySetError)
	}
	if x < 1 || y < 1 {
		return errors.New(DimensionError)
	}
	if x > MaxX || y > MaxY {
		return errors.New(DimensionError)
	}
	p.matrix = make([][]bool, y)

	for i := range p.matrix {
		p.matrix[i] = make([]bool, x)
	}
	p.x = x
	p.y = y
	p.init = true
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) getPixel(x, y int) bool {
	// x, y = p.scale(x, y)
	if !p.check(x, y) {
		return false
	}
	return p.matrix[y][x]
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) setPixel(x, y int, b bool) {
	if !p.check(x, y) {
		return
	}
	if p.toggle {
		p.matrix[y][x] = !p.matrix[y][x]
	} else {
		p.matrix[y][x] = b
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) Pixel(x, y int, b bool) {
	x, y = p.scale(x, y)
	if !p.check(x, y) {
		return
	}
	p.setPixel(x, y, b)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) LinePath(s string) {
	var x, y float64
	var lx, ly float64
	var ix, iy float64
	var c1x,c1y,c2x,c2y float64
	var err error
	sw := 0
	cmd := ""

	//var err error
	//s := "M 12 13 L 100 100 200 200 l 90 -90 M80 80"
	//r := "([MLml]?)\\D?([0-9.\\-]+)\\D?([0-9.\\-]+)"
	r := "[MmLlHhVvZzCc]|[+-]?\\d+\\.\\d+|[+-]?\\d+|[+-]?\\.\\d+"

	re := regexp.MustCompile(r)
	ss := re.FindAllString(s, -1)

	//fmt.Println("The Cow:", s)
	it := 0
	for _, ssSub := range ss {

		//fmt.Println(":::", ssSub, len(ssSub))
		it++

		_, err = strconv.ParseFloat(ssSub, 64)

		if err != nil { //Not numeric ...
			switch ssSub {
			case "M", "m", "L", "l", "V", "v", "H", "h":
				cmd = ssSub
				sw = 1
			case "C","c":
				cmd = ssSub
				sw =1
			case "Z", "z":
				cmd = ssSub
				x = ix
				y = iy
				sw = 9
			default:
				if p.debug {
					fmt.Println("Unknown :", ssSub)
				}
			}
		} else {
			if sw == 0 {
				sw = 1
			}
			switch {

			case (cmd == "M" || cmd == "L") && sw == 1:
				x, _ = strconv.ParseFloat(ssSub, 64)
				ix = x
				sw = 2
			case (cmd == "M" || cmd == "L") && sw == 2:
				y, _ = strconv.ParseFloat(ssSub, 64)
				iy = y
				sw = 9

			case (cmd == "m" || cmd == "l") && sw == 1:
				x, _ = strconv.ParseFloat(ssSub, 64)
				x = lx + x
				ix = x
				sw = 2
			case (cmd == "m" || cmd == "l") && sw == 2:
				y, _ = strconv.ParseFloat(ssSub, 64)
				y = ly + y
				iy = y
				sw = 9

			case cmd == "v" && sw == 1:
				y, _ = strconv.ParseFloat(ssSub, 64)
				y = ly + y
				sw = 9
			case cmd == "V" && sw == 1:
				y, _ = strconv.ParseFloat(ssSub, 64)
				sw = 9

			case cmd == "h" && sw == 1:
				x, _ = strconv.ParseFloat(ssSub, 64)
				x = lx + x
				sw = 9
			case cmd == "H" && sw == 1:
				x, _ = strconv.ParseFloat(ssSub, 64)
				sw = 9

			case cmd == "c" && sw == 1:
				c1x, _ = strconv.ParseFloat(ssSub, 64)
				c1x = lx + c1x
				sw = 2
			case cmd == "c" && sw == 2:
				c1y, _ = strconv.ParseFloat(ssSub, 64)
				c1y = ly + c1y
				sw = 3
			case cmd == "c" && sw == 3:
				c2x, _ = strconv.ParseFloat(ssSub, 64)
				c2x = lx + c2x
				sw = 4
			case cmd == "c" && sw == 4:
				c2y, _ = strconv.ParseFloat(ssSub, 64)
				c2y = ly + c2y
				sw = 5
			case cmd == "c" && sw == 5:
				x, _ = strconv.ParseFloat(ssSub, 64)
				x = lx + x
				sw = 6
			case cmd == "c" && sw == 6:
				y, _ = strconv.ParseFloat(ssSub, 64)
				y = ly + y
				sw = 9

			}
		}
		if sw == 9 {
			//fmt.Println(it, "(", cmd, ") lx", lx, "ly", ly, "x", x, "y", y)
			//fmt.Println("--------------------------------")
			switch cmd {
			case "L", "l":
				p.Line(int(lx), int(ly), int(x), int(y))
				lx = x
				ly = y
			case "M", "m":
				lx = x
				ly = y
			case "H", "h":
				p.Line(int(lx), int(ly), int(x), int(y))
				lx = x
			case "V", "v":
				p.Line(int(lx), int(ly), int(x), int(y))
				ly = y
			case "C", "c":
				p.Bezier(int(lx),int(ly),int(c1x),int(c1y),int(c2x),int(c2y),int(x),int(y))
				p.Line(int(lx), int(ly), int(x), int(y))
				lx = x
				ly = y
			case "Z", "z":
				p.Line(int(lx), int(ly), int(x), int(y))
				lx = x
				ly = y
			}
			sw = 0
		}

	}

}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) Bezier(x0, y0, x1, y1, x2, y2, x3, y3 int) {

	xu := 0.0
	yu := 0.0
	u := 0.0
	//i := 0

	x0, y0 = p.scale(x0, y0)
	x1, y1 = p.scale(x1, y1)
	x2, y2 = p.scale(x2, y2)
	x3, y3 = p.scale(x3, y3)

	for u = 0.0; u <= 1.0; u += 0.001 {
		xu = math.Pow(1-u, 3)*float64(x0) + 3*u*math.Pow(1-u, 2)*float64(x1) + 3*math.Pow(u, 2)*(1-u)*float64(x2) + math.Pow(u, 3)*float64(x3)
		yu = math.Pow(1-u, 3)*float64(y0) + 3*u*math.Pow(1-u, 2)*float64(y1) + 3*math.Pow(u, 2)*(1-u)*float64(y2) + math.Pow(u, 3)*float64(y3)
		p.setPixel(int(xu), int(yu), true)
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) Rectangle(x0, y0, x1, y1 int, b bool, f bool) {

	x0, y0 = p.scale(x0, y0)
	x1, y1 = p.scale(x1, y1)

	for i := x0; i <= x1; i++ {
		p.setPixel(i, y0, b)
	}

	for i := y0; i < y1; i++ {
		if f {
			for j := x0; j <= x1; j++ {
				p.setPixel(j, i, b)
			}
		} else {
			p.setPixel(x0, i, b)
			p.setPixel(x1, i, b)
		}
	}

	for i := x0; i <= x1; i++ {
		p.setPixel(i, y1, b)
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) floodFill(x, y int, prevC, newC bool) {
	if x < 0 || x >= p.x || y < 0 || y >= p.y {
		return
	}

	if prevC != p.getPixel(x, y) {
		return
	}
	// Replace the color at (x, y)
	p.setPixel(x, y, newC)
	// Recur for north, east, south and west
	p.floodFill(x+1, y, prevC, newC)
	p.floodFill(x-1, y, prevC, newC)
	p.floodFill(x, y+1, prevC, newC)
	p.floodFill(x, y-1, prevC, newC)
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) Fill(x, y int, newC bool) {
	x, y = p.scale(x, y)
	prevC := p.getPixel(x, y)
	if prevC == newC {
		return
	}
	p.floodFill(x, y, prevC, newC)
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) EllipseRect(x0, y0, x1, y1 int) {
	x0, y0 = p.scale(x0, y0)
	x1, y1 = p.scale(x1, y1)
	a := abs(x1 - x0)
	b := abs(y1 - y0)
	b1 := b & 1

	dx := 4 * (1 - a) * b * b
	dy := 4 * (b1 + 1) * a * a
	e := dx + dy + b1*a*a
	e2 := 0

	if x0 > x1 {
		x0 = x1
		x1 += a
	}
	if y0 > y1 {
		y0 = y1
	}
	y0 += (b + 1) / 2
	y1 = y0 - b1
	a *= 8 * a
	b1 = 8 * b * b
	for {
		p.setPixel(x1, y0, true) /*   I. Quadrant */
		p.setPixel(x0, y0, true) /*  II. Quadrant */
		p.setPixel(x0, y1, true) /* III. Quadrant */
		p.setPixel(x1, y1, true) /*  IV. Quadrant */
		e2 = 2 * e
		if e2 >= dx {
			x0++
			x1--
			dx += b1
			e += dx
		} /* x step */
		if e2 <= dy {
			y0++
			y1--
			dy += a
			e += dy
		} /* y step */
		if x0 > x1 {
			break
		}
	}

	for {
		if y0-y1 >= b {
			break
		}
		p.setPixel(x0-1, y0, true)
		p.setPixel(x1+1, y0, true)
		y0++
		p.setPixel(x0-1, y1, true)
		p.setPixel(x1+1, y1, true)
		y1--
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) Circle(x0, y0, r int) {
	x0, y0 = p.scale(x0, y0)
	r = p.sscale(r)
	x := -r
	y := 0
	e := 2 - 2*r
	for {
		p.setPixel(x0-x, y0+y, true)
		p.setPixel(x0-y, y0-x, true)
		p.setPixel(x0+x, y0-y, true)
		p.setPixel(x0+y, y0+x, true)
		r = e
		if r > x {
			x++
			e += x*2 + 1
		}
		if r <= y {
			y++
			e += y*2 + 1
		}
		if x >= 0 {
			break
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDing) Line(x0, y0, x1, y1 int) {
	x0, y0 = p.scale(x0, y0)
	x1, y1 = p.scale(x1, y1)
	var sx, sy int

	dx := abs(x1 - x0)
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}

	dy := -abs(y1 - y0)
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	e1 := dx + dy
	e2 := 0

	for {
		p.setPixel(x0, y0, true)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 = 2 * e1
		if e2 >= dy {
			e1 += dy
			x0 += sx
		}
		if e2 <= dx {
			e1 += dx
			y0 += sy
		}
	}
}
