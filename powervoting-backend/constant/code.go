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

// Status codes for API responses.
const (
	CodeOK             = 0    // CodeOK indicates a successful operation.
	CodeParamError     = 1000 // CodeParamError indicates an error due to invalid input parameters.
	CodeDataExistError = 1001 // CodeDataExistError indicates an error when the data already exists.
	CodeError          = 1002 // CodeError indicates an  error.
	CodeSystemError    = 9999 // CodeSystemError indicates an internal system error.
)

// Status messages for API responses.
const (
	CodeOKStr             = "success"              // CodeOKStr is the message for a successful operation.
	CodeParamErrorStr     = "param error"          // CodeParamErrorStr is the message for invalid input parameters.
	CodeSystemErrorStr    = "system error"         // CodeSystemErrorStr is the message for internal system errors.
	CodeDataExistErrorStr = "already exists error" // CodeDataExistErrorStr is the message when the data already exists.
)
