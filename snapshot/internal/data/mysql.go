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

package data

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"power-snapshot/config"
	models "power-snapshot/internal/model"
)

type Mysql struct {
	*gorm.DB
}

func NewMysql() *Mysql {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       fmt.Sprintf("%s:%s@tcp(%s)/power-voting-filecoin?charset=utf8&parseTime=True&loc=Local", config.Client.Mysql.Username, config.Client.Mysql.Password, config.Client.Mysql.Url),
		DefaultStringSize:         256,   // string size
		DisableDatetimePrecision:  true,  // datetime precision is disabled. Databases earlier than MySQL 5.6 do not support dateTime precision
		DontSupportRenameIndex:    true,  // Rename indexes by deleting and creating new indexes. Databases before MySQL 5.7 and MariaDB do not support rename indexes
		DontSupportRenameColumn:   true,  // Rename columns with 'change'. Databases prior to MySQL 8 and MariaDB do not support renaming columns
		SkipInitializeWithVersion: false, // This parameter is automatically configured based on the current MySQL version
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,   // Singular table
		},
	})

	if err != nil {
		panic(fmt.Errorf("mysql connect error: %v", err))
	}

	db.AutoMigrate(&models.SnapshotBackupTbl{})

	return &Mysql{db}
}
