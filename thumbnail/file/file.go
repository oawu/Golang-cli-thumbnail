/**
 * @author      OA Wu <oawu.tw@gmail.com>
 * @copyright   Copyright (c) 2015 - 2021
 * @license     http://opensource.org/licenses/MIT  MIT License
 * @link        https://www.ioa.tw/
 */

package file

import (
  "os"
  "fmt"
  "errors"
  "strings"
  _ "image/png"
  _ "image/jpeg"
  _fs "path/filepath"
  _img "image"
  _dimension "oago/thumbnail/dimension"
)

type Ext int

const (
  JPG Ext = iota
  PNG
)

func (ext Ext)Str() string {
  if ext == JPG { return ".jpg" }
  if ext == PNG { return ".png" }
  return ""
}

type File struct {
  Name string
  Path string
  Ext Ext

  BasePath string
  BaseName string

  Dimension *_dimension.Dimension
}

func New(name string) (*File, error, error) {
  if name == "" {
    return nil, errors.New("每有給予檔案名稱！"), nil
  }

  path, e := os.Getwd()
  if e != nil {
    return nil, errors.New("無法取得檔案目錄！"), e
  }
  
  if _fs.IsAbs(name) {
    path = name
  } else {
    path = _fs.Join(path, name)
  }

  basePath := _fs.Dir(path)
  name = _fs.Base(path)

  reader, e := os.Open(path)
  if e != nil {
    return nil, errors.New(fmt.Sprintf("無法讀取 %s 檔案(1)！", path)), e
  }
  defer reader.Close()
  
  stat, e := reader.Stat()
  if e != nil {
    return nil, errors.New(fmt.Sprintf("無法讀取 %s 檔案(2)！", path)), e
  }
  size := uint64(stat.Size())

  img, format, e := _img.DecodeConfig(reader)
  if e != nil {
    return nil, errors.New(fmt.Sprintf("無法讀取 %s 檔案(3)！", path)), e
  }

  var ext Ext
  switch {
  case format == "jpeg" || format == "jpg": ext = JPG
  case format == "png": ext = PNG
  default: return nil, errors.New("目前僅支援 .jpg、.png 格式的圖檔！"), nil
  }

  baseName := strings.TrimSuffix(name, ext.Str())
  dimension := &_dimension.Dimension{
    Width: uint(img.Width),
    Height: uint(img.Height),
    Gradient: float64(img.Height) / float64(img.Width),
    Size: size, }

  return &File{ Name: name, Path: path, Ext: ext, BasePath: basePath, BaseName: baseName, Dimension: dimension }, nil, nil
}