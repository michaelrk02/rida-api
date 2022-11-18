package model

import (
    "github.com/google/uuid"
    "github.com/michaelrk02/rida-api/resource"
    "gorm.io/gorm"
)

type Test struct {
    ID string `gorm:"primaryKey"`
    Foo string
    Bar string
}

func (Test) TableName() string {
    return "test"
}

func (t *Test) BeforeCreate(tx *gorm.DB) error {
    t.AssignID(uuid.New().String())

    return nil
}

func (t *Test) AssignID(id string) *Test {
    t.ID = id

    return t
}

func (t *Test) FromRequest(req *resource.TestRequest) *Test {
    t.Foo = req.Foo
    t.Bar = req.Bar

    return t
}

func (t Test) ToResponse() resource.TestResponse {
    return resource.TestResponse{
        ID: t.ID,
        Foo: t.Foo,
        Bar: t.Bar,
    }
}

func (t Test) ToResponseOnCreate() resource.TestResponseOnCreate {
    return resource.TestResponseOnCreate{
        ID: t.ID,
    }
}
