package models

import "errors"

var ErrUnsupportedVersion = errors.New("unsupported version")
var ErrUnsupportedProvider = errors.New("unsupported provider")
var ErrNoFields = errors.New("no fields")
var ErrNoLogFile = errors.New("no log file")
var ErrNoStorages = errors.New("no storages")
var ErrNoBindAddress = errors.New("no bind address")
var ErrUnsupportedDocvalueFields = errors.New("unsupported docvalue fields")
var ErrNoClickhouseDSN = errors.New("no clickhouse dsn")
