package util

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// error struct
type ApiError struct {
	Param string
	Message string
}

// custom validate error msg
func CustomMsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		msg := fmt.Sprintf("%s is required", fe.Field())
		return msg
  case "email": 
		return "Invalid email address"
	case "min":
		if fe.Field() == "Password" {
			msg := fmt.Sprintf("%s must be more than 6", fe.Field())
		return msg
		}
		msg := fmt.Sprintf("%s must be more than 3", fe.Field())
		return msg
		
	case "max":
		msg := fmt.Sprintf("%s must be less than 20", fe.Field())
		return msg
		
	}
	return fe.Error()
}

// return err array
func HandlerErrorMsg(err error) []ApiError{
	
	var ve validator.ValidationErrors
	if errors.As(err, &ve){
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			out[i]= ApiError{fe.Field(), CustomMsgForTag(fe)}
		}
		return out
	}
	return nil
}