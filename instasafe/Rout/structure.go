package Rout

import "time"

type Request struct {
	Amount    float64   `validate:"required"`
	Timestamp time.Time `validate:"required"`
}

type Stats struct {
	Sum   float64
	Avg   float64
	Max   float64
	Min   float64
	Count int
}
