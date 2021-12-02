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
	"bytes"
	"image"
	"io/ioutil"
	"path"
	"strings"
)

var (
	bgImgCache [][]byte //缓存背景图片
	bkImgCache [][]byte //缓存滑块模板图片
)

func LoadBackgroudImages(path string) (err error) {
	bgImgCache, err = loadImages(path)
	return
}

func LoadBlockImages(path string) (err error) {
	bkImgCache, err = loadImages(path)
	return
}

func loadImages(basePath string) ([][]byte, error) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}
	var fileArr [][]byte
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".png") {
			buf, err := ioutil.ReadFile(path.Join(basePath, f.Name()))
			if err != nil {
				return nil, err
			}
			fileArr = append(fileArr, buf)
		}
	}
	return fileArr, nil
}

// randBackgroudImage 随机抽取 背景图
func randBackgroudImage() (*ImageBuf, error) {
	n := r.Intn(len(bgImgCache))
	buf := bgImgCache[n]
	im, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	return &ImageBuf{
		w: im.Bounds().Dx(),
		h: im.Bounds().Dy(),
		i: im,
	}, nil
}

// randBlockImage 随机抽取 滑块图，和干扰图
func randBlockImage() (a *ImageBuf, b *ImageBuf, err error) {
	l := len(bkImgCache)
	n := r.Intn(len(bkImgCache))
	buf := bkImgCache[n]
	im, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, nil, err
	}
	var next = n + 1
	if next == l {
		next = 0
	}
	buf2 := bkImgCache[next]
	im2, _, err := image.Decode(bytes.NewReader(buf2))
	if err != nil {
		return nil, nil, err
	}
	a = &ImageBuf{
		w: im.Bounds().Dx(),
		h: im.Bounds().Dy(),
		i: im,
	}
	b = &ImageBuf{
		w: im2.Bounds().Dx(),
		h: im2.Bounds().Dy(),
		i: im2,
	}
	return
}