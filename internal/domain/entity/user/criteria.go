package user

import "time"

type Criteria struct {
	Tag              *string
	Email            *string
	PublicName       *string
	LastActiveAfter  *time.Time
	LastActiveBefore *time.Time
	UpdatedAfter     *time.Time
	UpdatedBefore    *time.Time
	CreatedAfter     *time.Time
	CreatedBefore    *time.Time
}
