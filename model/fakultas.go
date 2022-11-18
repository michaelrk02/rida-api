package model

import (
    //"github.com/michaelrk02/rida-api/resource"
    //"gorm.io/gorm"
)

type Fakultas struct {
    ID string `gorm:"primaryKey"`
    Nama string
}

func (Fakultas) TableName() string {
    return "fakultas"
}

/*func (t Test) ToResponse() resource.TestResponse {
    return resource.TestResponse{
        ID: t.ID,
        Nama: t.Nama,
    }
}

func (t Test) ToResponseOnCreate() resource.TestResponseOnCreate {
    return resource.TestResponseOnCreate{
        ID: t.ID,
    }
}*/
