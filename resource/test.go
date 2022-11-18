package resource

type TestRequest struct {
    ID *string `json:"id"`
    Foo string `json:"foo"`
    Bar string `json:"bar"`
}

type TestResponse struct {
    ID string `json:"id"`
    Foo string `json:"foo"`
    Bar string `json:"bar"`
}

type TestResponseCollection struct {
    Count int `json:"count"`
    Data []TestResponse `json:"data"`
}

type TestResponseOnCreate struct {
    ID string `json:"id"`
}
