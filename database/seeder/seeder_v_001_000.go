package seeder

import (
    "crypto/rand"
    "database/sql"
    "math/big"
    "strconv"

    "github.com/go-faker/faker/v4"
    "github.com/google/uuid"
    "github.com/michaelrk02/rida-api/model"
)

func (s *Seeder) RunV001000(localData, remoteData bool) {
    var random *big.Int

    s.DB.Exec("DELETE FROM `peneliti`")
    s.DB.Exec("DELETE FROM `admin`")

    var fakultasList []model.Fakultas
    s.DB.Find(&fakultasList)

    adminList := []model.Admin{
        // Password: superadmin
        {Nama: "Superadmin", Email: "root@localhost.localdomain", Password: "186cf774c97b60a1c106ef718d10970a6a06e06bef89553d9ae65d938a886eae", FakultasID: sql.NullString{"", false}},

        // Password: somepassword
        {Nama: "Alice", Email: "alice@localhost.localdomain", Password: "42a9798b99d4afcec9995e47a1d246b98ebc96be7a732323eee39d924006ee1d", FakultasID: sql.NullString{fakultasList[0].ID, true}},
        {Nama: "Bob", Email: "bob@localhost.localdomain", Password: "42a9798b99d4afcec9995e47a1d246b98ebc96be7a732323eee39d924006ee1d", FakultasID: sql.NullString{fakultasList[1].ID, true}},
        {Nama: "Charlie", Email: "charlie@localhost.localdomain", Password: "42a9798b99d4afcec9995e47a1d246b98ebc96be7a732323eee39d924006ee1d", FakultasID: sql.NullString{fakultasList[2].ID, true}},
    }

    s.DB.Create(&adminList)

    if localData {
        for i := range fakultasList {
            random, _ = rand.Int(rand.Reader, big.NewInt(40))
            penelitiCount := int(random.Int64()) + 10
            for j := 0; j < penelitiCount; j++ {
                peneliti := model.Peneliti{
                    Nama: faker.Name(),
                    ScopusAuthorID: uuid.New().String(),
                    GscholarAuthorID: uuid.New().String(),
                    FakultasID: fakultasList[i].ID,
                }

                random, _ = rand.Int(rand.Reader, big.NewInt(9000000))
                peneliti.Nidn = strconv.Itoa(1000000 + int(random.Int64()))

                random, _ = rand.Int(rand.Reader, big.NewInt(2))
                peneliti.JenisKelamin = []string{"Laki-Laki", "Perempuan"}[random.Int64()]

                random, _ = rand.Int(rand.Reader, big.NewInt(int64(len(adminList))))
                peneliti.DiciptakanOlehID = adminList[int(random.Int64())].ID

                random, _ = rand.Int(rand.Reader, big.NewInt(26))
                peneliti.HIndex = int(random.Int64())

                s.DB.Create(&peneliti)
            }
        }
    }

    if remoteData {
        penelitiList := []model.Peneliti{
            {Nidn: "0027038005", Nama: "Haryono Setiadi", JenisKelamin: "Laki-Laki", ScopusAuthorID: "-", GscholarAuthorID: "y-2QKVYAAAAJ", FakultasID: "b62d129e-3517-4c87-806c-6b77fb3512ff", DiciptakanOlehID: adminList[0].ID, IsRemote: true},
            {Nidn: "0002038307", Nama: "Heri Prasetyo", JenisKelamin: "Laki-Laki", ScopusAuthorID: "-", GscholarAuthorID: "vAk-cQ0AAAAJ", FakultasID: "b62d129e-3517-4c87-806c-6b77fb3512ff", DiciptakanOlehID: adminList[0].ID, IsRemote: true},
            {Nidn: "0002038307", Nama: "Dewi Wisnu Wardani", JenisKelamin: "Perempuan", ScopusAuthorID: "-", GscholarAuthorID: "jqTz_58AAAAJ", FakultasID: "b62d129e-3517-4c87-806c-6b77fb3512ff", DiciptakanOlehID: adminList[0].ID, IsRemote: true},
            {Nidn: "0002038307", Nama: "Wiharto", JenisKelamin: "Laki-Laki", ScopusAuthorID: "-", GscholarAuthorID: "ctkQIV0AAAAJ", FakultasID: "b62d129e-3517-4c87-806c-6b77fb3512ff", DiciptakanOlehID: adminList[0].ID, IsRemote: true},
        }

        s.DB.Create(&penelitiList)
    }
}
