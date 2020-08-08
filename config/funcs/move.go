package funcs

import (
	"fmt"
	"strings"

	"github.com/lucasepe/jumble"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

var MoveFunc = func(tiles map[string]jumble.Tile) function.Function {

	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "origin",
				Type: cty.String,
			},
			{
				Name: "steps",
				Type: cty.Number,
			},
			{
				Name: "dir",
				Type: cty.String,
			},
		},
		Type: function.StaticReturnType(cty.Number),
		Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
			var origin string
			if err := gocty.FromCtyValue(args[0], &origin); err != nil {
				return cty.NumberIntVal(-1), err
			}

			var steps int64
			if err := gocty.FromCtyValue(args[1], &steps); err != nil {
				return cty.NumberIntVal(-1), err
			}

			var dir string
			if err := gocty.FromCtyValue(args[2], &dir); err != nil {
				return cty.NumberIntVal(-1), err
			}

			t, ok := tiles[origin]
			if !ok {
				return cty.NumberIntVal(-1), fmt.Errorf("tile (ID: %s) not found", origin)
			}

			row, col := t.Location()

			switch dir = strings.ToLower(dir); dir {
			case "north":
				return cty.NumberIntVal(int64(row) - steps), nil
			case "east":
				return cty.NumberIntVal(int64(col) + steps), nil
			case "south":
				return cty.NumberIntVal(int64(row) + steps), nil
			case "west":
				return cty.NumberIntVal(int64(col) - steps), nil
			}

			return cty.NumberIntVal(1), fmt.Errorf("unknow direction: %s", dir)
		},
	})
}

var RowOfFunc = func(tiles map[string]jumble.Tile) function.Function {

	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "origin",
				Type: cty.String,
			},
		},
		Type: function.StaticReturnType(cty.Number),
		Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
			var origin string
			if err := gocty.FromCtyValue(args[0], &origin); err != nil {
				return cty.NumberIntVal(-1), err
			}

			t, ok := tiles[origin]
			if !ok {
				return cty.NumberIntVal(-1), fmt.Errorf("tile (ID: %s) not found", origin)
			}

			row, _ := t.Location()
			return cty.NumberIntVal(int64(row)), nil
		},
	})
}

var ColOfFunc = func(tiles map[string]jumble.Tile) function.Function {

	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "origin",
				Type: cty.String,
			},
		},
		Type: function.StaticReturnType(cty.Number),
		Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
			var origin string
			if err := gocty.FromCtyValue(args[0], &origin); err != nil {
				return cty.NumberIntVal(-1), err
			}

			t, ok := tiles[origin]
			if !ok {
				return cty.NumberIntVal(-1), fmt.Errorf("tile (ID: %s) not found", origin)
			}

			_, col := t.Location()
			return cty.NumberIntVal(int64(col)), nil
		},
	})
}
