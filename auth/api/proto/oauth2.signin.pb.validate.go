// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: oauth2.signin.proto

package proto

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on SignInRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SignInRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SignInRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SignInRequestMultiError, or
// nil if none found.
func (m *SignInRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *SignInRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetUsername()); l < 4 || l > 32 {
		err := SignInRequestValidationError{
			field:  "Username",
			reason: "value length must be between 4 and 32 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetPassword()); l < 8 || l > 64 {
		err := SignInRequestValidationError{
			field:  "Password",
			reason: "value length must be between 8 and 64 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return SignInRequestMultiError(errors)
	}

	return nil
}

// SignInRequestMultiError is an error wrapping multiple validation errors
// returned by SignInRequest.ValidateAll() if the designated constraints
// aren't met.
type SignInRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SignInRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SignInRequestMultiError) AllErrors() []error { return m }

// SignInRequestValidationError is the validation error returned by
// SignInRequest.Validate if the designated constraints aren't met.
type SignInRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SignInRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SignInRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SignInRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SignInRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SignInRequestValidationError) ErrorName() string { return "SignInRequestValidationError" }

// Error satisfies the builtin error interface
func (e SignInRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSignInRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SignInRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SignInRequestValidationError{}

// Validate checks the field values on SignInResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SignInResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SignInResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SignInResponseMultiError,
// or nil if none found.
func (m *SignInResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *SignInResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccessToken

	// no validation rules for RefreshToken

	if len(errors) > 0 {
		return SignInResponseMultiError(errors)
	}

	return nil
}

// SignInResponseMultiError is an error wrapping multiple validation errors
// returned by SignInResponse.ValidateAll() if the designated constraints
// aren't met.
type SignInResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SignInResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SignInResponseMultiError) AllErrors() []error { return m }

// SignInResponseValidationError is the validation error returned by
// SignInResponse.Validate if the designated constraints aren't met.
type SignInResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SignInResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SignInResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SignInResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SignInResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SignInResponseValidationError) ErrorName() string { return "SignInResponseValidationError" }

// Error satisfies the builtin error interface
func (e SignInResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSignInResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SignInResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SignInResponseValidationError{}
