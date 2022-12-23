package main

import (
    "context"
    "fmt"
    "flag"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
    "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
    "github.com/michaelrk02/rida-api/api/v1"
    "github.com/michaelrk02/rida-api/database/seeder"
    "github.com/michaelrk02/rida-api/service"
    "github.com/michaelrk02/rida-api/service/route"
    gorm_mysql "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func main() {
    var err error

    app := service.InitApplication()

    flag.BoolVar(&app.Params.Help, "help", false, "show help")
    flag.BoolVar(&app.Params.Seed, "seed", false, "perform database seeding")
    flag.BoolVar(&app.Params.Sync, "sync", false, "perform h-index synchronization")
    flag.BoolVar(&app.Params.Daemon, "daemon", false, "run HTTP daemon process")
    flag.BoolVar(&app.Params.Panic, "panic", false, "panics when a server error occurs")
    flag.Parse()

    if app.Params.Help {
        flag.PrintDefaults()
        return
    }

    err = godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env: %s", err)
    }

    cfg := mysql.NewConfig()
    cfg.User = os.Getenv("DB_USER")
    cfg.Passwd = os.Getenv("DB_PASS")
    cfg.Addr = os.Getenv("DB_HOST")
    cfg.DBName = os.Getenv("DB_NAME")

    customLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
        SlowThreshold: time.Second,
        LogLevel: logger.Error,
        IgnoreRecordNotFoundError: true,
        Colorful: true,
    })
    app.DB, err = gorm.Open(gorm_mysql.Open(cfg.FormatDSN()), &gorm.Config{Logger: customLogger})
    if err != nil {
        log.Fatalf("Error opening database connection: %s", err)
    }

    app.RouteV1 = &v1.RouteCollection{
        App: app,
    }

    if app.Params.Seed {
        s := seeder.Seeder{DB: app.DB}

        s.RunV001000()

        fmt.Println("Database seeding completed")
    }

    if app.Params.Sync {
        fmt.Println("Synchronizing H-Index values ...")

        app.SyncHIndex()

        fmt.Println("H-Index synchronization completed")
    }

    if app.Params.Daemon {
        r := chi.NewRouter()

        r.Use(middleware.Logger)
        r.Use(middleware.Recoverer)

        r.Use(cors.Handler(cors.Options{
            AllowedOrigins: []string{"https://*", "http://*"},
            AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
            AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        }))

        r.Use(func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                ctx := r.Context()
                ctx = context.WithValue(ctx, "app", app)

                next.ServeHTTP(w, r.WithContext(ctx))
            })
        })

        route.InitRoutesV1(r, app)

        fmt.Printf("Listening on port %s ...\n", os.Getenv("APP_PORT"))
        http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), r)
    }
}
