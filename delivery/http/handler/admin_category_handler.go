package handler

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"../../../entity"
	"../../../productpage"
)

// AdminCategoryHandler handles category handler admin requests
type AdminCategoryHandler struct {
	categorySrv productpage.CategoryService
}

// NewAdminCategoryHandler initializes and returns new AdminCateogryHandler
func NewAdminCategoryHandler(cs productpage.CategoryService) *AdminCategoryHandler {
	return &AdminCategoryHandler{categorySrv: cs}
}

// AdminCategories handle requests on route /admin/categories
func (ach *AdminCategoryHandler) AdminCategories(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	categories, errs := ach.categorySrv.Categories()

	if len(errs) > 0 {
		w.Header().Set("Content-type", "application/json")
		http.Error(w, http.StatusText(http.StatusSeeOther), 303)
		return
	}
	output, err := json.MarshalIndent(categories, "", "\t\t")

	if err != nil{
		w.Header().Set("Content-type", "application/json")
		http.Error(w, http.StatusText(http.StatusSeeOther), 303)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (ach *AdminCategoryHandler) AdminCategory(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	category, errs := ach.categorySrv.Category(uint(id))

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	output, err := json.MarshalIndent(category, "", "\t\t")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}
// AdminCategoriesNew hanlde requests on route /admin/categories/new
func (ach *AdminCategoryHandler) AdminCategoriesNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	l := r.ContentLength
	body := make([]byte, l)
	r.Body.Read(body)

	var ctg []entity.Category

	err := json.Unmarshal(body, &ctg)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	category, errs := ach.categorySrv.StoreCategory(&ctg[0])

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		//http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		http.Error(w, "first", http.StatusNotFound)
		return
	}

	p := fmt.Sprintf("localhost:8383/v1/categories/%d", category.ID)
	w.Header().Set("Location", p)
	w.WriteHeader(http.StatusCreated)
	return
}

// AdminCategoriesUpdate handle requests on /admin/categories/update
func (ach *AdminCategoryHandler) AdminCategoriesUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	category, errs := ach.categorySrv.Category(uint(id))

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	l := r.ContentLength

	body := make([]byte, l)
	r.Body.Read(body)

	var ctg []entity.Category
	err = json.Unmarshal(body, &ctg)
	if err != nil {
		panic(err.Error())
	}

	category, errs = ach.categorySrv.UpdateCategory(&ctg[0])

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(category, "", "\t\t")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

// AdminCategoriesDelete handle requests on route /admin/categories/delete
func (ach *AdminCategoryHandler) AdminCategoriesDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, errs := ach.categorySrv.DeleteCategory(uint(id))

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

func writeFile(mf *multipart.File, fname string) {

	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	path := filepath.Join(wd, "../", "../", "ui", "assets", "img", fname)
	image, err := os.Create(path)

	if err != nil {
		panic(err)
	}
	defer image.Close()
	io.Copy(image, *mf)
}

//AdminCategoryHandler handles category handler admin requests
//type AdminCategoryHandler struct {
//	tmpl        *template.Template
//	categorySrv productpage.CategoryService
//}
//
//// NewAdminCategoryHandler initializes and returns new AdminCateogryHandler
//func NewAdminCategoryHandler(t *template.Template, cs productpage.CategoryService) *AdminCategoryHandler {
//	return &AdminCategoryHandler{tmpl: t, categorySrv: cs}
//}
//
//// AdminCategories handle requests on route /admin/categories
//func (ach *AdminCategoryHandler) AdminCategories(w http.ResponseWriter, r *http.Request) {
//	categories, errs := ach.categorySrv.Categories()
//	if errs != nil {
//		panic(errs)
//	}
//	ach.tmpl.ExecuteTemplate(w, "admin.categ.layout", categories)
//}
//
//// AdminCategoriesNew hanlde requests on route /admin/categories/new
//func (ach *AdminCategoryHandler) AdminCategoriesNew(w http.ResponseWriter, r *http.Request) {
//
//	if r.Method == http.MethodPost {
//
//		ctg := entity.Category{}
//		ctg.Name = r.FormValue("name")
//		ctg.Description = r.FormValue("description")
//
//		mf, fh, err := r.FormFile("catimg")
//		if err != nil {
//			panic(err)
//		}
//		defer mf.Close()
//
//		ctg.Image = fh.Filename
//
//		writeFile(&mf, fh.Filename)
//
//		_, err = ach.categorySrv.StoreCategory(ctg)
//
//		if err != nil{
//			panic(err.Error())
//		}
//
//		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
//
//	} else {
//
//		ach.tmpl.ExecuteTemplate(w, "admin.categ.new.layout", nil)
//
//	}
//}
//
//// AdminCategoriesUpdate handle requests on /admin/categories/update
//func (ach *AdminCategoryHandler) AdminCategoriesUpdate(w http.ResponseWriter, r *http.Request) {
//
//	if r.Method == http.MethodGet {
//
//		idRaw := r.URL.Query().Get("id")
//		id, err := strconv.Atoi(idRaw)
//
//		if err != nil {
//			panic(err)
//		}
//
//		cat, err := ach.categorySrv.Category(uint(id))
//
//		if err != nil{
//			panic(err.Error())
//		}
//
//		ach.tmpl.ExecuteTemplate(w, "admin.categ.update.layout", cat)
//
//	} else if r.Method == http.MethodPost {
//
//		ctg := entity.Category{}
//		id, _ := strconv.Atoi(r.FormValue("id"))
//		ctg.ID = uint(id)
//		ctg.Name = r.FormValue("name")
//		ctg.Description = r.FormValue("description")
//		ctg.Image = r.FormValue("image")
//
//		mf, _, err := r.FormFile("catimg")
//
//		if err != nil {
//			panic(err)
//		}
//
//		defer mf.Close()
//
//		writeFile(&mf, ctg.Image)
//
//		err = ach.categorySrv.UpdateCategory(ctg)
//
//		if err != nil{
//			panic(err.Error())
//		}
//
//		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
//
//	} else {
//		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
//	}
//
//}
//
//// AdminCategoriesDelete handle requests on route /admin/categories/delete
//func (ach *AdminCategoryHandler) AdminCategoriesDelete(w http.ResponseWriter, r *http.Request) {
//
//	if r.Method == http.MethodGet {
//
//		idRaw := r.URL.Query().Get("id")
//
//		id, err := strconv.Atoi(idRaw)
//
//		if err != nil {
//			panic(err)
//		}
//
//		err = ach.categorySrv.DeleteCategory(uint(id))
//
//		if err != nil{
//			panic(err.Error())
//		}
//	}
//
//	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
//}
//
//func writeFile(mf *multipart.File, fname string) {
//
//	wd, err := os.Getwd()
//
//	if err != nil {
//		panic(err)
//	}
//
//	path := filepath.Join(wd, "../", "../", "ui", "assets", "img", fname)
//	image, err := os.Create(path)
//
//	if err != nil {
//		panic(err)
//	}
//	defer image.Close()
//	io.Copy(image, *mf)
//}
//
