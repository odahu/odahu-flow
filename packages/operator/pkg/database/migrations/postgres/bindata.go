package postgres

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var __000001_zero_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x55\xa8\x4a\x2d\xca\x57\xc8\xcd\x4c\x2f\x4a\x2c\xc9\xcc\xcf\x03\x04\x00\x00\xff\xff\x55\xb8\xd3\x29\x11\x00\x00\x00")

func _000001_zero_down_sql() ([]byte, error) {
	return bindata_read(
		__000001_zero_down_sql,
		"000001_zero.down.sql",
	)
}

var __000001_zero_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x55\xa8\x4a\x2d\xca\x57\xc8\xcd\x4c\x2f\x4a\x2c\xc9\xcc\xcf\x03\x04\x00\x00\xff\xff\x55\xb8\xd3\x29\x11\x00\x00\x00")

func _000001_zero_up_sql() ([]byte, error) {
	return bindata_read(
		__000001_zero_up_sql,
		"000001_zero.up.sql",
	)
}

var __000002_init_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x72\x75\xf7\xf4\xb3\xe6\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\xc8\x4f\x49\xcc\x28\x8d\xcf\x2f\x48\x2d\x4a\x2c\xc9\x2f\x8a\x2f\x29\x4a\xcc\xcc\xcb\xcc\x4b\x27\x4e\x75\x41\x62\x72\x76\x62\x3a\xd1\xca\x53\x52\x0b\x72\xf2\x2b\x73\x53\xf3\x4a\x48\x73\x4c\x7c\x66\x5e\x49\x6a\x7a\x51\x62\x49\x66\x7e\x1e\x89\x0e\x43\xd5\xea\xec\xef\xeb\xeb\x19\xa2\x60\x0d\x08\x00\x00\xff\xff\x80\xaa\x9b\xf0\x11\x01\x00\x00")

func _000002_init_down_sql() ([]byte, error) {
	return bindata_read(
		__000002_init_down_sql,
		"000002_init.down.sql",
	)
}

var __000002_init_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\xd1\xc1\x8a\x83\x30\x10\x06\xe0\x7b\x9e\x62\x8e\x0a\x7b\x5c\xf6\xe2\x29\x4a\x76\x37\x6d\xd5\x12\x43\xa9\x27\x19\x34\x68\xa8\x4d\x42\x1c\x0f\x7d\xfb\x42\xdb\x17\xe8\xa1\x78\x9c\xe1\x87\xef\x87\x3f\x17\x7f\xb2\xca\x18\x2b\x94\xe0\x5a\x80\xe6\xf9\x41\x80\xfc\x85\xaa\xd6\x20\xce\xb2\xd1\x0d\xf8\x01\xa7\xb5\xf3\xc1\x44\x24\x1f\x3b\x8a\x68\x9d\x75\x23\x4b\x18\x00\x80\x1d\x00\xe0\xc4\x55\xf1\xcf\x55\xf2\xf3\x9d\xc2\x51\xc9\x92\xab\x16\xf6\xa2\xfd\x7a\x24\x96\x60\x7a\xd8\x35\x75\x95\xbf\x6e\x42\x5a\x97\xe7\x87\xa5\x19\x7b\x07\x0f\xd8\x5f\x70\xdc\x4a\x1f\x4c\x98\xfd\xed\x6a\x1c\x6d\xc2\x93\xf7\x73\x3f\xa1\x75\x9d\x75\x64\xc6\x88\x64\xbd\xdb\x76\x86\xcf\x34\x29\xea\xb2\x94\x3a\xbb\x07\x00\x00\xff\xff\x09\x72\x36\x73\x9c\x02\x00\x00")

func _000002_init_up_sql() ([]byte, error) {
	return bindata_read(
		__000002_init_up_sql,
		"000002_init.up.sql",
	)
}

var __000003_route_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x72\x75\xf7\xf4\xb3\xe6\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\xc8\x4f\x49\xcc\x28\x8d\xcf\x2f\x48\x2d\x4a\x2c\xc9\x2f\x8a\x2f\xca\x2f\x2d\x49\xb5\xe6\x72\xf6\xf7\xf5\xf5\x0c\xb1\x06\x04\x00\x00\xff\xff\x96\x06\x58\x69\x39\x00\x00\x00")

func _000003_route_down_sql() ([]byte, error) {
	return bindata_read(
		__000003_route_down_sql,
		"000003_route.down.sql",
	)
}

var __000003_route_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2c\x8d\xbf\x0a\xc2\x30\x18\x07\xf7\x3c\xc5\x6f\x6c\xc1\x51\x5c\x3a\xa5\xe5\x53\xa3\x26\x91\xe4\x43\xec\x54\x82\x0d\xe8\x94\x92\x26\xef\x2f\xfe\x19\xef\x38\xb8\x9e\x0e\xca\x74\x42\x0c\x8e\x24\x13\x58\xf6\x17\x82\xda\xc3\x58\x06\xdd\x95\x67\x8f\x34\x87\x67\x9d\xd2\x12\x73\x28\x29\x4f\x39\xd5\x12\x45\x23\x00\xe0\x35\x03\xb8\x49\x37\x1c\xa5\x6b\x76\xdb\x16\x57\xa7\xb4\x74\x23\xce\x34\x6e\xbe\xc5\xba\xc4\x07\x4e\xde\x9a\xfe\xcf\x25\x94\xba\xfe\x8c\x68\x3f\x63\xab\xb5\xe2\xee\x1d\x00\x00\xff\xff\x1b\x91\x14\x0a\x88\x00\x00\x00")

func _000003_route_up_sql() ([]byte, error) {
	return bindata_read(
		__000003_route_up_sql,
		"000003_route.up.sql",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"000001_zero.down.sql": _000001_zero_down_sql,
	"000001_zero.up.sql": _000001_zero_up_sql,
	"000002_init.down.sql": _000002_init_down_sql,
	"000002_init.up.sql": _000002_init_up_sql,
	"000003_route.down.sql": _000003_route_down_sql,
	"000003_route.up.sql": _000003_route_up_sql,
}
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"000001_zero.down.sql": &_bintree_t{_000001_zero_down_sql, map[string]*_bintree_t{
	}},
	"000001_zero.up.sql": &_bintree_t{_000001_zero_up_sql, map[string]*_bintree_t{
	}},
	"000002_init.down.sql": &_bintree_t{_000002_init_down_sql, map[string]*_bintree_t{
	}},
	"000002_init.up.sql": &_bintree_t{_000002_init_up_sql, map[string]*_bintree_t{
	}},
	"000003_route.down.sql": &_bintree_t{_000003_route_down_sql, map[string]*_bintree_t{
	}},
	"000003_route.up.sql": &_bintree_t{_000003_route_up_sql, map[string]*_bintree_t{
	}},
}}
