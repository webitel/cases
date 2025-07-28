package model

type CaseCommunication struct {
	Id                int64          `json:"id" db:"id"`
	Ver               int32          `json:"ver" db:"ver"`
	CommunicationType *GeneralLookup `json:"communication_type" db:"communication_type"`
	CommunicationId   string         `json:"communication_id" db:"communication_id"`
}

// GeneralLookup type should be defined elsewhere in the model package.
