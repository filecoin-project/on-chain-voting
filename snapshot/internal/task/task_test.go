// Copyright (C) 2023-2024 StorSwift Inc.
// This file is part of the PowerVoting library.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package task

import (
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/zap"

	"power-snapshot/config"
)

func TestRepeatTask(t *testing.T) {
	isRunning := int32(0)
	config.InitLogger()
	for i := range 5 {
		if atomic.CompareAndSwapInt32(&isRunning, 0, 1) {
			
			zap.L().Info("Task is running", zap.Int("task id", i))
			go func() {
				time.Sleep(3 * time.Second)
				defer atomic.StoreInt32(&isRunning, 0)
			}()
		} else {
			zap.L().Info("Task is already running")
		}
		time.Sleep(time.Second)
	}
}
