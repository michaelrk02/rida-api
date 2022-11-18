package service

import (
    "errors"
    "fmt"
    "net/http"
    "strconv"
    "strings"

    "gorm.io/gorm"
)

type DataSource struct {
    Display int
    Page int
    Search string
    SortMode string
    SortColumn string

    MaxDisplay int
    SearchColumns []string
    SortColumns []string

    TotalItems int
    MaxPage int
}

func NewDataSource(maxDisplay int, searchColumns []string, sortColumns []string, sortColumn string) *DataSource {
    ds := new(DataSource)

    ds.Display = 15
    ds.Page = 1
    ds.Search = ""
    ds.SortMode = "asc"
    ds.SortColumn = sortColumn

    ds.MaxDisplay = maxDisplay
    ds.SearchColumns = searchColumns
    ds.SortColumns = sortColumns

    return ds
}

func (ds *DataSource) FromRequest(r *http.Request, prefix string) *DataSource {
    var err error
    var buf int64

    if r.URL.Query().Has(fmt.Sprintf("%sdisplay", prefix)) {
        buf, err = strconv.ParseInt(r.URL.Query().Get(fmt.Sprintf("%sdisplay", prefix)), 10, 64)
        if err == nil {
            ds.Display = int(buf)
        }
    }

    if r.URL.Query().Has(fmt.Sprintf("%spage", prefix)) {
        buf, err = strconv.ParseInt(r.URL.Query().Get(fmt.Sprintf("%spage", prefix)), 10, 64)
        if err == nil {
            ds.Page = int(buf)
        }
    }

    if r.URL.Query().Has(fmt.Sprintf("%ssearch", prefix)) {
        ds.Search = r.URL.Query().Get(fmt.Sprintf("%ssearch", prefix))
    }

    if r.URL.Query().Has(fmt.Sprintf("%ssort_mode", prefix)) {
        ds.SortMode = r.URL.Query().Get(fmt.Sprintf("%ssort_mode", prefix))
    }

    if r.URL.Query().Has(fmt.Sprintf("%ssort_column", prefix)) {
        ds.SortColumn = r.URL.Query().Get(fmt.Sprintf("%ssort_column", prefix))
    }

    return ds
}

func (ds *DataSource) Populate(totalItems int) {
    ds.TotalItems = totalItems

    if totalItems > 0 {
        ds.MaxPage = totalItems / ds.Display
        if totalItems % ds.Display > 0 {
            ds.MaxPage++
        }
    } else {
        ds.MaxPage = 1
    }
}

func (ds *DataSource) Validate() error {
    var ok bool

    if ds.Display <= 0 || ds.Display > ds.MaxDisplay {
        return errors.New("Invalid display length")
    }

    if ds.Page <= 0 || ds.Page > ds.MaxPage {
        return errors.New("Invalid page number")
    }

    if ds.SortMode != "asc" && ds.SortMode != "desc" {
        return errors.New("Invalid sort mode")
    }

    ok = false
    for i := range ds.SortColumns {
        if ds.SortColumn == ds.SortColumns[i] {
            ok = true
            break
        }
    }
    if !ok {
        return errors.New("Invalid sort column")
    }

    return nil
}

func (ds *DataSource) SearchClause() string {
    clauses := make([]string, len(ds.SearchColumns))
    for i := range ds.SearchColumns {
        clauses[i] = fmt.Sprintf("%s LIKE ?", ds.escapeColumn(ds.SearchColumns[i]))
    }
    return strings.Join(clauses, " OR ")
}

func (ds *DataSource) SearchArg() []interface{} {
    args := make([]interface{}, len(ds.SearchColumns))
    for i := range ds.SearchColumns {
        args[i] = fmt.Sprintf("%%%s%%", ds.escapeLike(ds.Search))
    }
    return args
}

func (ds *DataSource) SortClause() string {
    return fmt.Sprintf("%s %s", ds.escapeColumn(ds.SortColumn), ds.SortMode)
}

func (ds *DataSource) Limit() int {
    return ds.Display
}

func (ds *DataSource) Offset() int {
    return (ds.Page - 1) * ds.Display
}

func (ds *DataSource) EnumerationScope(db *gorm.DB) *gorm.DB {
    return db.Where(ds.SearchClause(), ds.SearchArg()...)
}

func (ds *DataSource) PopulationScope(db *gorm.DB) *gorm.DB {
    return db.Where(ds.SearchClause(), ds.SearchArg()...).Order(ds.SortClause()).Limit(ds.Limit()).Offset(ds.Offset())
}

func (ds *DataSource) escapeLike(pattern string) string {
    clean := pattern

    clean = strings.ReplaceAll(clean, "\\", "\\\\")
    clean = strings.ReplaceAll(clean, "%", "\\%")
    clean = strings.ReplaceAll(clean, "_", "\\_")

    return clean
}

func (ds *DataSource) escapeColumn(name string) string {
    parts := strings.Split(name, ".")
    for i := range parts {
        parts[i] = fmt.Sprintf("`%s`", parts[i])
    }
    return strings.Join(parts, ".")
}
