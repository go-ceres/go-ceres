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

package schedule

import (
	"github.com/go-ceres/go-ceres/errors"
	"github.com/robfig/cron/v3"
)

type Schedule struct {
	Cron   *cron.Cron
	config *Config
}

// newSchedule 创建任务调度
func newSchedule(c *Config) *Schedule {
	opts := c.Opts
	opts = append(opts, cron.WithLogger(c.log))
	return &Schedule{
		config: c,
		Cron:   cron.New(opts...),
	}
}

// Size 获取当前任务数量
func (s *Schedule) Size() int {
	return len(s.Cron.Entries())
}

// AddJob 添加任务
func (s *Schedule) AddJob(job Job) (int, error) {
	// 如果超出了任务数量
	if s.Size() >= s.config.Size {
		return 0, errors.New(errors.CodeAddScheduleToMaximum, errors.MsgAddScheduleToMaximum)
	}
	id, err := s.Cron.AddJob(job.Cron(), job)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// RemoveJob 移除任务
func (s *Schedule) RemoveJob(id int) {
	s.Cron.Remove(cron.EntryID(id))
}

// Get 根据id获取单个任务
func (s *Schedule) Get(id int) cron.Entry {
	return s.Cron.Entry(cron.EntryID(id))
}

// List 任务列表
func (s *Schedule) List() []cron.Entry {
	return s.Cron.Entries()
}

// Run 阻塞启动
func (s *Schedule) Run() {
	s.Cron.Run()
}

func (s *Schedule) Stop() {
	s.Cron.Stop()
}

// Start 启动
func (s *Schedule) Start() {
	s.Cron.Start()
}
