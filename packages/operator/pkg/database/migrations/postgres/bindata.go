// Code generated by go-bindata. DO NOT EDIT.
// sources:
// pkg/database/migrations/postgres/sources/000001_zero.down.sql (17B)
// pkg/database/migrations/postgres/sources/000001_zero.up.sql (17B)
// pkg/database/migrations/postgres/sources/000002_init.down.sql (273B)
// pkg/database/migrations/postgres/sources/000002_init.up.sql (668B)
// pkg/database/migrations/postgres/sources/000003_route.down.sql (57B)
// pkg/database/migrations/postgres/sources/000003_route.up.sql (136B)
// pkg/database/migrations/postgres/sources/000004_deletionmark.down.sql (209B)
// pkg/database/migrations/postgres/sources/000004_deletionmark.up.sql (290B)
// pkg/database/migrations/postgres/sources/000005_created_updated.down.sql (642B)
// pkg/database/migrations/postgres/sources/000005_created_updated.up.sql (722B)
// pkg/database/migrations/postgres/sources/000006_created_updated_delmark_for_route.down.sql (181B)
// pkg/database/migrations/postgres/sources/000006_created_updated_delmark_for_route.up.sql (226B)
// pkg/database/migrations/postgres/sources/000007_default_route.down.sql (71B)
// pkg/database/migrations/postgres/sources/000007_default_route.up.sql (98B)
// pkg/database/migrations/postgres/sources/000008_outbox.down.sql (49B)
// pkg/database/migrations/postgres/sources/000008_outbox.up.sql (210B)
// pkg/database/migrations/postgres/sources/000009_batch.down.sql (756B)
// pkg/database/migrations/postgres/sources/000009_batch.up.sql (1.313kB)

package postgres

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

var __000001_zeroDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x55\xa8\x4a\x2d\xca\x57\xc8\xcd\x4c\x2f\x4a\x2c\xc9\xcc\xcf\x03\x04\x00\x00\xff\xff\x55\xb8\xd3\x29\x11\x00\x00\x00")

func _000001_zeroDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_zeroDownSql,
		"000001_zero.down.sql",
	)
}

func _000001_zeroDownSql() (*asset, error) {
	bytes, err := _000001_zeroDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_zero.down.sql", size: 17, mode: os.FileMode(0664), modTime: time.Unix(1614256075, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xbb, 0x55, 0x97, 0xa7, 0x54, 0x23, 0x9, 0xa1, 0xb6, 0x33, 0x54, 0x59, 0x96, 0xaa, 0x3a, 0xff, 0xe4, 0x79, 0x88, 0xe0, 0xa, 0xf, 0x19, 0xfb, 0x5a, 0x88, 0xda, 0x8f, 0x61, 0xe3, 0xf8, 0xa5}}
	return a, nil
}

var __000001_zeroUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x55\xa8\x4a\x2d\xca\x57\xc8\xcd\x4c\x2f\x4a\x2c\xc9\xcc\xcf\x03\x04\x00\x00\xff\xff\x55\xb8\xd3\x29\x11\x00\x00\x00")

func _000001_zeroUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_zeroUpSql,
		"000001_zero.up.sql",
	)
}

func _000001_zeroUpSql() (*asset, error) {
	bytes, err := _000001_zeroUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_zero.up.sql", size: 17, mode: os.FileMode(0664), modTime: time.Unix(1614256075, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xbb, 0x55, 0x97, 0xa7, 0x54, 0x23, 0x9, 0xa1, 0xb6, 0x33, 0x54, 0x59, 0x96, 0xaa, 0x3a, 0xff, 0xe4, 0x79, 0x88, 0xe0, 0xa, 0xf, 0x19, 0xfb, 0x5a, 0x88, 0xda, 0x8f, 0x61, 0xe3, 0xf8, 0xa5}}
	return a, nil
}

var __000002_initDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x72\x75\xf7\xf4\xb3\xe6\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\xc8\x4f\x49\xcc\x28\x8d\xcf\x2f\x48\x2d\x4a\x2c\xc9\x2f\x8a\x2f\x29\x4a\xcc\xcc\xcb\xcc\x4b\x27\x4e\x75\x41\x62\x72\x76\x62\x3a\xd1\xca\x53\x52\x0b\x72\xf2\x2b\x73\x53\xf3\x4a\x48\x73\x4c\x7c\x66\x5e\x49\x6a\x7a\x51\x62\x49\x66\x7e\x1e\x89\x0e\x43\xd5\xea\xec\xef\xeb\xeb\x19\xa2\x60\x0d\x08\x00\x00\xff\xff\x80\xaa\x9b\xf0\x11\x01\x00\x00")

func _000002_initDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000002_initDownSql,
		"000002_init.down.sql",
	)
}

func _000002_initDownSql() (*asset, error) {
	bytes, err := _000002_initDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000002_init.down.sql", size: 273, mode: os.FileMode(0664), modTime: time.Unix(1614256075, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x2c, 0x7, 0xea, 0xf9, 0xa6, 0xde, 0xb6, 0x6d, 0x2, 0xc0, 0x74, 0xed, 0x89, 0x1b, 0x1, 0x5, 0xc8, 0x61, 0xd9, 0xe0, 0xe7, 0xe6, 0x30, 0x36, 0x75, 0xb5, 0x3a, 0x84, 0x91, 0x78, 0xa6, 0xb2}}
	return a, nil
}

var __000002_initUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\xd1\xc1\x8a\x83\x30\x10\x06\xe0\x7b\x9e\x62\x8e\x0a\x7b\x5c\xf6\xe2\x29\x4a\x76\x37\x6d\xd5\x12\x43\xa9\x27\x19\x34\x68\xa8\x4d\x42\x1c\x0f\x7d\xfb\x42\xdb\x17\xe8\xa1\x78\x9c\xe1\x87\xef\x87\x3f\x17\x7f\xb2\xca\x18\x2b\x94\xe0\x5a\x80\xe6\xf9\x41\x80\xfc\x85\xaa\xd6\x20\xce\xb2\xd1\x0d\xf8\x01\xa7\xb5\xf3\xc1\x44\x24\x1f\x3b\x8a\x68\x9d\x75\x23\x4b\x18\x00\x80\x1d\x00\xe0\xc4\x55\xf1\xcf\x55\xf2\xf3\x9d\xc2\x51\xc9\x92\xab\x16\xf6\xa2\xfd\x7a\x24\x96\x60\x7a\xd8\x35\x75\x95\xbf\x6e\x42\x5a\x97\xe7\x87\xa5\x19\x7b\x07\x0f\xd8\x5f\x70\xdc\x4a\x1f\x4c\x98\xfd\xed\x6a\x1c\x6d\xc2\x93\xf7\x73\x3f\xa1\x75\x9d\x75\x64\xc6\x88\x64\xbd\xdb\x76\x86\xcf\x34\x29\xea\xb2\x94\x3a\xbb\x07\x00\x00\xff\xff\x09\x72\x36\x73\x9c\x02\x00\x00")

func _000002_initUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000002_initUpSql,
		"000002_init.up.sql",
	)
}

func _000002_initUpSql() (*asset, error) {
	bytes, err := _000002_initUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000002_init.up.sql", size: 668, mode: os.FileMode(0664), modTime: time.Unix(1614256075, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x28, 0xf2, 0xca, 0x3a, 0xc5, 0x59, 0x0, 0xdc, 0x4e, 0x15, 0x7e, 0xf4, 0xed, 0xe5, 0x81, 0x61, 0xcd, 0x27, 0xbb, 0x60, 0xf0, 0xb1, 0x72, 0xa1, 0xb3, 0x1c, 0x91, 0x61, 0x9b, 0x90, 0x7a, 0x7f}}
	return a, nil
}

var __000003_routeDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x72\x75\xf7\xf4\xb3\xe6\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\xc8\x4f\x49\xcc\x28\x8d\xcf\x2f\x48\x2d\x4a\x2c\xc9\x2f\x8a\x2f\xca\x2f\x2d\x49\xb5\xe6\x72\xf6\xf7\xf5\xf5\x0c\xb1\x06\x04\x00\x00\xff\xff\x96\x06\x58\x69\x39\x00\x00\x00")

func _000003_routeDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000003_routeDownSql,
		"000003_route.down.sql",
	)
}

func _000003_routeDownSql() (*asset, error) {
	bytes, err := _000003_routeDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000003_route.down.sql", size: 57, mode: os.FileMode(0664), modTime: time.Unix(1614256075, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x77, 0x60, 0xcb, 0x71, 0x37, 0xc0, 0xdc, 0x41, 0xf0, 0x21, 0x61, 0x44, 0xcc, 0x9a, 0x1d, 0xd6, 0xd7, 0x62, 0xbf, 0x62, 0xfd, 0xf, 0xbd, 0xa5, 0xc1, 0xef, 0x49, 0x10, 0x49, 0xb9, 0x69, 0xf1}}
	return a, nil
}

var __000003_routeUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2c\x8d\xbf\x0a\xc2\x30\x18\x07\xf7\x3c\xc5\x6f\x6c\xc1\x51\x5c\x3a\xa5\xe5\x53\xa3\x26\x91\xe4\x43\xec\x54\x82\x0d\xe8\x94\x92\x26\xef\x2f\xfe\x19\xef\x38\xb8\x9e\x0e\xca\x74\x42\x0c\x8e\x24\x13\x58\xf6\x17\x82\xda\xc3\x58\x06\xdd\x95\x67\x8f\x34\x87\x67\x9d\xd2\x12\x73\x28\x29\x4f\x39\xd5\x12\x45\x23\x00\xe0\x35\x03\xb8\x49\x37\x1c\xa5\x6b\x76\xdb\x16\x57\xa7\xb4\x74\x23\xce\x34\x6e\xbe\xc5\xba\xc4\x07\x4e\xde\x9a\xfe\xcf\x25\x94\xba\xfe\x8c\x68\x3f\x63\xab\xb5\xe2\xee\x1d\x00\x00\xff\xff\x1b\x91\x14\x0a\x88\x00\x00\x00")

func _000003_routeUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000003_routeUpSql,
		"000003_route.up.sql",
	)
}

func _000003_routeUpSql() (*asset, error) {
	bytes, err := _000003_routeUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000003_route.up.sql", size: 136, mode: os.FileMode(0664), modTime: time.Unix(1614256075, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xdd, 0xf3, 0xbe, 0xb7, 0x9c, 0xb4, 0x26, 0xae, 0x6f, 0xc8, 0x16, 0x8c, 0xf9, 0xd8, 0xbc, 0xb5, 0x96, 0x49, 0x75, 0xd5, 0xda, 0x6d, 0xcf, 0x8c, 0x3f, 0x6c, 0x6f, 0x87, 0x8d, 0xb9, 0x71, 0x48}}
	return a, nil
}

var __000004_deletionmarkDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\xcc\xb1\x0d\xc2\x40\x0c\x05\xd0\x9e\x29\xac\xac\x91\x0e\x84\x50\x8a\x40\x43\x1f\x7d\x72\x56\x38\xc5\x67\x5b\x96\x53\xb0\x3d\x1b\x5c\x93\x01\xde\xbb\xde\x1f\xd3\x73\xbc\x40\x92\x83\x12\x1f\x61\xb2\x82\xef\xb1\x98\x73\x20\x2d\x96\x0c\x54\xad\xba\x51\x09\x73\x5a\x4d\x8e\xa6\x34\x14\x16\xce\x6a\xda\x10\xfb\xd0\xf5\x8e\x75\xc7\x76\x22\x28\xec\x62\xbf\xc6\x9a\xdd\xe1\xf6\x9a\xe7\xe9\x3d\xfe\x03\x00\x00\xff\xff\x28\x9d\x34\x19\xd1\x00\x00\x00")

func _000004_deletionmarkDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000004_deletionmarkDownSql,
		"000004_deletionmark.down.sql",
	)
}

func _000004_deletionmarkDownSql() (*asset, error) {
	bytes, err := _000004_deletionmarkDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000004_deletionmark.down.sql", size: 209, mode: os.FileMode(0664), modTime: time.Unix(1614256075, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x9f, 0x93, 0x32, 0x59, 0xc, 0xa5, 0xb5, 0xa9, 0x53, 0x9a, 0xd6, 0x2b, 0xe9, 0x35, 0x3f, 0x17, 0x93, 0x0, 0xad, 0x35, 0xc4, 0xb7, 0x16, 0x2a, 0x1, 0x31, 0xbe, 0xe2, 0x7f, 0x23, 0xda, 0x16}}
	return a, nil
}

var __000004_deletionmarkUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\xcc\x31\x0e\x82\x50\x0c\x06\xe0\x9d\x53\xfc\xe1\x1a\x4c\x6a\xd0\x90\x88\x0e\xba\x93\x62\x2b\xbe\x50\xda\x97\x97\xbe\xc1\xdb\x7b\x03\x07\x13\x2f\xf0\xed\xfb\xd3\x70\xe9\x1a\xd2\x90\x82\xa0\x59\x05\xce\xf4\xaa\x93\x67\x29\x14\x5e\xa6\x28\x94\x2c\xd9\xd2\x00\x00\x31\xa3\x65\x51\x89\xe4\xb6\x51\x59\x5b\xcc\xee\x2a\x64\x60\x79\x52\xd5\xc0\x71\x77\xbe\xf5\x30\x0f\x58\x55\xfd\x4a\x67\x7a\xac\xb4\xfc\xc7\x66\xc9\xea\xef\x4d\x2c\x7e\xc5\x0f\xd7\x71\x1c\xee\xdd\x27\x00\x00\xff\xff\x0a\x00\xba\xa4\x22\x01\x00\x00")

func _000004_deletionmarkUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000004_deletionmarkUpSql,
		"000004_deletionmark.up.sql",
	)
}

func _000004_deletionmarkUpSql() (*asset, error) {
	bytes, err := _000004_deletionmarkUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000004_deletionmark.up.sql", size: 290, mode: os.FileMode(0664), modTime: time.Unix(1614256075, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x45, 0x7e, 0xcc, 0x75, 0x79, 0xed, 0x92, 0x2e, 0x2a, 0xb9, 0xd0, 0x23, 0xb9, 0x25, 0x2c, 0xb, 0xdd, 0x45, 0xa6, 0x8f, 0x4e, 0xe5, 0xe2, 0x4e, 0x67, 0x85, 0xa9, 0x34, 0x31, 0x54, 0x55, 0x77}}
	return a, nil
}

var __000005_created_updatedDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x8f\x31\xae\xc2\x30\x10\x44\xfb\x7f\x0a\xdf\x23\xdd\x47\x08\xa5\x08\x34\xf4\xd6\x60\xaf\x1c\x0b\x67\xd7\x5a\xad\x0b\x6e\x4f\x47\x45\x94\x38\x1c\xe0\xcd\x7b\xf3\x7f\xbe\x8c\xd7\xe1\x0f\xc5\x48\x9d\xe1\x51\xc8\x49\xc4\xdc\xbc\x54\x52\x98\xa8\x37\x45\xe6\xcc\xc9\x45\x95\xea\x82\x94\xb6\xb0\x0b\x4a\x30\x8a\xfd\x60\xab\x71\x13\xac\x08\x4f\xa4\x23\xca\xef\xe4\x1e\x67\xa4\x5a\xe4\xb5\x10\x5b\xb7\x74\x05\xed\x7a\xea\x33\x1b\x25\x85\x65\xe1\xe3\xaf\x57\x57\xf6\xb4\x98\x48\x09\x33\x32\xff\xd4\xb2\xbd\xf2\x69\x39\xdd\xa6\x69\xbc\x0f\xef\x00\x00\x00\xff\xff\x0e\xd3\x2d\xae\x82\x02\x00\x00")

func _000005_created_updatedDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000005_created_updatedDownSql,
		"000005_created_updated.down.sql",
	)
}

func _000005_created_updatedDownSql() (*asset, error) {
	bytes, err := _000005_created_updatedDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000005_created_updated.down.sql", size: 642, mode: os.FileMode(0664), modTime: time.Unix(1624535277, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x21, 0x4b, 0x2b, 0x21, 0xee, 0x57, 0x80, 0x55, 0x10, 0xf, 0x1e, 0x71, 0x51, 0xa5, 0x41, 0x48, 0x2a, 0x4a, 0xd9, 0xb7, 0x52, 0x38, 0x94, 0x57, 0xa8, 0x73, 0x4, 0x1e, 0x39, 0xd1, 0x95, 0xdd}}
	return a, nil
}

var __000005_created_updatedUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\xd1\x41\x0a\xc2\x30\x10\x85\xe1\x7d\x4f\x31\xf7\xe8\x4e\x11\xe9\xa2\xba\x71\x5f\x9e\xcd\xd0\x06\xd3\x4c\x88\xaf\x0b\x3d\xbd\x07\xb0\x14\x12\xbd\xc0\xcf\x07\xff\xe1\x74\xee\x2e\x6d\x83\x40\xcd\x42\xdc\x83\x8a\x39\xcc\xeb\x60\x49\x33\x68\x79\x60\x86\x8f\x3e\x4e\x8d\x88\x08\x9c\x93\x31\x2b\xa8\x4e\xe8\x17\x7d\x12\x4b\xe2\xbb\x2c\xb0\x26\x57\x14\x48\x18\x1f\x98\x7e\x21\x7c\x17\x4a\x0d\x4e\x53\xb0\xd7\xa2\x91\xd5\x88\x8d\x44\xa9\x82\x66\x61\x9c\xe1\xe3\xe0\x23\x75\xca\xa0\xb7\x58\x3f\x66\xb7\x56\x7d\xe9\x2f\xb6\xfd\xda\xa6\xed\x78\xed\xfb\xee\xd6\x7e\x02\x00\x00\xff\xff\x16\x7c\x3b\xea\xd2\x02\x00\x00")

func _000005_created_updatedUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000005_created_updatedUpSql,
		"000005_created_updated.up.sql",
	)
}

func _000005_created_updatedUpSql() (*asset, error) {
	bytes, err := _000005_created_updatedUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000005_created_updated.up.sql", size: 722, mode: os.FileMode(0664), modTime: time.Unix(1624535161, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x87, 0xda, 0x61, 0x60, 0x70, 0x4b, 0xf2, 0xad, 0x51, 0x9f, 0x19, 0x24, 0x26, 0x63, 0xe0, 0xe5, 0x18, 0x2e, 0x93, 0xbb, 0x40, 0xe, 0x2e, 0x4b, 0x14, 0x7, 0xd2, 0xc2, 0xb0, 0xdc, 0xbc, 0x84}}
	return a, nil
}

var __000006_created_updated_delmark_for_routeDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\xcb\xc1\x09\x42\x31\x0c\x06\xe0\xfb\x9b\x22\x7b\xf4\xa6\x88\xf4\x50\xbd\x78\x2f\xb1\xf9\x41\x31\x6d\x4a\x48\xf6\x77\x06\x1d\xe0\x3b\x5d\xae\xf5\x56\x0e\xd6\x80\x53\xf0\x53\x41\x26\xfc\xca\x6e\x1b\xce\x61\xde\xdd\x32\x40\xe2\xb6\x69\x98\xe6\x5c\x34\x1c\x1c\x90\x1f\x55\x6e\xf9\x43\x09\x14\xf1\xb6\x35\xd9\x3f\xe5\x38\xdf\x5b\xab\x8f\xf2\x0d\x00\x00\xff\xff\x4a\x13\xb4\x4c\xb5\x00\x00\x00")

func _000006_created_updated_delmark_for_routeDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000006_created_updated_delmark_for_routeDownSql,
		"000006_created_updated_delmark_for_route.down.sql",
	)
}

func _000006_created_updated_delmark_for_routeDownSql() (*asset, error) {
	bytes, err := _000006_created_updated_delmark_for_routeDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000006_created_updated_delmark_for_route.down.sql", size: 181, mode: os.FileMode(0664), modTime: time.Unix(1622817356, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xc5, 0x55, 0xe, 0x83, 0xb3, 0xae, 0x46, 0x3d, 0xee, 0xd5, 0xd9, 0x9c, 0xb3, 0xa5, 0xb8, 0x42, 0xc, 0xd3, 0x3d, 0x89, 0x3d, 0x1b, 0xe, 0x12, 0xf, 0x27, 0xf5, 0x3a, 0xbd, 0xc2, 0xd, 0xe1}}
	return a, nil
}

var __000006_created_updated_delmark_for_routeUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\xcb\xc1\xad\xc2\x30\x0c\x06\xe0\x7b\xa6\xf8\xd5\x35\x72\x7a\x0f\x15\x54\x89\xc2\x01\xee\x95\x8b\x8d\xa8\x70\xe2\xc8\x38\x17\xa6\x67\x06\xc4\x00\xdf\xff\x78\x98\x4e\x39\x91\x86\x38\x82\x56\x15\x18\xd3\xa3\x2f\xd6\xc4\x29\xcc\x17\xb7\x1e\x92\x00\x80\x98\x71\x73\xa1\x10\x46\x6c\x45\x5e\x41\xa5\xc5\xfb\x0b\xdd\x1b\xff\xa0\x07\x16\x95\xd8\xac\x16\xf2\xe7\x80\xd5\x4c\x85\x2a\x58\xee\xd4\x35\xb0\xff\x3b\x5e\x46\x54\x0b\xd4\xae\x9a\xd3\xee\x3c\xcf\xd3\x35\x7f\x02\x00\x00\xff\xff\x09\x87\x16\xb4\xe2\x00\x00\x00")

func _000006_created_updated_delmark_for_routeUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000006_created_updated_delmark_for_routeUpSql,
		"000006_created_updated_delmark_for_route.up.sql",
	)
}

func _000006_created_updated_delmark_for_routeUpSql() (*asset, error) {
	bytes, err := _000006_created_updated_delmark_for_routeUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000006_created_updated_delmark_for_route.up.sql", size: 226, mode: os.FileMode(0664), modTime: time.Unix(1622817356, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x10, 0x2a, 0x71, 0x4f, 0xad, 0x73, 0xed, 0xae, 0x3b, 0x59, 0xde, 0x3, 0x38, 0xe3, 0xa6, 0x97, 0x1f, 0xa9, 0xd4, 0xd8, 0x7, 0xf6, 0x8e, 0xbd, 0x7f, 0xa9, 0xb3, 0x48, 0xc, 0x76, 0x67, 0x2f}}
	return a, nil
}

var __000007_default_routeDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x04\xc0\x31\x0e\x80\x30\x08\x00\xc0\xbd\xaf\xe0\x1f\xdd\x34\xc6\x74\xa8\x2e\xee\x0d\x0a\x46\x13\x94\x06\xe1\xff\xde\x30\xcd\x65\xc9\x09\xc5\xd9\xc0\x71\x17\x06\x25\xbc\xa2\x69\x67\x43\x57\x6b\xa6\xe1\x0c\x64\xda\xe1\x50\x89\xe7\x85\xfb\x6b\xc4\x27\x86\x78\x4e\xe3\x5a\x6b\xd9\xf2\x1f\x00\x00\xff\xff\xc4\x58\xde\x6c\x47\x00\x00\x00")

func _000007_default_routeDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000007_default_routeDownSql,
		"000007_default_route.down.sql",
	)
}

func _000007_default_routeDownSql() (*asset, error) {
	bytes, err := _000007_default_routeDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000007_default_route.down.sql", size: 71, mode: os.FileMode(0664), modTime: time.Unix(1622817356, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xab, 0xea, 0xa4, 0xd5, 0x68, 0x43, 0xaa, 0xba, 0xf7, 0xee, 0x10, 0x2a, 0x59, 0x86, 0x7b, 0x1e, 0xf3, 0xe4, 0x6a, 0x1f, 0xe7, 0xf0, 0xbb, 0x5f, 0xf2, 0x66, 0xd6, 0x3f, 0x62, 0x3e, 0x88, 0xcb}}
	return a, nil
}

var __000007_default_routeUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x34\xc7\x31\x0a\x42\x31\x0c\x06\xe0\xbd\xa7\xf8\xef\xd1\x49\xe5\x29\x05\xab\x83\xee\x25\x25\x11\x85\xd0\x48\x4c\xee\xef\xf4\xc6\xef\xb8\x5d\xda\xad\x16\xd2\x10\x47\xd0\x54\x81\x31\xbd\x73\xd8\x57\x9c\xc2\x7c\xb8\x65\x48\x01\x00\x62\xc6\xe7\x37\x58\x5e\x94\x1a\x98\x66\x2a\xb4\xb0\xfb\x7c\xb8\x3e\x36\x2c\x0b\xac\x54\xad\xe5\x74\xef\xbd\x3d\xeb\x3f\x00\x00\xff\xff\x15\x93\xd4\x27\x62\x00\x00\x00")

func _000007_default_routeUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000007_default_routeUpSql,
		"000007_default_route.up.sql",
	)
}

func _000007_default_routeUpSql() (*asset, error) {
	bytes, err := _000007_default_routeUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000007_default_route.up.sql", size: 98, mode: os.FileMode(0664), modTime: time.Unix(1622817356, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x58, 0x10, 0x24, 0xd8, 0xce, 0xbe, 0x3, 0x6, 0x60, 0x6e, 0xe5, 0x79, 0x9, 0xc8, 0x18, 0x61, 0x23, 0x7f, 0xde, 0x72, 0xdb, 0x35, 0x8, 0xad, 0x10, 0x9b, 0xcd, 0xf7, 0x80, 0xe1, 0x4b, 0xe1}}
	return a, nil
}

var __000008_outboxDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x72\x75\xf7\xf4\xb3\xe6\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\xc8\x4f\x49\xcc\x28\x8d\xcf\x2f\x2d\x49\xca\xaf\xb0\xe6\x72\xf6\xf7\xf5\xf5\x0c\xb1\x06\x04\x00\x00\xff\xff\x5e\x08\x04\x5c\x31\x00\x00\x00")

func _000008_outboxDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000008_outboxDownSql,
		"000008_outbox.down.sql",
	)
}

func _000008_outboxDownSql() (*asset, error) {
	bytes, err := _000008_outboxDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000008_outbox.down.sql", size: 49, mode: os.FileMode(0664), modTime: time.Unix(1622817356, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x64, 0x7c, 0xab, 0xee, 0x7d, 0xb1, 0x7, 0x20, 0xdc, 0x6f, 0xc8, 0x16, 0xd9, 0x76, 0x2, 0xaa, 0x53, 0xbb, 0xea, 0xee, 0x49, 0xcb, 0x2d, 0x33, 0xd2, 0xea, 0xc3, 0x2, 0xa5, 0x62, 0xba, 0x96}}
	return a, nil
}

var __000008_outboxUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\xcf\xcd\x8a\xc2\x30\x14\x40\xe1\x7d\x9e\xe2\x2e\x5b\x98\xcd\x0c\xc3\x30\xd0\x55\x52\x62\x8d\x34\xad\x24\x17\x11\x37\x25\x92\xa0\x01\x6d\x8a\xdc\x8a\x7d\x7b\xf1\x07\x41\x70\x7b\xbe\xd5\x11\xb2\x52\x4d\xc1\x58\x69\x24\x47\x09\xc8\x45\x2d\x41\xcd\xa0\x69\x11\xe4\x5a\x59\xb4\x90\xbc\xdb\x8f\x5d\x1a\x69\x9b\x2e\x2c\x63\x00\x00\xd1\x03\x08\x55\x59\x69\x14\xaf\xbf\xee\x29\xf4\x14\x69\xea\xa2\x87\x15\x37\xe5\x9c\x9b\xec\xef\x37\x7f\xd2\x39\xf4\xd4\xd1\x34\x84\x97\x7d\xff\xfc\xbf\xe1\xee\x94\xc6\xe1\x83\x7a\x47\x81\xe2\x31\x00\x2a\x2d\x2d\x72\xbd\xc4\xcd\x43\x06\x37\x1d\x92\xf3\xb0\xb0\x6d\x23\x58\x7e\x5b\x68\xb5\x56\x58\x5c\x03\x00\x00\xff\xff\xf0\xc0\x2d\x0c\xd2\x00\x00\x00")

func _000008_outboxUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000008_outboxUpSql,
		"000008_outbox.up.sql",
	)
}

func _000008_outboxUpSql() (*asset, error) {
	bytes, err := _000008_outboxUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000008_outbox.up.sql", size: 210, mode: os.FileMode(0664), modTime: time.Unix(1622817356, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x69, 0x14, 0x14, 0x36, 0x6f, 0x12, 0xe7, 0x97, 0x89, 0xef, 0x3, 0x23, 0xd8, 0x34, 0x33, 0xc0, 0x6b, 0x13, 0x87, 0x32, 0xe6, 0xda, 0x7c, 0xd3, 0xb2, 0x27, 0x69, 0xba, 0xaf, 0x50, 0x8c, 0x2b}}
	return a, nil
}

var __000009_batchDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x91\x41\x73\x9b\x3e\x10\xc5\xef\x7c\x8a\x37\x3e\xfd\xff\x1d\xd7\xa4\x3e\xd6\x27\xe2\x90\x56\xd3\x18\x32\x86\x34\xcd\xc9\x23\xc4\x02\xdb\xc1\x12\x95\x44\x08\xdf\xbe\x83\x63\x3a\xf6\xf4\xd4\x1d\x9d\xb4\x4f\xbf\x7d\x6f\x15\x7e\x08\x30\x1d\x4c\xb5\x35\xdd\x68\xb9\x6e\x3c\xd6\x37\xeb\x4f\x88\x1f\xa3\x1d\xb2\xd1\x79\x3a\xba\x0b\xd5\x03\x2b\xd2\x8e\x4a\xf4\xba\x24\x0b\xdf\x10\xa2\x4e\xaa\x86\xe6\xce\x12\xdf\xc9\x3a\x36\x1a\xeb\xd5\x0d\xfe\x9b\x04\x8b\x73\x6b\xf1\xff\x66\xc6\x8c\xa6\xc7\x51\x8e\xd0\xc6\xa3\x77\x04\xdf\xb0\x43\xc5\x2d\x81\xde\x14\x75\x1e\xac\xa1\xcc\xb1\x6b\x59\x6a\x45\x18\xd8\x37\xa7\x59\x67\xd2\x6a\xe6\xbc\x9c\x39\xa6\xf0\x92\x35\x24\x94\xe9\x46\x98\xea\x52\x0c\xe9\x2f\x02\x4c\xd5\x78\xdf\x7d\x0e\xc3\x61\x18\x56\xf2\x64\x7e\x65\x6c\x1d\xb6\xef\x72\x17\x3e\x88\x6d\x9c\x64\xf1\xc7\xf5\xea\xe6\xe2\xe1\x93\x6e\xc9\x39\x58\xfa\xd5\xb3\xa5\x12\xc5\x08\xd9\x75\x2d\x2b\x59\xb4\x84\x56\x0e\x30\x16\xb2\xb6\x44\x25\xbc\x99\x02\x0c\x96\x3d\xeb\x7a\x09\x67\x2a\x3f\x48\x4b\x33\xaa\x64\xe7\x2d\x17\xbd\xbf\xda\xe3\x6c\x97\xdd\x95\xc0\x68\x48\x8d\x45\x94\x41\x64\x0b\xdc\x46\x99\xc8\x96\x33\xe8\x59\xe4\x5f\xd3\xa7\x1c\xcf\xd1\x7e\x1f\x25\xb9\x88\x33\xa4\x7b\x6c\xd3\xe4\x4e\xe4\x22\x4d\x32\xa4\xf7\x88\x92\x17\x7c\x13\xc9\xdd\x12\xc4\xbe\x21\x0b\x7a\xeb\xec\x94\xc4\x58\xf0\xb4\x61\x2a\xff\xac\x33\x23\xba\xb2\x52\x99\x77\x6b\xae\x23\xc5\x15\x2b\xb4\x52\xd7\xbd\xac\x09\xb5\x79\x25\xab\x59\xd7\xe8\xc8\x1e\xd9\x4d\x3f\xee\x20\x75\x39\xa3\x5a\x3e\xb2\x97\xfe\x74\xfd\x57\xc6\x69\x60\x18\x04\xb7\xf1\x17\x91\x6c\x82\xe0\x6e\x9f\x3e\x22\x8f\x6e\x1f\x62\x88\x7b\xc4\x3f\x44\x96\x67\x30\xa5\x6c\xfa\x43\x21\xbd\x6a\x0e\xac\x2b\xb2\xa4\x15\x1d\x7e\x9a\x62\xf3\x2f\x7a\x47\xf6\x95\x15\x6d\x82\x60\x9b\xee\x76\x22\xdf\xfc\x0e\x00\x00\xff\xff\x2b\xec\xf8\x54\xf4\x02\x00\x00")

func _000009_batchDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000009_batchDownSql,
		"000009_batch.down.sql",
	)
}

func _000009_batchDownSql() (*asset, error) {
	bytes, err := _000009_batchDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000009_batch.down.sql", size: 756, mode: os.FileMode(0664), modTime: time.Unix(1622817356, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x76, 0x2e, 0x7d, 0xd3, 0x38, 0xd1, 0x5b, 0x6f, 0x9a, 0x37, 0x88, 0xe5, 0x1a, 0x2d, 0x9f, 0x1d, 0x30, 0x3d, 0xe4, 0xb4, 0x65, 0xdf, 0x60, 0x2f, 0xf6, 0x7, 0x30, 0xb5, 0x71, 0x85, 0x12, 0x7f}}
	return a, nil
}

var __000009_batchUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x53\x41\x6f\x83\x36\x18\xbd\xf3\x2b\x9e\x72\x6a\xa7\x2c\xe9\xa2\x69\x87\xf5\x44\x52\xba\xb2\x26\x50\x01\x6d\x97\x53\x64\xe0\x03\xdc\x82\xcd\x6c\xd3\x34\xfb\xf5\x93\x21\xac\x69\x56\x69\xd2\x6a\x71\xe1\x7b\xcf\xcf\xdf\x7b\x9f\x3d\xff\xc1\x81\xfd\x60\xd7\x4a\xb6\x07\xc5\xcb\xca\x60\x71\xb5\xf8\x09\xde\x83\xbb\x41\x7c\xd0\x86\x1a\x7d\xc2\x5a\xf3\x8c\x84\xa6\x1c\x9d\xc8\x49\xc1\x54\x04\xb7\x65\x59\x45\x23\x32\xc5\x13\x29\xcd\xa5\xc0\x62\x76\x85\x0b\x4b\x98\x1c\xa1\xc9\xe5\xf5\x28\x73\x90\x1d\x1a\x76\x80\x90\x06\x9d\x26\x98\x8a\x6b\x14\xbc\x26\xd0\x7b\x46\xad\x01\x17\xc8\x64\xd3\xd6\x9c\x89\x8c\xb0\xe7\xa6\xea\xcf\x3a\x2a\xcd\x46\x9d\xed\x51\x47\xa6\x86\x71\x01\x86\x4c\xb6\x07\xc8\xe2\x94\x0c\x66\x4e\x0c\xd8\x55\x19\xd3\xfe\x3a\x9f\xef\xf7\xfb\x19\xeb\x9b\x9f\x49\x55\xce\xeb\x81\xae\xe7\x6b\x7f\xe5\x05\xb1\xf7\xe3\x62\x76\x75\xb2\xf1\x51\xd4\xa4\x35\x14\xfd\xd9\x71\x45\x39\xd2\x03\x58\xdb\xd6\x3c\x63\x69\x4d\xa8\xd9\x1e\x52\x81\x95\x8a\x28\x87\x91\xd6\xc0\x5e\x71\xc3\x45\x39\x85\x96\x85\xd9\x33\x45\xa3\x54\xce\xb5\x51\x3c\xed\xcc\xa7\x1c\xc7\x76\xb9\xfe\x44\x90\x02\x4c\x60\xe2\xc6\xf0\xe3\x09\x96\x6e\xec\xc7\xd3\x51\xe8\xd9\x4f\xee\xc2\xc7\x04\xcf\x6e\x14\xb9\x41\xe2\x7b\x31\xc2\x08\xab\x30\xb8\xf1\x13\x3f\x0c\x62\x84\xb7\x70\x83\x2d\xee\xfd\xe0\x66\x0a\xe2\xa6\x22\x05\x7a\x6f\x95\x75\x22\x15\xb8\x4d\x98\xf2\x7f\xe2\x8c\x89\x3e\xb5\x52\xc8\xa1\x35\xdd\x52\xc6\x0b\x9e\xa1\x66\xa2\xec\x58\x49\x28\xe5\x1b\x29\xc1\x45\x89\x96\x54\xc3\xb5\x9d\xb8\x06\x13\xf9\x28\x55\xf3\x86\x1b\x66\xfa\xf2\xbf\x3c\xda\x03\xe7\x8e\xb3\xf4\x7e\xf3\x83\x6b\xc7\x59\x45\x9e\x9b\x78\x48\xdc\xe5\xda\x83\x7f\x8b\x20\x4c\xe0\xfd\xe1\xc7\x49\x0c\x99\xb3\xaa\xdb\xa5\xcc\x64\xd5\x8e\x8b\x82\x14\x89\x8c\x76\x9a\xd4\x1b\xcf\xc8\xb9\x70\xec\x51\x3c\x07\xf0\xe4\x46\xab\x3b\x37\xba\xf8\xe5\xe7\x4b\x3c\x44\xfe\xc6\x8d\xb6\xb8\xf7\xb6\xd3\x9e\x91\x29\x62\x36\x4b\xc3\x1b\xd2\x86\x35\xad\xf9\x6b\x00\xba\x36\xff\x1a\xc8\xa9\x26\xdb\x7c\xc3\xd4\x2b\x52\x29\x6b\x62\x02\x39\x15\xac\xab\x0d\x6e\xdd\x75\xec\xf5\x77\x57\x74\x75\x3d\x6c\xb0\x11\xe1\xf7\x38\x0c\x96\xc7\x7f\xc3\x4c\xa7\x87\x8a\x73\xf9\x7f\x4c\xbe\xc8\xf4\x3b\x06\xcf\xfa\xfb\xc2\xe9\x19\xe3\x1b\x96\xcf\x81\x13\xef\xe7\xd0\x30\xb9\xe1\x1d\xbe\x31\x95\x55\x4c\xf5\x96\x3e\xd6\xb8\x61\x30\x26\x85\x36\x8a\x71\x61\xc6\x94\xf8\xcb\x2e\xe5\x7a\x57\xbc\xf6\xb8\xa2\x63\x5c\xfa\x3f\xae\x8a\x25\x4b\x71\x8c\x01\x8a\xec\x03\xcb\x8c\x2d\xf5\xbe\x3f\x4a\xc3\xb0\xc2\xcd\xc6\x4f\xae\xff\x0e\x00\x00\xff\xff\x0e\x5e\xb5\xcd\x21\x05\x00\x00")

func _000009_batchUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000009_batchUpSql,
		"000009_batch.up.sql",
	)
}

func _000009_batchUpSql() (*asset, error) {
	bytes, err := _000009_batchUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000009_batch.up.sql", size: 1313, mode: os.FileMode(0664), modTime: time.Unix(1622817356, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x1b, 0xd7, 0xf4, 0xd0, 0x19, 0xb, 0x6f, 0x34, 0x84, 0x82, 0xc, 0x30, 0x48, 0xd4, 0xa2, 0x33, 0xd1, 0x8a, 0xe6, 0xe4, 0x4d, 0x93, 0x18, 0x75, 0xd0, 0xe9, 0xad, 0x86, 0x69, 0x37, 0x87, 0x7}}
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
	"000001_zero.down.sql":                              _000001_zeroDownSql,
	"000001_zero.up.sql":                                _000001_zeroUpSql,
	"000002_init.down.sql":                              _000002_initDownSql,
	"000002_init.up.sql":                                _000002_initUpSql,
	"000003_route.down.sql":                             _000003_routeDownSql,
	"000003_route.up.sql":                               _000003_routeUpSql,
	"000004_deletionmark.down.sql":                      _000004_deletionmarkDownSql,
	"000004_deletionmark.up.sql":                        _000004_deletionmarkUpSql,
	"000005_created_updated.down.sql":                   _000005_created_updatedDownSql,
	"000005_created_updated.up.sql":                     _000005_created_updatedUpSql,
	"000006_created_updated_delmark_for_route.down.sql": _000006_created_updated_delmark_for_routeDownSql,
	"000006_created_updated_delmark_for_route.up.sql":   _000006_created_updated_delmark_for_routeUpSql,
	"000007_default_route.down.sql":                     _000007_default_routeDownSql,
	"000007_default_route.up.sql":                       _000007_default_routeUpSql,
	"000008_outbox.down.sql":                            _000008_outboxDownSql,
	"000008_outbox.up.sql":                              _000008_outboxUpSql,
	"000009_batch.down.sql":                             _000009_batchDownSql,
	"000009_batch.up.sql":                               _000009_batchUpSql,
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
	"000001_zero.down.sql":                              {_000001_zeroDownSql, map[string]*bintree{}},
	"000001_zero.up.sql":                                {_000001_zeroUpSql, map[string]*bintree{}},
	"000002_init.down.sql":                              {_000002_initDownSql, map[string]*bintree{}},
	"000002_init.up.sql":                                {_000002_initUpSql, map[string]*bintree{}},
	"000003_route.down.sql":                             {_000003_routeDownSql, map[string]*bintree{}},
	"000003_route.up.sql":                               {_000003_routeUpSql, map[string]*bintree{}},
	"000004_deletionmark.down.sql":                      {_000004_deletionmarkDownSql, map[string]*bintree{}},
	"000004_deletionmark.up.sql":                        {_000004_deletionmarkUpSql, map[string]*bintree{}},
	"000005_created_updated.down.sql":                   {_000005_created_updatedDownSql, map[string]*bintree{}},
	"000005_created_updated.up.sql":                     {_000005_created_updatedUpSql, map[string]*bintree{}},
	"000006_created_updated_delmark_for_route.down.sql": {_000006_created_updated_delmark_for_routeDownSql, map[string]*bintree{}},
	"000006_created_updated_delmark_for_route.up.sql":   {_000006_created_updated_delmark_for_routeUpSql, map[string]*bintree{}},
	"000007_default_route.down.sql":                     {_000007_default_routeDownSql, map[string]*bintree{}},
	"000007_default_route.up.sql":                       {_000007_default_routeUpSql, map[string]*bintree{}},
	"000008_outbox.down.sql":                            {_000008_outboxDownSql, map[string]*bintree{}},
	"000008_outbox.up.sql":                              {_000008_outboxUpSql, map[string]*bintree{}},
	"000009_batch.down.sql":                             {_000009_batchDownSql, map[string]*bintree{}},
	"000009_batch.up.sql":                               {_000009_batchUpSql, map[string]*bintree{}},
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
