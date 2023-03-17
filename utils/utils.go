package utils

import (
	"os"

	"github.com/pkg/errors"
)

func WriteTextFile(filePath string, rows []string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(err, "could not create file in path "+filePath)
	}
	defer f.Close()

	for _, row := range rows {
		_, err = f.WriteString(row + "\n")
		if err != nil {
			return errors.Wrap(err, "could not write text to file "+filePath)
		}
	}
	return nil
}