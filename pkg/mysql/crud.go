package mysql

import (
	"errors"
	"gorm.io/gorm"
)

func AddDTO(db *gorm.DB, dto interface{}) error {
	if db == nil {
		return errors.New("db required")
	}

	return db.Create(dto).Error
}

func DelDTO(db *gorm.DB, dto interface{}, whereQuery interface{}, whereArgs ...interface{}) error {
	if db == nil {
		return errors.New("db required")
	}

	rowsAffected := db.Delete(dto, whereQuery, whereArgs).RowsAffected
	if rowsAffected != 0 {
		return nil
	}
	return errors.New("empty record")
}

func UpdateDTO(db *gorm.DB, dto interface{}) error {
	if db == nil {
		return errors.New("db required")
	}

	// if exists, update; else insert
	return db.Save(dto).Error
}

func QueryDTO(db *gorm.DB, dto interface{}, whereQuery interface{}, whereArgs ...interface{}) error {
	if db == nil {
		return errors.New("db required")
	}

	rowsAffected := db.Where(whereQuery, whereArgs...).First(dto).RowsAffected
	if rowsAffected != 0 {
		return nil
	}
	return errors.New("empty record")
}

func UpdateOneColumn(db *gorm.DB, dbDTO interface{}, columnName string, columnValue interface{},
	whereQuery interface{}, whereArgs ...interface{}) error {
	if db == nil {
		return errors.New("db required")
	}

	rowsAffected := db.Model(dbDTO).Where(whereQuery, whereArgs...).Update(columnName, columnValue).RowsAffected
	if rowsAffected != 0 {
		return nil
	}
	return errors.New("empty record")
}
