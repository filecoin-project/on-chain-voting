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

import "time"

// BaseField is the base field for all database models
type BaseField struct {
	ID        int64    `gorm:"primary_key;auto_increment;column:id;not null" json:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"not null,autoUpdateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null,autoUpdateTime"`
}
