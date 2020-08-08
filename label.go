package jumble

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

// Label wraps a string.
type Label struct {
	Row        int
	Col        int
	text       string
	fontSize   float64
	color      string
	background string
	angle      float64
}

// NewLabel returns a new label
func NewLabel(row, col int, text string, opts ...func(*Label)) Label {
	res := Label{
		Row: row, Col: col,
		text:  text,
		color: "#000000",
	}

	for _, opt := range opts {
		opt(&res)
	}

	return res
}

// LabelFontSize sets the label font size
func LabelFontSize(val float64) func(*Label) {
	return func(lab *Label) {
		lab.fontSize = val
	}
}

// LabelColor sets the label font color
func LabelColor(hex string) func(*Label) {
	return func(lab *Label) {
		lab.color = hex
	}
}

// LabelBackground sets the label background color
func LabelBackground(hex string) func(*Label) {
	return func(lab *Label) {
		lab.background = hex
	}
}

// LabelAngle sets the label rotation angle in degree
func LabelAngle(val float64) func(*Label) {
	return func(lab *Label) {
		lab.angle = val
	}
}

// Location returns the grid position (row, col)
func (lab *Label) Location() (int, int) {
	return lab.Row, lab.Col
}

// Plot draws a string in a cell
func (lab *Label) Plot(g *Grid) error {
	if err := g.VerifyInBounds(lab.Row, lab.Col); err != nil {
		return err
	}

	if lab.fontSize == 0 {
		lab.fontSize = 0.3 * float64(g.cellSize)
	}

	face := truetype.NewFace(g.font, &truetype.Options{Size: lab.fontSize})

	center := g.CellCenter(lab.Row, lab.Col)

	dc := g.Context()

	if lab.background != "" {
		sw, sh := g.ctx.MeasureString(lab.text)
		pad := lab.fontSize

		dc.Push()
		dc.SetHexColor(lab.background)
		dc.DrawRoundedRectangle(center.X-(pad+sw/2), center.Y-(pad+sh)/2, sw+2*pad, sh+pad, 2)
		dc.Fill()
		dc.Pop()
	}

	dc.Push()
	dc.SetFontFace(face)
	dc.SetHexColor(lab.color)
	dc.RotateAbout(gg.Radians(lab.angle), center.X, center.Y)
	dc.DrawStringAnchored(lab.text, center.X, center.Y, 0.5, 0.5)
	dc.Pop()

	return nil
}
