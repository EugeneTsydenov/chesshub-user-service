package enums

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBanned    UserStatus = "banned"
	UserStatusDeleted   UserStatus = "deleted"
)

type TimeControl string

const (
	TimeControlBullet         TimeControl = "bullet"
	TimeControlBlitz          TimeControl = "blitz"
	TimeControlRapid          TimeControl = "rapid"
	TimeControlClassical      TimeControl = "classical"
	TimeControlCorrespondence TimeControl = "correspondence"
)

type ChangeReason string

const (
	ChangeReasonGame        ChangeReason = "game"
	ChangeReasonAdjustment  ChangeReason = "adjustment"
	ChangeReasonSeasonReset ChangeReason = "season_reset"
	ChangeReasonPenalty     ChangeReason = "penalty"
)
