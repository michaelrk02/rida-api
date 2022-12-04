package resource

type AdminRequest struct {
    Nama string `json:"nama"`
    Email string `json:"email"`
    Password *string `json:"password"`
    FakultasID string `json:"fakultas_id"`
}

type AdminRequestOnLogin struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type AdminResponse struct {
    ID string `json:"id"`
    Nama string `json:"nama"`
    Email string `json:"email"`
    FakultasID *string `json:"fakultas_id"`
    FakultasNama *string `json:"fakultas_nama"`
}

type AdminResponseCollection struct {
    Count int `json:"count"`
    Data []AdminResponse `json:"data"`
}

type AdminResponseOnCreate struct {
    ID string `json:"id"`
}

type AdminResponseOnLogin struct {
    ID string `json:"id"`
    Role string `json:"role"`
    Token string `json:"token"`
}
