package query

import (
	"errors"
	"fmt"
	"github.com/quiz_be/services/core/helper"
	"gorm.io/gorm"
)

type BaseQuery interface {
	GetDB() *gorm.DB
	SetDB(*gorm.DB)
	Set(*gorm.DB) BaseQuery
}

func NewBQ(db *gorm.DB) BaseQuery {
	return &baseQuery{DB: db}
}

type baseQuery struct {
	DB *gorm.DB
}

func (b *baseQuery) GetDB() *gorm.DB {
	return b.DB
}

func (b *baseQuery) SetDB(db *gorm.DB) {
	b.DB = db
}

func (b *baseQuery) Set(db *gorm.DB) BaseQuery {
	b.DB = db
	return b
}

// insert funcs

func Upsert[S any, D any](db *gorm.DB, source *S, mapper func(*S) *D) error {
	return db.Save(mapper(source)).Error
}

func Insert[S any, D any](db *gorm.DB, source *S, mapper func(*S) *D) error {
	return db.Create(mapper(source)).Error
}

// query func
func Where[I BaseQuery](bq I, statement string, args ...any) I {
	bq.SetDB(bq.GetDB().Where(statement, args...))
	return bq
}

// order func
func OrderBy[I any](bq I, field string, desc bool) I {
	db := any(bq).(BaseQuery)
	db.SetDB(db.GetDB().Order(getOrderByStr(field, desc)))
	return bq
}

// result func
func Result[S any, D any](bq BaseQuery, mapper func(s *S) *D) (*D, error) {
	var rs S
	err := bq.GetDB().First(&rs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	d := mapper(&rs)
	return d, nil
}

func ResultList[S any, D any](bq BaseQuery, mapper func(s *S) *D) ([]*D, error) {
	var rs []*S
	err := bq.GetDB().Find(&rs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return helper.MapList(rs, mapper), nil
}

func Limit[I any](q I, limit int) I {
	if limit <= 0 {
		return q
	}
	db := any(q).(BaseQuery)
	db.SetDB(db.GetDB().Limit(limit))
	return q
}

func Offset[I any](q I, offset int) I {
	db := any(q).(BaseQuery)
	db.SetDB(db.GetDB().Offset(offset))
	return q
}

func Count(bq BaseQuery) (int, error) {
	var count int64
	err := bq.GetDB().Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func Delete(bq BaseQuery) error {
	return bq.GetDB().Delete(any(bq)).Error
}

func getOrderByStr(field string, desc bool) string {
	if desc {
		return fmt.Sprintf("%s desc NULLS LAST", field)
	}
	return fmt.Sprintf("%s asc NULLS FIRST", field)
}
