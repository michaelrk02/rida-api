package middleware

import (
    "context"
    "errors"
    "net/http"
    "os"
    "regexp"

    "github.com/golang-jwt/jwt/v4"
    "github.com/michaelrk02/rida-api/api"
    "github.com/michaelrk02/rida-api/model"
    "github.com/michaelrk02/rida-api/service"
    "gorm.io/gorm"
)

func AuthAdmin(mustSuperadmin bool) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler  {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            var err error
            var ok bool

            authz := regexp.MustCompile("Bearer (.+)").FindStringSubmatch(r.Header.Get("Authorization"))
            if authz == nil {
                api.Error{Message: api.ErrUnauthorized}.Send(w, 401, errors.New("Invalid JWT header"))
                return
            }

            token, err := jwt.Parse(authz[1], JwtKeyFunc)
            if err != nil {
                api.Error{Message: api.ErrUnauthorized}.Send(w, 401, err)
                return
            }

            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                api.Error{Message: api.ErrUnauthorized}.Send(w, 401, errors.New("Invalid JWT token"))
                return
            }

            adminID, ok := claims["sub"]
            if !ok {
                api.Error{Message: api.ErrUnauthorized}.Send(w, 401, errors.New("Invalid authentication guard"))
                return
            }

            app := r.Context().Value("app").(*service.Application)

            var admin model.Admin
            err = app.DB.Where("id = ?", adminID).First(&admin).Error
            if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                    api.Error{Message: api.ErrUnauthorized}.Send(w, 401, errors.New("Admin not found"))
                } else {
                    api.Error{Message: api.ErrServerSide}.Send(w, 500, err)
                }
                return
            }

            if mustSuperadmin {
                if admin.FakultasID.Valid {
                    api.Error{Message: api.ErrUnauthorized}.Send(w, 401, errors.New("Invalid access rights"))
                    return
                }
            }

            ctx := r.Context()
            ctx = context.WithValue(ctx, "auth_mode", "admin")
            ctx = context.WithValue(ctx, "auth_id", adminID)
            ctx = context.WithValue(ctx, "auth_is_superadmin", !admin.FakultasID.Valid)
            ctx = context.WithValue(ctx, "auth_entity", admin)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func JwtKeyFunc(token *jwt.Token) (interface{}, error) {
    var ok bool

    _, ok = token.Method.(*jwt.SigningMethodHMAC)
    if !ok {
        return nil, errors.New("Invalid JWT signing method")
    }

    return []byte(os.Getenv("APP_KEY")), nil
}
