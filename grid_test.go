package jumble

import (
	"bytes"
	"encoding/base64"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fogleman/gg"
	"github.com/stretchr/testify/assert"
)

func TestNewGrid(t *testing.T) {
	grid, err := NewGrid(12, 10, 64)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 12, grid.rows)
	assert.Equal(t, 10, grid.cols)
	assert.Equal(t, 64, grid.cellSize)

	assert.NotNil(t, grid.Context())
	assert.Equal(t, float64(64), grid.CellSize())
}

func TestAsset(t *testing.T) {

	LoadImage("assets://aws/aws_lambda.png")
}

func TestGridCellCenter(t *testing.T) {
	tests := []struct {
		row  int
		col  int
		want gg.Point
	}{
		{2, 2, gg.Point{X: 160, Y: 160}},
		{5, 4, gg.Point{X: 288, Y: 352}},
		{6, 7, gg.Point{X: 480, Y: 416}},
		{9, 3, gg.Point{X: 224, Y: 608}},
	}

	grid, err := NewGrid(12, 10, 64)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := grid.CellCenter(tt.row, tt.col)
			if got != tt.want {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}

func TestGridVerifyInBounds(t *testing.T) {
	tests := []struct {
		row  int
		col  int
		want string
	}{
		{2, 2, ""},
		{5, 4, "cell (5, 4) is out of bounds"},
		{6, 7, "cell (6, 7) is out of bounds"},
		{1, 3, ""},
	}

	grid, err := NewGrid(5, 5, 64)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := grid.VerifyInBounds(tt.row, tt.col)
			if got != nil && got.Error() != tt.want {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}

func TestGridLayout(t *testing.T) {
	grid, err := NewGrid(4, 4, 12, GridTitle("test"))
	if err != nil {
		t.Fatal(err)
	}

	grid.DrawGrid()
	grid.DrawCoords()
	grid.DrawBorder()

	var data bytes.Buffer
	if err := grid.EncodePNG(&data); err != nil {
		t.Fatal(err)
	}

	str := base64.StdEncoding.EncodeToString(data.Bytes())
	//t.Logf(str)
	assert.True(t, strings.HasPrefix(str, "iVBORw0KGgoAAAANSUhEUgAAAGAAAABgCAIAAABt+uBvAAAErElEQVR4Aeyb3U7qShiG"))
}

func TestGrid(t *testing.T) {
	grid, err := NewGrid(12, 10, 72, GridTitle("test"))
	if err != nil {
		t.Fatal(err)
	}

	iconsPath := "/home/lus/Pictures/AWS-Architecture-Icons/PNG"

	icons := []Icon{
		{Row: 9, Col: 4, Fit: true, URI: filepath.Join(iconsPath, "aws_api_gateway.png")},
		{Row: 8, Col: 2, Fit: true, URI: filepath.Join(iconsPath, "aws_lambda.png")},
		{Row: 10, Col: 2, Fit: true, URI: filepath.Join(iconsPath, "aws_lambda.png")},

		{Row: 6, Col: 2, Fit: true, URI: filepath.Join(iconsPath, "aws_rds_mysql_instance.png")},
		{Row: 5, Col: 2, Fit: true, URI: filepath.Join(iconsPath, "aws_simple_storage_service_s3_bucket.png")},
		{Row: 4, Col: 2, Fit: true, URI: filepath.Join(iconsPath, "aws_elasticache_for_redis.png")},

		{Row: 5, Col: 4, Fit: true, URI: filepath.Join(iconsPath, "aws_elastic_container_service.png")},
		{Row: 2, Col: 4, Fit: true, URI: filepath.Join(iconsPath, "aws_vpc_vpn_connection.png")},
		{Row: 5, Col: 8, Fit: true, URI: filepath.Join(iconsPath, "aws_simple_notification_service_sns.png")},
		{Row: 3, Col: 8, Fit: true, URI: filepath.Join(iconsPath, "aws_simple_notification_service_sns_topic.png")},
		{Row: 5, Col: 6, Fit: true, URI: filepath.Join(iconsPath, "aws_vpc_elastic_network_interface.png")},
	}

	connectors := []Connector{
		ElbowLeftDownConnector(8, 3, ConnectorArrowLeft()),
		TeeRightConnector(9, 3),
		ElbowLeftUpConnector(10, 3, ConnectorArrowLeft()),
		ElbowLeftDownConnector(4, 3, ConnectorArrowLeft()),
		CrossConnector(5, 3, ConnectorArrowLeft()),
		ElbowLeftUpConnector(6, 3, ConnectorArrowLeft()),

		VerticalConnector(8, 4),
		VerticalConnector(7, 4),
		VerticalConnector(6, 4, ConnectorArrowUp()),

		HorizontalConnector(5, 5, ConnectorArrowRight()),
		VerticalConnector(4, 4),
		VerticalConnector(3, 4, ConnectorArrowUp()),

		HorizontalConnector(5, 7, ConnectorArrowRight()),
		VerticalConnector(4, 8, ConnectorArrowUp()),

		TeeDownConnector(1, 4),
		ElbowRightUpConnector(1, 3, ConnectorArrowUp()),
		ElbowLeftUpConnector(1, 5, ConnectorArrowUp()),
	}

	frames := []Frame{
		NewFrame(2, 0, 11, 9),
		NewFrame(3, 1, 7, 6),
	}

	strings := []Label{
		NewLabel(3, 2, "VPC", LabelColor("#e9eed7"), LabelBackground("#2e422a"), LabelFontSize(12)),
		NewLabel(2, 1, "AWS", LabelColor("#e9eed7"), LabelBackground("#c8a60d"), LabelFontSize(12)),
		NewLabel(11, 7, "Token Manager Account", LabelColor("#fafafa"), LabelBackground("#00000088"), LabelFontSize(16)),
		NewLabel(0, 3, "SSO", LabelColor("#987634"), LabelFontSize(20)),
		NewLabel(0, 5, "MDW", LabelColor("#987634"), LabelFontSize(20)),
	}

	grid.DrawGrid()
	grid.DrawBorder()
	//grid.DrawCoords()

	for _, el := range frames {
		if err := el.Plot(grid); err != nil {
			panic(err)
		}
	}

	for _, el := range icons {
		el.Plot(grid)
	}

	for _, el := range connectors {
		el.Plot(grid)
	}

	for _, el := range strings {
		el.Plot(grid)
	}

	grid.SavePNG()
}
