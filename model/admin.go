package model

import (
    "crypto/sha256"
    "database/sql"
    "encoding/hex"
    "errors"

    "github.com/google/uuid"
    "github.com/michaelrk02/rida-api/resource"
    "gorm.io/gorm"
)

const (
    ADMIN_VALIDATE_DETAILS uint     = 1 << 0
    ADMIN_VALIDATE_PASSWORD         = 1 << 1
    ADMIN_VALIDATE_FAKULTAS         = 1 << 2

    ADMIN_VALIDATE_ALL              = ADMIN_VALIDATE_DETAILS | ADMIN_VALIDATE_PASSWORD | ADMIN_VALIDATE_FAKULTAS
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

func (a *Admin) Validate(flags uint) error {
    if flags & ADMIN_VALIDATE_DETAILS != 0 {
        if a.Nama == "" {
            return errors.New("Nama tidak boleh kosong")
        }

        if a.Email == "" {
            return errors.New("E-mail tidak boleh kosong")
        }
    }

    if flags & ADMIN_VALIDATE_PASSWORD != 0 {
        if a.Password == "" {
            return errors.New("Password tidak boleh kosong")
        }

        if len(a.Password) < 8 {
            return errors.New("Password harus minimal 8 karakter")
        }
    }

    if flags & ADMIN_VALIDATE_FAKULTAS != 0 {
        if !a.FakultasID.Valid {
            return errors.New("Fakultas tidak boleh kosong")
        }
    }

    return nil
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
    a.Password = req.Password

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
