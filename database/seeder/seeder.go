package seeder

import (
    "gorm.io/gorm"
)

type Seeder struct {
    DB *gorm.DB
}
