package pixelding

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/bits"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const MaxX = 4000
const MaxY = 2000
const DefStep = 15

const OutOfBoundsError = "out of bounds"
const AlreadySetError = "already set"
const DimensionError = "dimension error"

const RegSplitter = "[MmLlHhVvZzCcSsQqTtAa]|[+-]?\\d+\\.\\d+|[+-]?\\d+|[+-]?\\.\\d+"

type PixelDING struct {
	init         bool
	matrix       [][]bool
	sizeX, sizeY int
	msteps       int
	aspectX      int
	aspectY      int
	scalef       float64
	debug        bool
	invert       bool
	toggle       bool
	render       int
	LastError    error
	buffer       []string
	fonts        map[string]*PixelFont
	stamps       map[string]*PixelStamp
}

type PixelStamp struct {
	prepared bool     `json:"-"`
	Len      int      `json:"len"`
	Data     []uint64 `json:"data"`
}

type PixelFont struct {
	prepared bool              `json:"-"`
	Chars    map[int]PixelChar `json:"chars"`
}

type PixelChar struct {
	OffsetX int      `json:"offsetX"`
	OffestY int      `json:"offsetY"`
	SizeX   int      `json:"sizeX"`
	SizeY   int      `json:"sizeY"`
	Len     int      `json:"len"`
	GN      int      `json:"gn"`
	GA      int      `json:"ga"`
	Data    []uint64 `json:"data"`
}

//----------------------------------------------------------------------------------------------------------------------
func New(dimensions ...int) PixelDING {
	x := PixelDING{}
	if len(dimensions) > 1 {
		x.sizeX = dimensions[0]
		x.sizeY = dimensions[1]
		x.init = true
	}
	x.SetStep(0)
	x.fonts = make(map[string]*PixelFont)
	x.stamps = make(map[string]*PixelStamp)
	x.AddFont("__std", x.LoadStdFont())
	x.AddStamp("__std", x.LoadStdStamp())
	x.LastError = x.Dimensions(x.sizeX, x.sizeY)
	return x
}

//----------------------------------------------------------------------------------------------------------------------
func maxUint64(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}

//----------------------------------------------------------------------------------------------------------------------
func minUint64(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}

//----------------------------------------------------------------------------------------------------------------------
func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

//----------------------------------------------------------------------------------------------------------------------
func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

//----------------------------------------------------------------------------------------------------------------------
func leftBound(x []uint64) ([]uint64, int) {
	var max uint64
	var y []uint64
	for _, u := range x {
		max = maxUint64(max, u)
	}
	c := bits.LeadingZeros64(max)
	for _, u := range x {
		y = append(y, u<<c)
	}
	min := 128
	for _, u := range y {
		min = minInt(min, bits.TrailingZeros64(u))
	}
	return y, min
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) X() int {
	return p.sizeX
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Y() int {
	return p.sizeY
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) SaveFont(name string, font *PixelFont, perm os.FileMode) error {

	buf, err := json.Marshal(font)
	err = ioutil.WriteFile(name, buf, perm)
	if err != nil {
		p.LastError = err
		return err
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) LoadFont(name string) *PixelFont {
	x := PixelFont{}
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		p.LastError = err
		return nil
	}
	err = json.Unmarshal(buf, &x)
	if err != nil {
		p.LastError = err
		return nil
	}
	return &x
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) SaveStamp(name string, stamp *PixelStamp, perm os.FileMode) error {
	buf, err := json.Marshal(stamp)
	err = ioutil.WriteFile(name, buf, perm)
	if err != nil {
		p.LastError = err
		return err
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) LoadStamp(name string) *PixelStamp {
	x := PixelStamp{}
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		p.LastError = err
		return nil
	}
	err = json.Unmarshal(buf, &x)
	if err != nil {
		p.LastError = err
		return nil
	}
	return &x
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) FontPrint(font *PixelFont, x, y int, text string, set bool) {
	ls := 0
	sx := x
	sy := y
	v := 0
	//_, ok := p.Fonts[font]
	//if !ok {
	//	return
	//}
	for _, z := range text {
		v = 0
		if ls != 0 && font.Chars[int(z)].GA == ls {
			v = -1
		}
		p.fontStamp(sx+v, sy, font.Chars[int(z)].Data, set)
		sx = sx + font.Chars[int(z)].SizeX + 1 + v
		ls = font.Chars[int(z)].GN
	}
}

//----------------------------------------------------------------------------------------------------------------------
func prepareFont(x PixelFont) PixelFont {
	var max uint64
	for i, char := range x.Chars {
		ch := char
		c := 0
		max = 0
		if char.SizeX != 0 {
			continue
		}
		for _, datum := range char.Data {
			max = maxUint64(max, uint64(bits.Len64(datum)))
			c++
		}
		ch.SizeX = int(max)
		ch.SizeY = c
		ch.Data, ch.Len = leftBound(char.Data)
		x.Chars[i] = ch
	}
	return x
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) AddFont(name string, font *PixelFont) {
	p.fonts[name] = font
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) GetFont(name string) *PixelFont {
	return p.fonts[name]
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) AddStamp(name string, stamp *PixelStamp) {
	p.stamps[name] = stamp
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) GetStamp(name string) *PixelStamp {
	return p.stamps[name]
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) FontInfo(name string) {
	for i, char := range p.fonts[name].Chars {
		fmt.Println("I", i, "X:", char.SizeX, "Y:", char.SizeY)
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) SetStep(x int) {
	if x < 1 || x > 50 {
		p.msteps = DefStep
	} else {
		p.msteps = x
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Toggle(b bool) {
	p.toggle = b
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Invert(b bool) {
	p.invert = b
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Debug(b bool) {
	p.debug = b
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Scale(s float64) {
	p.scalef = s
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Aspect(x, y int) {
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
func (s *PixelStamp) X() int {
	if !s.prepared {
		s.Data, s.Len = leftBound(s.Data)
		s.prepared = true
	}
	return 64 - s.Len
}

//----------------------------------------------------------------------------------------------------------------------
func (s *PixelStamp) Y() int {
	return len(s.Data)
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Stamp(x, y int, stamp *PixelStamp, set bool, st bool) {
	var j int
	if !stamp.prepared {
		stamp.Data, stamp.Len = leftBound(stamp.Data)
		stamp.prepared = true
	}
	for i, v := range stamp.Data {
		j = 0
		for xx := uint64(0x8000000000000000); xx > 0; xx = xx >> 1 {
			if v&xx != 0 {
				p.setPixel(x+j, y+i, set)
			} else {
				if st {
					p.setPixel(x+j, y+i, !set)
				}
			}
			j++
			if j >= 64-stamp.Len {
				break
			}
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) fontStamp(x, y int, stamp []uint64, set bool) {
	var j int
	for i, v := range stamp {
		j = 0
		for xx := uint64(0x8000000000000000); xx > 0; xx = xx >> 1 {
			if v&xx != 0 {
				p.setPixel(x+j, y+i, set)
			}
			j++
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Display() {
	for _, b := range p.buffer {
		fmt.Println(b)
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) RenderSmallest() []string {
	var mix, max, miy, may int
	mix = 0xFFFF
	max = 0
	miy = 0xFFFF
	may = 0
	for y := 0; y < p.sizeY-1; y++ {
		for x := 0; x < p.sizeX-1; x++ {
			if p.matrix[y][x] {
				mix = minInt(mix, x)
				miy = minInt(miy, y)
				max = maxInt(max, x)
				may = maxInt(may, y)
			}
		}
	}
	mix--
	max++
	miy--
	may++
	if mix < 0 {
		mix = 0
	}
	if miy < 0 {
		miy = 0
	}
	if max > p.sizeX {
		max = p.sizeX
	}
	if may > p.sizeY {
		may = p.sizeY
	}
	return p.RenderXY(mix, miy, max, may)
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) RenderXY(x1, y1, x2, y2 int) []string {
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
	p.buffer = []string{}
	lo := ""
	cmp := true
	if p.invert {
		cmp = !cmp
	}
	//xtoggle :=0
	//ytoggle := 0

	for y := y1; y < y2; y = y + 2 - p.aspectY {
		lo = ""
		for x := x1; x < x2; x = x + 2 - p.aspectX { // = sizeX + 2 {
			bit := 0
			if p.getPixel(x, y) == cmp {
				bit += 8
			}
			if p.getPixel(x+1-p.aspectX, y) == cmp {
				bit += 4
			}
			if p.getPixel(x, y+1-p.aspectY) == cmp {
				bit += 2
			}
			if p.getPixel(x+1-p.aspectX, y+1-p.aspectY) == cmp {
				bit += 1
			}
			lo = lo + cox[bit]
		}
		p.buffer = append(p.buffer, lo)
	}
	return p.buffer
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Render() []string {
	return p.RenderXY(0, 0, p.sizeX, p.sizeY)
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) BufferAnalyse() {
	for i, s := range p.buffer {
		fmt.Println(i, len(s))
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Clear() {
	p.matrix = make([][]bool, p.sizeY)
	for i := range p.matrix {
		p.matrix[i] = make([]bool, p.sizeX)
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) check(x, y int) bool {
	if x < 0 || x > p.sizeX-1 || y < 0 || y > p.sizeY-1 {
		return false
	}
	return true
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) scale(x, y int) (int, int) {
	if p.scalef != 0.0 {
		x = int(float64(x) * p.scalef)
		y = int(float64(y) * p.scalef)
	}
	return x, y
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) sscale(x int) int {
	if p.scalef != 0.0 {
		x = int(float64(x) * p.scalef)
	}
	return x
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Dimensions(x, y int) error {
	/*
		if p.init {
			p.LastError = errors.New(AlreadySetError)
			return p.LastError
		}
	*/
	if x < 1 || y < 1 {
		p.LastError = errors.New(DimensionError)
		return p.LastError
	}
	if x > MaxX || y > MaxY {
		p.LastError = errors.New(DimensionError)
		return p.LastError
	}
	p.matrix = make([][]bool, y)

	for i := range p.matrix {
		p.matrix[i] = make([]bool, x)
	}
	p.sizeX = x
	p.sizeY = y
	p.init = true
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) getPixel(x, y int) bool {
	// sizeX, sizeY = p.scale(sizeX, sizeY)
	if !p.check(x, y) {
		return false
	}
	return p.matrix[y][x]
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) setPixel(x, y int, b bool) {
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
func (p *PixelDING) GetPixel(x, y int) bool {
	x, y = p.scale(x, y)
	/*	if !p.check(x, y) {
			return false
		}
	*/
	return p.getPixel(x, y)
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Pixel(x, y int, set bool) {
	x, y = p.scale(x, y)
	/*	if !p.check(x, y) {
			return
		}
	*/
	p.setPixel(x, y, set)
}

func swap(a, b int) (int, int) {
	return b, a
}

//----------------------------------------------------------------------------------------------------------------------
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) TextFrame(x1, y1, x2, y2 int, l string, bitmask int, scale ...bool) {
	noLineH := false
	noLineV := false
	sx := strings.Split(l, "")
	h1 := ""
	h2 := ""

	if x1 > x2 {
		x1, x2 = swap(x1, x2)
	}
	if y1 > y2 {
		y1, y2 = swap(y1, y2)
	}

	if x2-x1 < 1 || y2-y1 < 1 {
		return
	}
	if x2-x1 < 2 {
		noLineH = true
	}
	if y2-y1 < 2 {
		noLineV = true
	}

	if noLineV == false {
		for i := y1 + 1; i < y2; i++ {
			if bitmask&(1<<5) != 0 {
				if len(scale) > 0 {
					p.Text(x1, i, sx[3], scale[0])
				} else {
					p.Text(x1, i, sx[3])
				}
			}
			if bitmask&(1<<3) != 0 {
				if len(scale) > 0 {
					p.Text(x2, i, sx[5], scale[0])
				} else {
					p.Text(x2, i, sx[5])
				}
			}
		}
	}

	if noLineH == false {
		hs := x2 - x1 - 1
		if len(scale) > 0 {
			if scale[0] == true && p.aspectX==0 {
				hs = (x2 - x1) / 2
			}
		}
		if bitmask&(1<<7) != 0 {
			h1 = strings.Repeat(sx[1], hs)
			if len(scale) > 0 {
				p.Text(x1+1, y1, h1, scale[0])
			} else {
				p.Text(x1+1, y1, h1)
			}
		}
		if bitmask&(1<<1) != 0 {
			h2 = strings.Repeat(sx[7], hs)
			if len(scale) > 0 {
				p.Text(x1+1, y2, h2, scale[0])
			} else {
				p.Text(x1+1, y2, h2)
			}
		}
	}

	if len(scale) > 0 {
		if bitmask&(1<<8) != 0 {
			p.Text(x1, y1, sx[0], scale[0])
		}
		if bitmask&(1<<6) != 0 {
			p.Text(x2, y1, sx[2], scale[0])
		}
		if bitmask&(1<<2) != 0 {
			p.Text(x1, y2, sx[6], scale[0])
		}
		if bitmask&(1<<0) != 0 {
			p.Text(x2, y2, sx[8], scale[0])
		}
	} else {
		if bitmask&(1<<8) != 0 {
			p.Text(x1, y1, sx[0])
		}
		if bitmask&(1<<6) != 0 {
			p.Text(x2, y1, sx[2])
		}
		if bitmask&(1<<2) != 0 {
			p.Text(x1, y2, sx[6])
		}
		if bitmask&(1<<0) != 0 {
			p.Text(x2, y2, sx[8])
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) TextLineH(x1, y1, x2, y2 int, l string, set ...bool) {
	sx := strings.Split(l, "")
	hs := ""
	//lc := strings.Split(l,"")
	if x1 > x2 {
		x1, x2 = swap(x1, x2)
		y1, y2 = swap(y1, y2)
	}
	hs = strings.Repeat(sx[1], x2-x1)
	if len(set) > 0 {
		p.Text(x1, y1, hs, set[0])
	} else {
		p.Text(x1, y1, hs)
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) TextLineV(x1, y1, x2, y2 int, l string, set ...bool) {
	sx := strings.Split(l, "")
	//lc := strings.Split(l,"")
	if y1 > y2 {
		x1, x2 = swap(x1, x2)
		y1, y2 = swap(y1, y2)
	}
	for i := y1; i < y2; i++ {
		if len(set) > 0 {
			p.Text(x1, i, sx[3], set[0])
		} else {
			p.Text(x1, i, sx[3])
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Text(x, y int, text string, scale ...bool) {
	if len(scale) > 0 {
		if scale[0] == true {
			x, y = p.scale(x, y)
			if p.aspectX == 0 {
				x = x / 2
			}
			if p.aspectY == 0 {
				y = y / 2
			}
		}
	}
	var n string
	if len(p.buffer) < y+1 {
		return //Out of bounds
	}

	if x < 0 || y < 0 {
		return //Out of bounds
	}

	sx := strings.Split(text, "")
	l := len(sx)

	s := strings.Split(p.buffer[y], "")
	sl := len(s)

	if x > sl {
		return //Out of bounds
	}

	//	fmt.Println("sl", sl, "l", l)
	cs := 1
	for i := 0; i < x; i++ {
		n = n + s[i]
		cs++
	}
	//fmt.Println(cs)
	for _, t := range sx {
		if cs > sl {
			break
		}
		n = n + t
		cs++
	}
	//fmt.Println(cs)
	for i := x + l; i < sl; i++ {
		if cs > sl {
			break
		}
		n = n + s[i]
		cs++
	}
	//fmt.Println(cs)
	//
	/*
		for i, scale := range s {
			fmt.Print("[",i,"]",scale,string(scale),"-")
		}
		fmt.Println(len(s))
	*/
	p.buffer[y] = n

}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) HBar(size int) string {
	sx := strings.Split(HBar, "")
	e := size / 8
	r := size % 8
	s := strings.Repeat(sx[0], e)
	if r > 0 {
		s = s + sx[r]
	}
	return s
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) SVGPath(xo, yo float64, s string, set bool, fscale ...float64) {
	var x, y float64
	var lx, ly float64
	var ix, iy float64
	var c1x, c1y, c2x, c2y float64
	var lcx, lcy float64
	var err error
	sw := 0
	cmd := ""
	scale := 1.0
	for _, f := range fscale {
		scale = f
	}

	re := regexp.MustCompile(RegSplitter)
	ss := re.FindAllString(s, -1)

	it := 0
	for _, ssSub := range ss {

		if p.debug {
			fmt.Println(":::", ssSub, len(ssSub))
		}

		it++

		_, err = strconv.ParseFloat(ssSub, 64)

		if err != nil { //Not numeric ...
			switch ssSub {
			case "M", "m", "L", "l", "V", "v", "H", "h":
				cmd = ssSub
				sw = 1
			case "C", "c":
				cmd = ssSub
				sw = 1
			case "S", "s":
				cmd = ssSub
				sw = 1
			case "Q", "q":
				cmd = ssSub
				sw = 1
			case "T", "t":
				cmd = ssSub
				sw = 1
			case "Z", "z":
				cmd = ssSub
				x = ix
				y = iy
				sw = 9
			case "A", "a":
				cmd = ssSub
				sw = 1
			default:
				if p.debug {
					fmt.Println("Unknown :", ssSub)
				}
				p.LastError = errors.New("Unknown Command")
			}
		} else {
			if sw == 0 {
				sw = 1
			}
			switch {

			case cmd == "M" && sw == 1:
				x, _ = strconv.ParseFloat(ssSub, 64)
				ix = x
				sw = 2
			case cmd == "M" && sw == 2:
				y, _ = strconv.ParseFloat(ssSub, 64)
				iy = y
				sw = 9

			case cmd == "m" && sw == 1:
				x, _ = strconv.ParseFloat(ssSub, 64)
				x = lx + x
				ix = x
				sw = 2
			case cmd == "m" && sw == 2:
				y, _ = strconv.ParseFloat(ssSub, 64)
				y = ly + y
				iy = y
				sw = 9

			case cmd == "L" && sw == 1:
				x, _ = strconv.ParseFloat(ssSub, 64)
				sw = 2
			case cmd == "L" && sw == 2:
				y, _ = strconv.ParseFloat(ssSub, 64)
				sw = 9

			case cmd == "l" && sw == 1:
				x, _ = strconv.ParseFloat(ssSub, 64)
				x = lx + x
				sw = 2
			case cmd == "l" && sw == 2:
				y, _ = strconv.ParseFloat(ssSub, 64)
				y = ly + y
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

			case cmd == "C" && sw == 1:
				c1x, _ = strconv.ParseFloat(ssSub, 64)
				sw = 2
			case cmd == "C" && sw == 2:
				c1y, _ = strconv.ParseFloat(ssSub, 64)
				sw = 3
			case cmd == "C" && sw == 3:
				c2x, _ = strconv.ParseFloat(ssSub, 64)
				sw = 4
			case cmd == "C" && sw == 4:
				c2y, _ = strconv.ParseFloat(ssSub, 64)
				sw = 5
			case cmd == "C" && sw == 5:
				x, _ = strconv.ParseFloat(ssSub, 64)
				sw = 6
			case cmd == "C" && sw == 6:
				y, _ = strconv.ParseFloat(ssSub, 64)
				sw = 9

			case cmd == "S" && sw == 1:
				c2x, _ = strconv.ParseFloat(ssSub, 64)
				sw = 2
			case cmd == "S" && sw == 2:
				c2y, _ = strconv.ParseFloat(ssSub, 64)
				sw = 3
			case cmd == "S" && sw == 3:
				x, _ = strconv.ParseFloat(ssSub, 64)
				sw = 4
			case cmd == "S" && sw == 4:
				y, _ = strconv.ParseFloat(ssSub, 64)
				sw = 9

			case cmd == "s" && sw == 1:
				c2x, _ = strconv.ParseFloat(ssSub, 64)
				c2x = lx + c2x
				sw = 2
			case cmd == "s" && sw == 2:
				c2y, _ = strconv.ParseFloat(ssSub, 64)
				c2y = ly + c2y
				sw = 3
			case cmd == "s" && sw == 3:
				x, _ = strconv.ParseFloat(ssSub, 64)
				x = lx + x
				sw = 4
			case cmd == "s" && sw == 4:
				y, _ = strconv.ParseFloat(ssSub, 64)
				y = ly + y
				sw = 9

			case cmd == "Q" && sw == 1:
				c1x, _ = strconv.ParseFloat(ssSub, 64)
				sw = 2
			case cmd == "Q" && sw == 2:
				c1y, _ = strconv.ParseFloat(ssSub, 64)
				sw = 3
			case cmd == "Q" && sw == 3:
				x, _ = strconv.ParseFloat(ssSub, 64)
				sw = 4
			case cmd == "Q" && sw == 4:
				y, _ = strconv.ParseFloat(ssSub, 64)
				sw = 9

			case cmd == "q" && sw == 1:
				c1x, _ = strconv.ParseFloat(ssSub, 64)
				c1x = lx + c1x
				sw = 2
			case cmd == "q" && sw == 2:
				c1y, _ = strconv.ParseFloat(ssSub, 64)
				c1y = ly + c1y
				sw = 3
			case cmd == "q" && sw == 3:
				x, _ = strconv.ParseFloat(ssSub, 64)
				x = lx + x
				sw = 4
			case cmd == "q" && sw == 4:
				y, _ = strconv.ParseFloat(ssSub, 64)
				y = ly + y
				sw = 9

			case cmd == "T" && sw == 1:
				x, _ = strconv.ParseFloat(ssSub, 64)
				sw = 2
			case cmd == "T" && sw == 2:
				y, _ = strconv.ParseFloat(ssSub, 64)
				sw = 9

			case cmd == "t" && sw == 1:
				x, _ = strconv.ParseFloat(ssSub, 64)
				x = lx + x
				sw = 2
			case cmd == "t" && sw == 2:
				y, _ = strconv.ParseFloat(ssSub, 64)
				y = ly + y
				sw = 9

			case cmd == "A" && sw == 1:
				x, _ = strconv.ParseFloat(ssSub, 64)
				sw = 2
			case cmd == "A" && sw == 2:
				y, _ = strconv.ParseFloat(ssSub, 64)
				sw = 9

			case cmd == "a" && sw == 1:
				x, _ = strconv.ParseFloat(ssSub, 64)
				x = lx + x
				sw = 2
			case cmd == "a" && sw == 2:
				y, _ = strconv.ParseFloat(ssSub, 64)
				y = ly + y
				sw = 9

			}
		}
		if sw == 9 {
			//fmt.Println(it, "(", cmd, ") lx", lx, "ly", ly, "sizeX", sizeX, "sizeY", sizeY)
			//fmt.Println("--------------------------------")
			switch cmd {
			case "L", "l":
				p.Line(int((lx+xo)*scale), int((ly+yo)*scale), int((x+xo)*scale), int((y+yo)*scale), set)
				lx = x
				ly = y
			case "M", "m":
				lx = x
				ly = y
				lcx = x //Set the last Control Point to
				lcy = y
			case "H", "h":
				p.Line(int((lx+xo)*scale), int((ly+yo)*scale), int((x+xo)*scale), int((y+yo)*scale), set)
				lx = x
			case "V", "v":
				p.Line(int((lx+xo)*scale), int((ly+yo)*scale), int((x+xo)*scale), int((y+yo)*scale), set)
				ly = y

			case "C", "c":
				p.CBezier(int((lx+xo)*scale), int((ly+yo)*scale), int((c1x+xo)*scale), int((c1y+yo)*scale), int((c2x+xo)*scale), int((c2y+yo)*scale), int((x+xo)*scale), int((y+yo)*scale), set)
				//p.Line(int(lx), int(ly), int(sizeX), int(sizeY))
				lx = x
				ly = y
				lcx = c2x
				lcy = c2y

			case "S", "s":
				c1x = lx - -(lx - lcx)
				c1y = ly - -(ly - lcy)

				p.CBezier(int((lx+xo)*scale), int((ly+yo)*scale), int((c1x+xo)*scale), int((c1y+yo)*scale), int((c2x+xo)*scale), int((c2y+yo)*scale), int((x+xo)*scale), int((y+yo)*scale), set)
				//p.Line(int(lx), int(ly), int(sizeX), int(sizeY))
				lx = x
				ly = y
				lcx = c2x
				lcy = c2y

			case "Q", "q":

				//c2x = sizeX- -(sizeX-c1x)
				//c2y = sizeY- -(sizeY-c1y)
				p.QBezier(int((lx+xo)*scale), int((ly+yo)*scale), int((c1x+xo)*scale), int((c1y+yo)*scale), int((x+xo)*scale), int((y+yo)*scale), set)
				//p.Line(int(lx), int(ly), int(sizeX), int(sizeY))
				lx = x
				ly = y
				lcx = c1x
				lcy = c1y

			case "T", "t":
				c1x = lx - -(lx - lcx)
				c1y = ly - -(ly - lcy)

				p.QBezier(int((lx+xo)*scale), int((ly+yo)*scale), int((c1x+xo)*scale), int((c1y+yo)*scale), int((x+xo)*scale), int((y+yo)*scale), set)
				//p.Line(int(lx), int(ly), int(sizeX), int(sizeY))
				lx = x
				ly = y
				lcx = c1x
				lcy = c1y

			case "Z", "z":
				p.Line(int((lx+xo)*scale), int((ly+yo)*scale), int((x+xo)*scale), int((y+yo)*scale), set)
				lx = x
				ly = y

			case "A", "a":
				//simulate by just setting new ends
				lx = x
				ly = y

			}
			sw = 0
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) QBezier(x1, y1, x2, y2, x3, y3 int, set bool) {
	var px, py [50 + 1]int
	x1, y1 = p.scale(x1, y1)
	x2, y2 = p.scale(x2, y2)
	x3, y3 = p.scale(x3, y3)
	fx1, fy1 := float64(x1), float64(y1)
	fx2, fy2 := float64(x2), float64(y2)
	fx3, fy3 := float64(x3), float64(y3)
	for i := 0; i <= p.msteps; i++ { //range px {
		//for i := range px {
		c := float64(i) / float64(p.msteps)
		a := 1 - c
		a, b, c := a*a, 2*c*a, c*c
		px[i] = int(a*fx1 + b*fx2 + c*fx3)
		py[i] = int(a*fy1 + b*fy2 + c*fy3)
	}
	x0, y0 := px[0], py[0]
	for i := 1; i <= p.msteps; i++ {
		x1, y1 := px[i], py[i]
		p.Line(x0, y0, x1, y1, set)
		x0, y0 = x1, y1
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) CBezier(x1, y1, x2, y2, x3, y3, x4, y4 int, set bool) {
	var px, py [50 + 1]int
	x1, y1 = p.scale(x1, y1)
	x2, y2 = p.scale(x2, y2)
	x3, y3 = p.scale(x3, y3)
	x4, y4 = p.scale(x4, y4)
	fx1, fy1 := float64(x1), float64(y1)
	fx2, fy2 := float64(x2), float64(y2)
	fx3, fy3 := float64(x3), float64(y3)
	fx4, fy4 := float64(x4), float64(y4)
	for i := 0; i <= p.msteps; i++ { //range px {
		//for i := range px {
		d := float64(i) / float64(p.msteps)
		a := 1 - d
		b, c := a*a, d*d
		a, b, c, d = a*b, 3*b*d, 3*a*c, c*d
		px[i] = int(a*fx1 + b*fx2 + c*fx3 + d*fx4)
		py[i] = int(a*fy1 + b*fy2 + c*fy3 + d*fy4)
	}
	x0, y0 := px[0], py[0]
	for i := 1; i <= p.msteps; i++ {
		x1, y1 := px[i], py[i]
		p.Line(x0, y0, x1, y1, set)
		x0, y0 = x1, y1
	}
}

//----------------------------------------------------------------------------------------------------------------------
/*
func (p *PixelDING) BezierAlt(x0, y0, x1, y1, x2, y2, x3, y3 int) {

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
*/

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Rectangle(x0, y0, x1, y1 int, set bool, fill bool) {
	x0, y0 = p.scale(x0, y0)
	x1, y1 = p.scale(x1, y1)
	for i := x0; i <= x1; i++ {
		p.setPixel(i, y0, set)
	}
	for i := y0; i < y1; i++ {
		if fill {
			for j := x0; j <= x1; j++ {
				p.setPixel(j, i, set)
			}
		} else {
			p.setPixel(x0, i, set)
			p.setPixel(x1, i, set)
		}
	}
	for i := x0; i <= x1; i++ {
		p.setPixel(i, y1, set)
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) floodFill(x, y int, prevC, newC bool) {
	if x < 0 || x >= p.sizeX || y < 0 || y >= p.sizeY {
		return
	}
	if prevC != p.getPixel(x, y) {
		return
	}
	// Replace the color at (sizeX, sizeY)
	p.setPixel(x, y, newC)
	// Recur for north, east, south and west
	p.floodFill(x+1, y, prevC, newC)
	p.floodFill(x-1, y, prevC, newC)
	p.floodFill(x, y+1, prevC, newC)
	p.floodFill(x, y-1, prevC, newC)
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Fill(x, y int, newC bool) {
	x, y = p.scale(x, y)
	prevC := p.getPixel(x, y)
	if prevC == newC {
		return
	}
	p.floodFill(x, y, prevC, newC)
}

func toRadian(angle int) float64 {
	return float64(angle) * (math.Pi / 180.0)
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) DotArc(x0, y0, r int, a1, a2, step int, set bool) { //wieso

	x0, y0 = p.scale(x0, y0)
	r = p.sscale(r)

	if a1 == a2 {
		return
	}

	if a1 < 0 || a2 < 0 || a1 > 360 || a2 > 360 {
		return
	}
	if a1 > a2 {
		a2 += 360
	}

	for {
		if a1 >= a2 {
			break
		}

		xo := int(math.Round(float64(r) * math.Sin(toRadian(a1%360))))
		yo := int(math.Round(float64(r) * math.Cos(toRadian(a1%360))))

		p.setPixel(x0+xo, y0-yo, set)

		a1 += step

	}

}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) LineArc(x0, y0, r int, a1, a2, step int, set bool) { //wieso

	x0, y0 = p.scale(x0, y0)
	r = p.sscale(r)

	fromx := 0
	fromy := 0
	firstfrom := true

	if a1 == a2 {
		return
	}

	if a1 < 0 || a2 < 0 || a1 > 360 || a2 > 360 {
		return
	}
	if a1 > a2 {
		a2 += 360
	}

	for {
		if a1 >= a2 {
			break
		}

		xo := int(math.Round(float64(r) * math.Sin(toRadian(a1%360))))
		yo := int(math.Round(float64(r) * math.Cos(toRadian(a1%360))))

		if !firstfrom {
			p.Line(x0+fromx, y0-fromy, x0+xo, y0-yo, set)

		}
		firstfrom = false
		fromx = xo
		fromy = yo

		a1 += step

	}

	xo := int(math.Round(float64(r) * math.Sin(toRadian(a2%360))))
	yo := int(math.Round(float64(r) * math.Cos(toRadian(a2%360))))
	p.Line(x0+fromx, y0-fromy, x0+xo, y0-yo, set)
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) LineRadius(x0, y0, r1, r2, a1 int, set bool) {
	x0, y0 = p.scale(x0, y0)
	r1 = p.sscale(r1)
	r2 = p.sscale(r2)

	x1 := int(math.Round(float64(r1) * math.Sin(toRadian(a1))))
	y1 := int(math.Round(float64(r1) * math.Cos(toRadian(a1))))
	x2 := int(math.Round(float64(r2) * math.Sin(toRadian(a1))))
	y2 := int(math.Round(float64(r2) * math.Cos(toRadian(a1))))

	p.Line(x0+x1, y0-y1, x0+x2, y0-y2, set)
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) EllipseRect(x0, y0, x1, y1 int, set bool) {
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
		p.setPixel(x1, y0, set)
		p.setPixel(x0, y0, set)
		p.setPixel(x0, y1, set)
		p.setPixel(x1, y1, set)
		e2 = 2 * e
		if e2 >= dx {
			x0++
			x1--
			dx += b1
			e += dx
		} /* sizeX step */
		if e2 <= dy {
			y0++
			y1--
			dy += a
			e += dy
		} /* sizeY step */
		if x0 > x1 {
			break
		}
	}

	for {
		if y0-y1 >= b {
			break
		}
		p.setPixel(x0-1, y0, set)
		p.setPixel(x1+1, y0, set)
		y0++
		p.setPixel(x0-1, y1, set)
		p.setPixel(x1+1, y1, set)
		y1--
	}
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Circle(x0, y0, r int, set bool) {
	x0, y0 = p.scale(x0, y0)
	r = p.sscale(r)
	x := -r
	y := 0
	e := 2 - 2*r
	for {
		p.setPixel(x0-x, y0+y, set)
		p.setPixel(x0-y, y0-x, set)
		p.setPixel(x0+x, y0-y, set)
		p.setPixel(x0+y, y0+x, set)
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
func (p *PixelDING) Line(x0, y0, x1, y1 int, set bool) {
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
		p.setPixel(x0, y0, set)
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

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) DotLine(x0, y0, x1, y1 int, set bool, pattern ...uint8) {
	var pat uint8
	if len(pattern) == 0 {
		pat = Dot1x1Pattern
	} else {
		pat = pattern[0]
	}

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
		if (pat & 0x01)==0x01 {
			p.setPixel(x0, y0, set)
			pat = (pat >> uint8(1)) + 0x80
		} else {
			pat = pat >> 1
		}

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
