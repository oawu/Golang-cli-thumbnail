/**
 * @author      OA Wu <oawu.tw@gmail.com>
 * @copyright   Copyright (c) 2015 - 2021
 * @license     http://opensource.org/licenses/MIT  MIT License
 * @link        https://www.ioa.tw/
 */

package dimension

import (
  "math"
  "errors"
  "strconv"
)

type Dimension struct {
  Width uint
  Height uint
  Gradient float64
  Size uint64
}

func New(w string, h string, s string, g float64) (*Dimension, error, error) {
  width, _ := strconv.ParseInt(w, 10, 64)
  height, _ := strconv.ParseInt(h, 10, 64)

  if width > 0 && height > 0 {
    return &Dimension{ Width: uint(width), Height: uint(height), Gradient: float64(height) / float64(width) }, nil, nil
  }

  if width > 0 && height <= 0 {
    return &Dimension{ Width: uint(width), Height: uint(math.Round(g * float64(width))), Gradient: g }, nil, nil
  }

  if height > 0 && width <= 0 {
    return &Dimension{ Width: uint(math.Round(float64(height) / g)), Height: uint(height), Gradient: g }, nil, nil
  }
  
  size, _ := strconv.ParseInt(s, 10, 64)

  if size > 0 {
    return &Dimension{ Width: uint(size), Height: uint(size), Gradient: 1 }, nil, nil
  }

  return nil, errors.New("沒有給予要轉換的尺寸！"), nil
}