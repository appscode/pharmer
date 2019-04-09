// Code generated by go-bindata.
// sources:
// cloud.json
// DO NOT EDIT!

package lightsail

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
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
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

var _cloudJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x97\xd1\x8e\xea\x36\x10\x86\xef\x79\x0a\x2b\xd7\x0b\x3d\x44\xad\x74\xc4\x1d\xa2\x6c\xb5\xda\x56\xbb\x2a\xab\x56\x55\xb5\x8a\x06\x67\x08\x2e\x8e\x9d\xb5\x1d\xb6\xe9\x8a\xa7\xe9\xa3\xf4\xc5\xaa\x38\x90\x38\xc1\x81\x50\x6d\x6f\xce\x59\xcd\xfc\x1e\xbe\xf9\xe3\x19\x25\x1f\x23\x42\x02\x01\x29\x06\x33\x12\x70\x96\x6c\x8d\x06\xc6\x83\xbb\x32\x8c\x62\xaf\x83\x19\xf9\x7d\x44\x08\x21\x41\x8c\x7b\x1b\x26\x24\x78\x83\xd3\x5f\x99\x92\x71\x30\x22\xe4\xd5\x1e\x50\x98\x30\x29\x9a\x33\x1f\xf6\x5f\x42\x02\x2e\x29\x18\x26\x45\xf9\x23\xbf\x30\x95\x30\xc1\x4e\x25\xea\x63\x65\x2e\xd7\x63\x04\x6d\xc6\xd3\x26\xf9\x97\x14\xd8\x54\xb4\xa1\x5a\xd5\xd4\x70\xa3\x6b\x6f\x94\x7a\xa3\xb1\x37\x8a\xc1\x31\xf8\x6a\xff\x3f\xdc\xf5\x77\xf3\xb4\x65\xf2\x52\x27\xe1\xa0\x4e\x42\x6f\x27\xa1\xb7\x93\x90\x0e\xa7\x5b\x00\x67\x1b\xa9\xfa\xdd\x7e\xc7\x21\x6e\x57\xaa\x2e\x63\x15\xed\x32\x56\xd1\x1b\x18\x9f\x14\x26\x52\x5c\xe2\xbb\xea\x61\xa5\xf2\xf2\x9d\x79\x58\x45\x6f\xe0\x7b\x50\xc8\x41\xc4\x3e\x40\xcc\x87\x18\x58\xab\x5a\x80\x75\x74\xed\x8d\xde\x00\xf8\xa3\x14\xb1\xdf\xc0\x53\xb5\xcb\x06\xfe\xdf\x7c\x2b\x26\x12\xc8\xa4\x42\x1f\x22\x64\x63\x2d\x73\xb3\x1d\x30\xf6\x6d\x69\x0b\xb6\x9d\x5a\xf7\xa7\x6e\xc0\x7e\x91\xbb\xc2\x3b\xda\x90\x8d\x85\x54\x43\x91\x1d\x69\x17\xd9\x49\x75\x91\x9d\xd4\x15\xe4\x1e\x2b\x9d\x47\xde\x7a\x16\x45\x2c\xb0\x18\x6e\x73\xd8\x6f\x73\xd8\x6f\xf3\xb5\xf1\x6a\x98\x35\x9c\x2d\xfc\x16\xee\x3f\x7f\x4b\xf2\x0c\x39\x97\x17\x91\xeb\x2a\x2d\xda\x3a\xba\xf6\x46\x6f\xf6\xb5\x0f\xf2\xa7\x3c\x5d\x03\x1b\xe4\x69\xcf\xb5\xed\xb9\xb2\x37\x30\x62\x3e\xa6\x28\x8c\x02\xde\x47\x79\xaf\x40\xec\x36\xb9\x32\xd7\x76\x41\x5d\xa7\xbb\x0e\xea\x44\x77\x23\xd4\x89\x2e\x6e\xfd\x62\xc0\x84\x36\x20\x28\xbe\x14\x19\x7a\x5e\x0f\xf4\x2e\x2f\x11\x05\x08\x19\x4d\xa3\x2f\x0d\x61\x8c\x9a\x2a\x96\x9d\x7a\x40\x1a\xce\x48\xa9\x6a\x14\x14\x0c\x26\x52\x15\x65\xfa\x07\x14\xa8\x80\x93\xe7\x5c\x65\x52\x3b\xeb\x86\x66\x65\xfd\x69\x33\xcb\x90\x06\x33\xf2\x65\xf2\x9d\xd7\xd5\x23\x4d\xca\xa8\x1a\x80\x63\x65\x9f\xc2\x33\xbd\x44\xa3\x53\xe0\xfc\x3a\x8d\x95\x7d\x0a\x4d\x78\xd1\x1b\x8c\x59\x9e\x0e\x30\xc7\xea\xfe\x13\x4f\xd8\xe1\xf9\xf6\x12\x0f\x07\x95\xe0\x75\x1c\x2b\xfb\x14\x9a\xaf\x9d\x0b\x4e\x15\xc6\x28\x0c\x03\xee\xb9\xde\x99\x92\x7b\x16\xa3\x2a\x7f\x67\xfe\xeb\xca\x61\x64\x3a\xe3\x50\xdc\x4b\x95\x82\x29\xb3\x1b\x86\xdc\x79\xd3\x00\x21\xa4\xb1\x13\x5c\x56\xfd\x68\x86\x8e\x72\x99\xc7\x93\x6c\x0b\x2a\x45\x35\x61\xf2\x1b\xca\x73\x6d\x50\x8d\x1b\x8e\xb2\x9c\x3b\xa7\x67\x47\x62\xa1\x6f\x91\x6b\x23\x15\x24\xd8\x3d\x72\x3c\x71\xa8\x99\x6d\x0b\xed\xad\xd2\x80\x57\x1f\x15\x54\x8a\x0d\x4b\x8e\x6e\x44\xf3\xc5\x62\xb9\x5a\x45\x8f\xcb\xdf\xa2\x87\xef\x1d\x84\xb2\x96\x54\xa9\xdd\xc0\xef\x3a\x02\x4a\x51\xeb\x68\x87\x45\xc4\xe2\xb6\xec\x0f\x7d\x5c\xd4\x56\xf2\x88\x45\xb7\x0e\x87\x35\x5a\xdc\xb9\x55\x90\x47\x2c\xc8\x43\xa7\x08\x13\x59\x6e\x9f\x81\xc1\x3f\x4d\x50\x67\x0e\x77\x03\xdb\x58\x2d\x17\x3f\x2f\x5f\x9c\x6e\xfa\x5b\xd1\x48\x15\x1a\xa7\x23\x7f\x3b\x95\x6c\x7e\x6a\xaa\xa7\xa5\x95\x55\x91\xa6\xb3\x9e\xb6\x32\xd0\xfa\x5d\xaa\xd8\x69\xad\x67\x59\xef\xf2\x35\x2a\x81\xc6\xb7\xa9\xf7\xa8\xf4\x71\xa2\xa6\x93\xaf\x93\xfe\x79\xeb\x64\x8f\x9f\x92\xce\x15\x2e\x3f\x27\x67\xc4\xa8\x1c\x9d\x5b\xf7\x06\xe7\x31\xfb\x89\x59\x45\x47\x2e\xb8\x05\x1e\x1d\xfe\x0d\x00\x00\xff\xff\xcb\x0e\x06\xde\xc1\x0e\x00\x00")

func cloudJsonBytes() ([]byte, error) {
	return bindataRead(
		_cloudJson,
		"cloud.json",
	)
}

func cloudJson() (*asset, error) {
	bytes, err := cloudJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cloud.json", size: 3777, mode: os.FileMode(420), modTime: time.Unix(1453795200, 0)}
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
	"cloud.json": cloudJson,
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
	"cloud.json": {cloudJson, map[string]*bintree{}},
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
