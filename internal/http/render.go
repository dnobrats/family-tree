package httpx

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strconv"
	"time"

	"familytree/internal/domain"
)

// Renderer loads each page once (base layout + page content) and executes it.
type Renderer struct {
	pages map[string]*template.Template
}

// pageTemplates maps a logical page name to its content template file.
var pageTemplates = map[string]string{
	"tree":        "templates/public/tree.html",
	"person":      "templates/public/person.html",
	"error":       "templates/public/error.html",
	"login":       "templates/admin/login.html",
	"dashboard":   "templates/admin/dashboard.html",
	"persons":     "templates/admin/persons.html",
	"person_form": "templates/admin/person_form.html",
	"spouses":     "templates/admin/spouses.html",
	"clans":       "templates/admin/clans.html",
}

// NewRenderer parses every page template against templates/layouts/base.html.
func NewRenderer(fsys fs.FS) (*Renderer, error) {
	funcs := template.FuncMap{
		"deref":       deref,
		"genderLabel": func(g domain.Gender) string { return g.Label() },
		"year":        func() int { return time.Now().Year() },
		"inc":         func(n int) int { return n + 1 },
		"dec":         func(n int) int { return n - 1 },
	}

	pages := make(map[string]*template.Template, len(pageTemplates))
	for name, file := range pageTemplates {
		t, err := template.New("base").Funcs(funcs).ParseFS(fsys,
			"templates/layouts/base.html", file)
		if err != nil {
			return nil, fmt.Errorf("parse template %q: %w", name, err)
		}
		pages[name] = t
	}
	return &Renderer{pages: pages}, nil
}

// Render executes a page into a buffer first, so a template error never emits
// a half-written HTML page.
func (r *Renderer) Render(w http.ResponseWriter, status int, name string, data any) {
	t, ok := r.pages[name]
	if !ok {
		http.Error(w, "unknown page: "+name, http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "base", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = buf.WriteTo(w)
}

// deref is a template helper that prints an optional value safely.
func deref(v any) string {
	switch x := v.(type) {
	case *string:
		if x == nil {
			return ""
		}
		return *x
	case *int:
		if x == nil {
			return ""
		}
		return strconv.Itoa(*x)
	case *int64:
		if x == nil {
			return ""
		}
		return strconv.FormatInt(*x, 10)
	}
	return ""
}

// ---------------------------------------------------------------------------
// View models passed to templates. Field names match the template variables.
// ---------------------------------------------------------------------------

type pageBase struct {
	Title     string
	CSRFToken string
	Flash     string
	Year      int
	Admin     bool
}

type treePageData struct {
	pageBase
	Clans  []domain.Clan
	RootID int64
	MaxGen int
}

type personPageData struct {
	pageBase
	Detail *domain.PersonDetail
}

type loginPageData struct {
	pageBase
	Error string
}

type dashboardData struct {
	pageBase
	PersonCount int
	ClanCount   int
}

type personsData struct {
	pageBase
	Page *personPage
}

// personPage mirrors service.PersonPage so templates stay decoupled from it.
type personPage struct {
	Items      []domain.Person
	Total      int
	Page       int
	PageSize   int
	TotalPages int
}

type personFormData struct {
	pageBase
	PersonID int64
	Values   personFormValues
	People   []domain.PersonBrief
	Clans    []domain.Clan
	IsEdit   bool
	Error    string
}

type spousesData struct {
	pageBase
	Person  domain.Person
	Spouses []domain.Spouse
	People  []domain.PersonBrief
	Error   string
}

type clansData struct {
	pageBase
	Clans        []domain.Clan
	People       []domain.PersonBrief
	Edit         *domain.Clan
	EditID       int64
	EditParentID int64
	EditRootID   int64
	Error        string
}

type errorPageData struct {
	pageBase
	Message string
}

// personFormValues holds raw form strings so a form can be re-rendered with the
// user's input after a validation error.
type personFormValues struct {
	FullName       string
	Gender         string
	BirthYear      string
	BirthDateSolar string
	BirthDateLunar string
	IsAlive        bool
	DeathYear      string
	DeathDateSolar string
	DeathDateLunar string
	FatherID       string
	MotherID       string
	ClanID         string
	Address        string
	Phone          string
	Occupation     string
	AvatarURL      string
	GraveLocation  string
	Note           string
}
