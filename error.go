package main

import "errors"

// Custom errors.
var (
	ErrMissingSignature = errors.New("Signature is missing")
	ErrEmptyBody        = errors.New("Body is empty")
)
