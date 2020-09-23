package identicon

import (
	"crypto/md5"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"image/png"
	"io"
)

type Identicon struct {
	hash       [16]byte
	color      [3]byte
	grid       []byte // New property to hold the grid
	gridPoints []GridPoint
	pixelMap   []DrawingPoint // pixelMap for drawing
}

type GridPoint struct {
	value byte
	index int
}

type Point struct {
	x, y int
}

type DrawingPoint struct {
	topLeft     Point
	bottomRight Point
}

type Apply func(Identicon) Identicon

func hashInput(input []byte) Identicon {
	checkSum := md5.Sum(input) // generate checksum from input

	return Identicon{
		hash: checkSum,
	}
}

func pickColor(identicon Identicon) Identicon {
	rgb := [3]byte{}                 // first we make a byte array with length 3
	copy(rgb[:], identicon.hash[:3]) // next we copy the first 3 values from the hash to the rgb array
	identicon.color = rgb            // we than assign it to the color value

	return identicon
}

func buildGrid(identicon Identicon) Identicon {

	var grid []byte // Create empty grid
	// Loop over the hash from the identicon
	// Increment with 3 (Chunk the array in 3 parts)
	// this ensures we wont get array out of bounds error and will retrieve exactly 5 chunks of 3
	for i := 0; i < len(identicon.hash) && i+3 <= len(identicon.hash)-1; i += 3 {
		chunk := make([]byte, 5)           // Create a placeholder for the chunk
		copy(chunk, identicon.hash[i:i+3]) // Copy the items from the old array to the new array
		chunk[3] = chunk[1]                // mirror the second value in the chunk
		chunk[4] = chunk[0]                // mirror the first value in the chunk
		grid = append(grid, chunk...)      // append the chunk to the grid
	}

	identicon.grid = grid // set the grid property on the identicon
	return identicon
}

func filterOddSquares(identicon Identicon) Identicon {
	var grid []GridPoint                  // create a placeholder to hold the values of the loop
	for i, code := range identicon.grid { // loop over the grid
		if code%2 == 0 { // check if the value is odd or not
			// create a new Gridpoint where we save the value and the index in the grid
			point := GridPoint{
				value: code,
				index: i,
			}

			grid = append(grid, point) // append the item to the new grid
		}
	}
	// set the property
	identicon.gridPoints = grid
	return identicon
}

func buildPixelMap(identicon Identicon) Identicon {
	var drawingPoints []DrawingPoint // define placeholder for drawingpoints

	// Closure, this function returns a Drawingpoint
	pixelFunc := func(p GridPoint) DrawingPoint {
		horizontal := (p.index % 5) * 50 // This is the formula, we use the index from the gridpoint to calculate the horizontal dimension
		vertical := (p.index / 5) * 50   // This is the formula, we use the index from the gridpoint to calculate the vertical dimension

		topLeft := Point{horizontal, vertical}               // this is the topleft point with x and the y
		bottomRight := Point{horizontal + 50, vertical + 50} // the bottom right point is just the topleft point +50 because 1 block in the grid is 50x50

		return DrawingPoint{
			topLeft,
			bottomRight,
		}
	}

	for _, gridPoint := range identicon.gridPoints {
		// for every gridPoint we calculate the drawingpoints and we add them to the array
		drawingPoints = append(drawingPoints, pixelFunc(gridPoint))
	}
	identicon.pixelMap = drawingPoints // set the drawingpoint value on the identicon
	return identicon                   // return the modified identicon
}

func rect(img *image.RGBA, col color.Color, x1, y1, x2, y2 float64) {
	gc := draw2dimg.NewGraphicContext(img) // Prepare new image context
	gc.SetFillColor(col)                   // set the color
	gc.MoveTo(x1, y1)                      // move to the topleft in the image
	// Draw the lines for the dimensions
	gc.LineTo(x1, y1)
	gc.LineTo(x1, y2)
	gc.MoveTo(x2, y1) // move to the right in the image
	// Draw the lines for the dimensions
	gc.LineTo(x2, y1)
	gc.LineTo(x2, y2)
	// Set the linewidth to zero
	gc.SetLineWidth(0)
	// Fill the stroke so the rectangle will be filled
	gc.FillStroke()
}

func pipe(identicon Identicon, funcs ...Apply) Identicon {
	for _, applyer := range funcs {
		identicon = applyer(identicon)
	}
	return identicon
}

func Generate(input []byte) Identicon {
	identicon := hashInput(input)
	identicon = pipe(identicon, pickColor, buildGrid, filterOddSquares, buildPixelMap)

	return identicon
}

func (i Identicon) WriteImage(w io.Writer) error {
	// We create our default image containing a 250x250 rectangle
	var img = image.NewRGBA(image.Rect(0, 0, 250, 250))
	// We retrieve the color from the color property on the identicon
	col := color.RGBA{R: i.color[0], G: i.color[1], B: i.color[2], A: 255}

	// Loop over the pixelmap and call the rect function with the img, color and the dimensions
	for _, pixel := range i.pixelMap {
		rect(
			img,
			col,
			float64(pixel.topLeft.x),
			float64(pixel.topLeft.y),
			float64(pixel.bottomRight.x),
			float64(pixel.bottomRight.y),
		)
	}

	return png.Encode(w, img)
}
