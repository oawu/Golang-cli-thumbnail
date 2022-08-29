/**
 * @author      OA Wu <oawu.tw@gmail.com>
 * @copyright   Copyright (c) 2015 - 2022
 * @license     http://opensource.org/licenses/MIT  MIT License
 * @link        https://www.ioa.tw/
 */

package multiple

import (
	"fmt"
	_dimension "oago/thumbnail/dimension"
	_file "oago/thumbnail/file"
	_fs "path/filepath"
	"strconv"
	"strings"
)

type Multiple struct {
	Title string
	File  *_file.File
}

func trim(strs []string) []string {
	tmps := []string{}
	for _, str := range strs {
		tmps = append(tmps, strings.Trim(str, " "))
	}
	return tmps
}

func toUint8(strs []string) []uint8 {
	tmps := []uint8{}
	for _, str := range strs {
		val, err := strconv.ParseInt(str, 10, 16)
		if err != nil || val > 255 {
			continue
		}

		tmps = append(tmps, uint8(val))
	}
	return tmps
}

func New(multiple string, file *_file.File, dimension *_dimension.Dimension) ([]Multiple, error, error) {
	if multiple == "" {
		multiple = "1,2,3"
	}

	multiples := []Multiple{}

	nums := toUint8(trim(strings.Split(multiple, ",")))

	for _, num := range nums {
		baseName := fmt.Sprintf("%s@%dx", file.BaseName, num)
		basePath := file.BasePath
		ext := file.Ext
		name := fmt.Sprintf("%s%s", baseName, ext.Str())
		path := _fs.Join(basePath, name)

		multiple := uint(num)
		width := multiple * dimension.Width
		height := multiple * dimension.Height

		multiples = append(multiples, Multiple{Title: fmt.Sprintf("@%dx", num), File: &_file.File{Name: name, Path: path, Ext: ext, BasePath: basePath, BaseName: baseName, Dimension: &_dimension.Dimension{Width: width, Height: height, Gradient: float64(height) / float64(width)}}})
	}

	return multiples, nil, nil
}
