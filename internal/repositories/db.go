package repositories

import "gorm.io/gorm"

type DB interface {
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Model(value interface{}) *gorm.DB
	Update(column string, value interface{}) *gorm.DB
	Count(count *int64) *gorm.DB
	Pluck(column string, dest interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	Error() error
}
