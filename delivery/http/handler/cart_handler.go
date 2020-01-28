package handler
//
//
//import (
//	"html/template"
//	"net/http"
//	"strconv"
//
//	"../../../entity"
//	"../../../productpage"
//)
//
//// AdminCategoryHandler handles category handler admin requests
//type CartHandler struct {
//	tmpl        *template.Template
//	carSrv .ItemService
//}
//
//// NewAdminCategoryHandler initializes and returns new AdminCateogryHandler
//func NewCartHandler(t *template.Template, is productpage.ItemService) *SellerProductHandler {
//	return &SellerProductHandler{tmpl: t, productSrv: is}
//}
//
//
//func (sph *SellerProductHandler) ProductDetail(w http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodGet {
//
//		idRaw := r.URL.Query().Get("id")
//		id, err := strconv.Atoi(idRaw)
//
//		if err != nil {
//			panic(err)
//		}
//
//		pro, errs := sph.productSrv.Item(uint(id))
//
//		if len(errs) > 0 {
//			panic(errs)
//		}
//
//		_ = sph.tmpl.ExecuteTemplate(w, "productdetail.layout", pro)
//	}
//}
//
