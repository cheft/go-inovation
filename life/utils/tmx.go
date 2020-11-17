// © 2016 Steve McCoy under the MIT license. See LICENSE for details.

/*
Package tmx provides a struct for parsing the Tiled map editor's TMX map format. See
http://doc.mapeditor.org/reference/tmx-map-format/ .
https://github.com/mccoyst/tmx/blob/master/struct.go
*/
package utils

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"io"
	"strconv"
	"strings"
)

func Decode(r io.Reader) (*Map, error) {
	m := new(Map)
	x := xml.NewDecoder(r)
	err := x.Decode(m)
	if err != nil {
		return nil, err
	}

	for i := range m.Layers {
		err = m.Layers[i].decodeIDs()
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

type Map struct {
	XMLName         xml.Name `xml:"map"`
	Version         string   `xml:"version,attr"`
	Orientation     string   `xml:"orientation,attr"`
	Width           int      `xml:"width,attr"`
	Height          int      `xml:"height,attr"`
	TileWidth       int      `xml:"tilewidth,attr"`
	TileHeight      int      `xml:"tileheight,attr"`
	BackgroundColor string   `xml:"backgroundcolor,attr"`

	Properties   []Property    `xml:"properties>property"`
	Tilesets     []Tileset     `xml:"tileset"`
	Layers       []Layer       `xml:"layer"`
	ObjectGroups []ObjectGroup `xml:"objectgroup"`
	ImageLayers  []ImageLayer  `xml:"imagelayer"`
}

type Tileset struct {
	FirstGID   int32  `xml:"firstgid,attr"`
	Source     string `xml:"source,attr"`
	Name       string `xml:"name,attr"`
	TileWidth  int    `xml:"tilewidth,attr"`
	TileHeight int    `xml:"tileheight,attr"`
	Spacing    int    `xml:"spacing,attr"`
	Margin     int    `xml:"margin,attr"`

	TileOffset   TileOffset `xml:"tileoffset"`
	Properties   []Property `xml:"properties>property"`
	Image        Image      `xml:"image"`
	TerrainTypes []Terrain  `xml:"terraintypes>terrain"`
}

type TileOffset struct {
	X int `xml:"x,attr"`
	Y int `xml:"y,attr"`
}

type Image struct {
	Format string `xml:"format,attr"`
	Source string `xml:"source,attr"`
	Trans  string `xml:"trans,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`

	Data Data `xml:"data"`
}

type Terrain struct {
	Name string `xml:"name,attr"`
	Tile int    `xml:"tile,attr"`

	Properties []Property `xml:"properties>property"`
}

type Tile struct {
	ID          int32   `xml:"id,attr"`
	Terrain     string  `xml:"terrain,attr"`
	Probability float32 `xml:"probability,attr"`

	Properties []Property `xml:"properties>property"`
	Image      Image      `xml:"image"`
}

type Layer struct {
	Name    string  `xml:"name,attr"`
	Opacity float32 `xml:"opacity,attr"`
	Visible bool    `xml:"visible,attr"`

	Properties []Property `xml:"properties>property"`
	Data       Data       `xml:"data"`

	// The GID of each tile, in order. Use this instead of raw Data, which is cleaned up by Decode.
	GIDs []int32 `xml:"-"`
}

func (l *Layer) decodeIDs() error {
	d := &l.Data
	if d.Encoding == "csv" {
		r := bufio.NewScanner(strings.NewReader(d.Text))
		for r.Scan() {
			line := r.Text()
			parts := strings.Split(line, ",")
			for _, p := range parts {
				if p == "" {
					continue
				}
				n, err := strconv.ParseInt(p, 10, 32)
				if err != nil {
					return err
				}
				l.GIDs = append(l.GIDs, int32(n))
			}
		}
	} else if d.Encoding == "base64" {
		raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(d.Text))
		if err != nil {
			return err
		}
		var r io.Reader = bytes.NewReader(raw)
		if d.Compression == "zlib" {
			r, err = zlib.NewReader(r)
		} else if d.Compression == "gzip" {
			r, err = gzip.NewReader(r)
		}
		if err != nil {
			return err
		}
		for {
			var n int32
			err = binary.Read(r, binary.LittleEndian, &n)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			l.GIDs = append(l.GIDs, n)
		}
	} else {
		for i := range d.Tiles {
			l.GIDs = append(l.GIDs, d.Tiles[i].GID)
		}
	}

	l.Data = Data{}
	return nil
}

type Data struct {
	Encoding    string `xml:"encoding,attr"`
	Compression string `xml:"compression,attr"`

	Text  string       `xml:",chardata"`
	Tiles []SingleTile `xml:"tile"`
}

type SingleTile struct {
	GID int32 `xml:"gid,attr"`
}

type ObjectGroup struct {
	Name    string  `xml:"name,attr"`
	Color   string  `xml:"color,attr"`
	Opacity float32 `xml:"opacity,attr"`
	Visible bool    `xml:"visible,attr"`

	Properties []Property `xml:"properties>property"`
	Objects    []Object   `xml:"object"`
}

type Object struct {
	Name     string  `xml:"name,attr"`
	Type     string  `xml:"type,attr"`
	X        int     `xml:"x,attr"`
	Y        int     `xml:"y,attr"`
	Width    int     `xml:"width,attr"`
	Height   int     `xml:"height,attr"`
	Rotation float32 `xml:"rotation,attr"`
	GID      int32   `xml:"gid,attr"`
	Visible  bool    `xml:"visible,attr"`

	Properties []Property `xml:"properties>property"`
	Ellipse    *Ellipse   `xml:"ellipse"`
	Polygon    string     `xml:"polygon>points"`
	Polylines  string     `xml:"polyline>points"`
}

type Ellipse struct{}

type ImageLayer struct {
	Name    string  `xml:"name,attr"`
	Opacity float32 `xml:"opacity,attr"`
	Visible bool    `xml:"visible,attr"`

	Properties []Property `xml:"properties>property"`
	Image      Image      `xml:"image"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}
