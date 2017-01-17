package handler
import (
	"errors"
)

var (
	InvalidMagicError = errors.New("Invalid Magic Error")

	NotEnoughDataError = errors.New("Not Enough Data Error")
)
