package config

import "errors"

type OutputModeType string

const (
	File   OutputModeType = "file"
	Stdout OutputModeType = "stdout"
)

func (o *OutputModeType) String() string {
	return string(*o)
}

func (o *OutputModeType) Set(v string) error {
	switch v {
	case "file", "stdout":
		*o = OutputModeType(v)
		return nil
	}

	return errors.New(`must be one of "file" or "stdout"`)
}

func (o *OutputModeType) Type() string {
	return "OutputModeType"
}
