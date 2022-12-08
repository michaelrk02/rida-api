package resource

type AdminRequest struct {
    Nama string `json:"nama"`
    Email string `json:"email"`
    Password string `json:"password"`
    FakultasID string `json:"fakultas_id"`
}

type AdminRequestOnLogin struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type AdminRequestOnUpdatePassword struct {
    OldPassword string `json:"old_password"`
    NewPassword string `json:"new_password"`
}

type AdminResponse struct {
    ID string `json:"id"`
    Nama string `json:"nama"`
    Email string `json:"email"`
    FakultasID *string `json:"fakultas_id"`
    FakultasNama *string `json:"fakultas_nama"`
}

type AdminResponseCollection struct {
    Population int `json:"population"`
    Display int `json:"display"`
    Page int `json:"page"`
    MaxPage int `json:"max_page"`
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
