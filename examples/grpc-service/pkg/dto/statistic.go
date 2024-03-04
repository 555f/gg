package dto

import "time"

type Statistic struct {
	ProfileID int
	Sum       float64
	CreatedAt time.Time
}
