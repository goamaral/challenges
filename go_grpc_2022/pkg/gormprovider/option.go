package gormprovider

import "gorm.io/gorm"

type Option interface {
	Apply(db *gorm.DB) *gorm.DB
}

func ApplyOptions(qry *gorm.DB, opts ...Option) *gorm.DB {
	for _, opt := range opts {
		qry = opt.Apply(qry)
	}
	return qry
}
