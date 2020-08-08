package funcs

import (
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

// MakeURIFunc constructs a function that joins a baseUri with a path.
var MakeURIFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "elem1",
			Type: cty.String,
		},
		{
			Name: "elem2",
			Type: cty.String,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		var elem1, elem2 string
		if err := gocty.FromCtyValue(args[0], &elem1); err != nil {
			return cty.UnknownVal(cty.String), err
		}
		if err := gocty.FromCtyValue(args[1], &elem2); err != nil {
			return cty.UnknownVal(cty.String), err
		}

		if strings.HasPrefix(elem1, "http") {
			u, err := url.Parse(elem1)
			if err != nil {
				return cty.UnknownVal(cty.String), err
			}
			u.Path = path.Join(u.Path, elem2)

			return cty.StringVal(u.String()), nil
		}

		return cty.StringVal(filepath.Join(elem1, elem2)), nil
	},
})
