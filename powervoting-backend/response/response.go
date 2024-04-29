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

package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"powervoting-server/model"
)

func R(code int, data any, message string, c *gin.Context) {
	c.JSON(http.StatusOK, model.Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Success return success
func Success(c *gin.Context) {
	R(1, nil, "ok", c)
}

func SuccessWithMsg(msg string, c *gin.Context) {
	R(1, nil, msg, c)
}

func SuccessWithData(data any, c *gin.Context) {
	R(1, data, "ok", c)
}

// Fail return fail
func Fail(code int, message string, c *gin.Context) {
	R(code, nil, message, c)
}

// Error return error
func Error(err error, c *gin.Context) {
	R(0, nil, err.Error(), c)
}

// ParamError param error
func ParamError(c *gin.Context) {
	R(0, nil, "param error", c)
}

func SystemError(c *gin.Context) {
	R(0, nil, "system error", c)
}
