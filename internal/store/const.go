package store

import (
	errors "github.com/webitel/cases/internal/errors"
	"google.golang.org/grpc/codes"
)

// Case communication types
const (
	CommunicationChat  = "Messaging"
	CommunicationCall  = "Phone"
	CommunicationEmail = "Email"
)

// error types
var (
	ErrInternal = errors.Internal("internal server error")

	ErrNoRows              = errors.NotFound("entity does not exists or you do not have enough permissions to perform the operation")
	ErrUniqueViolation     = errors.New("invalid input: entity already exists", errors.WithCode(codes.AlreadyExists))
	ErrForeignKeyViolation = errors.Aborted("invalid input: violates foreign key constraint")
	ErrCheckViolation      = errors.Aborted("invalid input: violates check constraint")
	ErrNotNullViolation    = errors.Aborted("invalid input: violates not null constraint: column can not be null")
	ErrEntityConflict      = errors.Aborted("invalid input: found more then one requested entity")
)
