package model

import (
    "github.com/google/uuid"
    "github.com/michaelrk02/rida-api/resource"
    "gorm.io/gorm"
)

type Peneliti struct {
    ID string `gorm:"primaryKey"`
    Nidn string
    Nama string
    JenisKelamin string
    ScopusAuthorID string
    GscholarAuthorID string
    FakultasID string
    DiciptakanOlehID string

    Fakultas *Fakultas
    DiciptakanOleh *Admin
}

func (Peneliti) TableName() string {
    return "peneliti"
}

func (p *Peneliti) BeforeCreate(tx *gorm.DB) error {
    p.AssignID(uuid.New().String())

    return nil
}

func (p *Peneliti) AssignID(id string) *Peneliti {
    p.ID = id

    return p
}

func (p *Peneliti) FromRequest(req *resource.PenelitiRequest) *Peneliti {
    p.Nidn = req.Nidn
    p.Nama = req.Nama
    p.JenisKelamin = req.JenisKelamin
    p.ScopusAuthorID = req.ScopusAuthorID
    p.GscholarAuthorID = req.GscholarAuthorID

    return p
}

func (p Peneliti) ToResponse() resource.PenelitiResponse {
    return resource.PenelitiResponse{
        ID: p.ID,
        Nidn: p.Nidn,
        Nama: p.Nama,
        JenisKelamin: p.JenisKelamin,
        ScopusAuthorID: p.ScopusAuthorID,
        GscholarAuthorID: p.GscholarAuthorID,
        FakultasID: p.FakultasID,
        DiciptakanOlehID: p.DiciptakanOlehID,
        FakultasNama: p.Fakultas.Nama,
        DiciptakanOlehNama: p.DiciptakanOleh.Nama,
    }
}

func (p Peneliti) ToResponseOnCreate() resource.PenelitiResponseOnCreate {
    return resource.PenelitiResponseOnCreate{
        ID: p.ID,
    }
}
