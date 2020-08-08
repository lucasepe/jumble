package config

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/lucasepe/jumble"
	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"

	"github.com/lucasepe/jumble/config/data"
	"github.com/lucasepe/jumble/config/funcs"
)

// Config defines a grid layout with all the tiles.
type Config struct {
	Rows       int
	Cols       int
	Margin     int
	Background string
	Grid       bool
	Border     bool
	Hints      bool

	Tiles map[string]jumble.Tile
}

// rootHCL is the helper struct for parsing our HCL file.
type rootHCL struct {
	Rows       int    `hcl:"rows"`
	Cols       int    `hcl:"cols"`
	Margin     int    `hcl:"margin,optional"`
	Background string `hcl:"background,optional"`
	Grid       bool   `hcl:"grid,optional"`
	Border     bool   `hcl:"border,optional"`
	Hints      bool   `hcl:"hints,optional"`

	Variables []*struct {
		Name  string         `hcl:"name,label"`
		Value hcl.Attributes `hcl:"value,remain"`
	} `hcl:"var,block"`

	Tiles []*struct {
		Kind    string   `hcl:"type,label"`
		ID      string   `hcl:"id,label"`
		HCLBody hcl.Body `hcl:",remain"`
	} `hcl:"tile,block"`
}

// DecodeURI parses the given uri with our HCL content.
func DecodeURI(uri string) (Config, error) {
	const bytesLimit = 100 * 1024

	if strings.HasPrefix(uri, "http") {
		body, err := data.FetchFromURI(uri, bytesLimit)
		if err != nil {
			return Config{}, errors.Wrapf(err, "fetching '%s'", uri)
		}

		return Decode(body, uri)
	}

	body, err := data.FetchFromFile(uri, bytesLimit)
	if err != nil {
		return Config{}, errors.Wrapf(err, "fetching '%s'", uri)
	}

	return Decode(body, uri)
}

// Decode parses the given buffer with our HCL content.
// The `uri` string is just for debugging purposes.
// On success this function returns a Config struct.
func Decode(data []byte, uri string) (Config, error) {

	// Instantiate an HCL parser with the source byte slice.
	parser := hclparse.NewParser()

	srcHCL, diags := parser.ParseHCL(data, uri)
	if diags.HasErrors() {
		return Config{}, fmt.Errorf("error parsing HCL file: %w", diags)
	}

	// Start the first pass of decoding
	var root rootHCL
	if diags := gohcl.DecodeBody(srcHCL.Body, nil, &root); diags.HasErrors() {
		return Config{}, fmt.Errorf("error decoding HCL configuration: %w", diags)
	}

	// Decode all variables
	variables := map[string]cty.Value{}
	for _, v := range root.Variables {
		if len(v.Value) == 0 {
			continue
		}

		val, diags := v.Value["value"].Expr.Value(nil)
		if diags.HasErrors() {
			return Config{}, fmt.Errorf("error decoding HCL configuration: %w", diags)
		}

		variables[v.Name] = val
	}

	cfg := Config{
		Rows:       root.Rows,
		Cols:       root.Cols,
		Background: root.Background,
		Margin:     root.Margin,
		Grid:       root.Grid,
		Border:     root.Border,
		Hints:      root.Hints,
		Tiles:      map[string]jumble.Tile{},
	}

	// Call a helper function which creates an HCL context for use in
	// decoding the parsed HCL.
	evalContext, err := createContext(variables, cfg.Tiles)
	if err != nil {
		return Config{}, fmt.Errorf("error creating HCL evaluation context: %w", err)
	}

	// Start decoding
	for _, tile := range root.Tiles {
		if len(strings.TrimSpace(tile.ID)) == 0 {
			if tile.ID, err = shortid.Generate(); err != nil {
				return Config{}, err
			}
		}

		switch t := tile.Kind; t {
		case "icon":
			el, err := decodeIcon(tile.HCLBody, evalContext)
			if err != nil {
				return Config{}, err
			}
			cfg.Tiles[tile.ID] = &el

		case "label":
			el, err := decodeLabel(tile.HCLBody, evalContext)
			if err != nil {
				return Config{}, err
			}
			cfg.Tiles[tile.ID] = &el

		case "frame":
			el, err := decodeFrame(tile.HCLBody, evalContext)
			if err != nil {
				return Config{}, err
			}
			cfg.Tiles[tile.ID] = &el

		default:
			el, err := decodeConnector(tile.HCLBody, evalContext, tile.Kind)
			if err != nil {
				if _, ok := err.(*unknowTileTypeError); ok {
					continue
				}
				return Config{}, err
			}
			cfg.Tiles[tile.ID] = &el
		}
	}

	return cfg, nil
}

// createContext is a helper function that creates an *hcl.EvalContext to be
// used in decoding HCL. It add all variables to the eval context.
// It also creates custom functions.
func createContext(vars map[string]cty.Value, tiles map[string]jumble.Tile) (*hcl.EvalContext, error) {
	return &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"var": cty.ObjectVal(vars),
		},
		Functions: map[string]function.Function{
			"mkURI":    funcs.MakeURIFunc,
			"row":      funcs.RowOfFunc(tiles),
			"col":      funcs.ColOfFunc(tiles),
			"add":      stdlib.AddFunc,
			"subtract": stdlib.SubtractFunc,
		},
	}, nil
}

// decodeIcon decode the HCL 'icon' block
func decodeIcon(body hcl.Body, ctx *hcl.EvalContext) (jumble.Icon, error) {
	var tmp struct {
		Row int    `hcl:"row"`
		Col int    `hcl:"col"`
		Fit bool   `hcl:"fit,optional"`
		URI string `hcl:"uri"`
	}
	tmp.Fit = true

	if diags := gohcl.DecodeBody(body, ctx, &tmp); diags.HasErrors() {
		return jumble.Icon{}, fmt.Errorf("error decoding HCL configuration: %w", diags)
	}

	return jumble.Icon{
		Row: tmp.Row, Col: tmp.Col,
		Fit: tmp.Fit,
		URI: tmp.URI,
	}, nil
}

// decodeFrame decode the HCL 'frame' block
func decodeFrame(body hcl.Body, ctx *hcl.EvalContext) (jumble.Frame, error) {
	var tmp struct {
		Left        int     `hcl:"left"`
		Top         int     `hcl:"top"`
		Right       int     `hcl:"right"`
		Bottom      int     `hcl:"bottom"`
		Dashes      float64 `hcl:"dashes,optional"`
		Color       string  `hcl:"color,optional"`
		Stroke      bool    `hcl:"stroke,optional"`
		StrokeWidth float64 `hcl:"stroke_width,optional"`
		Oval        bool    `hcl:"oval,optional"`
	}

	if diags := gohcl.DecodeBody(body, ctx, &tmp); diags.HasErrors() {
		return jumble.Frame{}, fmt.Errorf("error decoding HCL configuration: %w", diags)
	}

	return jumble.NewFrame(tmp.Left, tmp.Top, tmp.Right, tmp.Bottom,
		jumble.FrameOval(tmp.Oval),
		jumble.FrameColor(tmp.Color),
		jumble.FrameStroke(tmp.Stroke),
		jumble.FrameDashes(tmp.Dashes),
		jumble.FrameStrokeWidth(tmp.StrokeWidth),
	), nil
}

// decodeLabel decode the HCL 'label' block
func decodeLabel(body hcl.Body, ctx *hcl.EvalContext) (jumble.Label, error) {
	var tmp struct {
		Row        int     `hcl:"row"`
		Col        int     `hcl:"col"`
		Text       string  `hcl:"text"`
		FontSize   float64 `hcl:"font_size,optional"`
		Color      string  `hcl:"color,optional"`
		Background string  `hcl:"background,optional"`
		Angle      float64 `hcl:"angle,optional"`
	}

	if diags := gohcl.DecodeBody(body, ctx, &tmp); diags.HasErrors() {
		return jumble.Label{}, fmt.Errorf("error decoding HCL configuration: %w", diags)
	}

	return jumble.NewLabel(tmp.Row, tmp.Col, tmp.Text,
		jumble.LabelColor(tmp.Color),
		jumble.LabelBackground(tmp.Background),
		jumble.LabelFontSize(tmp.FontSize),
		jumble.LabelAngle(tmp.Angle)), nil
}

// decodeConnectors decode all the HCL connector block
func decodeConnector(body hcl.Body, ctx *hcl.EvalContext, kind string) (jumble.Connector, error) {
	var tmp struct {
		Row        int    `hcl:"row"`
		Col        int    `hcl:"col"`
		Color      string `hcl:"color,optional"`
		ArrowUp    bool   `hcl:"arrow_up,optional"`
		ArrowRight bool   `hcl:"arrow_right,optional"`
		ArrowDown  bool   `hcl:"arrow_down,optional"`
		ArrowLeft  bool   `hcl:"arrow_left,optional"`
	}

	if diags := gohcl.DecodeBody(body, ctx, &tmp); diags.HasErrors() {
		return jumble.Connector{}, fmt.Errorf("error decoding HCL configuration: %w", diags)
	}

	var res jumble.Connector

	switch kind {
	case "cross":
		res = jumble.CrossConnector(tmp.Row, tmp.Col)
	case "horizontal_line":
		res = jumble.HorizontalConnector(tmp.Row, tmp.Col)
	case "vertical_line":
		res = jumble.VerticalConnector(tmp.Row, tmp.Col)
	case "elbow_right_up":
		res = jumble.ElbowRightUpConnector(tmp.Row, tmp.Col)
	case "elbow_right_down":
		res = jumble.ElbowRightDownConnector(tmp.Row, tmp.Col)
	case "elbow_left_down":
		res = jumble.ElbowLeftDownConnector(tmp.Row, tmp.Col)
	case "elbow_left_up":
		res = jumble.ElbowLeftUpConnector(tmp.Row, tmp.Col)
	case "tee_left":
		res = jumble.TeeLeftConnector(tmp.Row, tmp.Col)
	case "tee_down":
		res = jumble.TeeDownConnector(tmp.Row, tmp.Col)
	case "tee_right":
		res = jumble.TeeRightConnector(tmp.Row, tmp.Col)
	case "tee_up":
		res = jumble.TeeUpConnector(tmp.Row, tmp.Col)
	default:
		return jumble.Connector{}, &unknowTileTypeError{
			err: fmt.Errorf("unknown type: %s", kind),
		}
	}

	if tmp.ArrowUp {
		jumble.ConnectorArrowUp()(&res)
	}
	if tmp.ArrowRight {
		jumble.ConnectorArrowRight()(&res)
	}
	if tmp.ArrowDown {
		jumble.ConnectorArrowDown()(&res)
	}
	if tmp.ArrowLeft {
		jumble.ConnectorArrowLeft()(&res)
	}

	return res, nil
}

// unknowTileTypeError custom error
//to identify unknown tile types
type unknowTileTypeError struct {
	err error
}

// Error implement the standard library
// interface type for errors.
func (e *unknowTileTypeError) Error() string {
	return e.err.Error()
}
