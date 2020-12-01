// Code generated by go-bindata. (@generated) DO NOT EDIT.

 //Package bindata generated by go-bindata.// sources:
// pkg/deployment/assets/core.rego
// pkg/deployment/assets/mapper.rego
// pkg/deployment/assets/roles.rego
package bindata

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// ModTime return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _coreRego = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x92\xc1\x4e\xc4\x20\x10\x86\xef\x3c\xc5\x04\xaf\x1b\xfa\x04\x7b\x54\x8f\x9a\x75\x6f\x4d\xd3\x4c\x60\xb4\xa4\xd0\x21\x40\xad\xb1\xe9\xbb\x9b\x82\x89\x6e\xb2\xf1\xa0\x07\x8f\xfc\x3f\x7c\xf9\x80\x09\xa8\x47\x7c\x21\x60\x83\xc3\xac\x34\x47\x12\xc2\xfa\xc0\x31\x83\xc1\x8c\xaa\xe6\x1e\x43\xa0\x78\xa5\x88\xec\x28\x09\x81\xce\xf1\x02\xab\x00\xa8\x3b\x15\xea\x6c\x79\x82\xe3\x11\xe4\xfd\xed\x59\x7e\x15\x91\x12\xcf\x51\x53\xa9\x9a\x81\xd0\xe5\xe1\x5d\x8a\xed\xf7\x0c\x4f\x39\x5a\x9d\xfe\xc4\xa8\x1e\x7a\x20\x3d\x16\xce\x0d\x9c\x70\x81\xfd\x72\x3f\x21\xdb\xc2\x3c\x80\x7c\x7c\x78\x3a\xcb\xae\xed\xbb\xeb\xfc\x56\x36\x18\x6c\xe3\xd9\x90\x6b\xec\xf4\xcc\xfb\x99\x8b\xe8\x95\x47\xaa\x80\x6f\x04\x5c\xfa\xf2\xbc\x6d\xdf\x15\xcb\x75\x55\x27\x76\xb4\x6d\x9f\x8a\x77\xf6\x8d\x0c\xd4\x1f\xf8\x37\xcb\x39\x51\xbc\xd0\x6c\xcb\x42\xed\x33\xd2\x27\x6d\x69\xca\x36\xe5\x43\xd5\x54\x68\xbc\x9d\x0a\x61\x13\x1f\x01\x00\x00\xff\xff\x8b\xcd\xc6\x89\x7a\x02\x00\x00")

func coreRegoBytes() ([]byte, error) {
	return bindataRead(
		_coreRego,
		"core.rego",
	)
}

func coreRego() (*asset, error) {
	bytes, err := coreRegoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "core.rego", size: 634, mode: os.FileMode(420), modTime: time.Unix(1604043088, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _mapperRego = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x91\xdf\x6a\xeb\x30\x0c\xc6\xaf\xeb\xa7\x30\xb9\x3e\xf8\x01\x0e\xf4\x09\xb6\xbb\xc1\x6e\x42\x30\x9a\xa3\x36\xea\xe2\x3f\xb3\xe4\x66\xa5\xec\xdd\x87\xdd\x36\xdb\xa0\x63\x37\x16\xfa\x7e\x9f\x14\x49\x49\xe0\x5e\x61\x8f\x3a\x8e\x30\x15\xe3\x21\x25\xcc\x4a\x91\x4f\x31\x8b\x1e\x41\xc0\x5c\x48\x8e\x33\xb2\x52\x2d\x58\x0f\x49\x6f\xf5\x59\x6d\xba\x06\x2d\x8c\x9e\x42\xf7\x5f\x37\x6a\x5a\xf6\x4f\x69\x7d\xa5\xb5\x8b\x65\x47\x18\x84\x58\x56\xdb\x4f\xf9\x9b\xff\x48\xb8\x60\x5e\x7d\x97\x54\x7d\x28\x75\x58\x44\x6f\x35\x85\x54\xc4\x80\x48\xa6\x97\x22\xc8\xc6\xa3\x40\x6b\xe6\x62\x10\x7c\x17\xb3\xa3\x59\x30\xdb\x9b\xde\x77\x18\x8e\xf1\x74\x95\xd9\x4c\x22\xc9\x1c\x16\xb1\x50\x64\x0a\xdd\x60\x76\x84\xf3\xc8\x4d\x4a\x70\x9a\x23\x8c\x4a\x65\x58\x6c\xfb\x7e\x5f\xdf\xe1\xac\x36\x35\xea\xad\x3e\x2c\x62\x1e\x28\x8c\xe6\x49\x72\x71\xf2\x0c\x73\xc1\x5b\x87\x8c\x30\x7b\x0b\xce\x21\xf3\xef\xa6\xb6\x55\xa3\x8f\xc4\x57\x76\xac\x2f\xf7\x76\x58\xcb\x28\xec\x1b\xaa\x8b\x17\xc6\x7c\x7f\x9a\xf5\x7f\xf4\x5f\x03\xdb\x61\xa8\x45\xe0\x84\x62\xb8\x77\xb0\x8c\x6f\x05\x59\x2e\x87\xf0\x28\x53\xac\x0b\x23\xc7\x92\x1d\xfe\x59\x90\x40\x26\xf5\x19\x00\x00\xff\xff\x35\xc2\xe9\x5d\x37\x02\x00\x00")

func mapperRegoBytes() ([]byte, error) {
	return bindataRead(
		_mapperRego,
		"mapper.rego",
	)
}

func mapperRego() (*asset, error) {
	bytes, err := mapperRegoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "mapper.rego", size: 567, mode: os.FileMode(420), modTime: time.Unix(1604043088, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _rolesRego = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2a\x48\x4c\xce\x4e\x4c\x4f\x55\xc8\x4f\x49\xcc\x28\xd5\x2b\xca\xcf\x49\x2d\xe6\xe2\x4a\x4c\xc9\xcd\xcc\x53\xb0\xb2\x55\x50\x02\xb3\x94\xb8\x52\x12\x4b\x12\xe3\x8b\x93\x33\x53\xf3\x4a\x32\x8b\x4b\xc0\x32\xa8\x42\x4a\x5c\x65\x99\xa9\xe5\xa9\x45\x60\x29\x08\x53\x09\x10\x00\x00\xff\xff\x95\x1e\x87\xee\x5b\x00\x00\x00")

func rolesRegoBytes() ([]byte, error) {
	return bindataRead(
		_rolesRego,
		"roles.rego",
	)
}

func rolesRego() (*asset, error) {
	bytes, err := rolesRegoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "roles.rego", size: 91, mode: os.FileMode(420), modTime: time.Unix(1604043088, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
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
var _bindata = map[string]func() (*asset, error){
	"core.rego":   coreRego,
	"mapper.rego": mapperRego,
	"roles.rego":  rolesRego,
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
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"core.rego":   &bintree{coreRego, map[string]*bintree{}},
	"mapper.rego": &bintree{mapperRego, map[string]*bintree{}},
	"roles.rego":  &bintree{rolesRego, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
