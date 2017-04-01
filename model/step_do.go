package model

import (
	"time"
)

// Do collection
type Do struct {
	// Steps in the collection
	Do ISteps

	// Time duration in which the execution of the do steps will be timed-out
	Timeout time.Duration
}
