package model

import "context"

type CreateOptions struct {
	context.Context
	Fields []string
}
