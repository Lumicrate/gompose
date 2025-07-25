package postgres

import (
	"fmt"
	"github.com/Lumicrate/gompose/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
)

type PostgresAdapter struct {
	dsn string
	db  *gorm.DB
}

func New(dsn string) *PostgresAdapter {
	return &PostgresAdapter{dsn: dsn}
}

func (p *PostgresAdapter) Init() error {
	var err error
	p.db, err = gorm.Open(postgres.Open(p.dsn), &gorm.Config{})
	return err
}

func (p *PostgresAdapter) Migrate(entities []any) error {
	for _, entity := range entities {
		if err := p.db.AutoMigrate(entity); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}
	return nil
}

func (p *PostgresAdapter) Create(entity any) error {
	return p.db.Create(entity).Error
}

func (p *PostgresAdapter) Update(entity any) error {
	return p.db.Save(entity).Error
}

func (p *PostgresAdapter) Delete(id string, entity any) error {
	return p.db.Delete(entity, "id = ?", id).Error
}

func (p *PostgresAdapter) FindAll(entity any, filters map[string]any, pagination db.Pagination, sort []db.Sort) (any, error) {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	sliceType := reflect.SliceOf(entityType)

	resultValue := reflect.New(sliceType) // *([]Entity)

	tx := p.db.Model(entity)

	for key, val := range filters {
		tx = tx.Where(fmt.Sprintf("%s = ?", key), val)
	}

	for _, s := range sort {
		tx = tx.Order(fmt.Sprintf("%s %s", s.Field, s.Direction))
	}

	if pagination.Limit > 0 {
		tx = tx.Limit(pagination.Limit)
	}

	if pagination.Offset > 0 {
		tx = tx.Offset(pagination.Offset)
	}

	if err := tx.Find(resultValue.Interface()).Error; err != nil {
		return nil, err
	}

	result := resultValue.Elem().Interface()

	return result, nil
}

func (p *PostgresAdapter) FindByID(id string, entity any) (any, error) {
	err := p.db.First(entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}
