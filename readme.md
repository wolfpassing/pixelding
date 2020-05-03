# PixelDING
![pixelding](screenshots/console.png "The PixelDING")

## What is PixelDING ?
It's not easy to describe the purpose of PixelDING. It started with straight lines, in an SVG like path string as output to the **console**.
Over the time more and more things got implemented and it is still not finished.

Lets have a look of some samples.

### Simple things
![circle](screenshots/circle.png "The PixelDING")
![rectangle](screenshots/rectngle.png "The PixelDING")
![line](screenshots/line.png "The PixelDING")
![ellipse](screenshots/ellipse.png "The PixelDING")
Circles, Rectangles, Lines and Ellipses

### Bezier things
![bezier](screenshots/bezier1.png "The PixelDING")
![bezier](screenshots/bezier2.png "The PixelDING")
Simple and complex Bezier curves

### SVG Path things
![svg1](screenshots/svg1.png "The PixelDING")
![svg2](screenshots/svg2.png "The PixelDING")
SVG compatible and scalable paths

### Font things
![font](screenshots/font.png "The PixelDING")
Fonts for Text

### Stamps
![stamps](screenshots/stamp.png "The PixelDING")
Stamps (Bitmaps)

# Lets have a general look into

## Index

func New() PixelDING
func (p *PixelDing) SaveFont(name string, font PixelFont) error
func (p *PixelDing) LoadFont(name string) PixelFont
func (p *PixelDing) SaveStamp(name string, stamp PixelStamp) error
func (p *PixelDing) LoadStamp(name string) PixelStamp
func (p *PixelDing) FontPrint(font string, x, y int, text string, set bool)
func (p *PixelDing) AddFont(name string, font PixelFont)
func (p *PixelDing) AddStamp(name string, stamp PixelStamp)
func (p *PixelDing) FontInfo(name string)
func (p *PixelDing) SetStep(x int)
func (p *PixelDing) Toggle(b bool)
func (p *PixelDing) Invert(b bool)
func (p *PixelDing) Debug(b bool)
func (p *PixelDing) Scale(s float64)
func (p *PixelDing) Aspect(x, y int)
func (s *PixelStamp) X() int
func (p *PixelDing) Stamp(x, y int, stamp *PixelStamp, set bool, st bool)
func (p *PixelDing) Display()
func (p *PixelDing) RenderSmallest() []string
func (p *PixelDing) RenderXY(x1, y1, x2, y2 int) []string
func (p *PixelDing) Render() []string
func (p *PixelDing) Clear()
func (p *PixelDing) Dimensions(x, y int) error

func (p *PixelDing) Pixel(x, y int, b bool)
func (p *PixelDing) SVGPath(xo,yo float64, s string, fscale ...float64)
func (p *PixelDing) Bezier2(x1, y1, x2, y2, x3, y3 int)
func (p *PixelDing) Bezier(x1, y1, x2, y2, x3, y3, x4, y4 int)
func (p *PixelDing) Rectangle(x0, y0, x1, y1 int, b bool, f bool)
func (p *PixelDing) Fill(x, y int, newC bool)
func (p *PixelDing) EllipseRect(x0, y0, x1, y1 int)
func (p *PixelDing) Circle(x0, y0, r int)
func (p *PixelDing) Line(x0, y0, x1, y1 int)




