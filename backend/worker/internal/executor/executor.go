package executor

import (
	"errors"
	"time"
)

func Execute(jobType string) error {
	switch jobType {
	case "email":
		time.Sleep(2 * time.Second)
		return nil
	case "fail":
		time.Sleep(1 * time.Second)
		return errors.New("simulated job failure")
	default:
		time.Sleep(1 * time.Second)
		return nil
	}
}
