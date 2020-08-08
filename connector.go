package jumble

import (
	"math"
)

// Connector connect two or more tiles.
type Connector struct {
	Row int
	Col int

	color string

	strokeUp    bool
	strokeRight bool
	strokeDown  bool
	strokeLeft  bool

	arrowUp    bool
	arrowRight bool
	arrowDown  bool
	arrowLeft  bool
}

// ConnectorColor sets the connector strokes color
func ConnectorColor(hex string) func(c *Connector) {
	return func(c *Connector) {
		c.color = hex
	}
}

// ConnectorArrowUp enable the arrow up on the connector
func ConnectorArrowUp() func(c *Connector) {
	return func(c *Connector) {
		c.arrowUp = true
	}
}

// ConnectorArrowDown enable the arrow down on the connector
func ConnectorArrowDown() func(c *Connector) {
	return func(c *Connector) {
		c.arrowDown = true
	}
}

// ConnectorArrowLeft enable the arrow left on the connector
func ConnectorArrowLeft() func(c *Connector) {
	return func(c *Connector) {
		c.arrowLeft = true
	}
}

// ConnectorArrowRight enable the arrow right on the connector
func ConnectorArrowRight() func(c *Connector) {
	return func(c *Connector) {
		c.arrowRight = true
	}
}

// VerticalConnector returns a vertical connector
func VerticalConnector(row, col int, opts ...func(*Connector)) Connector {
	con := Connector{
		Row:        row,
		Col:        col,
		strokeUp:   true,
		strokeDown: true,
		color:      "#000000",
	}

	for _, opt := range opts {
		opt(&con)
	}

	return con
}

// HorizontalConnector returns a hotizontal connector
func HorizontalConnector(row, col int, opts ...func(*Connector)) Connector {
	con := Connector{
		Row:         row,
		Col:         col,
		strokeLeft:  true,
		strokeRight: true,
		color:       "#000000",
	}

	for _, opt := range opts {
		opt(&con)
	}

	return con
}

// ElbowRightDownConnector returns a right (→) down (↓) elbow connector
func ElbowRightDownConnector(row, col int, opts ...func(*Connector)) Connector {
	con := Connector{
		Row:         row,
		Col:         col,
		strokeRight: true,
		strokeDown:  true,
		color:       "#000000",
	}

	for _, opt := range opts {
		opt(&con)
	}

	return con
}

// ElbowRightUpConnector returns a right (→) up (↑) elbow connector
func ElbowRightUpConnector(row, col int, opts ...func(*Connector)) Connector {
	con := Connector{
		Row:         row,
		Col:         col,
		strokeRight: true,
		strokeUp:    true,
		color:       "#000000",
	}

	for _, opt := range opts {
		opt(&con)
	}

	return con
}

// ElbowLeftUpConnector returns left (←) up (↑) elbow connector
func ElbowLeftUpConnector(row, col int, opts ...func(*Connector)) Connector {
	con := Connector{
		Row:        row,
		Col:        col,
		strokeLeft: true,
		strokeUp:   true,
		color:      "#000000",
	}

	for _, opt := range opts {
		opt(&con)
	}

	return con
}

// ElbowLeftDownConnector returns left (←) down (↓) elbow connector
func ElbowLeftDownConnector(row, col int, opts ...func(*Connector)) Connector {
	con := Connector{
		Row:        row,
		Col:        col,
		strokeLeft: true,
		strokeDown: true,
		color:      "#000000",
	}

	for _, opt := range opts {
		opt(&con)
	}

	return con
}

// TeeUpConnector returns a tee connector (← ↑ →)
func TeeUpConnector(row, col int, opts ...func(*Connector)) Connector {
	con := Connector{
		Row:         row,
		Col:         col,
		strokeUp:    true,
		strokeLeft:  true,
		strokeRight: true,
		color:       "#000000",
	}

	for _, opt := range opts {
		opt(&con)
	}

	return con
}

// TeeDownConnector returns a tee connector (← ↓ →)
func TeeDownConnector(row, col int, opts ...func(*Connector)) Connector {
	con := Connector{
		Row:         row,
		Col:         col,
		strokeDown:  true,
		strokeLeft:  true,
		strokeRight: true,
		color:       "#000000",
	}

	for _, opt := range opts {
		opt(&con)
	}

	return con
}

// TeeLeftConnector returns a tee connector (← ↑ ↓)
func TeeLeftConnector(row, col int, opts ...func(*Connector)) Connector {
	con := Connector{
		Row:        row,
		Col:        col,
		strokeLeft: true,
		strokeDown: true,
		strokeUp:   true,
		color:      "#000000",
	}

	for _, opt := range opts {
		opt(&con)
	}

	return con
}

// TeeRightConnector returns a tee connector (↑ ↓ →)
func TeeRightConnector(row, col int, opts ...func(*Connector)) Connector {
	con := Connector{
		Row:         row,
		Col:         col,
		strokeRight: true,
		strokeDown:  true,
		strokeUp:    true,
		color:       "#000000",
	}

	for _, opt := range opts {
		opt(&con)
	}

	return con
}

// CrossConnector returns a tee connector (← ↑ ↓ →)
func CrossConnector(row, col int, opts ...func(*Connector)) Connector {
	con := Connector{
		Row:         row,
		Col:         col,
		strokeRight: true,
		strokeDown:  true,
		strokeUp:    true,
		strokeLeft:  true,
		color:       "#000000",
	}

	for _, opt := range opts {
		opt(&con)
	}

	return con
}

// Location returns the grid position (row, col)
func (c *Connector) Location() (int, int) {
	return c.Row, c.Col
}

// Plot draws a connector on the grid.
func (c *Connector) Plot(g *Grid) error {
	if err := g.VerifyInBounds(c.Row, c.Col); err != nil {
		return err
	}

	const strokeMultiplier float64 = 0.03

	w, h := g.CellSize(), g.CellSize()
	lw := strokeMultiplier * math.Max(w, h)

	center := g.CellCenter(c.Row, c.Col)

	// Draw the shape.
	dc := g.Context()
	dc.Push()
	dc.SetLineWidth(lw)
	dc.SetHexColor(c.color)
	dc.SetLineCapRound()
	dc.Translate(center.X, center.Y)

	if c.strokeUp {
		dc.MoveTo(0.0, 0.0)
		dc.LineTo(0.0, -0.5*h)
		dc.Stroke()
	}

	if c.arrowUp {
		as := 0.15 * w
		x, y := 0.0, -0.5*h
		dc.MoveTo(x, y)
		dc.LineTo(x-as, y+as)
		dc.LineTo(x+as, y+as)
		dc.LineTo(x, y)
		dc.Fill()
	}

	if c.strokeRight {
		dc.MoveTo(0.0, 0.0)
		dc.LineTo(0.5*w, 0.0)
		dc.Stroke()
	}
	if c.arrowRight {
		as := 0.15 * w
		x, y := 0.5*w, 0.0
		dc.MoveTo(x, y)
		dc.LineTo(x-as, y-as)
		dc.LineTo(x-as, y+as)
		dc.LineTo(x, y)
		dc.Fill()
	}

	if c.strokeDown {
		dc.MoveTo(0.0, 0.0)
		dc.LineTo(0.0, 0.5*h)
		dc.Stroke()
	}
	if c.arrowDown {
		as := 0.15 * w
		x, y := 0.0, 0.5*h
		dc.MoveTo(x, y)
		dc.LineTo(x-as, y-as)
		dc.LineTo(x+as, y-as)
		dc.LineTo(x, y)
		dc.Fill()
	}

	if c.strokeLeft {
		dc.MoveTo(0.0, 0.0)
		dc.LineTo(-0.5*w, 0.0)
		dc.Stroke()
	}
	if c.arrowLeft {
		as := 0.15 * w
		x, y := -0.5*w, 0.0
		dc.MoveTo(x, y)
		dc.LineTo(x+as, y-as)
		dc.LineTo(x+as, y+as)
		dc.LineTo(x, y)
		dc.Fill()
	}

	dc.Pop()

	return nil
}
