package database

import "gorm.io/gorm"

// GORM/real implementation of DB operations
type GormDB struct {
	DB *gorm.DB
}

func (g *GormDB) Create(value interface{}) error {
	return g.DB.Create(value).Error
}

func (g *GormDB) Where(query interface{}, args ...interface{}) Database {
	return &GormDB{DB: g.DB.Where(query, args...)}
}

func (g *GormDB) First(out interface{}, where ...interface{}) error {
	return g.DB.First(out, where...).Error
}

func (g *GormDB) Save(value interface{}) error {
	return g.DB.Save(value).Error
}

func (g *GormDB) Model(value interface{}) Database {
	g.DB = g.DB.Model(value)
	return g
}

func (g *GormDB) Take(out interface{}, where ...interface{}) error {
	return g.DB.Take(out, where...).Error
}

func (g *GormDB) Delete(value interface{}, where ...interface{}) error {
	return g.DB.Delete(value, where...).Error
}

func (g *GormDB) Find(out interface{}, where ...interface{}) error {
	return g.DB.Find(out, where...).Error
}

func (g *GormDB) Updates(values interface{}) error {
	return g.DB.Updates(values).Error
}
