package entity

import "time"

type TimeFormat time.Time

func (t *TimeFormat) MarshalJSON() ([]byte, error) {
	if time.Time(*t).IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + time.Time(*t).Format("2006-01-02 15:04:05") + `"`), nil
}

type Segments struct {
	BizTag    string    `xorm:"not null pk VARCHAR(32) 'biz_tag'" json:"biz_tag" binding:"required,max=32"`
	MaxId     int64     `xorm:"BIGINT(20) 'max_id'" json:"max_id" binding:"required"`
	Step      int64     `xorm:"INT(11) 'step'" json:"step" binding:"required"`
	Remark    string    `xorm:"VARCHAR(200) 'remark'" json:"remark"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

func (s *Segments) TableName() string {
	return "segments"
}
