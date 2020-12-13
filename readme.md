# PixelDING
![pixelding](screenshots/console.png "The PixelDING")

## What is PixelDING ?
It's not easy to describe the complete purpose of PixelDING. It started with the idea to have in our microservices nice looking text statistic in the console logs.
First there where just straight lines - in an SVG like path string as output to the **console**.
Over the time more and more things got implemented and it is still not finished.

![pixelding](screenshots/charts.png "The PixelDING")

Lets have a look of some samples.
----
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

### Arc und Lines
![arc](screenshots/arclines1.png "The PixelDING")
![arc](screenshots/arclines2.png "The PixelDING")

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
That should be the result. A perfect circle Yin and Yang showing up in yor console.
>Note: Depending on the console font size the result could be squeezed horizontal or vertical, see the Aspect command for compensating this.

----
## Preparations and Information

### Coordinates
Whenever there is a X and Y coordinate used, it is zero based. On an paint area of 200 by 200 dimensions the X and Y coordinate range is from 0 to 199.


#### New(dimensions... int) PixelDING
Initialize and return a PixelDING object. If you want you can specify the dimensions X and Y for the paint area.
````GO
pd = pixelding.New(100,100) //Create and set dimensions to 100 by 100
pd = pixelding.New()        //Just Create an PixelDING object
````

----
#### Dimensions(x, y int) error
Set the dimension for the paint area
````GO
pd.Dimensions(100,100) //Set dimensions to 100 by 100
````

----
#### Clear()
Clear the painting area
````GO
pd.Clear() //Clear the painting area
````

----
#### Aspect(x, y int)
Set the aspect ratio. 0 = normal, 1 = double. Due to different font metrics on different font sizes in the console, this could help you keep the aspect ratio near 1:1. The perfect font metrics would be a of square size (same height an width). If you paint a circle and the circle is squeezed horizontal, then use X aspect 1. If it is squeezed vertically use 1 at Y aspect.
````GO
pd.Aspect(1,0) //Change the x aspect to double to reduce horizontal squeeze
````

----
#### Invert(b bool)
Enable or disable the invert mode for rendering
````GO
pd.Invert(false)    //Normal output mode
````

----
#### SetStep(x int)
Set the steps for Bezier or other curves. Default value is 15, the lesser the value the faster, but also getting rougher. The maximum Value supported is 50.
````GO
pd.SetStep(7)
````

----
#### Render() []string
#### RenderSmallest() []string
#### RenderXY(x1, y1, x2, y2 int) []string
Render the internal paint area. The result is returned as an string slice and internal stored in the PixelDING object saved for the Display command.

For the RenderSmallest all leading and trailing empty pixels are not rendered. If there is just one pixel painted at 75,25 on a 100 by 100 paint area there is just one pixel renderd and the string slice result consists of on char.

````GO
pd.Render()                 //Simple Render, full dimension as set in the PixelDING object
pd.RenderSmallest()         //Render exactly what is needed.
pd.RenderXY(10,10,50,50)    //Render from 10,10 to 50,50
buffer := pd.Render()       //Render AND return the result into buffer variable
````

----
#### Display()
Display the rendered paint area on console output
````GO
pd.Display()        //prints out the rendered buffer from the PixelDING object.
````

## Drawing

#### Pixel(x, y int, set bool)
Put a pixel on the paint area at x,y if set=true. If set=false the pixel is cleared
````GO
pd.Pixel(25,50, true)   //Paint a Pixel
pd.Pixel(50,20, false)  //Clear the Pixel
````

----
#### GetPixel(x, y int) bool
Returns true if the pixel on the paint area at x,y if set, otherwise false
````GO
result := pd.GetPixel(25,50)   //Check if the pixel is set
````

----
#### Line(x0, y0, x1, y1 int, set bool)
Paint a simple line from x0,y0 to x1,y1. Set the pixels on set=true otherwise clear them.
````GO
pd.Line(10,10,50,50,true)       //paint a line
````

----
#### DotLine(x0, y0, x1, y1 int, set bool, pattern uint8)
Paint a simple line from x0,y0 to x1,y1. Set the pixels on set=true otherwise clear them and use the bit pattern for dotted line variations.
````GO
pd.Line(10,10,50,50,true,pixelDing.Dot1x1Pattern)       //paint a line with standard dots
pd.Line(50,50,10,10,true,0B10010110)       //paint a line with uggly dots :
````
>Dot1x1Pattern   `- - - - - - - - - -`   
>Dot2x2Pattern   `--  --  --  --  --`   
>Dot4x4Pattern   `----    ----    ----`   
>Dot1x3Pattern   `-   -   -   -   -   -`   
>Dot3x2x1Pattern `---  -  ---  -  ---  -`   
>Dot6x2Pattern   `------  ------  ------`   
>Dot7x1Pattern   `------- ------- -------`   
>Dot5x1x1Pattern `----- - ----- - ----- -`   


----
#### QBezier(x1, y1, x2, y2, x3, y3 int, set bool)
Quadratic Bezier Curve from x1,y1 to x3,y3, control point x2,y2. Set the pixels on set=true otherwise clear them.
````GO
pd.QBezier(5,50,30,0,50,50)
````

----
#### CBezier(x1, y1, x2, y2, x3, y3, x4, y4 int, set bool)
Cubic Bezier from x1,y1 to x4,y4. Control points x2,y2 and x3,y3. Set the pixels on set=true otherwise clear them.
````GO
pd.CBezier(5,50,5,25,50,25,50,50)
````

----
#### Rectangle(x0, y0, x1, y1 int, set bool, fill bool)
Draw a rectangle from x0,y0 to x1,y1. Set the pixels on set=true otherwise clear them. Fill the rectangle on fill=true
````GO
pd.Rectangle(10,10,50,50,true)          //Rectangle
pd.Rectangle(10,10,50,50,true,true)     //Rectangle filled
````

----
#### Fill(x, y int, newC bool)
Floodfill, start at x,y
````GO
pd.Floodfill(22,22)     //Floodfill starting add 22,22
````

----
#### EllipseRect(x0, y0, x1, y1 int, set bool)
Draw ellipse in the given box defined by x0,y0 to x1,y1. Set the pixels on set=true otherwise clear them.
````GO
pd.EllipseRect(10,10,100,50)
````

----
#### Circle(x0, y0, r int, set bool)
Draw a circle at x0,y0 with the radius r. Set the pixels on set=true otherwise clear them.
````GO
pd.Circle(25,25,10,true)
````

----
#### DotArc(x0, y0, r int, a1, a2, step int, set bool)
Draw an dotted arc at x0,y0 with the radius r. a1 and a2 specify the degrees from and to. Set the pixels on set=true otherwise clear them.
The arc is drawn allways clockwise. So the example should show the upper half circle. (9 to 3 on the clock) 
````GO
pd.DotArc(100,100,25,270,90,true)
````

----
#### LineArc(x0, y0, r int, a1, a2, step int, set bool)
Draw an arc at x0,y0 with the radius r. a1 and a2 specify the degrees from and to. Set the pixels on set=true otherwise clear them.
The arc is drawn allways clockwise. So the example should show the lower half circle. (3 to 9 on the clock)
````GO
pd.LineArc(100,100,25,90,270,true)
````

----
#### LineRadius(x0,y0,r1,r2,a1 int, set bool)
Draw a line from an circle center point x0,y0 to the given direction in degrees a1.
The line is starting at the first radius r1 and ends at the second radius r2. Set the pixels on b=true otherwise clear them.
````GO
pd.LineRadius(50,50,20,50,45,true)
````

LineRadius(x0,y0,r1,r2,a1 int, bs bool)


----
#### SVGPath(xo,yo float64, s string, set bool, fscale ...float64)
Interprete and draw the path in s. Set the pixels on set=true otherwise clear them. Scale it by fscale.
````GO
pd.SVGPath(0, 0, "M 5, 60 c 25, -25 50, 25 75, 12 s 50, 50 75, -15", true)
````


## Font and Stamp Loading

----
#### FontPrint(font *PixelFont, x, y int, text string, set bool)
Print out the text string at x,y on the paint area, Set the pixels on set=true otherwise clear them.
````GO
pd.FontPrint(myFont,20,20,"This is a test", true)       //Write on the paint area, mode pixel set
pd.FontPrint(pd.GetFont("copper"),20,50,"(a+b)=x",true) //Use a stored font "copper"
````

----
#### Stamp(x, y int, stamp *PixelStamp, set bool, st bool)
Set the pixels on set=true otherwise clear them. Stamp mode if st=true. Stamp mode means, that the bitmap is transfered to the paint area without blending it together with the background.
````GO
pd.Stamp(myStamp,20,20,true)                    //Stamp in set Mode
pd.Stamp(pd.GetStamp("flower"),20,50,true,true) //Stamp a stored stamp with name "flower" in stamp mode
````

----
#### LoadFont(name string) *PixelFont
Load a front by the given name (including path). This is returning a PixelFont object which can be added via AddFont function. LoadFont as a other Load functions does return **nil** on file errors
````GO
myFont := pd.LoadFont("/home/user/fonts/smallfont.fnt")
if myFont == nil {
    return errors.New("File not found.")
}
````

----
#### SaveFont(name string, font *PixelFont, perm os.FileMode) error
Generated Fonts can be saved too.
````GO
pd.SaveFont("C:\\FONTS\\SMALLFONT.FNT")
````

----
#### LoadStamp(name string) *PixelStamp
Load a stamp by the given name (including path). This is returning a PixelStamp object which can be added via AddStamp function. LoadStamp as a other Load functions does return **nil** on file errors
````GO
myStamp := pd.LoadStamp("/home/user/stamps/duck.stp")
if myStamp == nil {
    return errors.New("File not found.")
}
````

----
#### SaveStamp(name string, stamp *PixelStamp, perm os.FileMode) error
````GO
pd.SaveStamp("C:\\STAMPS\\DUCK.STP")
````

----
#### AddFont(name string, font *PixelFont)
Adds the font to the PixelDing object
````GO
pd.AddFont("copper",myFont)     //Add the font object to the PixelDING object
````

----
#### GetFont(name string) *PixelFont
Get the font from the PixelDing object
````GO
myFont := pd.GetFont("copper")
````

----
#### AddStamp(name string, stamp *PixelStamp)
Add the stamp to the PixelDing object
````GO
pd.AddStamp("duck", myStamp)    //Add a stamp with the name "duck" to the PixelDING object
````

----
#### GetStamp(name string) *PixelStamp
Get the the stamp from the PixelDing object
````GO
myStamp := pd.GetStamp("duck")
````

----
#### Toggle(b bool)
Toggle the pixelmode, from set to clear. WARNING Deprecated ! If toggle is enabled, all pixel set operations are inverting the pixel on the position. If a pixel is already set, the pixel is cleared and vice versa.
````GO
pd.Toggle(true)
````

----

#### Scale(s float64)
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
#### X() int
Get the X size of the paint area. Together with the stamp dimensions better positioning calculation is possible
````GO
pd := pixelding.New(100,100)
cat := pd.LoadStamp("cat.stp")
pd.Stamp(pd.X()/2-cat.X()/2, pd.Y()/2-cat.Y()/2, myStamp, true)  //Stamp cat in centered position
````

----
#### Y() int
Get the Y size of the paint area, for details see X() function

----
#### X() int
Get the X size of the stamp
````GO
stampWidth := myStamp.X()
````

----
#### Y() int
Get the Y size of the stamp
````GO
stampHeight := myStamp.Y()
````

----
## Bonus Text Functions
Additional there are some bonus text functions. Frames, Lines and normal Text.
![pixelding](screenshots/frames.png "The PixelDING")

### Things good to know...
There are a few **predefined** frame strings as constants available
````GO
pixelding.SingleFrame   //Single Line unicode chars
pixelding.DoubleFrame   //Double Line unicode chars
pixelding.RoundFrame    //Round Line unicode chars
pixelding.BlockFrame    //Block Char unicode chars
pixelding.TextFrame     //Text Chars
````
The strings have a total of nine chars in the following oder:
> Left upper corner; Upper line; Right upper corner; Left line; Center char (unused); Right line; Left lower corner; Lower line; Right lower corner

The TextFrame is using a **bitmask** to specify the frame elements used. That allows L-shaped or U-shaped, etc. frames. The easiest way is to use this as a binary representation.line
````GO
lshaped := 0B000100110
````
That draws no upper corners and upper lines, but the left line, the lower left corner and the bottom line. See the chart picture as example.

All functions has the options to use two coordinates systems. If the paint area has the size of 100 by 100 the external buffer representation is using 50x50 chars strings.
Therefor a circle which has as center 30,30, that center for the text functions would be 15,15. That translation could be a mess and so there is a boolean parameter that enable translation.
If enabled, the coordinates used are the same as for the painting area, additional in that mode the scaling will be used too.

#### TextFrame(x1, y1, x2, y2 int, l string, bitmask int, btranslate ...bool)
Draw a Textframe.
````GO
pd.TextFrame(10,10,40,30,pixelding.SingleFrame,0B111111111)
pd.TextFrame(1,1,40,30,pixelding.SingleFrame,0B111111111,true)
````

#### TextLineH(x1, y1, x2, y2 int, l string, b ...bool)
Draw a horizontal line.
````GO
pd.TextLineH(15,15,30,15,pixelding.SingleFrame,0B111111111)
````
>Even if the parameters suggest the possibility of diagonal lines, that is not (yet) supported. The essential y coordinate is y1, y2 is not used.

#### TextLineV(x1, y1, x2, y2 int, l string, b ...bool)
Draw a horizontal line.
````GO
pd.TextLineV(15,15,30,15,pixelding.SingleFrame,0B111111111)
````
>Even if the parameters suggest the possibility of diagonal lines, that is not (yet) supported. The essential y coordinate is x1, x2 is not used.

#### Text(x, y int, text string, b ...bool)
````GO
pd.Text(11,11,"The PixelDING")          //Put the text on 11,11
pd.Text(22,22,"The PixelDING", true)    //Using paint area coordinates,
                                        //which is the same as the line above
````





