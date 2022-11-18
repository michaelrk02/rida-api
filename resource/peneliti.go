package resource

type PenelitiRequest struct {
    ID string `json:"id"`
    Nidn string `json:"nidn"`
    Nama string `json:"nama"`
    JenisKelamin string `json:"jenis_kelamin"`
    ScopusAuthorID string `json:"scopus_author_id"`
    GscholarAuthorID string `json:"gscholar_author_id"`
    FakultasID *string `json:"fakultas_id"`
}

type PenelitiResponse struct {
    ID string `json:"id"`
    Nidn string `json:"nidn"`
    Nama string `json:"nama"`
    JenisKelamin string `json:"jenis_kelamin"`
    ScopusAuthorID string `json:"scopus_author_id"`
    GscholarAuthorID string `json:"gscholar_author_id"`
    FakultasID string `json:"fakultas_id"`
    FakultasNama string `json:"fakultas_nama"`
    DiciptakanOlehID string `json:"diciptakan_oleh_id"`
    DiciptakanOlehNama string `json:"diciptakan_oleh_nama"`
}

type PenelitiResponseCollection struct {
    Population int `json:"population"`
    Display int `json:"display"`
    Page int `json:"page"`
    MaxPage int `json:"max_page"`
    Data []PenelitiResponse `json:"data"`
}

type PenelitiResponseOnCreate struct {
    ID string `json:"id"`
}
