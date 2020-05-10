# PixelDING
![pixelding](screenshots/console.png "The PixelDING")

## What is PixelDING ?
It's not easy to describe the purpose of PixelDING. It started with straight lines, in an SVG like path string as output to the **console**.
Over the time more and more things got implemented and it is still not finished.

![pixelding](screenshots/charts.png "The PixelDING")

Lets have a look of some samples.
----
### Simple things
![circle](screenshots/circle.png) <!-- .element height="50%" width="50%" -->
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

# Lets have a more detailed look into

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
That should be the result. A perfect circle yinyang showing up in yor console.
>Note: Depending on the console font size the result could be squeezed horizontal or vertical, see the Aspect command for compensating this.

----
## Preparations and Information

### Coordinates
Whenever there is a X and Y coordinate used, it is zero based. On an paint area of 200 by 200 dimensions the X and Y coordinate range is from 0 to 199.


#### func New(dimensions... int) PixelDING
Initialize and return a PixelDING object. If you want you can specify the dimensions X and Y for the paint area.
````GO
pd = pixelding.New(100,100) //Create and set dimensions to 100 by 100
pd = pixelding.New()        //Just Create an PixelDING object
````

----
#### func (p *PixelDING) Dimensions(x, y int) error
Set the dimension for the paint area
````GO
pd.Dimensions(100,100) //Set dimensions to 100 by 100
````

----
#### func (p *PixelDING) Clear()
Clear the painting area
````GO
pd.Clear() //Clear the painting area
````

----
#### func (p *PixelDING) Aspect(x, y int)
Set the aspect ratio. 0 = normal, 1 = double. Due to different font metrics on different font sizes in the console, this could help you keep the aspect ratio near 1:1. The perfect font metrics would be a of square size (same height an width). If you paint a circle and the circle is squeezed horizontal, then use X aspect 1. If it is squeezed vertically use 1 at Y aspect.
````GO
pd.Aspect(1,0) //Change the x aspect to double to reduce horizontal squeeze
````

----
#### func (p *PixelDING) Invert(b bool)
Enable or disable the invert mode for rendering
````GO
pd.Invert(false)    //Normal output mode
````

----
#### func (p *PixelDING) SetStep(x int)
Set the steps for Bezier or other curves. Default value is 15, the lesser the value the faster, but also getting rougher. The maximum Value supported is 50.
````GO
pd.SetStep(7)
````

----
#### func (p *PixelDING) Render() []string
#### func (p *PixelDING) RenderSmallest() []string
#### func (p *PixelDING) RenderXY(x1, y1, x2, y2 int) []string
Render the internal paint area. The result is returned as an string slice and internal stored in the PixelDING object saved for the Display command.

For the RenderSmallest all leading and trailing empty pixels are not rendered. If there is just one pixel painted at 75,25 on a 100 by 100 paint area there is just one pixel renderd and the string slice result consists of on char.

````GO
pd.Render()                 //Simple Render, full dimension as set in the PixelDING object
pd.RenderSmallest()         //Render exactly what is needed.
pd.RenderXY(10,10,50,50)    //Render from 10,10 to 50,50
buffer := pd.Render()       //Render AND return the result into buffer variable
````

----
#### func (p *PixelDING) Display()
Display the rendered paint area on console output
````GO
pd.Display()        //prints out the rendered buffer from the PixelDING object.
````

## Drawing

#### func (p *PixelDING) Pixel(x, y int, b bool)
Put a pixel on the paint area at x,y if b=true. If b=false the pixel is cleared
````GO
pd.Pixel(25,50, true)   //Paint a Pixel
pd.Pixel(50,20, false)  //Clear the Pixel
````

----
#### func (p *PixelDING) GetPixel(x, y int) bool
Returns true if the pixel on the paint area at x,y if set, otherwise false
````GO
result := pd.GetPixel(25,50)   //Check if the pixel is set
````

----
#### func (p *PixelDING) Line(x0, y0, x1, y1 int, b bool)
Paint a simple line from x0,y0 to x1,y1. Set the pixels on b=true otherwise clear them.
````GO
pd.Line(10,10,50,50,true)       //paint a line
````

----
#### func (p *PixelDING) QBezier(x1, y1, x2, y2, x3, y3 int, b bool)
Quadratic Bezier Curve from x1,y1 to x3,y3, control point x2,y2. Set the pixels on b=true otherwise clear them.
````GO
pd.QBezier(5,50,30,0,50,50)
````

----
#### func (p *PixelDING) CBezier(x1, y1, x2, y2, x3, y3, x4, y4 int, b bool)
Cubic Bezier from x1,y1 to x4,y4. Control points x2,y2 and x3,y3. Set the pixels on b=true otherwise clear them.
````GO
pd.CBezier(5,50,5,25,50,25,50,50)
````

----
#### func (p *PixelDING) Rectangle(x0, y0, x1, y1 int, b bool, f bool)
Draw a rectangle from x0,y0 to x1,y1. Set the pixels on b=true otherwise clear them. Fill the rectangle on f=true
````GO
pd.Rectangle(10,10,50,50,true)          //Rectangle
pd.Rectangle(10,10,50,50,true,true)     //Rectangle filled
````

----
#### func (p *PixelDING) Fill(x, y int, newC bool)
Floodfill, start at x,y
````GO
pd.Floodfill(22,22)     //Floodfill starting add 22,22
````

----
#### func (p *PixelDING) EllipseRect(x0, y0, x1, y1 int, b bool)
Draw ellipse in the given box defined by x0,y0 to x1,y1. Set the pixels on b=true otherwise clear them.
````GO
pd.EllipseRect(10,10,100,50)
````

----
#### func (p *PixelDING) Circle(x0, y0, r int, b bool)
Draw a circle at x0,y0 with the radius r. Set the pixels on b=true otherwise clear them.
````GO
pd.Circle(25,25,10,true)
````

----
#### func (p *PixelDING) SVGPath(xo,yo float64, s string, b bool, fscale ...float64)
Interprete and draw the path in s. Set the pixels on b=true otherwise clear them. Scale it by fscale.
````GO
pd.SVGPath(0, 0, "M 5, 60 c 25, -25 50, 25 75, 12 s 50, 50 75, -15", true)
````


## Font and Stamp Loading

----
#### func (p *PixelDING) FontPrint(font *PixelFont, x, y int, text string, b bool)
Print out the text string at x,y on the paint area, Set the pixels on b=true otherwise clear them.
````GO
pd.FontPrint(myFont,20,20,"This is a test", true)       //Write on the paint area, mode pixel set
pd.FontPrint(pd.GetFont("copper"),20,50,"(a+b)=x",true) //Use a stored font "copper"
````

----
#### func (p *PixelDING) Stamp(x, y int, stamp *PixelStamp, set bool, st bool)
Set the pixels on set=true otherwise clear them. Stamp mode if st=true. Stamp mode means, that the bitmap is transfered to the paint area without blending it together with the background.
````GO
pd.Stamp(myStamp,20,20,true)                    //Stamp in set Mode
pd.Stamp(pd.GetStamp("flower"),20,50,true,true) //Stamp a stored stamp with name "flower" in stamp mode
````

----
#### func (p *PixelDING) LoadFont(name string) *PixelFont
Load a front by the given name (including path). This is returning a PixelFont object which can be added via AddFont function. LoadFont as a other Load functions does return **nil** on file errors
````GO
myFont := pd.LoadFont("/home/user/fonts/smallfont.fnt")
if myFont == nil {
    return errors.New("File not found.")
}
````

----
#### func (p *PixelDING) SaveFont(name string, font *PixelFont) error
Generated Fonts can be saved too.
````GO
pd.SaveFont("C:\\FONTS\\SMALLFONT.FNT")
````

----
#### func (p *PixelDING) LoadStamp(name string) *PixelStamp
Load a stamp by the given name (including path). This is returning a PixelStamp object which can be added via AddStamp function. LoadStamp as a other Load functions does return **nil** on file errors
````GO
myStamp := pd.LoadStamp("/home/user/stamps/duck.stp")
if myStamp == nil {
    return errors.New("File not found.")
}
````

----
#### func (p *PixelDING) SaveStamp(name string, stamp *PixelStamp) error
````GO
pd.SaveStamp("C:\\STAMPS\\DUCK.STP")
````

----
#### func (p *PixelDING) AddFont(name string, font *PixelFont)
Adds the font to the PixelDing object
````GO
pd.AddFont("copper",myFont)     //Add the font object to the PixelDING object
````

----
#### func (p *PixelDING) GetFont(name string) *PixelFont
Get the font from the PixelDing object
````GO
myFont := pd.GetFont("copper")
````

----
#### func (p *PixelDING) AddStamp(name string, stamp *PixelStamp)
Add the stamp to the PixelDing object
````GO
pd.AddStamp("duck", myStamp)    //Add a stamp with the name "duck" to the PixelDING object
````

----
#### func (p *PixelDING) GetStamp(name string) *PixelStamp
Get the the stamp from the PixelDing object
````GO
myStamp := pd.GetStamp("duck")
````

----
#### func (p *PixelDING) Toggle(b bool)
Toggle the pixelmode, from set to clear. WARNING Deprecated ! If toggle is enabled, all pixel set operations are inverting the pixel on the position. If a pixel is already set, the pixel is cleared and vice versa.
````GO
pd.Toggle(true)
````

----

#### func (p *PixelDING) Scale(s float64)
Global scale the coordinates.
````GO
pd.Scale(0.5)   //Set scale to 50%
````
Here a more detailed example
````GO
pd := pixelding.New(100,100)
pd.Circle(99,50,20)            //Half of the circle is visible
pd.Render()
pd.Display()                   //String array 50 by 50 chars
````
````GO
pd := pixelding.New(100,100)
pd.Scale(0.5)                  //Set scale to 50%
pd.Circle(99,50,20)            //Full circle visible but half size
pd.Render()
pd.Display()                   //String Arry 50 by 50 no size change on scale
````

----
#### func (p *PixelDING) X() int
Get the X size of the paint area. Together with the stamp dimensions better positioning calulations are possible
````GO
pd := pixelding.New(100,100)
cat := pd.LoadStamp("cat.stp")
pd.Stamp(pd.X()/2-cat.X()/2, pd.Y()/2-cat.Y()/2, myStamp, true)  //Stamp cat in centered position
````

----
#### func (p *PixelDING) Y() int
Get the Y size of the paint area, for details see X() function

----
#### func (p *PixelStamp) X() int
Get the X size of the stamp
````GO
stampWidth := myStamp.X()
````

----
#### func (p *PixelStamp) Y() int
Get the Y size of the stamp
````GO
stampHeight := myStamp.Y()
````





