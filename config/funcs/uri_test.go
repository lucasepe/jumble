package funcs

import (
	"fmt"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestMakeURI(t *testing.T) {
	tests := []struct {
		Base cty.Value
		Path cty.Value
		Want cty.Value
		Err  bool
	}{
		{
			cty.StringVal("http://www.google.com/images/"),
			cty.StringVal("2020/camera/pic1.png"),
			cty.StringVal("http://www.google.com/images/2020/camera/pic1.png"),
			false,
		},
		{
			cty.StringVal("/Pictures/AWS-Architecture-Icons/PNG"),
			cty.StringVal("aws_api_gateway.png"),
			cty.StringVal("/Pictures/AWS-Architecture-Icons/PNG/aws_api_gateway.png"),
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("makeURI(%#v)", test.Base), func(t *testing.T) {
			got, err := MakeURIFunc.Call([]cty.Value{test.Base, test.Path})

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}
