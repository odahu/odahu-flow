// Code generated by go-bindata. DO NOT EDIT.
// sources:
// pkg/deployment/assets/mapper.rego (567B)
// pkg/deployment/assets/roles.rego (91B)
// pkg/deployment/assets/ml_servers/odahu_ml_server.rego (634B)
// pkg/deployment/assets/ml_servers/triton.rego (896B)

package bindata

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
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
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
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

	info := bindataFileInfo{name: "mapper.rego", size: 567, mode: os.FileMode(0664), modTime: time.Unix(1634714528, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x34, 0xf8, 0x76, 0x31, 0xbf, 0x92, 0xdc, 0x47, 0x7a, 0xc2, 0x28, 0xfd, 0xa6, 0xeb, 0x1c, 0xe1, 0x16, 0x4e, 0xc2, 0x21, 0x82, 0x97, 0x9e, 0x79, 0xd7, 0x9b, 0xc8, 0x71, 0xa4, 0x8e, 0xc4, 0x28}}
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

	info := bindataFileInfo{name: "roles.rego", size: 91, mode: os.FileMode(0664), modTime: time.Unix(1634714528, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x63, 0xe5, 0x4f, 0x68, 0xc0, 0x5d, 0xd9, 0x2b, 0xed, 0xa8, 0xd3, 0xb5, 0xef, 0x72, 0x43, 0xd3, 0x9f, 0x6a, 0x5, 0x11, 0x88, 0x3, 0x23, 0xb0, 0x4a, 0xa8, 0x86, 0xd0, 0xb9, 0x71, 0x4, 0x9f}}
	return a, nil
}

var _ml_serversOdahu_ml_serverRego = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x92\xc1\x4e\xc4\x20\x10\x86\xef\x3c\xc5\x04\xaf\x1b\xfa\x04\x7b\x54\x8f\x9a\x75\x6f\x4d\xd3\x4c\x60\xb4\xa4\xd0\x21\x40\xad\xb1\xe9\xbb\x9b\x82\x89\x6e\xb2\xf1\xa0\x07\x8f\xfc\x3f\x7c\xf9\x80\x09\xa8\x47\x7c\x21\x60\x83\xc3\xac\x34\x47\x12\xc2\xfa\xc0\x31\x83\xc1\x8c\xaa\xe6\x1e\x43\xa0\x78\xa5\x88\xec\x28\x09\x81\xce\xf1\x02\xab\x00\xa8\x3b\x15\xea\x6c\x79\x82\xe3\x11\xe4\xfd\xed\x59\x7e\x15\x91\x12\xcf\x51\x53\xa9\x9a\x81\xd0\xe5\xe1\x5d\x8a\xed\xf7\x0c\x4f\x39\x5a\x9d\xfe\xc4\xa8\x1e\x7a\x20\x3d\x16\xce\x0d\x9c\x70\x81\xfd\x72\x3f\x21\xdb\xc2\x3c\x80\x7c\x7c\x78\x3a\xcb\xae\xed\xbb\xeb\xfc\x56\x36\x18\x6c\xe3\xd9\x90\x6b\xec\xf4\xcc\xfb\x99\x8b\xe8\x95\x47\xaa\x80\x6f\x04\x5c\xfa\xf2\xbc\x6d\xdf\x15\xcb\x75\x55\x27\x76\xb4\x6d\x9f\x8a\x77\xf6\x8d\x0c\xd4\x1f\xf8\x37\xcb\x39\x51\xbc\xd0\x6c\xcb\x42\xed\x33\xd2\x27\x6d\x69\xca\x36\xe5\x43\xd5\x54\x68\xbc\x9d\x0a\x61\x13\x1f\x01\x00\x00\xff\xff\x8b\xcd\xc6\x89\x7a\x02\x00\x00")

func ml_serversOdahu_ml_serverRegoBytes() ([]byte, error) {
	return bindataRead(
		_ml_serversOdahu_ml_serverRego,
		"ml_servers/odahu_ml_server.rego",
	)
}

func ml_serversOdahu_ml_serverRego() (*asset, error) {
	bytes, err := ml_serversOdahu_ml_serverRegoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "ml_servers/odahu_ml_server.rego", size: 634, mode: os.FileMode(0664), modTime: time.Unix(1634714528, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x7f, 0x2e, 0x86, 0x3e, 0x25, 0x95, 0xdd, 0x3e, 0xde, 0x9d, 0xe4, 0x3, 0x6b, 0xea, 0x2d, 0x3d, 0x77, 0x11, 0xc6, 0x36, 0x26, 0x82, 0x79, 0x53, 0x26, 0x8e, 0x3f, 0xd, 0x1a, 0xb1, 0x47, 0xfd}}
	return a, nil
}

var _ml_serversTritonRego = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x92\x41\x6f\xd3\x40\x10\x85\xef\xfb\x2b\x46\x2e\x87\x44\x4d\x6d\x89\x23\x52\x14\xf5\x50\xe0\x82\xa8\x4a\x6e\xc6\xb8\xa3\xdd\x49\xbd\xea\x7a\xd7\x9a\x9d\x24\x40\xc8\x7f\x47\xbb\x36\x6a\x28\xa8\x87\xe6\x66\xcd\x1b\xbd\xf7\xcd\xf3\x0e\xa8\x1f\xf1\x81\x20\x18\xec\xb6\xa5\x0e\x4c\x4a\xd9\x7e\x08\x2c\x60\x50\xb0\x1c\xe7\x3d\x0e\x03\xf1\x7f\x04\x0e\x8e\xa2\x52\x17\x70\xe3\xcd\x10\xac\x97\x08\x61\x03\xd1\x1a\xd2\xc8\xa0\x83\x17\xb4\x9e\x38\x2a\x74\x2e\xec\xe1\xa0\x00\x46\xaf\x12\xb5\xd8\xe0\x61\xb9\x84\xe2\xc3\xcd\xba\x78\x12\x98\x62\xd8\xb2\xa6\x2c\x55\x1d\xa1\x93\xee\x67\xa1\x8e\xea\xd5\x1e\x3d\x09\x5b\x1d\xcf\xf2\x18\x39\x74\x47\xfa\x31\xfb\x5c\xc0\x9a\xad\x04\x0f\x1f\xb3\x00\xf4\xe7\xfe\x57\x27\xec\xde\x4e\x21\x95\xb3\x3b\x3a\x0b\xf6\xc9\x8a\x09\xcd\x8f\x09\xf8\xfd\xd6\x39\x40\xad\x29\x46\xd8\x04\x86\x6b\xd3\x5b\xff\x6f\xc6\x36\x12\xb7\xf9\xbf\xd6\x6d\x93\xfc\xf2\x77\x89\x79\xfb\xf4\xf2\x4f\xc1\x90\x3b\x39\x5c\xf2\xf8\x16\xa5\xbb\xa3\x07\xfa\x0e\xef\x96\x70\xff\x2d\xb1\xf4\x69\x31\x56\xf5\xd7\xfd\x55\x73\x39\xab\x76\xc4\xd1\x06\x9f\x06\xa6\xb9\x9c\xaf\x66\x95\xf5\x1b\xe2\x5f\x23\xec\x7c\x55\xad\xde\xdc\xa7\x98\xeb\x11\x55\x02\x64\x9d\xbc\xa6\xcc\x9d\x9e\x1f\x44\x6d\xc9\x8b\xf5\x51\x22\xa0\x37\x90\xa8\x23\xec\xad\x74\x30\x10\x5f\xe5\xc8\x4c\xfe\x52\x89\x75\x6e\x71\x01\xc5\xed\xe7\x2f\xeb\xa2\xa9\xdb\x46\x01\x30\xb5\x3d\x8a\xee\x66\xcf\xee\x59\x3c\xef\x7a\xae\x4e\xfa\xc7\xfd\x5f\x9d\xd5\xc5\xe1\x50\xde\x05\x47\xc7\x63\xb1\x98\x1a\x4c\xdc\xed\xc4\x1d\x25\xa7\x1d\xd5\xef\x00\x00\x00\xff\xff\x1d\xdb\xd2\x0d\x80\x03\x00\x00")

func ml_serversTritonRegoBytes() ([]byte, error) {
	return bindataRead(
		_ml_serversTritonRego,
		"ml_servers/triton.rego",
	)
}

func ml_serversTritonRego() (*asset, error) {
	bytes, err := ml_serversTritonRegoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "ml_servers/triton.rego", size: 896, mode: os.FileMode(0664), modTime: time.Unix(1634714528, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x6c, 0x7e, 0x46, 0x52, 0x62, 0x67, 0x88, 0x60, 0xd7, 0xa0, 0x3c, 0xcc, 0x19, 0xfe, 0x8e, 0x8, 0xa0, 0x3, 0x3f, 0x6d, 0xca, 0x9d, 0x9e, 0x74, 0xcc, 0x94, 0xc0, 0x5, 0xe4, 0xa6, 0xbd, 0x44}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
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

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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
	"mapper.rego":                     mapperRego,
	"roles.rego":                      rolesRego,
	"ml_servers/odahu_ml_server.rego": ml_serversOdahu_ml_serverRego,
	"ml_servers/triton.rego":          ml_serversTritonRego,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	"mapper.rego": {mapperRego, map[string]*bintree{}},
	"ml_servers": {nil, map[string]*bintree{
		"odahu_ml_server.rego": {ml_serversOdahu_ml_serverRego, map[string]*bintree{}},
		"triton.rego":          {ml_serversTritonRego, map[string]*bintree{}},
	}},
	"roles.rego": {rolesRego, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
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
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
