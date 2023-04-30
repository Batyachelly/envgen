package envgen

import "errors"

// ErrInvalidASTObjectType ast object has a type different from the type declaration.
var ErrInvalidASTObjectType = errors.New("invalid type")
