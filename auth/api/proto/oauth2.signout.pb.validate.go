// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: oauth2.signout.proto

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

// Validate checks the field values on SignOutRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SignOutRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SignOutRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SignOutRequestMultiError,
// or nil if none found.
func (m *SignOutRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *SignOutRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccessToken

	if len(errors) > 0 {
		return SignOutRequestMultiError(errors)
	}

	return nil
}

// SignOutRequestMultiError is an error wrapping multiple validation errors
// returned by SignOutRequest.ValidateAll() if the designated constraints
// aren't met.
type SignOutRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SignOutRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SignOutRequestMultiError) AllErrors() []error { return m }

// SignOutRequestValidationError is the validation error returned by
// SignOutRequest.Validate if the designated constraints aren't met.
type SignOutRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SignOutRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SignOutRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SignOutRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SignOutRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SignOutRequestValidationError) ErrorName() string { return "SignOutRequestValidationError" }

// Error satisfies the builtin error interface
func (e SignOutRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSignOutRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SignOutRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SignOutRequestValidationError{}

// Validate checks the field values on SignOutResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *SignOutResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SignOutResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SignOutResponseMultiError, or nil if none found.
func (m *SignOutResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *SignOutResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return SignOutResponseMultiError(errors)
	}

	return nil
}

// SignOutResponseMultiError is an error wrapping multiple validation errors
// returned by SignOutResponse.ValidateAll() if the designated constraints
// aren't met.
type SignOutResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SignOutResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SignOutResponseMultiError) AllErrors() []error { return m }

// SignOutResponseValidationError is the validation error returned by
// SignOutResponse.Validate if the designated constraints aren't met.
type SignOutResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SignOutResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SignOutResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SignOutResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SignOutResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SignOutResponseValidationError) ErrorName() string { return "SignOutResponseValidationError" }

// Error satisfies the builtin error interface
func (e SignOutResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSignOutResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SignOutResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SignOutResponseValidationError{}
