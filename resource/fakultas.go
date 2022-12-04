package resource

type FakultasResponse struct {
    ID string `json:"id"`
    Nama string `json:"nama"`
}

type FakultasResponseCollection struct {
    Count int `json:"count"`
    Data []FakultasResponse `json:"data"`
}
