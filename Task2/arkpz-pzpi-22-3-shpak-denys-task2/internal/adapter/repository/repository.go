package repository

import (
	"context"
	"wayra/internal/core/domain/interfaces"

	"gorm.io/gorm"
)

type GenericRepository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{db: db}
}

func (r *GenericRepository[T]) Add(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		return err
	}

	if loader, ok := any(entity).(interfaces.RelationLoader); ok {
		return loader.LoadRelations(r.db.WithContext(ctx)).First(entity).Error
	}

	return nil
}

func (r *GenericRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	query := r.db.WithContext(ctx).Model(&entity).Where("id = ?", id)

	if loader, ok := any(&entity).(interfaces.RelationLoader); ok {
		query = loader.LoadRelations(query)
	}

	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *GenericRepository[T]) Where(ctx context.Context, params *T) ([]T, error) {
	var entities []T
	query := r.db.WithContext(ctx).Where(params)

	err := query.Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *GenericRepository[T]) Update(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Model(entity).Updates(entity).Error
	if err != nil {
		return err
	}
	if loader, ok := any(entity).(interfaces.RelationLoader); ok {
		return loader.LoadRelations(r.db.WithContext(ctx)).First(entity).Error
	}
	return nil
}

func (r *GenericRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GenericRepository[T]) SkipTake(ctx context.Context, skip int, take int) (*[]T, error) {
	var entities []T
	query := r.db.WithContext(ctx).Offset(skip).Limit(take)

	err := query.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

func (r *GenericRepository[T]) CountWhere(ctx context.Context, params *T) int64 {
	var count int64
	query := r.db.WithContext(ctx).Model(new(T)).Where(params)

	if loader, ok := any(params).(interfaces.RelationLoader); ok {
		query = loader.LoadRelations(query)
	}

	query.Count(&count)
	return count
}
