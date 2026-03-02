package database

import (
	"github.com/ladecadence/MapaLabs/pkg/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

type SQLite struct {
	db *gorm.DB
}

func (s *SQLite) Open(fileName string) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	s.db = database
	return s.db, nil
}

func (s *SQLite) Init() error {
	err := s.db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	err = s.db.AutoMigrate(&models.Lab{})
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLite) UpsertUser(u models.User) error {
	result := s.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&u)
	return result.Error
}

func (s *SQLite) DeleteUser(models.User) error {
	// TODO
	return nil
}

func (s *SQLite) GetUsers() ([]models.User, error) {
	var users []models.User
	result := s.db.Find(&users)
	return users, result.Error
}

func (s *SQLite) GetUser(name string) (models.User, error) {
	var user models.User
	result := s.db.Where("name=?", name).First(&user)
	return user, result.Error
}

func (s *SQLite) UpsertLab(lab models.Lab) error {
	result := s.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&lab)
	return result.Error
}

func (s *SQLite) GetLabs() ([]models.Lab, error) {
	var labs []models.Lab
	result := s.db.Find(&labs)
	return labs, result.Error
}

func (s *SQLite) GetLab(id int) (models.Lab, error) {
	var lab models.Lab
	result := s.db.Where("id=?", id).First(&lab)
	return lab, result.Error
}
