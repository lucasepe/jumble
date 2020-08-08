package jumble

import (
	"github.com/disintegration/imaging"
)

// Icon wraps an image.
type Icon struct {
	Row int
	Col int
	Fit bool
	URI string
}

// NewIcon returns a new icon from the specified uri
func NewIcon(row, col int, uri string) Icon {
	return Icon{
		Row: row, Col: col,
		URI: uri,
		Fit: true,
	}
}

// Location returns the grid position (row, col)
func (ic *Icon) Location() (int, int) {
	return ic.Row, ic.Col
}

// Plot draws a image (eventually rescaling) in a cell.
func (ic *Icon) Plot(g *Grid) error {
	if err := g.VerifyInBounds(ic.Row, ic.Col); err != nil {
		return err
	}

	im, err := LoadImage(ic.URI)
	if err != nil {
		return err
	}

	if ic.Fit {
		b := im.Bounds()
		size := b.Max.X
		if b.Max.Y > size {
			size = b.Max.Y
		}

		if g.cellSize < size {
			im = imaging.Resize(im, g.cellSize, g.cellSize, imaging.Lanczos)
		}
	}

	center := g.CellCenter(ic.Row, ic.Col)

	dc := g.Context()
	dc.Push()
	//g.ctx.RotateAbout(gg.Radians(alpha), center.X, center.Y)
	dc.DrawImageAnchored(im, int(center.X), int(center.Y), 0.5, 0.5)
	dc.Pop()

	return nil
}
