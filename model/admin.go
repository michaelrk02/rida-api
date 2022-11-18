package model

import (
    "crypto/sha256"
    "database/sql"
    "encoding/hex"

    "github.com/google/uuid"
    "github.com/michaelrk02/rida-api/resource"
    "gorm.io/gorm"
)

type Admin struct {
    ID string `gorm:"primaryKey"`
    Nama string
    Email string
    Password string
    FakultasID sql.NullString

    Fakultas *Fakultas
}

func (Admin) TableName() string {
    return "admin"
}

func (a *Admin) BeforeCreate(tx *gorm.DB) error {
    a.AssignID(uuid.New().String())

    return nil
}

func (a *Admin) AssignID(id string) *Admin {
    a.ID = id

    return a
}

func (a *Admin) HashPassword() *Admin {
    passwordHash := sha256.New()
    passwordHash.Write([]byte(a.Password))
    a.Password = hex.EncodeToString(passwordHash.Sum(nil))

    return a
}

func (a *Admin) FromRequest(req *resource.AdminRequest) *Admin {
    a.Nama = req.Nama
    a.Email = req.Email

    if req.Password != nil {
        a.Password = *req.Password
    }

    a.FakultasID.Scan(req.FakultasID)

    return a
}

func (a Admin) ToResponse() resource.AdminResponse {
    resp := resource.AdminResponse{
        ID: a.ID,
        Nama: a.Nama,
        Email: a.Email,
    }

    if a.FakultasID.Valid {
        resp.FakultasID = &a.FakultasID.String
        resp.FakultasNama = &a.Fakultas.Nama
    }

    return resp
}

func (a Admin) ToResponseOnCreate() resource.AdminResponseOnCreate {
    return resource.AdminResponseOnCreate{
        ID: a.ID,
    }
}
