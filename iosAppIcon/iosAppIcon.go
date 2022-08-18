/**
 * @author      OA Wu <oawu.tw@gmail.com>
 * @copyright   Copyright (c) 2015 - 2021
 * @license     http://opensource.org/licenses/MIT  MIT License
 * @link        https://www.ioa.tw/
 */

package iosAppIcon

import (
  "os"
  "io"
  "fmt"
  "sync"
  "errors"
  "strings"
  "oago/lib"
  _fs "path/filepath"
  _file "oago/thumbnail/file"
  _lang "golang.org/x/text/language"
  _message "golang.org/x/text/message"
  _thumbnail "oago/thumbnail"
)

type Size struct {
  size uint16
  path string
  multiples []uint8
  file *_file.File
}
type Result struct {
  status bool
  message string
  error error
}

func copy(wg *sync.WaitGroup, mux *sync.Mutex, result chan<- *Result, src string, dest string, size *Size) {
  defer wg.Done()
  
  stat, e := os.Stat(src)
  if e != nil {
    result <- &Result{ status: false, message: fmt.Sprintf("無法取得 %s 檔案狀態！", src), error: e }
    return
  }
  if !stat.Mode().IsRegular() {
    result <- &Result{ status: false, message: fmt.Sprintf("%s 不是常規的檔案！", src) }
    return
  }

  mux.Lock()
  defer mux.Unlock()

  source, e := os.Open(src)
  if e != nil {
    result <- &Result{ status: false, message: fmt.Sprintf("無法開啟 %s 檔案！", src), error: e }
    return
  }
  defer source.Close()
  
  destination, e := os.Create(dest)
  if e != nil {
    result <- &Result{ status: false, message: fmt.Sprintf("無法建立 %s 檔案！", dest), error: e }
    return
  }
  defer destination.Close()

  _, e = io.Copy(destination, source)
  if e != nil {
    result <- &Result{ status: false, message: fmt.Sprintf("複製 %s ➜ %s 失敗！", src, dest), error: e }
    return
  }
  
  file, err, _ := _file.New(dest)
  if file == nil {
    result <- &Result{ status: false, message: fmt.Sprintf("無法取得 %s 檔案資訊！", dest), error: err }
  } else {
    size.file = file
    result <- &Result{ status: true }
  }
}
func createSizes(file *_file.File) ([]Size, error) {
  sizes := []Size{
    Size{ size: 20, multiples: []uint8{1, 2, 3} },
    Size{ size: 29, multiples: []uint8{1, 2, 3} },
    Size{ size: 40, multiples: []uint8{1, 2, 3} },
    Size{ size: 60, multiples: []uint8{1, 2, 3} },
    Size{ size: 76, multiples: []uint8{1, 2} },
    Size{ size: 167, multiples: []uint8{1} },
    Size{ size: 1024, multiples: []uint8{1} },
  }

  for i, size := range sizes {
    sizes[i].path = _fs.Join(file.BasePath, fmt.Sprintf("%s-%d%s", file.BaseName, size.size, file.Ext.Str()))
  }

  total := len(sizes)
  wg  := new(sync.WaitGroup)
  mux := new(sync.Mutex)
  ous := make(chan *Result, total)

  wg.Add(total)
  for i, size := range sizes {
    go copy(wg, mux, ous, file.Path, size.path, &sizes[i])
  }
  wg.Wait()

  for i := 1; i <= total; i++ {
    ou := <-ous
    if !ou.status {
      return sizes, errors.New(ou.message)
    }
  }

  return sizes, nil
}
func joinInt2Str(nums []uint8, str string) string {
  tmp := ""
  for _, num := range nums {
    if tmp == "" {
      tmp = fmt.Sprintf("%d", num)
    } else {
      tmp = tmp + str + fmt.Sprintf("%d", num)
    }
  }
  return tmp
}
func modify(file *_file.File) (special1 *_file.File, special2 *_file.File, err error) {
  baseName1 := fmt.Sprintf("%s-167@1x", file.BaseName)
  name1 := fmt.Sprintf("%s%s", baseName1, file.Ext.Str())
  special1 = &_file.File { Name: name1, Path: _fs.Join(file.BasePath, name1), Ext: file.Ext, BasePath: file.BasePath, BaseName: baseName1 }

  baseName2 := fmt.Sprintf("%s-83.5@2x", file.BaseName)
  name2 := fmt.Sprintf("%s%s", baseName2, file.Ext.Str())
  special2 = &_file.File { Name: name2, Path: _fs.Join(file.BasePath, name2), Ext: file.Ext, BasePath: file.BasePath, BaseName: baseName2 }

  err = os.Rename(special1.Path, special2.Path)
  return
}
func Run(arguments []string) {

  // 檢查參數
  if len(arguments) == 0 {
    lib.ShowHelpIosAppIcon()
    return
  }

  // 取得檔案
  file, err, _ := _file.New(lib.Args(arguments, "-P", "--pic"))
  if file == nil {
    lib.ShowHelpIosAppIcon(err)
    return
  }

  // 複製檔案
  sizes, err := createSizes(file)
  if err != nil {
    lib.ShowHelpIosAppIcon(err)
    return
  }

  // 各別縮圖
  total := len(sizes)
  wg := new(sync.WaitGroup)
  ins := make(chan []string, total)
  ous := make(chan *_thumbnail.Thumbnail, total)

  wg.Add(total)

  for i := 0; i < total; i++ {
    go _thumbnail.RunSync(wg, ins, ous)
  }

  for _, size := range sizes {
    ins <- []string{ "-P", size.file.Path, "-S", fmt.Sprintf("%d", size.size), "-M", joinInt2Str(size.multiples, ",") }
  }
  close(ins)
  wg.Wait()

  thumbnails := []*_thumbnail.Thumbnail{}
  for i := 0; i < total; i++ {
    thumbnails = append(thumbnails, <-ous)
  }

  close(ous)

  // 結果
  results := []_thumbnail.Result{}
  for _, thumbnail := range thumbnails {
    if thumbnail.Error != nil {
      lib.ShowHelpIosAppIcon(thumbnail.Error)
      return
    }

    for _, result := range thumbnail.Results {
      results = append(results, result)
    }
  }

  // 刪除複製檔
  wg = new(sync.WaitGroup)
  wg.Add(total)
  for _, size := range sizes {
    go func (wg *sync.WaitGroup, path string) {
      defer wg.Done()
      os.Remove(path)
    }(wg, size.path)
  }
  wg.Wait()

  // 處理 83.5@2x 檔案問題
  s1, s2, err := modify(file)
  if err != nil {
    lib.ShowHelpIosAppIcon(err)
    return
  }

  type Tmp struct{
    path string
    size string
    status string
  }

  tmps := []Tmp{}

  for _, result := range results {
    path := result.Multiple.File.Path
    if path == s1.Path {
      path = s2.Path
    }

    if result.Status {
      tmps = append(tmps, Tmp{ path: path, size: fmt.Sprintf("%d byte(%.2f%%)", result.Multiple.File.Dimension.Size, float64(result.Multiple.File.Dimension.Size) / float64(file.Dimension.Size) * 100), status: "\x1b[38;5;2m成功\x1b[0m" })
    } else {
      tmps = append(tmps, Tmp{ path: path, size: fmt.Sprintf("%d byte(%.2f%%)", result.Multiple.File.Dimension.Size, float64(result.Multiple.File.Dimension.Size) / float64(file.Dimension.Size) * 100), status: "\x1b[38;5;1m失敗\x1b[0m" })
    }
  }

  // 處理輸出排版
  w1 := uint64(0)
  w2 := uint64(0)
  for _, tmp := range tmps {
    t1 := lib.StrWidth(tmp.path)
    t2 := lib.StrWidth(tmp.size)
    if t1 > w1 {
      w1 = t1
    }
    if t2 > w2 {
      w2 = t2
    }
  }

  formater := _message.NewPrinter(_lang.English)

  fmt.Println("")
  fmt.Println("  原始圖檔")
  fmt.Println("")
  size  := fmt.Sprintf("%d byte(100%%)", file.Dimension.Size)
  t1 := strings.Repeat(" ", int(w1 - lib.StrWidth(file.Path)))
  t2 := strings.Repeat(" ", int(w2 - lib.StrWidth(size)))
  fmt.Println(formater.Sprintf("    %s%s ｜ %s%s", file.Path, t1, size, t2))
  fmt.Println("")
  fmt.Println("  轉換結果")
  fmt.Println("")

  for _, tmp := range tmps {
    t1 := strings.Repeat(" ", int(w1 - lib.StrWidth(tmp.path)))
    t2 := strings.Repeat(" ", int(w2 - lib.StrWidth(tmp.size)))
    fmt.Println(formater.Sprintf("    %s%s ｜ %s%s ─ %s", tmp.path, t1, tmp.size, t2, tmp.status))
  }

  fmt.Println("")
}