package tmx

import (
	"encoding/xml"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Map struct {
	Version         string  `xml:"version,attr"`                   // The TMX format version. Was “1.0” so far, and will be incremented to match minor Tiled releases.
	TiledVersion    string  `xml:"tiledversion,attr"`              // The Tiled version used to save the file (since Tiled 1.0.1). May be a date (for snapshot builds).
	Orientation     string  `xml:"orientation,attr"`               // Map orientation. Tiled supports “orthogonal”, “isometric”, “staggered” and “hexagonal” (since 0.11).
	RenderOrder     *string `xml:"renderorder,attr,omitempty"`     // The order in which tiles on tile layers are rendered. Valid values are right-down (the default), right-up, left-down and left-up. In all cases, the map is drawn row-by-row. (only supported for orthogonal maps at the moment)
	Width           uint32  `xml:"width,attr"`                     // The map width in tiles.
	Height          uint32  `xml:"height,attr"`                    // The map height in tiles.
	TileWidth       uint32  `xml:"tilewidth,attr"`                 // The width of a tile.
	TileHeight      uint32  `xml:"tileheight,attr"`                // The height of a tile.
	HexSideLength   *uint32 `xml:"hexsidelength,attr,omitempty"`   // Only for hexagonal maps. Determines the width or height (depending on the staggered axis) of the tile’s edge, in pixels.
	StaggerAxis     *string `xml:"staggeraxis,attr,omitempty"`     // For staggered and hexagonal maps, determines which axis (“x” or “y”) is staggered. (since 0.11)
	StaggerIndex    *string `xml:"staggerindex,attr,omitempty"`    // For staggered and hexagonal maps, determines whether the “even” or “odd” indexes along the staggered axis are shifted. (since 0.11)
	BackgroundColor *string `xml:"backgroundcolor,attr,omitempty"` // The background color of the map. (optional, may include alpha value since 0.15 in the form #AARRGGBB
	NextObjectId    uint32  `xml:"nextobjectid,attr"`              // Stores the next available ID for new objects. This number is stored to prevent reuse of the same ID after objects have been removed. (since 0.11)

	Properties *Properties `xml:"properties,omitempty"`
	TileSets   []*TileSet  `xml:"tileset,omitempty"`
	Layers     []*Layer    `xml:",any"`
}

type TileSet struct {
	FirstGID   uint32 `xml:"firstgid,attr"`   // The first global tile ID of this tileset (this global ID maps to the first tile in this tileset).
	Source     string `xml:"source,attr"`     // If this tileset is stored in an external TSX (Tile Set XML) file, this attribute refers to that file. That TSX file has the same structure as the <tileset> element described here. (There is the firstgid attribute missing and this source attribute is also not there. These two attributes are kept in the TMX map, since they are map specific.)
	Name       string `xml:"name,attr"`       // The name of this tileset.
	TileWidth  uint32 `xml:"tilewidth,attr"`  // The (maximum) width of the tiles in this tileset.
	TileHeight uint32 `xml:"tileheight,attr"` // The (maximum) height of the tiles in this tileset.
	Spacing    uint32 `xml:"spacing,attr"`    // The spacing in pixels between the tiles in this tileset (applies to the tileset image).
	Margin     uint32 `xml:"margin,attr"`     // The margin around the tiles in this tileset (applies to the tileset image).
	TileCount  uint32 `xml:"tilecount,attr"`  // The number of tiles in this tileset (since 0.13)
	Columns    uint32 `xml:"columns,attr"`    // The number of tile columns in the tileset. For image collection tilesets it is editable and is used when displaying the tileset. (since 0.15)

	Offset       *TileOffset `xml:"tileoffset,omitempty"`
	Properties   *Properties `xml:"properties,omitempty"`
	Image        *Image      `xml:"image,omitempty"`
	TerrainTypes []*Terrain  `xml:"terraintypes>terrain,omitempty"`
	Tiles        []*Tile     `xml:"tile,omitempty"`
	// TODO: Add Wangsets
}

type TileOffset struct {
	X int32 `xml:"x,attr"` // Horizontal offset in pixels
	Y int32 `xml:"y,attr"` // Vertical offset in pixels (positive is down)
}

type Terrain struct {
	Name string `xml:"name,attr"` // The name of the terrain type.
	Tile uint32 `xml:"tile,attr"` // The local tile-id of the tile that represents the terrain visually.

	Properties *Properties `xml:"properties,omitempty"`
}

type Tile struct {
	ID          uint32   `xml:"id,attr"`                    // The local tile ID within its tileset.
	Type        *string  `xml:"type,attr,omitempty"`        // The type of the tile. Refers to an object type and is used by tile objects. (optional) (since 1.0)
	Terrain     *string  `xml:"terrain,attr,omitempty"`     // Defines the terrain type of each corner of the tile, given as comma-separated indexes in the terrain types array in the order top-left, top-right, bottom-left, bottom-right. Leaving out a value means that corner has no terrain. (optional)
	Probability *float64 `xml:"probability,attr,omitempty"` // A percentage indicating the probability that this tile is chosen when it competes with others while editing with the terrain tool. (optional)

	Properties  *Properties `xml:"properties,omitempty"`
	Image       *Image      `xml:"image,omitempty"`
	ObjectGroup *Layer      `xml:"objectgroup,omitempty"`
	Animation   []*Frame    `xml:"animation,omitempty"`
}

type Frame struct {
	TileID   uint32  `xml:"tileid,attr"`   // The local ID of a tile within the parent <tileset>.
	Duration float64 `xml:"duration,attr"` // How long (in milliseconds) this frame should be displayed before advancing to the next frame.
}

type Layer struct {
	XMLName xml.Name

	Name      string   `xml:"name,attr"`                // The name of the layer.
	Width     *uint32  `xml:"width,attr,omitempty"`     // The width of the layer in tiles. Always the same as the map width for fixed-size maps.
	Height    *uint32  `xml:"height,attr,omitempty"`    // The height of the layer in tiles. Always the same as the map height for fixed-size maps.
	Color     *string  `xml:"color,attr,omitempty"`     // The color used to display the objects in this group.
	Opacity   *float64 `xml:"opacity,attr,omitempty"`   // The opacity of the layer as a value from 0 to 1. Defaults to 1.
	Visible   *int     `xml:"visible,attr,omitempty"`   // Whether the layer is shown (1) or hidden (0). Defaults to 1.
	OffsetX   *float64 `xml:"offsetx,attr,omitempty"`   // Rendering offset for this layer in pixels. Defaults to 0. (since 0.14)
	OffsetY   *float64 `xml:"offsety,attr,omitempty"`   // Rendering offset for this layer in pixels. Defaults to 0. (since 0.14)
	DrawOrder *string  `xml:"draworder,attr,omitempty"` // Whether the objects are drawn according to the order of appearance (“index”) or sorted by their y-coordinate (“topdown”). Defaults to “topdown”.

	Properties *Properties `xml:"properties,omitempty"`
	Data       *Data       `xml:"data,omitempty"`
	Objects    []*Object   `xml:"object,omitempty"`
	Image      *Image      `xml:"image,omitempty"`
	Layers     []*Layer    `xml:",any"`
}

type Data struct {
	Encoding    *string `xml:"encoding,attr,omitempty"`    // The encoding used to encode the tile layer data. When used, it can be “base64” and “csv” at the moment.
	Compression *string `xml:"compression,attr,omitempty"` // The compression used to compress the tile layer data. Tiled supports “gzip” and “zlib”.

	TileData []TileData `xml:"data,omitempty"`
	Chunks   []Chunk    `xml:"chunk,omitempty"`
	Data     []byte     `xml:",innerxml"`
}

// This should probably not be used, rather use raw encoding
type TileData struct {
	GID uint32 `xml:"gid"`
}

type Chunk struct {
	X      float64 `xml:"x,attr"`      // The x coordinate of the chunk in tiles.
	Y      float64 `xml:"y,attr"`      // The y coordinate of the chunk in tiles.
	Width  int     `xml:"width,attr"`  // The width of the chunk in tiles.
	Height int     `xml:"height,attr"` // The height of the chunk in tiles.

	TileData []TileData `xml:"data,omitempty"`
	Data     []byte     `xml:",innerxml"`
}

type Object struct {
	ID       uint32   `xml:"id,attr"`                 // Unique ID of the object. Each object that is placed on a map gets a unique id. Even if an object was deleted, no object gets the same ID. Can not be changed in Tiled. (since Tiled 0.11)
	Name     string   `xml:"name,attr"`               // The name of the object. An arbitrary string.
	Type     *string  `xml:"type,attr,omitempty"`     // The type of the object. An arbitrary string.
	X        float64  `xml:"x,attr"`                  // The x coordinate of the object in pixels.
	Y        float64  `xml:"y,attr"`                  // The y coordinate of the object in pixels.
	Width    *float64 `xml:"width,attr,omitempty"`    // The width of the object in pixels (defaults to 0).
	Height   *float64 `xml:"height,attr,omitempty"`   // The height of the object in pixels (defaults to 0).
	Rotation *float64 `xml:"rotation,attr,omitempty"` // The rotation of the object in degrees clockwise (defaults to 0).
	GID      *uint32  `xml:"gid,attr,omitempty"`      // A reference to a tile (optional).
	Visible  *int     `xml:"visible,attr,omitempty"`  // Whether the object is shown (1) or hidden (0). Defaults to 1.
	TID      *uint32  `xml:"tid,attr,omitempty"`      // A reference to a template (optional).

	Ellipse  *Ellipse  `xml:"ellipse,omitempty"`
	Point    *Point    `xml:"point,omitempty"`
	Polygon  *Polygon  `xml:"polygon,omitempty"`
	Polyline *Polyline `xml:"polyline,omitempty"`
	Text     *Text     `xml:"text,omitempty"`
}

type Ellipse struct{}

type Point struct{}

type Polygon struct {
	Points string `xml:"points,attr"` // A list of x,y coordinates in pixels
}

type Polyline struct {
	Points string `xml:"points,attr"`
}

type Text struct {
	FontFamily *string `xml:"fontfamily,attr,omitempty"` // The font family used (default: “sans-serif”)
	PixelSize  *int    `xml:"pixelsize,attr,omitempty"`  // The size of the font in pixels (not using points, because other sizes in the TMX format are also using pixels) (default: 16)
	Wrap       *int    `xml:"wrap,attr,omitempty"`       // Whether word wrapping is enabled (1) or disabled (0). Defaults to 0.
	Color      *string `xml:"color,attr,omitempty"`      // Color of the text in #AARRGGBB or #RRGGBB format (default: #000000)
	Bold       *int    `xml:"bold,attr,omitempty"`       // Whether the font is bold (1) or not (0). Defaults to 0.
	Italic     *int    `xml:"italic,attr,omitempty"`     // Whether the font is italic (1) or not (0). Defaults to 0.
	Underline  *int    `xml:"underline,attr,omitempty"`  // Whether a line should be drawn below the text (1) or not (0). Defaults to 0.
	Strikeout  *int    `xml:"strikeout,attr,omitempty"`  // Whether a line should be drawn through the text (1) or not (0). Defaults to 0.
	Kerning    *int    `xml:"kerning,attr,omitempty"`    // Whether kerning should be used while rendering the text (1) or not (0). Default to 1.
	HAlign     *string `xml:"halign,attr,omitempty"`     // Horizontal alignment of the text within the object (left (default), center or right)
	VAlign     *string `xml:"valign,attr,omitempty"`     // Vertical alignment of the text within the object (top (default), center or bottom)

	Text string `xml:",innerxml"`
}

type Image struct {
	Format string  `xml:"format,attr,omitempty"` // Used for embedded images, in combination with a data child element. Valid values are file extensions like png, gif, jpg, bmp, etc.
	Source string  `xml:"source,attr"`           // The reference to the tileset image file (Tiled supports most common image formats).
	Trans  *string `xml:"trans,attr,omitempty"`  // Defines a specific color that is treated as transparent (example value // “#FF00FF” for magenta). Up until Tiled 0.12, this value is written out without a # but this is planned to change.
	Width  *int    `xml:"width,attr,omitempty"`  // The image width in pixels (optional, used for tile index correction when the image changes)
	Height *int    `xml:"height,attr,omitempty"` // The image height in pixels (optional)

	Data *Data `xml:"data,omitempty"`
}

type Properties struct {
	Properties []Property `xml:"property"`
}

type Property struct {
	Name  string  `xml:"name,attr"`           // The name of the property.
	Type  *string `xml:"type,attr,omitempty"` // The type of the property. Can be string (default), int, float, bool, color or file (since 0.16, with color and file added in 0.17).
	Value string  `xml:"value,attr"`          // The value of the property.
}

type Template struct {
	TileSet *TileSet `xml:"tileset,omitempty"`
	Object  *Object  `xml:"object,omitempty"`
}

// Load parses a tmx file into a new tmx.Map object, it will also parse any
// tsx tileset files that are referenced.
func Load(file *os.File) (*Map, error) {
	decoder := xml.NewDecoder(file)
	tmxMap := &Map{}
	err := decoder.Decode(tmxMap)
	for _, ts := range tmxMap.TileSets {
		if ts.Source != "" {
			tsxFile, err := os.Open(filepath.Join(filepath.Dir(file.Name()), ts.Source))
			if err != nil {
				return nil, errors.Wrap(err, "unable to open tileset source file")
			}
			d := xml.NewDecoder(tsxFile)
			err = d.Decode(ts)
			tsxFile.Close()
			if err != nil {
				return nil, errors.Wrap(err, "unable to decode tileset source file")
			}
		}
	}
	return tmxMap, errors.Wrap(err, "unable to decode tmx map")
}
