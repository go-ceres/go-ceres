//    Copyright 2021. Go-Ceres
//    Author https://github.com/go-ceres/go-ceres
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package old

import (
	"image"
	"image/color"
)

type ImageBuf struct {
	i image.Image
	w int
	h int
}

func (i *ImageBuf) getHeight() int {
	return i.h
}

func (i *ImageBuf) getWidth() int {
	return i.w
}

func (i *ImageBuf) getRGBA(x, y int) color.RGBA64 {
	r, g, b, a := i.i.At(x, y).RGBA()
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func (i *ImageBuf) setRGBA(x, y int, c color.Color) {
	switch i.i.(type) {
	case *image.RGBA:
		i.i.(*image.RGBA).Set(x, y, c)
	case *image.NRGBA:
		i.i.(*image.NRGBA).Set(x, y, c)
	}
}
