package ydapp

import (
	"errors"
)

func Error(desc string, err error) error {
	if err == nil {
		return errors.New(desc)
	}
	return errors.New(desc + ": " + err.Error())
}
