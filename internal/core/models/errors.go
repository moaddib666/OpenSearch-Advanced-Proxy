package models

import "errors"

var ErrUnsupportedVersion = errors.New("unsupported version")
var ErrUnsupportedProvider = errors.New("unsupported provider")
var ErrNoFields = errors.New("no fields")
var ErrNoLogFile = errors.New("no log file")
