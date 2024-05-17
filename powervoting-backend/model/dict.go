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

package model

// Dict records the ID of the synchronized vote.
type Dict struct {
	Id    int64  `json:"id"`                    // ID of the synchronized vote
	Name  string `json:"name" gorm:"not null"`  // Key representing the name
	Value string `json:"value" gorm:"not null"` // Value associated with the name key
}
