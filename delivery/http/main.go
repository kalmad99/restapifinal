package main

import (
	"../../productpage/prorepository"
	"../../productpage/proservice"
	"../../user/repository"
	"../../user/service"
	"github.com/julienschmidt/httprouter"

	//"database/sql"
	"github.com/jinzhu/gorm"
	//"../../entity"
	"../http/handler"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"html/template"
	"net/http"
)

//func createTables(dbconn *gorm.DB) []error {
//	errs := dbconn.CreateTable(&entity.Product{}, &entity.AddToCart{}, &entity.Category{}, &entity.User{}, &entity.Role{}).GetErrors()
//	if errs != nil {
//		return errs
//	}
//	return nil
//}
func main() {

	dbconn, err := gorm.Open("postgres", "postgres://postgres:password@localhost/restaurantdb2?sslmode=disable")
	//dbconn, err := sql.Open("postgres", "postgres://postgres:password@localhost/restaurantdb2?sslmode=disable")

	if err != nil {
		panic(err)
	}

	defer dbconn.Close()

	//wd, err := os.Getwd()
	//path := filepath.Join(wd, "../", "../")
	//
	//log.Println(path)

	//createTables(dbconn)
	tmpl := template.Must(template.ParseGlob("../../ui/templates/*"))

	categoryRepo := prorepository.NewCategoryGormRepo(dbconn)
	//categoryRepo := prorepository.NewCategoryRepositoryImpl(dbconn)
	categoryServ := proservice.NewCategoryService(categoryRepo)

	productRepo := prorepository.NewItemGormRepo(dbconn)
	//productRepo := prorepository.NewProductRepositoryImpl(dbconn)
	productServ := proservice.NewItemService(productRepo)

	userRepo := repository.NewUserGormRepo(dbconn)
	userServ := service.NewUserService(userRepo)

	roleRepo := repository.NewRoleGormRepo(dbconn)
	roleSrv := service.NewRoleService(roleRepo)


	adminCatgHandler := handler.NewAdminCategoryHandler(categoryServ)
	sellerProHandler := handler.NewSellerProductHandler(productServ)
	userHandler := handler.NewUserHandler(userServ)

	menuHandler := handler.NewMenuHandler(tmpl, productServ)
	adminRoleHandler := handler.NewAdminRoleHandler(roleSrv)

	fs := http.FileServer(http.Dir("ui/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	router := httprouter.New()

	router.GET("/v1/admin/roles", adminRoleHandler.GetRoles)


	router.GET("/v1/users", userHandler.Users)
	router.GET("/v1/users/:id", userHandler.User)
	router.GET("/v1/login", userHandler.UserLogin)
	router.PUT("/v1/user/update/:id", userHandler.UserUpdate)
	router.POST("/v1/users", userHandler.UserNew)
	router.POST("/v1/usersvalid", userHandler.UserNewValidity)
	router.DELETE("/v1/user/delete/:id", userHandler.UserDelete)

	router.GET("/v1/products",sellerProHandler.SellerProducts)
	router.GET("/v1/searchProducts/:search", sellerProHandler.SearchProducts)
	router.GET("/v1/detail/:id", sellerProHandler.ProductDetail)
	router.POST("/v1/products", sellerProHandler.SellerProductsNew)
	router.POST("/v1/rate/:id", sellerProHandler.RateProduct)
	router.PATCH("/v1/products/update/:id", sellerProHandler.SellerProductsUpdate)
	router.DELETE("/v1/products/delete/:id", sellerProHandler.SellerProductsDelete)

	router.GET("/v1/categories", adminCatgHandler.AdminCategories)
	router.GET("/v1/categories/:id", adminCatgHandler.AdminCategory)
	router.POST("/v1/categories", adminCatgHandler.AdminCategoriesNew)
	router.PUT("/v1/categories/update/:id", adminCatgHandler.AdminCategoriesUpdate)
	router.DELETE("/v1/categories/delete/:id", adminCatgHandler.AdminCategoriesDelete)

	router.GET("/", menuHandler.Index)
	router.GET("/about", menuHandler.About)
	router.GET("/contact", menuHandler.Contact)
	router.GET("/menu", menuHandler.Menu)
	router.GET("/admin", menuHandler.Admin)
	router.GET("/Loginpage", menuHandler.LoginPage)
	router.GET("/Registpage", menuHandler.RegistPage)
	router.GET("/Registpage2", menuHandler.RegistPage2)


	//http.HandleFunc("/users", userHandler.Users)
	//http.HandleFunc("/users/success", userHandler.UserUpdate)
	////http.HandleFunc("/admin/users",userHandler.Users)
	//http.HandleFunc("/registrationprocess1", userHandler.UserNew)
	//http.HandleFunc("/user/update", userHandler.UserUpdate)
	//http.HandleFunc("/user/delete", userHandler.UserDelete)
	////http.HandleFunc("/user/changepass", userHandler.UserChangePassword)
	//http.HandleFunc("/login", userHandler.UserLogin)

	//http.HandleFunc("/admin/roles/new", roleHandler.AdminRolesNew)

	err = http.ListenAndServe(":8383", router)
	if err!=nil{
		panic(err.Error())
	}

}