//go:build !windows

package main

import "errors"

var errUnsupported = errors.New("windows service not supported")

func isWindowsService() (bool, error)                                 { return false, errUnsupported }
func runService(name, host, port, dsn string, origins []string) error { return errUnsupported }
func installService(name, desc string) error                          { return errUnsupported }
func removeService(name string) error                                 { return errUnsupported }
