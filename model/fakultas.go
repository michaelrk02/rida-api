package model

import (
    "github.com/michaelrk02/rida-api/resource"
)

type Fakultas struct {
    ID string `gorm:"primaryKey"`
    Nama string
}

func (Fakultas) TableName() string {
    return "fakultas"
}

func (f Fakultas) ToResponse() resource.FakultasResponse {
    return resource.FakultasResponse{
        ID: f.ID,
        Nama: f.Nama,
    }
}
