package models

type Org struct {
	OrgId   int32  `json:"org_id,omitempty" db:"org_id"`
	OrgName string `json:"org_name,omitempty" db:"org_name"`
}
