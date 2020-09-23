package identicon

import (
	"crypto/sha256"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"image/png"
	"io"
)

type Identicon struct {
	hash       [32]byte
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

func Generate(input []byte) Identicon {
	identicon := hashInput(input)
	identicon = pipe(identicon, buildGrid, pickColor, filterOddSquares, buildPixelMap)

	return identicon
}

func hashInput(input []byte) Identicon {
	checkSum := sha256.Sum256(input)
	return Identicon{
		hash: checkSum,
	}
}

func pipe(identicon Identicon, Applies ...Apply) Identicon {
	for _, f := range Applies {
		identicon = f(identicon)
	}
	return identicon
}

// 產生網格
func buildGrid(identicon Identicon) Identicon {
	var grid []byte
	// Loop over the hash from the identicon, and increment with 3 (Chunk the array in 3 parts)
	for i := 0; i < len(identicon.hash) && i+3 <= len(identicon.hash)-1; i += 3 {
		chunk := make([]byte, 5)
		copy(chunk, identicon.hash[i:i+3]) // Copy the items from the old array to the new array
		chunk[3] = chunk[1]                // mirror the second value in the chunk
		chunk[4] = chunk[0]                // mirror the first value in the chunk
		grid = append(grid, chunk...)
	}

	identicon.grid = grid // set the grid property on the identicon
	return identicon
}

// 設置填充顏色
func pickColor(identicon Identicon) Identicon {
	rgb := [3]byte{}
	copy(rgb[:], identicon.hash[:3]) // copy first 3 values to the rgb array from identicon hash
	identicon.color = rgb            // rgb array assign to the color value of Identicon

	return identicon
}

// 濾掉奇數網格
func filterOddSquares(identicon Identicon) Identicon {
	var grid []GridPoint
	for i, code := range identicon.grid {
		if code%2 == 0 { // check is even or odd
			point := GridPoint{
				value: code,
				index: i,
			}

			grid = append(grid, point)
		}
	}
	// set the property
	identicon.gridPoints = grid
	return identicon
}

// 建立圖示
func buildPixelMap(identicon Identicon) Identicon {
	var drawingPoints []DrawingPoint

	// Closure, this function returns a Drawingpoint
	pixelFunc := func(p GridPoint) DrawingPoint {
		horizontal := (p.index % 5) * 50 // use the index from the gridpoint to calculate the horizontal dimension
		vertical := (p.index / 5) * 50   // use the index from the gridpoint to calculate the vertical dimension

		topLeft := Point{horizontal, vertical}               // top left
		bottomRight := Point{horizontal + 50, vertical + 50} // bottom right（top left point +50)

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
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFillColor(col) // 設值顏色
	gc.MoveTo(x1, y1)    // 讓圖片靠左上角
	gc.LineTo(x1, y1)
	gc.LineTo(x1, y2)

	gc.MoveTo(x2, y1) // move to the right in the image
	gc.LineTo(x2, y1)
	gc.LineTo(x2, y2)
	gc.SetLineWidth(0)
	gc.FillStroke()
}

// 產生Identicon 圖示
func (i Identicon) WriteImage(w io.Writer) error {
	var img = image.NewRGBA(image.Rect(0, 0, 250, 250))
	col := color.RGBA{R: i.color[0], G: i.color[1], B: i.color[2], A: 255}

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
