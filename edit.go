package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type File struct {
	Name   string
	file   os.File
	Logger zerolog.Logger
}

func (f *File) Open() error {
	fsFile, err := os.Create(f.Name)
	if err != nil {
		f.Logger.Error().Err(err).Msg(fmt.Sprintf("Error creating file \"%s\".", f.Name))
		return err
	}
	defer fsFile.Close()

	// Implement editor to write to files. //

	return nil
}
