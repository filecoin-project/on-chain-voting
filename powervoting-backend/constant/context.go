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

package constant

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Context struct {
	*gin.Context
	*validator.Validate
}

// BindAndValidate binds and validates request params.
func (c *Context) BindAndValidate(req any) error {
	// 1. BindUri binds the passed struct pointer using the specified binding engine.
	if err := c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("router parameter error: %v", err)
	}
	
	// 2. BindQuery binds the passed struct pointer using the specified binding engine.
	if err := c.ShouldBindQuery(req); err != nil {
		return fmt.Errorf("query parameter error: %v", err)
	}

	// 3. BindJSON binds the passed struct pointer using the specified binding engine.
	if err := c.ShouldBind(req); err != nil {
		return fmt.Errorf("body parameter error: %v", err)
	}

	// 4. Validate the request params.
	if err := c.Validate.Struct(req); err != nil {
		zap.L().Error("Param validate error", zap.Error(err))
		return fmt.Errorf("param validate error: %v", err)
	}

	return nil
}
