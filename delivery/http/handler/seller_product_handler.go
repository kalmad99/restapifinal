package handler

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"

	"../../../entity"
	"../../../productpage"
)

// AdminCategoryHandler handles category handler admin requests
type SellerProductHandler struct {
	productSrv productpage.ItemService
}

// NewAdminCategoryHandler initializes and returns new AdminCateogryHandler
func NewSellerProductHandler(is productpage.ItemService) *SellerProductHandler {
	return &SellerProductHandler{productSrv: is}
}

// AdminCategories handle requests on route /admin/categories
func (sph *SellerProductHandler) SellerProducts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	products, errs := sph.productSrv.Items()
	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	output, err := json.MarshalIndent(products, "", "\t\t")

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
func (sph *SellerProductHandler) SellerProductsNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	l := r.ContentLength
	log.Println("content length is", l)
	body := make([]byte, l)
	r.Body.Read(body)

	var pro []entity.Product
	//
	err := json.Unmarshal(body, &pro)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		//http.Error(w, "first", http.StatusNotFound)
		return
	}
	product, errs := sph.productSrv.StoreItem(&pro[0])

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		//http.Error(w, "second", http.StatusNotFound)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	p := fmt.Sprintf("/v1/seller/products/%d", product.ID)
	log.Println(product.Name)
	log.Println(p)
	w.Header().Set("Location", p)
	w.WriteHeader(http.StatusCreated)
	return
}

// AdminCategoriesUpdate handle requests on /admin/categories/update
func (sph *SellerProductHandler) SellerProductsUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	product, errs := sph.productSrv.Item(uint(id))

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	l := r.ContentLength

	body := make([]byte, l)
	r.Body.Read(body)
	for _, b := range body{
		log.Print(string(b))
	}

	var pro []entity.Product

	err = json.Unmarshal(body, &pro)
	if err != nil {
		panic(err.Error())
	}

	product, errs = sph.productSrv.UpdateItem(&pro[0])

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(product, "", "\t\t")

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
func (sph *SellerProductHandler) SellerProductsDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, errs := sph.productSrv.DeleteItem(uint(id))

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

func (sph *SellerProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := ps.ByName("search")

	products, err := sph.productSrv.SearchProduct(res)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(products, "", "\t\t")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

func (sph *SellerProductHandler) ProductDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	prod, errs := sph.productSrv.Item(uint(id))

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	output, err := json.MarshalIndent(prod, "", "\t\t")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

// AdminCategoriesUpdate handle requests on /admin/categories/update
func (sph *SellerProductHandler) RateProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	product, errs := sph.productSrv.Item(uint(id))

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	prod := &entity.Product{}
	id, _ = strconv.Atoi(r.FormValue("id"))
	prod.ID = uint(id)
	prod.Rating, _ = strconv.ParseFloat(r.FormValue("rating"), 64)

	log.Println("Rating is ", prod.Rating)
	product, errs = sph.productSrv.RateProduct(prod)

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(product, "", "\t\t")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}