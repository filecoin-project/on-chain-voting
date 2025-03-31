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

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"powervoting-server/constant"
	"powervoting-server/model/api"
)

// R is response function that formats and sends a JSON.
func R(code int, data any, message string, c *gin.Context) {
	c.JSON(http.StatusOK, api.Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Success sends a success response with status code 1 and message "ok".
func Success(c *gin.Context) {
	R(constant.CodeOK, nil, constant.CodeOKStr, c)
}

// SuccessWithData sends a success response with status code 1, data, and message "ok".
func SuccessWithData(c *gin.Context, data any) {
	R(constant.CodeOK, data, constant.CodeOKStr, c)
}

// Error sends an error response with status code 0 and the error message.
func Error(c *gin.Context, err error) {
	R(constant.CodeError, nil, err.Error(), c)
}

// ParamError sends a response indicating a parameter error.
func ParamError(c *gin.Context) {
	R(constant.CodeParamError, nil, constant.CodeParamErrorStr, c)
}

// SystemError sends a response indicating a system error.
func SystemError(c *gin.Context) {
	R(constant.CodeSystemError, nil, constant.CodeSystemErrorStr, c)
}
