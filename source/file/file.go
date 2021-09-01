//  Copyright 2020 Go-Ceres
//  Author https://github.com/go-ceres/go-ceres
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package file

import (
	"github.com/fsnotify/fsnotify"
	"github.com/go-ceres/go-ceres/config"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type fileSource struct {
	path      string
	dir       string
	unmarshal string
	watcher   *fsnotify.Watcher
	changed   chan struct{}
}

func (fs *fileSource) Read() (*config.DataSet, error) {
	fh, err := os.Open(fs.path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = fh.Close()
	}()
	b, err := ioutil.ReadAll(fh)
	if err != nil {
		return nil, err
	}
	info, err := fh.Stat()
	if err != nil {
		return nil, err
	}
	cs := &config.DataSet{
		Format:    fs.getUnmarshal(),
		Source:    fs.String(),
		Timestamp: info.ModTime(),
		Data:      b,
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

func (fs *fileSource) Write(dataSet *config.DataSet) error {
	return nil
}

func (fs *fileSource) IsChanged() <-chan struct{} {
	return fs.changed
}

func (fs *fileSource) Watch() {
	if fs.watcher == nil {
		var err error
		fs.watcher, err = fsnotify.NewWatcher()
		fs.changed = make(chan struct{}, 1)
		if err != nil {
			log.Fatal("new watch", err)
		}
	}
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event := <-fs.watcher.Events:
					const writeOrCreateMask = fsnotify.Write | fsnotify.Create
					if event.Op&writeOrCreateMask != 0 && filepath.Clean(event.Name) == filepath.Clean(fs.path) {
						log.Println("modified file: ", event.Name)
						select {
						case fs.changed <- struct{}{}:
						default:
						}
					}
				case err := <-fs.watcher.Errors:
					log.Printf("watcher error: %v\n", err)
				}
			}
		}()
		err := fs.watcher.Add(fs.dir)
		if err != nil {
			log.Fatal("add file to fsnotify error:", err)
		}
		initWG.Done()
		eventsWG.Wait()
	}()
	initWG.Wait()
}

func (fs *fileSource) UnWatch() {
	close(fs.changed)
	fs.changed = nil
	_ = fs.watcher.Close()
	fs.watcher = nil
}

func (fs *fileSource) String() string {
	return "file"
}

func (fs *fileSource) getUnmarshal() string {
	parts := strings.Split(fs.path, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return fs.unmarshal
}

func NewSource(file string, opts ...Option) config.Source {
	file, err := filepath.Abs(file)
	if err != nil {
		panic("new file source：" + err.Error())
	}
	eng := fileSource{
		path: file,
		dir:  filepath.Dir(file),
	}
	for _, o := range opts {
		o(&eng)
	}
	return &eng
}
