package repositories

import "gorm.io/gorm"

// GormDB interface defines the required database operations
type GormDB interface {
	Create(value interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Joins(query string, args ...interface{}) *gorm.DB
	Model(value interface{}) *gorm.DB
	Count(value *int64) *gorm.DB
	Offset(offset int) *gorm.DB
	Limit(limit int) *gorm.DB
	Preload(query string, args ...interface{}) *gorm.DB
	Pluck(column string, value interface{}) *gorm.DB
	Update(column string, value interface{}) *gorm.DB
	Error() error
}

// Ensure *gorm.DB implements GormDB interface
type gormDB struct {
	*gorm.DB
}

func (db *gormDB) Error() error {
	return db.DB.Error
}

// WrapDB wraps a *gorm.DB to implement the GormDB interface
func WrapDB(db *gorm.DB) GormDB {
	return &gormDB{db}
}
