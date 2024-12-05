package interfaces

import "gorm.io/gorm"

type RelationLoader interface {
	LoadRelations(db *gorm.DB) *gorm.DB
}
