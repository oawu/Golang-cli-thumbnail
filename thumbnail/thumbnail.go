/**
 * @author      OA Wu <oawu.tw@gmail.com>
 * @copyright   Copyright (c) 2015 - 2021
 * @license     http://opensource.org/licenses/MIT  MIT License
 * @link        https://www.ioa.tw/
 */

package thumbnail

import (
	"errors"
	"fmt"
	_resize "github.com/nfnt/resize"
	_lang "golang.org/x/text/language"
	_message "golang.org/x/text/message"
	"image"
	_jpeg "image/jpeg"
	_png "image/png"
	"oago/lib"
	_dimension "oago/thumbnail/dimension"
	_file "oago/thumbnail/file"
	_multiple "oago/thumbnail/multiple"
	"os"
	"strings"
	"sync"
)

type Result struct {
	Status   bool
	Multiple _multiple.Multiple
}

type Thumbnail struct {
	Results []Result
	file    *_file.File
	Error   error
}

func thumbnail(img image.Image, ins <-chan _multiple.Multiple, ous chan<- *Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for multiple := range ins {

		img = _resize.Resize(multiple.File.Dimension.Width, multiple.File.Dimension.Height, img, _resize.NearestNeighbor)

		src, err := os.Create(multiple.File.Path)
		if err != nil {
			ous <- &Result{Status: false, Multiple: multiple}
			return
		}

		if multiple.File.Ext == _file.JPG {
			_jpeg.Encode(src, img, &_jpeg.Options{100})
		} else {
			_png.Encode(src, img)
		}

		src.Close()

		multiple.File.Dimension.Size = lib.FileSize(multiple.File.Path)
		ous <- &Result{Status: true, Multiple: multiple}
	}
}

func cover(file *_file.File, dimension *_dimension.Dimension, multiples []_multiple.Multiple) ([]Result, error, error) {
	returns := []Result{}
	reader, e := os.Open(file.Path)
	if e != nil {
		return returns, errors.New(fmt.Sprintf("無法讀取 %s 檔案(1)！", file.Path)), e
	}
	defer reader.Close()

	img, _, e := image.Decode(reader)
	if e != nil {
		return returns, errors.New(fmt.Sprintf("無法讀取 %s 檔案(2)！", file.Path)), e
	}

	total := len(multiples)

	wg := new(sync.WaitGroup)
	ins := make(chan _multiple.Multiple, total)
	ous := make(chan *Result, total)

	for i := 0; i < total; i++ {
		wg.Add(1)
		go thumbnail(img, ins, ous, wg)
	}

	for _, multiple := range multiples {
		ins <- multiple
	}

	close(ins)
	wg.Wait()

	for i := 1; i <= total; i++ {
		returns = append(returns, *<-ous)
	}
	close(ous)

	return returns, nil, nil
}

func RunSync(wg *sync.WaitGroup, ins <-chan []string, ous chan<- *Thumbnail) {
	defer wg.Done()

	for arguments := range ins {
		results := []Result{}

		// 檢查參數
		if len(arguments) == 0 {
			ous <- &Thumbnail{Results: results, file: nil, Error: errors.New("參數錯誤！")}
			return
		}

		// 取得檔案
		file, err, _ := _file.New(lib.Args(arguments, "-P", "--pic"))
		if file == nil {
			ous <- &Thumbnail{Results: results, file: file, Error: err}
			return
		}

		// 取得尺寸
		dimension, err, _ := _dimension.New(lib.Args(arguments, "-W", "--width"), lib.Args(arguments, "-H", "--height"), lib.Args(arguments, "-S", "--size"), file.Dimension.Gradient)
		if dimension == nil {
			ous <- &Thumbnail{Results: results, file: file, Error: err}
			return
		}

		// 取得倍率
		multiples, err, _ := _multiple.New(lib.Args(arguments, "-M", "--multiple"), file, dimension)
		if len(multiples) <= 0 {
			ous <- &Thumbnail{Results: results, file: file, Error: err}
			return
		}

		// 轉換
		results, err, _ = cover(file, dimension, multiples)
		if len(results) <= 0 {
			ous <- &Thumbnail{Results: results, file: file, Error: err}
			return
		}
		ous <- &Thumbnail{Results: results, file: file, Error: nil}
	}
}
func Run(arguments []string) {

	wg := new(sync.WaitGroup)
	ins := make(chan []string, 1)
	ous := make(chan *Thumbnail, 1)

	wg.Add(1)
	go RunSync(wg, ins, ous)
	ins <- arguments
	close(ins)

	wg.Wait()

	ou := <-ous
	close(ous)

	results := ou.Results
	file := ou.file
	err := ou.Error

	if err != nil {
		lib.ShowHelpThumbnail(err)
		return
	}

	// 結果
	formater := _message.NewPrinter(_lang.English)
	width := uint64(0)
	for _, result := range results {
		tmp := lib.StrWidth(result.Multiple.Title)
		if tmp > width {
			width = tmp
		}
	}

	tmp1 := strings.Repeat(" ", int(width-lib.StrWidth("@0x")))
	tmp2 := strings.Repeat(" ", int(width))

	fmt.Println("")
	fmt.Println("  原始圖檔")
	fmt.Println("")
	fmt.Println(formater.Sprintf("    %s@0x ｜ %s%s ｜ %d byte(100%%)", tmp1, tmp2, file.Path, file.Dimension.Size))
	fmt.Println("")
	fmt.Println("  轉換結果")
	fmt.Println("")
	for _, result := range results {
		tmp := strings.Repeat(" ", int(width-lib.StrWidth(result.Multiple.Title)))
		if result.Status {
			fmt.Println(formater.Sprintf("    %s%s ｜ %s%s ｜ %d byte(%.2f%%) ─ %s", tmp, result.Multiple.Title, tmp, result.Multiple.File.Path, result.Multiple.File.Dimension.Size, float64(result.Multiple.File.Dimension.Size)/float64(file.Dimension.Size)*100, "\x1b[38;5;2m成功\x1b[0m"))
		} else {
			fmt.Println(formater.Sprintf("    %s%s ｜ %s%s ｜ %d byte(%.2f%%) ─ %s", tmp, result.Multiple.Title, tmp, result.Multiple.File.Path, result.Multiple.File.Dimension.Size, float64(result.Multiple.File.Dimension.Size)/float64(file.Dimension.Size)*100, "\x1b[38;5;1m失敗\x1b[0m"))
		}
	}
	fmt.Println("")
}
