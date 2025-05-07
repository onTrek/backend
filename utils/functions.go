package utils

import (
	"fmt"
	"os"
)

func DeleteFile(path any) error {
	err := os.Remove(path.(string))
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
