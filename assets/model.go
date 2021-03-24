package assets

import (
	"database/sql/driver"
	"fmt"
	"go-xops/pkg/common"
	"time"

	"gorm.io/gorm"
)

// 由于gorm提供的base model没有json tag, 使用自定义
type Model struct {
	gorm.Model
	Id        uint      `gorm:"primary_key;comment:'自增编号'" json:"id"`
	CreatedAt LocalTime `gorm:"column:created_at;comment:'创建时间'" json:"created_at"`
	UpdatedAt LocalTime `gorm:"column:updated_at;comment:'更新时间'" json:"updated_at"`
	DeletedAt LocalTime `gorm:"column:deleted_at;null;comment:'软删除'" sql:"index" json:"deleted_at"`
}

// 表名设置
func (Model) TableName(name string) string {
	// 添加表前缀
	return fmt.Sprintf("%s_%s", common.Conf.Mysql.TablePrefix, name)
}

// 本地时间
type LocalTime struct {
	time.Time
}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	// ""空值不进行解析
	if len(data) == 2 {
		*t = LocalTime{Time: time.Time{}}
		return
	}

	// 指定解析的格式
	now, err := time.Parse(`"`+common.SecLocalTimeFormat+`"`, string(data))
	*t = LocalTime{Time: now}
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.Format(common.SecLocalTimeFormat))
	return []byte(output), nil
}

// gorm 写入 mysql 时调用
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// gorm 检出 mysql 时调用
func (t *LocalTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to LocalTime", v)
}

// 用于 fmt.Println 和后续验证场景
func (t LocalTime) String() string {
	return t.Format(common.SecLocalTimeFormat)
}

// 只需要日期
func (t LocalTime) DateString() string {
	return t.Format(common.DateLocalTimeFormat)
}
