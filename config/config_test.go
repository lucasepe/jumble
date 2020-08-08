package config

import (
	"testing"
)

func TestConfigParser(t *testing.T) {
	demo := `
rows = 10
cols = 10

var "baseUri" {
	value = "/home/lus/Pictures/AWS-Architecture-Icons/PNG"
}


tile "icon" "agw" {
	row = 3
	col = 4
	fit = true
	uri = "${mkURI(var.baseUri, "aws_api_gateway.png")}"
}	


tile "horizontal_line" "hl1" {
	row = "${move("agw", 2, "north")}"
	col = 4
}	

`
	cfg, err := Decode([]byte(demo), "demo.hcl")
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range cfg.Tiles {
		t.Logf("%s: %v\n", k, v)
	}

}
