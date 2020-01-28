package handler

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"

	"../../../productpage"
)

// MenuHandler handles menu related requests
type MenuHandler struct {
	tmpl        *template.Template
	productSrv productpage.ItemRepository
}

// NewMenuHandler initializes and returns new MenuHandler
func NewMenuHandler(T *template.Template, CS productpage.ItemService) *MenuHandler {
	return &MenuHandler{tmpl: T, productSrv: CS}
}

// Index handles request on route /
func (mh *MenuHandler) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	products, err := mh.productSrv.Items()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/html")
	err2 := mh.tmpl.ExecuteTemplate(w, "index.layout", products)
	if err2 != nil {
		panic(err2)
	}
}

// About handles requests on route /about
func (mh *MenuHandler) About(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	mh.tmpl.ExecuteTemplate(w, "about.layout", nil)
}

// Menu handle request on route /menu
func (mh *MenuHandler) Menu(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	mh.tmpl.ExecuteTemplate(w, "menu.layout", nil)
}

// Contact handle request on route /Contact
func (mh *MenuHandler) Contact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	mh.tmpl.ExecuteTemplate(w, "contact.layout", nil)
}

// Admin handle request on route /admin
func (mh *MenuHandler) Admin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	mh.tmpl.ExecuteTemplate(w, "admin.index.layout", nil)
}

func (mh *MenuHandler) RegistPage(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	mh.tmpl.ExecuteTemplate(w, "Registrationform.html", nil)
}

func (mh *MenuHandler) RegistPage2(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	mh.tmpl.ExecuteTemplate(w, "Registrationformpart2.html", nil)
}


func (mh *MenuHandler) LoginPage(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	mh.tmpl.ExecuteTemplate(w, "login.html", nil)
}