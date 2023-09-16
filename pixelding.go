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

const ModeTrueColor = 3
const ModeNoColor = 0
const Mode16Color = 1
const ModePaletteColor = 2
const (
	ColorBlack = 30 + iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorDefault
	ColorReset
)

const ESCHome = "\033[0;0H"
const ESCClear = "\033[J"

const OutOfBoundsError = "out of bounds"
const AlreadySetError = "already set"
const DimensionError = "dimension error"
const ColormodeError = "colormode error"

const RegSplitter = "[MmLlHhVvZzCcSsQqTtAa]|[+-]?\\d+\\.\\d+|[+-]?\\d+|[+-]?\\.\\d+"

type PixelDING struct {
	init           bool
	matrix         [][]uint32
	tmatrix        [][]rune
	sizeX, sizeY   int
	clipsx, clipsy int
	clipex, clipey int
	clipping       bool
	msteps         int
	aspectX        int
	aspectY        int
	faspectX       int
	faspectY       int
	scalef         float64
	debug          bool
	invert         bool
	toggle         bool
	acolor         uint32
	bcolor         uint32
	colorrender    int
	LastError      error
	buffer         []string
	fonts          map[string]*PixelFont
	stamps         map[string]*PixelStamp
	pics           map[string]*PixelPicture
}

type PixelStamp struct {
	Prepared bool     `json:"prepared"`
	Len      int      `json:"len"`
	Data     []uint64 `json:"data"`
}

type PixelPicture struct {
	Mode     int      `json:"mode"`
	ColorKey uint32   `json:"colorKey"`
	SizeX    int      `json:"sizeX"`
	SizeY    int      `json:"sizeY"`
	SegX     int      `json:"segX"`
	SegY     int      `json:"segY"`
	Data     []uint32 `json:"data"`
}

type PixelFont struct {
	Prepared bool              `json:"prepared"`
	sizex    int               `json:"-"`
	sizey    int               `json:"-"`
	numchar  int               `json:"-"`
	Chars    map[int]PixelChar `json:"chars"`
}

type PixelFontInfo struct {
	MaxX  int
	MaxY  int
	Chars int
}

type PixelChar struct {
	OffsetX int      `json:"OffsetX"`
	OffsetY int      `json:"OffsetY"`
	SizeX   int      `json:"sizeX"`
	SizeY   int      `json:"sizeY"`
	Len     int      `json:"len"`
	GN      int      `json:"gn"`
	GA      int      `json:"ga"`
	Data    []uint64 `json:"data"`
}

// New create a new PixelDING with optional size parameter (x,y), and attaches the
// standard Font and Stamp, return pixelDING object
// ----------------------------------------------------------------------------------------------------------------------
func New(dimensions ...int) PixelDING {
	x := PixelDING{}
	if len(dimensions) > 1 {
		x.sizeX = dimensions[0]
		x.sizeY = dimensions[1]
		x.init = true
	}
	x.SetStep(0)
	x.acolor = 1
	x.bcolor = 0
	x.fonts = make(map[string]*PixelFont)
	x.stamps = make(map[string]*PixelStamp)
	x.pics = make(map[string]*PixelPicture)
	x.AddFont("__std", x.LoadStdFont())
	x.AddStamp("__std", x.LoadStdStamp())
	x.LastError = x.Dimensions(x.sizeX, x.sizeY)
	return x
}

// ----------------------------------------------------------------------------------------------------------------------
func maxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

// ----------------------------------------------------------------------------------------------------------------------
func minUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

// ----------------------------------------------------------------------------------------------------------------------
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ----------------------------------------------------------------------------------------------------------------------
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ----------------------------------------------------------------------------------------------------------------------
func leftBound(x []uint64, fixsize int) ([]uint64, int) {
	var max uint64
	var y []uint64
	var c int
	if fixsize == 0 {
		for _, u := range x {
			max = maxUint64(max, u)
		}
		c = bits.LeadingZeros64(max)
	} else {
		c = 64 - fixsize
	}
	for _, u := range x {
		y = append(y, u<<c)
	}
	min := 128
	for _, u := range y {
		min = minInt(min, bits.TrailingZeros64(u))
	}
	return y, min
}

// X returns pixelDING maximum X
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) X() int {
	return p.sizeX
}

// Y returns pixelDING maximum Y
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Y() int {
	return p.sizeY
}

// SaveFont save a font to disk
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) SaveFont(name string, font *PixelFont, permissions os.FileMode) error {

	buf, err := json.Marshal(font)
	err = ioutil.WriteFile(name, buf, permissions)
	if err != nil {
		p.LastError = err
		return err
	}
	return nil
}

// LoadFont load a font into pixelDING font object
// ----------------------------------------------------------------------------------------------------------------------
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

// SaveStamp save a stamp to disk
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) SaveStamp(name string, stamp *PixelStamp, permissions os.FileMode) error {
	buf, err := json.Marshal(stamp)
	err = ioutil.WriteFile(name, buf, permissions)
	if err != nil {
		p.LastError = err
		return err
	}
	return nil
}

// LoadStamp load a stamp into pixelDING stamp object
// ----------------------------------------------------------------------------------------------------------------------
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

// SavePicture Save a picture that was created ot converted into memory
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) SavePicture(name string, picture *PixelPicture, permissions os.FileMode) error {
	buf, err := json.Marshal(picture)
	err = ioutil.WriteFile(name, buf, permissions)
	if err != nil {
		p.LastError = err
		return err
	}
	return nil
}

// LoadPicture Load a picture into pixelDING picture object
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) LoadPicture(name string) *PixelPicture {
	x := PixelPicture{}
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

// FontPrint Print a text into pixelDING with given font at x,y
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) FontPrint(font *PixelFont, x0, y0 int, text string, set bool, param ...int) {
	ls := 0
	sx := x0
	sy := y0
	v := 0
	offset := 0
	if len(param) > 0 {
		offset = param[0]
	}
	//_, ok := p.Fonts[font]
	//if !ok {
	//	return
	//}
	for _, z := range text {
		v = 0
		if ls != 0 && font.Chars[int(z)].GA == ls {
			v = -1
		}
		p.fontStamp(sx+v, sy, font.Chars[int(z)].Data, set, p.faspectX, p.faspectY)
		sx = sx + font.Chars[int(z)].SizeX + 1 + v + offset
		if p.faspectX > 0 {
			sx = sx + font.Chars[int(z)].SizeX + 1 + v
		}
		ls = font.Chars[int(z)].GN
	}
}

// Prepare This is a compression option to reduce the saved size on disk
// ----------------------------------------------------------------------------------------------------------------------
func (f *PixelChar) Prepare() {
	c := 0
	var max uint64
	for _, datum := range f.Data {
		max = maxUint64(max, uint64(bits.Len64(datum)))
		c++
	}
	if f.SizeX == 0 {
		f.SizeX = int(max)
	}
	f.SizeY = c
	f.Data, f.Len = leftBound(f.Data, f.SizeX)
}

// AddChar adds a char to a pixelDING font object
// ----------------------------------------------------------------------------------------------------------------------
func (f *PixelFont) AddChar(ix int, char PixelChar) {
	f.Chars[ix] = char
}

// PrepareFont this is a compression option to reduce the saved size on disk
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) PrepareFont(x PixelFont) *PixelFont {
	var max uint64
	if x.Prepared {
		return &x
	}
	for i, char := range x.Chars {
		ch := char
		c := 0
		max = 0
		//if char.SizeX != 0 {
		//	continue
		//}
		for _, datum := range char.Data {
			max = maxUint64(max, uint64(bits.Len64(datum)))
			c++
		}
		if char.SizeX == 0 {
			ch.SizeX = int(max)
		}
		ch.SizeY = c
		ch.Data, ch.Len = leftBound(char.Data, ch.SizeX)
		x.Chars[i] = ch
	}
	x.Prepared = true
	return &x
}

// AddFont adds a font object to the pixelDING object, replaces existing one
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) AddFont(name string, font *PixelFont) {
	p.fonts[name] = font
}

// RemoveFont removes a font object from pixelDING object
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) RemoveFont(name string) {
	_, ok := p.fonts[name]
	if ok {
		delete(p.fonts, name)
	}
}

// GetFont returns pixelDING font object of existing font, nil if not exist
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) GetFont(name string) *PixelFont {
	_, ok := p.fonts[name]
	if ok {
		return p.fonts[name]
	}
	return nil
}

// AddStamp adds a stamp obect to the pixelDING object, replaces existing one
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) AddStamp(name string, stamp *PixelStamp) {
	p.stamps[name] = stamp
}

// RemoveStamp removes a stamp object from pixelDING object
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) RemoveStamp(name string) {
	_, ok := p.stamps[name]
	if ok {
		delete(p.stamps, name)
	}
}

// GetStamp returns stamp object from pixelDING object
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) GetStamp(name string) *PixelStamp {
	_, ok := p.stamps[name]
	if ok {
		return p.stamps[name]
	}
	return nil
}

// AddPicture adds a picture object to the pixelDING object
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) AddPicture(name string, picture *PixelPicture) {
	p.pics[name] = picture
}

// RemovePicture removes a picture from the PixelDING object
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) RemovePicture(name string) {
	_, ok := p.pics[name]
	if ok {
		delete(p.pics, name)
	}
}

// GetPicture returns a picture object from pixelDING or nil if not exist
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) GetPicture(name string) *PixelPicture {
	_, ok := p.pics[name]
	if ok {
		return p.pics[name]
	}
	return nil
}

// FontInfo returns a font info structure from font object
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) FontInfo(font *PixelFont) (*PixelFontInfo, error) {
	maxX := 0
	maxY := 0
	fi := PixelFontInfo{}
	if font == nil {
		return nil, nil
	}

	if !font.Prepared {
		return &fi, errors.New("font not Prepared")
	} else {
		if font.sizex != 0 || font.sizey != 0 {
			fi.MaxX = font.sizex
			fi.MaxY = font.sizey
			fi.Chars = font.numchar
			return &fi, nil
		}
	}
	for _, char := range font.Chars {
		maxX = maxInt(maxX, char.SizeX)
		maxY = maxInt(maxY, char.SizeY)
	}
	font.sizex = maxX
	font.sizex = maxY
	font.numchar = len(font.Chars)
	fi.MaxX = maxX
	fi.MaxY = maxY
	fi.Chars = len(font.Chars)
	return &fi, nil
}

// SetStep sets the maximum steps for some curves and other functions
// use a higher number (up to 50) if you need more quality on bigger curves
// reduce the steps if you need more performance or drawing smaller curves
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) SetStep(steps int) {
	if steps < 1 || steps > 50 {
		p.msteps = DefStep
	} else {
		p.msteps = steps
	}
}

// RGBMul experimental color manipulation to dim or brighten colors
// ----------------------------------------------------------------------------------------------------------------------
func RGBMul(color uint32, mod float64) uint32 {
	var r, g, b float64

	r = float64((color>>16)&0xff) * mod
	g = float64((color>>8)&0xff) * mod
	b = float64(color&0xff) * mod

	if r > 255.0 {
		r = 255.0
	}
	if g > 255.0 {
		g = 255.0
	}
	if b > 255.0 {
		b = 255.0
	}

	return uint32(r)<<16 + uint32(g)<<8 + uint32(b)
}

// RGB helper to construct a 32bit color value from R,G,B values
// ----------------------------------------------------------------------------------------------------------------------
func RGB(r, g, b uint8) uint32 {
	var c uint32
	c = uint32(r)
	c <<= 8
	c += uint32(g)
	c <<= 8
	c += uint32(b)
	return c
}

// Color set the desired color for drawing
// Color(a) sets foreground to a
// Color(a,b) sets foreground to a and background to b
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Color(c ...uint32) {
	if len(c) == 1 {
		p.acolor = c[0]
	} else {
		if len(c) == 2 {
			p.acolor = c[0]
			p.bcolor = c[1]
		}
	}
}

// Toggle experimental pixel toggle
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Toggle(b bool) {
	p.toggle = b
}

// ColorMode set the desired color mode
// need to be one of : ModeNoColor, ModePaletteColor, Mode16Color, ModeTrueColor
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) ColorMode(mode int) error {
	switch mode {
	case ModeNoColor, ModePaletteColor, Mode16Color, ModeTrueColor:
		p.colorrender = mode
		return nil
	default:
		return errors.New(ColormodeError)
	}
}

// Invert experimental pixel inverting
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Invert(b bool) {
	p.invert = b
}

// Debug switches some debug messages on (experimental)
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Debug(b bool) {
	p.debug = b
}

// Scale set a scalefactor (experimental)
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Scale(s float64) {
	p.scalef = s
}

// Aspect sets the aspect ratio for the rendering. On no Colormode the X dimension is
// doubled. To draw still "perfect" circles you need to set the x aspect ration to 1
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Aspect(x0, y0 int) {
	p.aspectX = 0
	p.aspectY = 0
	if x0 > 0 {
		p.aspectX = 1
	}
	if y0 > 0 {
		p.aspectY = 1
	}
}

// FontAspect doubles font size X and/or Y (experimental)
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) FontAspect(x0, y0 int) {
	p.faspectX = 0
	p.faspectY = 0
	if x0 > 0 {
		p.faspectX = 1
	}
	if y0 > 0 {
		p.faspectY = 1
	}
}

// X returns maximum X from stamp object
// ----------------------------------------------------------------------------------------------------------------------
func (s *PixelStamp) X() int {
	if !s.Prepared {
		s.Data, s.Len = leftBound(s.Data, 0)
		s.Prepared = true
	}
	return 64 - s.Len
}

// Y return maximum Y from stamp object
// ----------------------------------------------------------------------------------------------------------------------
func (s *PixelStamp) Y() int {
	return len(s.Data)
}

// Picture draws a picture obeject at the x,y coordinates given.
// if a segment is specified the function is painting only that segment
// The calculation of the segment is done only via the segment size of the picture
// SegX and SegY which need to be specified in the picture itself
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Picture(picture *PixelPicture, x0, y0 int, segment int) {
	xdehn := 0
	ydehn := 0
	xstart := 0
	ystart := 0
	if picture.SegX > 0 {
		xdehn = picture.SizeX / picture.SegX
	}
	if picture.SegY > 0 {
		ydehn = picture.SizeY / picture.SegY
	}

	if segment > xdehn*ydehn {
		segment = 1
	}

	if segment > 0 {
		//zero based calculation
		segment--
		xstart = segment % xdehn
		ystart = segment / xdehn

		ix := (xstart * picture.SegX) + (ystart * picture.SizeX * picture.SegY)
		jx := ix
		/*
			fmt.Println("PicX", picture.SizeX)
			fmt.Println("PicY", picture.SizeY)
			fmt.Println("SegX", picture.SegX)
			fmt.Println("SegY", picture.SegY)

			fmt.Println("Segment", segment)
			fmt.Println("xdehn", xdehn)
			fmt.Println("ydehn", ydehn)
			fmt.Println("xstart", xstart)
			fmt.Println("ystart", ystart)
			fmt.Println("ix", ix)
		*/
		for i := 0; i < picture.SegY; i++ {
			for j := 0; j < picture.SegX; j++ {
				p.PixelC(x0+j, y0+i, picture.Data[ix])
				ix++
			}
			ix = jx + picture.SizeX
			jx = ix
		}

	} else {

		ix := 0

		for i := 0; i < picture.SizeY; i++ {
			for j := 0; j < picture.SizeX; j++ {
				p.PixelC(x0+j, y0+i, picture.Data[ix])
				ix++
			}
		}
	}

}

// Stamp stamps a stamp object at the given x,y coordinates
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Stamp(stamp *PixelStamp, x0, y0 int, set bool, st bool) {
	var j int
	if !stamp.Prepared {
		stamp.Data, stamp.Len = leftBound(stamp.Data, 0)
		stamp.Prepared = true
	}
	for i, v := range stamp.Data {
		j = 0
		for xx := uint64(0x8000000000000000); xx > 0; xx = xx >> 1 {
			if v&xx != 0 {
				p.setPixel(x0+j, y0+i, set)
			} else {
				if st {
					p.setPixel(x0+j, y0+i, !set)
				}
			}
			j++
			if j >= 64-stamp.Len {
				break
			}
		}
	}
}

// fontStamp internal
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) fontStamp(x0, y0 int, stamp []uint64, set bool, ax, ay int) {
	var jx int
	var jy int
	jy = 0
	for _, v := range stamp {
		jx = 0
		for xx := uint64(0x8000000000000000); xx > 0; xx = xx >> 1 {
			if v&xx != 0 {
				p.setPixel(x0+jx, y0+jy, set)
				if ax > 0 {
					p.setPixel(x0+jx+1, y0+jy, set)
				}
				if ay > 0 {
					p.setPixel(x0+jx, y0+jy+1, set)
					if ax > 0 {
						p.setPixel(x0+jx+1, y0+jy+1, set)
					}
				}
			}
			jx = jx + 1 + ax
		}
		jy = jy + 1 + ay
	}
}

// Display prints the rendered display buffer to the console
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Display() {
	for _, b := range p.buffer {
		fmt.Println(b)
	}
}

// Display prints the rendered display buffer to the console with suffix for raw modus
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) DisplaySuff(suffix string) {
	for _, b := range p.buffer {
		fmt.Println(b, suffix)
	}
}

// RenderSmallest calculates the minimal output and renders only that area
// Note: on multicolor mode the background color specifie the background color
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) RenderSmallest(color ...uint32) []string {
	var mix, max, miy, may int
	cc := uint32(0)
	if len(color) > 0 {
		cc = color[0]
	}

	mix = 0xFFFF
	max = 0
	miy = 0xFFFF
	may = 0
	for y := 0; y < p.sizeY-1; y++ {
		for x := 0; x < p.sizeX-1; x++ {
			if p.matrix[y][x] != cc {
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

// setFG internal
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) setFG(c uint32) string {
	switch p.colorrender {
	case Mode16Color:
		return fmt.Sprint("\033[1;", c&0xff, "m")
	case ModePaletteColor:
		return fmt.Sprint("\033[38;5;", c&0xff, "m")
	case ModeTrueColor:
		return fmt.Sprint("\033[38;2;", (c>>16)&0xff, ";", (c>>8)&0xff, ";", c&0xff, "m")
	}
	return ""
}

// setBG internal
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) setBG(c uint32) string {
	switch p.colorrender {
	case Mode16Color:
		return fmt.Sprint("\033[1;", c&0xff+10, "m")
	case ModePaletteColor:
		return fmt.Sprint("\033[48;5;", c&0xff, "m")
	case ModeTrueColor:
		return fmt.Sprint("\033[48;2;", (c>>16)&0xff, ";", (c>>8)&0xff, ";", c&0xff, "m")
	}
	return ""
}

// RenderXY renders a given rectangle from pixelDING, given by x1,y1 to x2,y2
// ----------------------------------------------------------------------------------------------------------------------
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
	coy := string(0x2584)
	var afg uint32
	var abg uint32
	afg = math.MaxInt32
	abg = math.MaxInt32

	p.buffer = []string{}
	lo := ""

	switch p.colorrender {
	case ModeTrueColor, ModePaletteColor, Mode16Color:
		for y := y1; y < y2; y = y + 2 {
			lo = ""
			for x := x1; x < x2; x++ {

				c1 := p.getPixelC(x, y)
				c2 := p.getPixelC(x, y+1)

				if p.tmatrix[y/2][x] > 0 {

					if abg != c1 {
						//lo = lo + fmt.Sprint("\033[48;2;", (c1>>16)&0xff, ";", (c1>>8)&0xff, ";", c1&0xff, "m")
						lo = lo + p.setBG(c1) //  fmt.Sprint("\033[48;2;", (c1>>16)&0xff, ";", (c1>>8)&0xff, ";", c1&0xff, "m")
						abg = c1
					}
					if afg != c2 {
						lo = lo + p.setFG(c2) //fmt.Sprint("\033[38;2;", (c2>>16)&0xff, ";", (c2>>8)&0xff, ";", c2&0xff, "m")
						afg = c2
					}
					lo = lo + string(p.tmatrix[y/2][x])

				} else {

					ucoy := true
					if c1 == c2 {
						ucoy = false
					}
					if abg != c1 {
						lo = lo + p.setBG(c1) //fmt.Sprint("\033[48;2;", (c1>>16)&0xff, ";", (c1>>8)&0xff, ";", c1&0xff, "m")
						abg = c1
					}
					if afg != c2 {
						lo = lo + p.setFG(c2) //fmt.Sprint("\033[38;2;", (c2>>16)&0xff, ";", (c2>>8)&0xff, ";", c2&0xff, "m")
						afg = c2
					}
					if ucoy {
						lo = lo + coy
					} else {
						lo = lo + " "
					}
				}
			}
			lo = lo + "\033[0m"
			afg = math.MaxInt32
			abg = math.MaxInt32
			p.buffer = append(p.buffer, lo)
		}

	case ModeNoColor:
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
	}
	return p.buffer
}

// Render renders a pixelDING object
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Render() []string {
	return p.RenderXY(0, 0, p.sizeX, p.sizeY)
}

// bufferAnalyse internal
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) bufferAnalyse() {
	for i, s := range p.buffer {
		fmt.Println(i, len(s))
	}
}

// Clear empty the pixelDING drawing buffer
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Clear() {
	p.matrix = make([][]uint32, p.sizeY)
	p.tmatrix = make([][]rune, (p.sizeY/2)+1)
	for i := range p.matrix {
		p.matrix[i] = make([]uint32, p.sizeX)
	}
	for i := range p.tmatrix {
		p.tmatrix[i] = make([]rune, p.sizeX)
	}
}

// SetClipping experimental clipping
// SetClipping (true, 100,100,200,200) set cliping to arectangle 100,100 to 200,200
// SetClipping (false) deactivate clipping
// SetClipping (true) activate clipping with the last given area
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) SetClipping(clp bool, xyxy ...int) {
	p.clipping = clp
	if len(xyxy) > 3 {
		p.clipsx = xyxy[0]
		p.clipsy = xyxy[1]
		p.clipex = xyxy[2]
		p.clipey = xyxy[3]
	}
}

// check internal
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) check(x0, y0 int) bool {
	if p.clipping {
		if x0 < p.clipsx || x0 > p.clipex || y0 < p.clipsy || y0 > p.clipey {
			return false
		}
	}
	if x0 < 0 || x0 > p.sizeX-1 || y0 < 0 || y0 > p.sizeY-1 {
		return false
	}
	return true
}

// scale internal
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) scale(x0, y0 int) (int, int) {
	if p.scalef != 0.0 {
		x0 = int(float64(x0) * p.scalef)
		y0 = int(float64(y0) * p.scalef)
	}
	return x0, y0
}

// sscale internal
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) sscale(x0 int) int {
	if p.scalef != 0.0 {
		x0 = int(float64(x0) * p.scalef)
	}
	return x0
}

// Dimensions sets the dimensions of a pixelDING by x,y
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Dimensions(x0, y0 int) error {
	/*
		if p.init {
			p.LastError = errors.New(AlreadySetError)
			return p.LastError
		}
	*/
	if x0 < 1 || y0 < 1 {
		p.LastError = errors.New(DimensionError)
		return p.LastError
	}
	if x0 > MaxX || y0 > MaxY {
		p.LastError = errors.New(DimensionError)
		return p.LastError
	}
	p.matrix = make([][]uint32, y0)
	p.tmatrix = make([][]rune, (y0/2)+1)

	for i := range p.matrix {
		p.matrix[i] = make([]uint32, x0)
	}
	for i := range p.tmatrix {
		p.tmatrix[i] = make([]rune, x0)
	}
	p.sizeX = x0
	p.sizeY = y0
	p.init = true
	return nil
}

// getPixel internal
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) getPixel(x0, y0 int) bool {
	// sizeX, sizeY = p.scale(sizeX, sizeY)
	if !p.check(x0, y0) {
		return false
	}
	if p.matrix[y0][x0] != 0 {
		return true
	}
	return false
}

// getPixelC internal
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) getPixelC(x0, y0 int) uint32 {
	// sizeX, sizeY = p.scale(sizeX, sizeY)
	if !p.check(x0, y0) {
		return 0
	}
	return p.matrix[y0][x0]
}

// setPixel internal
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) setPixel(x0, y0 int, b bool) {
	if !p.check(x0, y0) {
		return
	}
	if b {
		p.matrix[y0][x0] = p.acolor
	} else {
		p.matrix[y0][x0] = p.bcolor
	}
	/*
		if p.toggle {
			p.matrix[y0][x0] = !p.matrix[y0][x0]
		} else {
			p.matrix[y0][x0] = b
		}
	*/
}

// setPixelC internal
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) setPixelC(x0, y0 int, color uint32) {
	if !p.check(x0, y0) {
		return
	}
	p.matrix[y0][x0] = color
}

// GetPixelC gets the color of the Pixel at x,y
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) GetPixelC(x0, y0 int) uint32 {
	x0, y0 = p.scale(x0, y0)
	/*	if !p.check(x0, y0) {
			return false
		}
	*/
	return p.getPixelC(x0, y0)
}

// PixelC set the Pixel at x,y to color
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) PixelC(x0, y0 int, color uint32) {
	x0, y0 = p.scale(x0, y0)
	/*	if !p.check(x0, y0) {
			return
		}
	*/
	p.setPixelC(x0, y0, color)
}

// GetPixel return true if the pixel is set (none color mode)
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) GetPixel(x0, y0 int) bool {
	x0, y0 = p.scale(x0, y0)
	/*	if !p.check(x0, y0) {
			return false
		}
	*/
	return p.getPixel(x0, y0)
}

// Pixel set the Pixel on x,y as set (none colore mode)
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Pixel(x0, y0 int, set bool) {
	x0, y0 = p.scale(x0, y0)
	/*	if !p.check(x0, y0) {
			return
		}
	*/
	p.setPixel(x0, y0, set)
}

// abs internal
// ----------------------------------------------------------------------------------------------------------------------
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// TextFrame set various text frames
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) TextFrame(x1, y1, x2, y2 int, l string, bitmask int) {
	noLineH := false
	noLineV := false
	sx := strings.Split(l, "")
	h1 := ""
	h2 := ""

	if x1 > x2 {
		x1, x2 = x2, x1 //swap(x1, x2)
	}
	if y1 > y2 {
		y1, y2 = y2, y1 //swap(y1, y2)
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
				p.Text(x1, i, sx[3])
			}
			if bitmask&(1<<3) != 0 {
				p.Text(x2, i, sx[5])
			}
		}
	}

	if noLineH == false {
		hs := x2 - x1 - 1

		if bitmask&(1<<7) != 0 {
			h1 = strings.Repeat(sx[1], hs)
			p.Text(x1+1, y1, h1)
		}
		if bitmask&(1<<1) != 0 {
			h2 = strings.Repeat(sx[7], hs)
			p.Text(x1+1, y2, h2)
		}
	}

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

// TextFrameBuffer Put a text Frame in the rendered buffer
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) TextFrameBuffer(x1, y1, x2, y2 int, l string, bitmask int, scale ...bool) {
	noLineH := false
	noLineV := false
	sx := strings.Split(l, "")
	h1 := ""
	h2 := ""

	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
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
					p.TextBuffer(x1, i, sx[3], scale[0])
				} else {
					p.TextBuffer(x1, i, sx[3])
				}
			}
			if bitmask&(1<<3) != 0 {
				if len(scale) > 0 {
					p.TextBuffer(x2, i, sx[5], scale[0])
				} else {
					p.TextBuffer(x2, i, sx[5])
				}
			}
		}
	}

	if noLineH == false {
		hs := x2 - x1 - 1
		if len(scale) > 0 {
			if scale[0] == true && p.aspectX == 0 {
				hs = (x2 - x1) / 2
			}
		}
		if bitmask&(1<<7) != 0 {
			h1 = strings.Repeat(sx[1], hs)
			if len(scale) > 0 {
				p.TextBuffer(x1+1, y1, h1, scale[0])
			} else {
				p.TextBuffer(x1+1, y1, h1)
			}
		}
		if bitmask&(1<<1) != 0 {
			h2 = strings.Repeat(sx[7], hs)
			if len(scale) > 0 {
				p.TextBuffer(x1+1, y2, h2, scale[0])
			} else {
				p.TextBuffer(x1+1, y2, h2)
			}
		}
	}

	if len(scale) > 0 {
		if bitmask&(1<<8) != 0 {
			p.TextBuffer(x1, y1, sx[0], scale[0])
		}
		if bitmask&(1<<6) != 0 {
			p.TextBuffer(x2, y1, sx[2], scale[0])
		}
		if bitmask&(1<<2) != 0 {
			p.TextBuffer(x1, y2, sx[6], scale[0])
		}
		if bitmask&(1<<0) != 0 {
			p.TextBuffer(x2, y2, sx[8], scale[0])
		}
	} else {
		if bitmask&(1<<8) != 0 {
			p.TextBuffer(x1, y1, sx[0])
		}
		if bitmask&(1<<6) != 0 {
			p.TextBuffer(x2, y1, sx[2])
		}
		if bitmask&(1<<2) != 0 {
			p.TextBuffer(x1, y2, sx[6])
		}
		if bitmask&(1<<0) != 0 {
			p.TextBuffer(x2, y2, sx[8])
		}
	}
}

// TextLineHBuffer splits a text horizontal in the rendered buffer
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) TextLineHBuffer(x1, y1, x2, y2 int, l string, set ...bool) {
	sx := strings.Split(l, "")
	hs := ""
	//lc := strings.Split(l,"")
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	hs = strings.Repeat(sx[1], x2-x1)
	if len(set) > 0 {
		p.TextBuffer(x1, y1, hs, set[0])
	} else {
		p.TextBuffer(x1, y1, hs)
	}
}

// TextLineVBuffer splits a text vertical into the rendered buffer
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) TextLineVBuffer(x1, y1, x2, y2 int, l string, set ...bool) {
	sx := strings.Split(l, "")
	//lc := strings.Split(l,"")
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	for i := y1; i < y2; i++ {
		if len(set) > 0 {
			p.TextBuffer(x1, i, sx[3], set[0])
		} else {
			p.TextBuffer(x1, i, sx[3])
		}
	}
}

// Text put a text on x,y to the pixelDING
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Text(x0, y0 int, text string) {
	if x0 < 0 || x0 > p.sizeX-1 || y0 < 0 || y0 > p.sizeY-1 {
		return
	}

	xx := x0
	xy := y0 / 2
	rc := []rune(text)

	for i := range rc {

		if xx > p.sizeX-1 {
			break
		}
		////rs := []rune("\u2220")
		//r := rune(text[i])
		p.tmatrix[xy][xx] = rc[i]
		p.Pixel(xx, xy*2, false)
		p.Pixel(xx, xy*2+1, true)
		xx++
	}
}

// TextBuffer put a text on x,y in the rendered buffer
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) TextBuffer(x0, y0 int, text string, scale ...bool) {
	if len(scale) > 0 {
		if scale[0] == true {
			x0, y0 = p.scale(x0, y0)
			if p.aspectX == 0 {
				x0 = x0 / 2
			}
			if p.aspectY == 0 {
				y0 = y0 / 2
			}
		}
	}
	var n string
	if len(p.buffer) < y0+1 {
		return //Out of bounds
	}

	if x0 < 0 || y0 < 0 {
		return //Out of bounds
	}

	sx := strings.Split(text, "")
	l := len(sx)

	s := strings.Split(p.buffer[y0], "")
	sl := len(s)

	if x0 > sl {
		return //Out of bounds
	}

	//	fmt.Println("sl", sl, "l", l)
	cs := 1
	for i := 0; i < x0; i++ {
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
	for i := x0 + l; i < sl; i++ {
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
	p.buffer[y0] = n

}

// HBar return a horizontal bar string
// ----------------------------------------------------------------------------------------------------------------------
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

// SVGPath draws a standard SVG path at x,y
// There is at the moment the arc function missing (work in progress)
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) SVGPath(x0, y0 float64, s string, set bool, fscale ...float64) {
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
				p.Line(int((lx+x0)*scale), int((ly+y0)*scale), int((x+x0)*scale), int((y+y0)*scale), set)
				lx = x
				ly = y
			case "M", "m":
				lx = x
				ly = y
				lcx = x //Set the last Control Point to
				lcy = y
			case "H", "h":
				p.Line(int((lx+x0)*scale), int((ly+y0)*scale), int((x+x0)*scale), int((y+y0)*scale), set)
				lx = x
			case "V", "v":
				p.Line(int((lx+x0)*scale), int((ly+y0)*scale), int((x+x0)*scale), int((y+y0)*scale), set)
				ly = y

			case "C", "c":
				p.CBezier(int((lx+x0)*scale), int((ly+y0)*scale), int((c1x+x0)*scale), int((c1y+y0)*scale), int((c2x+x0)*scale), int((c2y+y0)*scale), int((x+x0)*scale), int((y+y0)*scale), set)
				//p.Line(int(lx), int(ly), int(sizeX), int(sizeY))
				lx = x
				ly = y
				lcx = c2x
				lcy = c2y

			case "S", "s":
				c1x = lx - -(lx - lcx)
				c1y = ly - -(ly - lcy)

				p.CBezier(int((lx+x0)*scale), int((ly+y0)*scale), int((c1x+x0)*scale), int((c1y+y0)*scale), int((c2x+x0)*scale), int((c2y+y0)*scale), int((x+x0)*scale), int((y+y0)*scale), set)
				//p.Line(int(lx), int(ly), int(sizeX), int(sizeY))
				lx = x
				ly = y
				lcx = c2x
				lcy = c2y

			case "Q", "q":

				//c2x = sizeX- -(sizeX-c1x)
				//c2y = sizeY- -(sizeY-c1y)
				p.QBezier(int((lx+x0)*scale), int((ly+y0)*scale), int((c1x+x0)*scale), int((c1y+y0)*scale), int((x+x0)*scale), int((y+y0)*scale), set)
				//p.Line(int(lx), int(ly), int(sizeX), int(sizeY))
				lx = x
				ly = y
				lcx = c1x
				lcy = c1y

			case "T", "t":
				c1x = lx - -(lx - lcx)
				c1y = ly - -(ly - lcy)

				p.QBezier(int((lx+x0)*scale), int((ly+y0)*scale), int((c1x+x0)*scale), int((c1y+y0)*scale), int((x+x0)*scale), int((y+y0)*scale), set)
				//p.Line(int(lx), int(ly), int(sizeX), int(sizeY))
				lx = x
				ly = y
				lcx = c1x
				lcy = c1y

			case "Z", "z":
				p.Line(int((lx+x0)*scale), int((ly+y0)*scale), int((x+x0)*scale), int((y+y0)*scale), set)
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

// QBezier plots a quadratic Bezier in pixelDING
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) QBezier(x1, y1, cx1, cy1, x2, y2 int, set bool) {
	var px, py int
	x0 := x1
	y0 := y1
	x1, y1 = p.scale(x1, y1)
	cx1, cy1 = p.scale(cx1, cy1)
	x2, y2 = p.scale(x2, y2)
	fx1, fy1 := float64(x1), float64(y1)
	fx2, fy2 := float64(cx1), float64(cy1)
	fx3, fy3 := float64(x2), float64(y2)
	for i := 0; i <= p.msteps; i++ { //range px {
		//for i := range px {
		c := float64(i) / float64(p.msteps)
		a := 1 - c
		a, b, c := a*a, 2*c*a, c*c

		px = int(a*fx1 + b*fx2 + c*fx3)
		py = int(a*fy1 + b*fy2 + c*fy3)

		p.Line(x0, y0, px, py, set)
		x0, y0 = px, py
	}
}

// CBezier plots a Bezier from x,y(1) to x,y(2) with two power lines x,y(3) x,y(4)
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) CBezier(x1, y1, cx1, cy1, cx2, cy2, x2, y2 int, set bool) {
	var px, py int
	x0 := x1
	y0 := y1
	x1, y1 = p.scale(x1, y1)
	cx1, cy1 = p.scale(cx1, cy1)
	cx2, cy2 = p.scale(cx2, cy2)
	x2, y2 = p.scale(x2, y2)
	fx1, fy1 := float64(x1), float64(y1)
	fx2, fy2 := float64(cx1), float64(cy1)
	fx3, fy3 := float64(cx2), float64(cy2)
	fx4, fy4 := float64(x2), float64(y2)
	for i := 0; i <= p.msteps; i++ { //range px {
		//for i := range px {
		d := float64(i) / float64(p.msteps)
		a := 1 - d
		b, c := a*a, d*d
		a, b, c, d = a*b, 3*b*d, 3*a*c, c*d

		px = int(a*fx1 + b*fx2 + c*fx3 + d*fx4)
		py = int(a*fy1 + b*fy2 + c*fy3 + d*fy4)

		p.Line(x0, y0, px, py, set)
		x0, y0 = px, py
	}
}

// QBezier plots a quadratic Bezier in pixelDING
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) GetQBezierXY(x1, y1, cx1, cy1, x2, y2 int, f float64) (float64, float64) {
	//	var px, py [50 + 1]int
	var px, py float64
	x1, y1 = p.scale(x1, y1)
	cx1, cy1 = p.scale(cx1, cy1)
	x2, y2 = p.scale(x2, y2)
	fx1, fy1 := float64(x1), float64(y1)
	fx2, fy2 := float64(cx1), float64(cy1)
	fx3, fy3 := float64(x2), float64(y2)
	a := 1 - f
	a, b, c := a*a, 2*f*a, f*f
	px = a*fx1 + b*fx2 + c*fx3
	py = a*fy1 + b*fy2 + c*fy3
	return px, py
}

// CBezier plots a Bezier from x,y(1) to x,y(2) with two power lines x,y(3) x,y(4)
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) GetCBezierXY(x1, y1, cx1, cy1, cx2, cy2, x2, y2 int, f float64) (float64, float64) {
	var px, py float64
	x1, y1 = p.scale(x1, y1)
	cx1, cy1 = p.scale(cx1, cy1)
	cx2, cy2 = p.scale(cx2, cy2)
	x2, y2 = p.scale(x2, y2)
	fx1, fy1 := float64(x1), float64(y1)
	fx2, fy2 := float64(cx1), float64(cy1)
	fx3, fy3 := float64(cx2), float64(cy2)
	fx4, fy4 := float64(x2), float64(y2)
	d := f
	a := 1 - d
	b, c := a*a, d*d
	a, b, c, d = a*b, 3*b*d, 3*a*c, c*d
	px = a*fx1 + b*fx2 + c*fx3 + d*fx4
	py = a*fy1 + b*fy2 + c*fy3 + d*fy4
	return px, py
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
// Rectangle plots a rectange x,y(1) to x,y(2) filled or unfilled
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

// floodfill internal
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) floodFill(x0, y0 int, prevC, newC bool) {
	if x0 < 0 || x0 >= p.sizeX || y0 < 0 || y0 >= p.sizeY {
		return
	}
	if prevC != p.getPixel(x0, y0) {
		return
	}
	// Replace the color at (sizeX, sizeY)
	p.setPixel(x0, y0, newC)
	// Recur for north, east, south and west
	p.floodFill(x0+1, y0, prevC, newC)
	p.floodFill(x0-1, y0, prevC, newC)
	p.floodFill(x0, y0+1, prevC, newC)
	p.floodFill(x0, y0-1, prevC, newC)
}

// Fill floodfills the area, starting with the pixel color at x,y
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) Fill(x0, y0 int, newC bool) {
	x0, y0 = p.scale(x0, y0)
	prevC := p.getPixel(x0, y0)
	if prevC == newC {
		return
	}
	p.floodFill(x0, y0, prevC, newC)
}

func toRadian(angle int) float64 {
	return float64(angle) * (math.Pi / 180.0)
}

// DotArc plot a dotted Arc at x,y with radius r, from degree a1 to degree a2
// NOTE: the angle 0° is at 12 o'clock, 90° at 3 o'clock, 180° at 6 o'clock, 270° at 9 o'clock
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) DotArcClock(x0, y0, r int, a1, a2, step int, set bool) { //wieso

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

// LineArc plot a Arc at x,y with radius r, from degree a1 to degree a2
// NOTE: the angle 0° is at 12 o'clock, 90° at 3 o'clock, 180° at 6 o'clock, 270° at 9 o'clock
// ----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) LineArcClock(x0, y0, r int, a1, a2, step int, set bool) { //wieso

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

// DotArc plot a dotted Arc at x,y with radius r, from degree a1 to degree a2
// ----------------------------------------------------------------------------------------------------------------------
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

		yo := int(math.Round(float64(r) * math.Sin(toRadian(a1%360))))
		xo := int(math.Round(float64(r) * math.Cos(toRadian(a1%360))))

		p.setPixel(x0+xo, y0-yo, set)

		a1 += step

	}

}

// LineArc plot a Arc at x,y with radius r, from degree a1 to degree a2
// ----------------------------------------------------------------------------------------------------------------------
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

		yo := int(math.Round(float64(r) * math.Sin(toRadian(a1%360))))
		xo := int(math.Round(float64(r) * math.Cos(toRadian(a1%360))))

		if !firstfrom {
			p.Line(x0+fromx, y0-fromy, x0+xo, y0-yo, set)

		}
		firstfrom = false
		fromx = xo
		fromy = yo

		a1 += step

	}

	yo := int(math.Round(float64(r) * math.Sin(toRadian(a2%360))))
	xo := int(math.Round(float64(r) * math.Cos(toRadian(a2%360))))
	p.Line(x0+fromx, y0-fromy, x0+xo, y0-yo, set)
}

// LineRadius draw a line from center x,y(0) to degree a1 from radius r1 to r2
// r1 is the distance from center to r1. If you want o draw from center r1 = 0.0
// and r2 is the max radius.
// ----------------------------------------------------------------------------------------------------------------------
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

// EllipseRect draw a Elipse which fits into the box given by x,y(0) to x,y(1)
// ----------------------------------------------------------------------------------------------------------------------
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

// Circle draw a Circle
// ----------------------------------------------------------------------------------------------------------------------
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

// Line draw a line
// ----------------------------------------------------------------------------------------------------------------------
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

// DotLine draw a line, specified by the pattern
// ----------------------------------------------------------------------------------------------------------------------
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
		if (pat & 0x01) == 0x01 {
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
