package grpc

import (
	"fmt"
	"log/slog"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/pratchaya-maneechot/service-exchange/libs/bus/command"
	"github.com/pratchaya-maneechot/service-exchange/libs/bus/query"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandlerOption struct {
	Command   command.CommandBus
	Query     query.QueryBus
	Logger    *slog.Logger
	Validator *validator.Validate
}

func NewGrpcHandlerOption(cb command.CommandBus, qb query.QueryBus, log *slog.Logger, v *validator.Validate) GrpcHandlerOption {
	return GrpcHandlerOption{
		Command:   cb,
		Query:     qb,
		Logger:    log,
		Validator: v,
	}
}

func ProvideValidator() *validator.Validate {
	v := validator.New()

	v.RegisterValidation("password_strength", validatePasswordStrength)
	v.RegisterValidation("timezone", validateTimezone)

	return v
}

func (h GrpcHandlerOption) ValidationErrors(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	st := status.New(codes.InvalidArgument, "Request validation failed")

	fieldViolations := make([]*epb.BadRequest_FieldViolation, 0, len(validationErrors))

	for _, fieldErr := range validationErrors {
		fieldViolations = append(fieldViolations, &epb.BadRequest_FieldViolation{
			Field:       getFieldName(fieldErr),
			Description: getValidationMessage(fieldErr),
		})
	}

	withDetails, detailErr := st.WithDetails(&epb.BadRequest{
		FieldViolations: fieldViolations,
	})
	if detailErr != nil {
		return st.Err()
	}

	return withDetails.Err()
}

func getFieldName(fieldErr validator.FieldError) string {
	// Try to get JSON tag name first
	if jsonTag := fieldErr.Tag(); jsonTag != "" {
		// You might want to parse the actual JSON tag from the struct field
		// For now, return the field name as is
	}
	return fieldErr.Field()
}

func getValidationMessage(fieldErr validator.FieldError) string {
	field := fieldErr.Field()
	tag := fieldErr.Tag()
	value := fieldErr.Value()
	param := fieldErr.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("Field '%s' is required", field)
	case "email":
		return fmt.Sprintf("Field '%s' must be a valid email address", field)
	case "min":
		return fmt.Sprintf("Field '%s' must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("Field '%s' must be at most %s characters long", field, param)
	case "len":
		return fmt.Sprintf("Field '%s' must be exactly %s characters long", field, param)
	case "gt":
		return fmt.Sprintf("Field '%s' must be greater than %s", field, param)
	case "gte":
		return fmt.Sprintf("Field '%s' must be greater than or equal to %s", field, param)
	case "lt":
		return fmt.Sprintf("Field '%s' must be less than %s", field, param)
	case "lte":
		return fmt.Sprintf("Field '%s' must be less than or equal to %s", field, param)
	case "oneof":
		return fmt.Sprintf("Field '%s' must be one of: %s", field, param)
	case "uuid":
		return fmt.Sprintf("Field '%s' must be a valid UUID", field)
	case "url":
		return fmt.Sprintf("Field '%s' must be a valid URL", field)
	case "alpha":
		return fmt.Sprintf("Field '%s' must contain only alphabetic characters", field)
	case "alphanum":
		return fmt.Sprintf("Field '%s' must contain only alphanumeric characters", field)
	case "numeric":
		return fmt.Sprintf("Field '%s' must be numeric", field)
	default:
		return fmt.Sprintf("Field '%s' failed validation '%s' with value '%v'", field, tag, value)
	}
}

func validatePasswordStrength(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasNumber {
		count++
	}
	if hasSpecial {
		count++
	}

	return count >= 3
}

func validateTimezone(fl validator.FieldLevel) bool {
	timezone := fl.Field().String()
	_, err := time.LoadLocation(timezone)
	return err == nil
}
