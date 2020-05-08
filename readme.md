# PixelDING
![pixelding](screenshots/console.png "The PixelDING")

## What is PixelDING ?
It's not easy to describe the purpose of PixelDING. It started with straight lines, in an SVG like path string as output to the **console**.
Over the time more and more things got implemented and it is still not finished.

Lets have a look of some samples.

### Simple things
![circle](screenshots/circle.png "The PixelDING")
![rectangle](screenshots/rectangle.png "The PixelDING")
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

## Small Example
````GO
pd := pixelding.New(100,100)
pd.Circle(50,50,45)
pd.CBezier(5,50,5,25,50,25,50,50)
pd.CBezier(50,50,50,75,95,75,95,50)
pd.Fill(49,51,true)
pd.Render()
pd.Display()

````
![yinyang](screenshots/yinyang.png "The PixelDING")
That should be the result. A perfect circle yinyang.
>Note: Depending on the console font size the result could be squeezed horizontal or vertical, see the Aspect command for compensating this.


## Preparations

#### func New(dimensions... int) PixelDING
Initialize and return a PixelDING object. If you want you can specify the dimensions X and Y for the paint area.

#### func (p *PixelDING) Dimensions(x, y int) error
Set the dimension for the paint area

#### func (p *PixelDING) Clear()
Clear the painting area

#### func (p *PixelDING) Aspect(x, y int)
Set the aspect ratio. 0 = normal, 1 = double

#### func (p *PixelDING) Invert(b bool)
Enable or disable the invert mode for rendering

#### func (p *PixelDING) SetStep(x int)
Set the steps for Bezier or other curves. Default value is 15, the lesser the value the faster, but also getting rougher. The maximum Value supported is 50.

#### func (p *PixelDING) Render() []string
Render the internal paint area. The result is returned but also internaly saved for the Display command.

#### func (p *PixelDING) RenderSmallest() []string
#### func (p *PixelDING) RenderXY(x1, y1, x2, y2 int) []string

#### func (p *PixelDING) Display()
Display the rendered paint area on console output

## Drawing

#### func (p *PixelDING) Pixel(x, y int, b bool)
Put a pixel on the paint area at x,y if b=true. If b=false the pixel is cleared

#### func (p *PixelDING) Line(x0, y0, x1, y1 int, b bool)
Paint a simple line from x0,y0 to x1,y1. Set the pixels on b=true otherwise clear them.

#### func (p *PixelDING) QBezier(x1, y1, x2, y2, x3, y3 int, b bool)
Quadratic Bezier Curve from x1,y1 to x3,y3, control point x2,y2. Set the pixels on b=true otherwise clear them.

#### func (p *PixelDING) CBezier(x1, y1, x2, y2, x3, y3, x4, y4 int, b bool)
Cubic Bezier from x1,y1 to x4,y4. Control points x2,y2 and x3,y3. Set the pixels on b=true otherwise clear them.

#### func (p *PixelDING) Rectangle(x0, y0, x1, y1 int, b bool, f bool)
Draw a rectangle from x0,y0 to x1,y1. Set the pixels on b=true otherwise clear them. Fill the rectangle on f=true

#### func (p *PixelDING) Fill(x, y int, newC bool)
Floodfill, start at x,y

#### func (p *PixelDING) EllipseRect(x0, y0, x1, y1 int, b bool)
Draw ellipse in the given box defined by x0,y0 to x1,y1. Set the pixels on b=true otherwise clear them.

#### func (p *PixelDING) Circle(x0, y0, r int, b bool)
Draw a circle at x0,y0 with the radius r. Set the pixels on b=true otherwise clear them.

#### func (p *PixelDING) SVGPath(xo,yo float64, s string, b bool, fscale ...float64)
Interprete and draw the path in s. Set the pixels on b=true otherwise clear them. Scale it by fscale.


## Font and Stamp Loading

#### func (p *PixelDING) FontPrint(font *PixelFont, x, y int, text string, b bool)
Print out the text string at x,y on the paint area, Set the pixels on b=true otherwise clear them.

#### func (p *PixelDING) Stamp(x, y int, stamp *PixelStamp, set bool, st bool)
Set the pixels on set=true otherwise clear them. Stamp mode if st=true. Stamp mode means, that the bitmap is transfered to the paint area without blending it together with the background.

#### func (p *PixelDING) LoadFont(name string) PixelFont
Load a front by the given name (including path). This is returning a PixelFont object which can be added via AddFont function.

#### func (p *PixelDING) SaveFont(name string, font PixelFont) error
Generated Fotns can be saved too.

#### func (p *PixelDING) LoadStamp(name string) PixelStamp
Load a stamp by the given name (including path). This is returning a PixelStamp object which can be added via AddStamp function.

#### func (p *PixelDING) SaveStamp(name string, stamp PixelStamp) error

#### func (p *PixelDING) AddFont(name string, font PixelFont)
#### func (p *PixelDING) AddStamp(name string, stamp PixelStamp)

#### func (p *PixelDING) Toggle(b bool)
#### func (p *PixelDING) Scale(s float64)
#### func (p *PixelDING) X() int
#### func (p *PixelDING) Y() int
#### func (p *PixelStamp) X() int
#### func (p *PixelStamp) Y() int




