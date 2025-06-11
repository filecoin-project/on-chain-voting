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
	"errors"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Context struct {
	*gin.Context
	*validator.Validate
}

type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationErrorResponse

func (v ValidationErrors) Errors() []error {
	var errs []error
	for _, e := range v {
		errs = append(errs, fmt.Errorf("%s: %s", e.Field, e.Message))
	}

	return errs
}

func (v ValidationErrors) Error() error {
	var err error
	for _, e := range v {
		err = errors.Join(err, fmt.Errorf("%s: %s", e.Field, e.Message))
	}

	return err
}

// BindAndValidate binds and validates request params.
func (c *Context) BindAndValidate(req any) ValidationErrors {
	// 1. BindUri binds the passed struct pointer using the specified binding engine.
	if err := c.ShouldBindUri(req); err != nil {
		return []ValidationErrorResponse{{
			Field:   "uri",
			Message: fmt.Sprintf("uri parameter error: %v", err),
		}}
	}

	// 2. BindQuery binds the passed struct pointer using the specified binding engine.
	if err := c.ShouldBindQuery(req); err != nil {
		return []ValidationErrorResponse{{
			Field:   "query",
			Message: fmt.Sprintf("query parameter error: %v", err),
		}}
	}

	// 3. BindJSON binds the passed struct pointer using the specified binding engine.
	if err := c.ShouldBind(req); err != nil {
		return []ValidationErrorResponse{{
			Field:   "body",
			Message: fmt.Sprintf("body parameter error: %v", err),
		}}
	}

	// 4. Validate the request params.
	var errors []ValidationErrorResponse
	if err := c.Validate.Struct(req); err != nil {
		zap.L().Error("Param validate error", zap.Error(err))
		dataType := reflect.TypeOf(req)
		if dataType.Kind() == reflect.Ptr {
			dataType = dataType.Elem()
		}

		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationErrorResponse
			field, _ := dataType.FieldByName(err.Field())

			element.Field = err.Field()
			element.Message = getValidationMessage(field, err.Tag())

			errors = append(errors, element)
		}
	}

	return errors
}

func getValidationMessage(field reflect.StructField, tag string) string {
	if msg := field.Tag.Get("msg"); msg != "" {
		return msg
	}

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field.Name)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field.Name, field.Tag.Get("min"))
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field.Name, field.Tag.Get("max"))
	case "oneof":
		return fmt.Sprintf("%s must be one of %s", field.Name, field.Tag.Get("oneof"))
	default:
		return fmt.Sprintf("Validation failed on %s for %s", field.Name, tag)
	}
}
