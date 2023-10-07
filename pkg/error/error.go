package error

import "errors"

var (
	TableDoesNotExist  = errors.New("Table does not exist")
	ColumnDoesNotExist = errors.New("Column does not exist")
	InvalidSelectItem  = errors.New("Select item is not valid")
	InvalidDatatype    = errors.New("Invalid datatype")
	MissingValues      = errors.New("Missing values")
)
