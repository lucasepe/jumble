package jumble

// FrameOval sets the frame shape to
// oval (default is rectangle)
func FrameOval(val bool) func(f *Frame) {
	return func(fr *Frame) {
		fr.oval = val
	}
}

// FrameColor sets the frame color
func FrameColor(hex string) func(f *Frame) {
	return func(fr *Frame) {
		fr.color = hex
	}
}

// FrameStroke enables the frame stroke
func FrameStroke(val bool) func(f *Frame) {
	return func(fr *Frame) {
		fr.stroke = val
	}
}

// FrameDashes sets the frame dashes
func FrameDashes(val float64) func(f *Frame) {
	return func(fr *Frame) {
		fr.dashes = val
	}
}

// FrameStrokeWidth sets the stroke width
func FrameStrokeWidth(val float64) func(f *Frame) {
	return func(fr *Frame) {
		fr.strokeWidth = val
	}
}

// Frame represents a frame on the grid.
type Frame struct {
	Left        int
	Top         int
	Right       int
	Bottom      int
	dashes      float64
	color       string
	stroke      bool
	strokeWidth float64
	oval        bool
}

// NewFrame returns a new frame
func NewFrame(l, t int, r, b int, opts ...func(*Frame)) Frame {
	res := Frame{
		Left: l, Top: t,
		Right: r, Bottom: b,
		oval:        false,
		color:       "#000000",
		stroke:      true,
		strokeWidth: 1,
		dashes:      5,
	}

	for _, opt := range opts {
		opt(&res)
	}

	return res
}

// Location returns the grid position (row, col)
func (fr *Frame) Location() (int, int) {
	return fr.Left, fr.Right
}

// Plot draws a frame accross the specified cells
func (fr *Frame) Plot(g *Grid) error {
	if err := g.VerifyInBounds(fr.Left, fr.Top); err != nil {
		return err
	}

	p1 := g.CellCenter(fr.Left, fr.Top)
	p2 := g.CellCenter(fr.Right, fr.Bottom)
	dx, dy := p2.X-p1.X, p2.Y-p1.Y

	dc := g.Context()

	dc.Push()
	if fr.dashes > 0 {
		dc.SetDash(fr.dashes)
	} else {
		dc.SetDash()
	}

	dc.SetLineWidth(fr.strokeWidth)
	dc.SetHexColor(fr.color)

	if fr.oval {
		rx, ry := 0.5*dx, 0.5*dy
		cx, cy := p1.X+rx, p1.Y+ry
		dc.DrawEllipse(cx, cy, rx, ry)
	} else {
		dc.DrawRectangle(p1.X, p1.Y, dx, dy)
	}

	if fr.stroke {
		dc.Stroke()
	} else {
		dc.Fill()
	}
	dc.Pop()

	return nil
}
