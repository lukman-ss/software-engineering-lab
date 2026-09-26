package dberr

import (
	"errors"
	"fmt"
)

type SQLState string

const (
	SQLStateRestrictViolation   SQLState = "23001"
	SQLStateNotNullViolation    SQLState = "23502"
	SQLStateForeignKeyViolation SQLState = "23503"
	SQLStateUniqueViolation     SQLState = "23505"
	SQLStateCheckViolation      SQLState = "23514"
	SQLStateExclusionViolation  SQLState = "23P01"
	SQLStateSerializationFail   SQLState = "40001"
)

type ConstraintError struct {
	Code           SQLState
	ConstraintName string
	TableName      string
	Detail         string
	Err            error
}

func (e *ConstraintError) Error() string {
	if e.ConstraintName != "" {
		return fmt.Sprintf("db constraint error [%s] violated on table '%s' (SQLSTATE %s): %s", e.ConstraintName, e.TableName, e.Code, e.Detail)
	}
	return fmt.Sprintf("db error on table '%s' (SQLSTATE %s): %s", e.TableName, e.Code, e.Detail)
}

func (e *ConstraintError) Unwrap() error {
	return e.Err
}

func NewNotNullViolation(table, column, constraintName string) *ConstraintError {
	return &ConstraintError{
		Code:           SQLStateNotNullViolation,
		ConstraintName: constraintName,
		TableName:      table,
		Detail:         fmt.Sprintf("null value in column %q violates not-null constraint", column),
	}
}

func NewUniqueViolation(table, constraintName, detail string) *ConstraintError {
	return &ConstraintError{
		Code:           SQLStateUniqueViolation,
		ConstraintName: constraintName,
		TableName:      table,
		Detail:         detail,
	}
}

func NewCheckViolation(table, constraintName, detail string) *ConstraintError {
	return &ConstraintError{
		Code:           SQLStateCheckViolation,
		ConstraintName: constraintName,
		TableName:      table,
		Detail:         detail,
	}
}

func NewForeignKeyViolation(table, constraintName, detail string) *ConstraintError {
	return &ConstraintError{
		Code:           SQLStateForeignKeyViolation,
		ConstraintName: constraintName,
		TableName:      table,
		Detail:         detail,
	}
}

func IsConstraintViolation(err error, code SQLState) bool {
	var cErr *ConstraintError
	if errors.As(err, &cErr) {
		return cErr.Code == code
	}
	return false
}

func MapToDomainError(err error) error {
	var cErr *ConstraintError
	if errors.As(err, &cErr) {
		switch cErr.Code {
		case SQLStateUniqueViolation:
			return fmt.Errorf("conflict: resource with this unique attribute already exists (rule: %s)", cErr.ConstraintName)
		case SQLStateNotNullViolation:
			return fmt.Errorf("invalid input: mandatory field is missing (rule: %s)", cErr.ConstraintName)
		case SQLStateCheckViolation:
			return fmt.Errorf("validation failed: value outside permissible boundary (rule: %s)", cErr.ConstraintName)
		case SQLStateForeignKeyViolation:
			return fmt.Errorf("reference error: referenced entity does not exist (rule: %s)", cErr.ConstraintName)
		default:
			return fmt.Errorf("database integrity violation: %s", cErr.Detail)
		}
	}
	return err
}
