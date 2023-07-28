package config

import (
	"errors"
	"danfwing.com/m/zhansheng/utils/sqlWithTransaction"
	"github.com/jinzhu/gorm"
	"time"
)

type SystemConfig struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	Key          string    `json:"key" gorm:"type:varchar(64);not null;unique"`                //key
	Name         string    `json:"name" gorm:"type:varchar(255)"`                              //名称
	Value        string    `json:"value" gorm:"type:varchar(255)"`                             //值
	ValueExplain string    `json:"value_explain" gorm:"type:varchar(255);not null;default:''"` //值说明
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:datetime"`
}

type SystemConfigCreate struct {
	Name         string `json:"name" validate:"required" label:"名称"`
	Key          string `json:"key" validate:"required" label:"键名称"`
	Value        string `json:"value"`
	ValueExplain string `json:"value_explain"`
}

type SystemConfigUpdate struct {
	ID           uint   `json:"id" validate:"required"`
	Name         string `json:"name" validate:"required" label:"名称"`
	Value        string `json:"value" validate:"required"`
	ValueExplain string `json:"value_explain"`
}

func (c *SystemConfig) Create() (id int, err error) {
	err = sqlWithTransaction.D_Pool(func(DB **gorm.DB) (err error) {
		conn := *DB
		return conn.Create(&c).Error
	})
	if err != nil {
		return 0, err
	}
	return int(c.ID), err
}

func (c *SystemConfig) UpdateById(param map[string]interface{}) (id int, err error) {
	err = sqlWithTransaction.D_Pool(func(DB **gorm.DB) (err error) {
		conn := *DB
		return conn.Model(&c).Update(param).Error
	})
	if err != nil {
		return 0, err
	}
	return int(c.ID), err
}

func (c SystemConfig) Getlist(page, page_size int) (rs []SystemConfig, err error) {
	err = sqlWithTransaction.D_Pool(func(DB **gorm.DB) (err error) {
		conn := *DB
		return conn.Model(&SystemConfig{}).Offset((page - 1) * page_size).Limit(page_size).Order("created_at desc").Scan(&rs).Error
	})
	if err != nil {
		return nil, err
	}
	return rs, err
}

func (c SystemConfig) Exit() (err error) {
	err = sqlWithTransaction.D_Pool(func(db **gorm.DB) (err error) {
		conn := *db
		if c.Key != "" {
			conn = conn.Where("`key` = ?", c.Key)
		}

		err = conn.Select("id").First(&c).Error
		return
	})
	return
}

func (c *SystemConfig) Delete() (id int, err error) {
	err = sqlWithTransaction.D_Pool(func(DB **gorm.DB) (err error) {
		conn := *DB
		return conn.Where("id = ? ", c.ID).Delete(&SystemConfig{}).Error
	})
	if err != nil {
		return id, err
	}
	return int(c.ID), err
}

func (c SystemConfig) GetInfoByKey() (value SystemConfig, err error) {
	err = sqlWithTransaction.D_Pool(func(DB **gorm.DB) (err error) {
		conn := *DB
		return conn.Model(&SystemConfig{}).Where("key = ?", c.Key).First(&SystemConfig{}).Scan(&value).Error
	})
	if err == gorm.ErrRecordNotFound {
		return value, errors.New("key not found")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return value, err
	}
	return value, err
}

func (c SystemConfig) GetInfoById() (value SystemConfig, err error) {
	err = sqlWithTransaction.D_Pool(func(DB **gorm.DB) (err error) {
		conn := *DB
		return conn.Model(&SystemConfig{}).Where("id = ?", c.ID).First(&SystemConfig{}).Scan(&value).Error
	})
	if err == gorm.ErrRecordNotFound {
		return value, errors.New("id not found")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return value, err
	}
	return value, err
}

func (c SystemConfig) GetCount() (count int, err error) {
	err = sqlWithTransaction.D_Pool(func(DB **gorm.DB) (err error) {
		conn := *DB
		return conn.Model(&SystemConfig{}).Select("id").Count(&count).Error
	})
	if err != nil {
		return count, err
	}
	return count, err
}
