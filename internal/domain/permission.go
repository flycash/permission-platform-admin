package domain

type Permission struct {
	ID          int64
	BizID       int64
	Name        string
	Description string
	Resource    Resource
	Action      string
	Metadata    string
}
type Effect string

const (
	EffectAllow Effect = "allow"
	EffectDeny  Effect = "deny"
)

func (e Effect) String() string {
	return string(e)
}

func (e Effect) IsAllow() bool {
	return e.String() == "allow"
}

func (e Effect) IsDeny() bool {
	return e.String() == "deny"
}
type UserPermission struct {
	ID         int64      `json:"id,omitzero"`
	BizID      int64      `json:"bizId,omitzero"`
	UserID     int64      `json:"userId,omitzero"`
	Permission Permission `json:"permission,omitzero"`
	StartTime  int64      `json:"startTime,omitzero"` // 权限生效时间
	EndTime    int64      `json:"endTime,omitzero"`   // 权限失效时间
	Effect     Effect     `json:"effect,omitzero"`
	Ctime      int64      `json:"cTime,omitzero"`
	Utime      int64      `json:"uTime,omitzero"`
}
