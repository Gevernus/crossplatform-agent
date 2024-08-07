// Code generated for package main by go-bindata DO NOT EDIT. (@generated)
// sources:
// tray_icon_24x24.png
package assets

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

// Mode return file modify time
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

var _tray_icon_24x24Png = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x00\xda\x03\x25\xfc\x89\x50\x4e\x47\x0d\x0a\x1a\x0a\x00\x00\x00\x0d\x49\x48\x44\x52\x00\x00\x00\x18\x00\x00\x00\x18\x08\x02\x00\x00\x00\x6f\x15\xaa\xaf\x00\x00\x03\xa1\x49\x44\x41\x54\x78\x9c\x4d\x55\xcb\x8e\x1c\x45\x10\x8c\xc8\xaa\x7e\xcc\x32\xcc\x6a\xd9\x65\xbd\x27\x04\x48\xac\x05\x48\xb6\x10\x70\xe7\x64\x24\xe0\x0b\x6c\x03\xf2\x09\xf0\xbf\x98\xa3\x91\x40\xbe\x98\xbf\x30\x47\x73\x40\xe6\xe0\x05\xdb\x08\x89\xe5\x61\x2f\x12\x62\x25\xe4\x07\xdb\xf3\xe8\x9e\xca\xe0\x50\x35\x63\xe6\x30\xd3\xd3\x99\x15\x95\x95\x19\x11\xc5\x8b\xe7\xcf\x57\x75\x05\x01\x70\x09\x20\xe5\x2e\x09\x24\x24\x90\x04\x05\x41\x22\x09\x80\x24\x49\x49\x80\x48\x03\x20\xc9\xcc\xac\x69\x1a\x49\x66\x2c\x39\x25\x93\xc8\x20\x92\xa4\x0c\x90\x51\x00\xe4\xfd\x32\x0a\x49\x33\x73\xb9\x9d\x9c\x9c\x54\x55\x1c\x86\x25\x00\x01\x2e\x97\x24\x21\x6f\x05\x40\xf9\x57\x10\xb0\x0a\x09\x39\x03\x14\x24\x88\x30\x6b\x9a\x66\x36\x9d\x55\x75\x95\x13\xf2\x17\x59\x1e\x32\x44\x5e\x8a\x35\xfe\x3a\x54\x92\x00\x20\x02\x68\x9a\xc6\x53\x62\xd9\x5c\x24\x87\x61\x90\x04\x10\x2b\x10\x10\x00\xcd\x48\x32\x86\x28\x09\xc4\xba\x77\x82\x22\xcd\xe6\xf3\x39\xc9\xba\x6e\xdc\x93\xdc\x17\xfd\x30\xd9\xdc\xcc\x00\xb9\x65\xca\x58\x80\x52\x12\x70\x32\xed\xda\xba\x29\x71\x30\x97\x16\x25\xb5\xa3\xd1\x6c\x3a\x5d\x2c\xe6\x4d\xd3\x0c\xc3\xf0\xde\xfb\x1f\x4c\x36\x27\x7d\xdf\x87\x10\xd6\x73\x92\x14\x43\x98\x6c\x4d\xfe\xfa\xe7\xf8\xee\xad\x1f\x7e\xfb\xfd\xd7\xb6\x6d\x33\x7a\x9e\x44\x1c\x86\x81\x64\xdb\x8e\xfa\xc5\x02\xd2\x33\xe3\xf1\xee\xde\xa9\x6b\x5f\x7d\x69\x16\x48\xe6\x06\x08\x80\x54\x5b\x75\xe6\x9d\xb7\x67\xa3\xe5\x5b\x6f\xbc\xf9\xd3\xcf\xf7\x46\xa3\x91\xe4\x85\x1f\x92\x79\x4a\x80\xdc\x53\xdd\xd4\xf3\xc5\x22\xc6\xd8\xf7\x8b\x10\x62\x5d\xd7\x31\xc6\x18\xab\x18\x63\x55\x55\x4d\xd3\x24\xea\xfe\x9d\xc3\xd3\xf5\x0b\x0b\xef\x0d\xb6\x6e\xbb\x24\xc1\xad\x6d\x47\xb9\xbc\xe4\x1e\x42\x48\x29\x19\x99\x07\xbd\xa6\x82\x24\x77\x55\x55\x3c\x3a\x7a\x70\xf0\xe3\x81\xd1\x56\x27\xca\x25\xcb\x68\xb6\x4c\x99\x41\x65\xb5\x99\x91\x96\x73\x8c\x86\xdc\x6d\xe4\x20\x9a\xa6\x19\x8d\x46\x34\x62\xf5\x51\x66\x93\x14\x73\x23\x09\x89\xcc\x92\xe8\xfb\xbe\xeb\xba\xb6\x6d\x52\x72\xe4\x28\x61\xa4\x59\x98\x4e\xa7\x5d\xd7\xa1\xb0\x55\x92\x93\x96\xf9\x1a\x49\x64\x35\x19\x00\x22\xb9\xef\xed\xed\xed\x9f\xde\x27\x2d\x73\xae\x69\x9b\xbe\x1f\x52\x4a\x21\x84\x7f\x9f\x3c\x3e\x73\xf6\x6c\x8c\x55\x4a\x89\xa4\x8a\xe8\x00\x30\xca\x85\x00\x82\x00\x82\x85\xae\xeb\x5e\xd9\xdf\xbf\x72\xe5\xf3\x75\xf1\xdf\xdd\xfa\xfe\xf5\xd7\x5e\x7d\x76\x3c\x5e\xbf\xb9\xfe\xf5\x75\x0b\x01\xca\x42\xc9\xea\x53\x2c\x27\x25\x05\x58\xb0\xd9\x49\xf7\xc5\xd5\xab\xe7\xde\x3d\xd7\xf7\x7d\x4a\xc9\x42\xec\xfb\xe5\xed\x83\x83\xb4\x1c\x48\x6b\xdb\xf6\xc9\x93\xc7\xdf\xde\xbc\xb9\xb1\xb1\xe1\xee\x00\x4a\x09\x24\x3f\xba\x78\x31\x37\x9f\x34\xc9\x25\x74\x5d\xd7\xf7\x7d\xac\xaa\x17\x5f\x7a\xb9\x9b\x2d\x2e\x7f\xf6\xc9\xed\xdb\x07\xdf\xdc\xb8\xb1\xb3\xbd\x75\x78\xf8\x4b\x55\xd5\x1b\x1b\xa3\xaa\xaa\xb8\x82\x61\x39\x1a\x64\xa4\xc0\x52\x97\xfb\x78\x3c\x9e\xcd\x67\x93\xc9\xe6\xc7\x97\x2e\x3d\x38\x3a\xfa\xfb\xf8\xf8\xf9\xdd\x9d\x4f\x2f\x5f\xfe\xf3\xfe\x1f\xf7\xee\xde\x39\xb5\xbb\x9b\x6b\x21\x59\x14\x48\xb8\xc4\x0f\x2f\x5c\x08\x31\x40\x4f\x19\x9c\x52\x4a\xc9\x1f\x3e\x7a\x04\xda\xe6\xd6\x56\x5a\xa6\xec\x71\x8f\x1f\x3d\xdc\x7e\x6e\x2b\x98\x85\x10\xcc\x98\x09\x43\x98\x20\xb9\xc7\xe2\x55\x28\x36\x96\xad\x2a\x04\xec\xec\x6c\xa7\x65\x5a\xf6\xf3\xcc\x19\x82\xbb\x3b\xdb\x20\xcc\x6c\x6d\x0a\xc8\x7a\x06\x48\x8b\x85\xe2\x85\x5e\xc5\x63\x32\x7a\x56\x88\xd6\x64\x95\x42\x08\x6b\x73\x41\x11\x62\x59\x66\xeb\x97\xe5\x74\x65\xa2\x5c\xe5\x14\xee\x81\x74\xc9\x57\xce\xfd\xff\x5d\xb3\x0f\x1a\x57\x4e\xfd\xd4\x11\x08\x40\x79\x1e\xe5\x1f\x00\xa1\x8a\x55\xe9\x02\xb8\xf6\xde\xd5\xc4\x69\xc5\xe4\xb1\xb2\xf7\xd5\x3d\x51\x04\x06\x00\x34\xb2\xae\xab\x94\x96\xa4\x01\xc2\x6a\xf2\x66\xc1\xdd\x69\xe6\x2e\x13\x28\x57\xbe\x82\x04\x64\xbd\x4b\x70\x77\x01\xee\x72\x79\x5d\xd7\x5d\x37\xa5\x99\xcb\xb5\xba\x03\xf2\x4d\xb1\x1c\x96\xd9\x76\xff\x03\x69\x41\x53\x07\xb2\xb8\x61\x63\x00\x00\x00\x00\x49\x45\x4e\x44\xae\x42\x60\x82\x01\x00\x00\xff\xff\xb0\x05\xfb\x80\xda\x03\x00\x00")

func tray_icon_24x24PngBytes() ([]byte, error) {
	return bindataRead(
		_tray_icon_24x24Png,
		"tray_icon_24x24.png",
	)
}

func tray_icon_24x24Png() (*asset, error) {
	bytes, err := tray_icon_24x24PngBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "tray_icon_24x24.png", size: 986, mode: os.FileMode(420), modTime: time.Unix(1723037139, 0)}
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
	"tray_icon_24x24.png": tray_icon_24x24Png,
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
	"tray_icon_24x24.png": &bintree{tray_icon_24x24Png, map[string]*bintree{}},
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
