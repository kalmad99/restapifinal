package handler

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"

	"../../../entity"
	"../../../user"
)

var name, email, phone, pass string
var id uint
// AdminCategoryHandler handles category handler admin requests
type UserHandler struct {
	userSrv user.UserService
}

// NewAdminCategoryHandler initializes and returns new AdminCateogryHandler
func NewUserHandler(us user.UserService) *UserHandler {
	return &UserHandler{userSrv: us}
}

// AdminCategories handle requests on route /admin/categories
func (uh *UserHandler) Users(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	users, errs := uh.userSrv.Users()
	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	output, err := json.MarshalIndent(users, "", "\t\t")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

func (uh *UserHandler) User(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	user, errs := uh.userSrv.User(uint(id))

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	output, err := json.MarshalIndent(user, "", "\t\t")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

// AdminCategoriesNew handle requests on route /admin/categories/new
func (uh *UserHandler) UserNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//decoder := json.NewDecoder(r.Body)
	//usr := entity.User{}
	//err := decoder.Decode(&usr)

	l := r.ContentLength
	log.Println("content length is", l)
	body := make([]byte, l)
	r.Body.Read(body)

	//usr := &entity.User{}
	var usr []entity.User
	//
	err := json.Unmarshal(body, &usr)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		//http.Error(w, "first", http.StatusNotFound)
		return
	}

	name = usr[0].Name
	email = usr[0].Email
	phone = usr[0].Phone
	pass = usr[0].Password
	id = usr[0].ID

	log.Println("Name", name)
	log.Println("Email", email)
	log.Println("Phone", phone)
	log.Println("Password", pass)
	log.Println("ID", id)

	hostURL := "smtp.gmail.com"
	hostPort := "587"
	emailSender := "kalemesfin12go@gmail.com"
	password := "qnzfgwbnaxykglvu"
	emailReceiver := usr[0].Email

	emailAuth := smtp.PlainAuth(
		"",
		emailSender,
		password,
		hostURL,
	)

	msg := []byte("To: " + emailReceiver + "\r\n" +
		"Subject: " + "Hello " + usr[0].Name + "\r\n" +
		"This is your OTP. 123456789")

	err = smtp.SendMail(
		hostURL + ":" + hostPort,
		emailAuth,
		emailSender,
		[]string{emailReceiver},
		msg,
	)

	if err != nil{
		fmt.Print("Error: ", err)
	}
	fmt.Print("Email Sent")

	p := fmt.Sprintf("localhost:8383/Registpage2")
	log.Println(usr[0].Name)
	log.Println(p)
	w.Header().Set("Location", p)
	text := fmt.Sprintf("To verify your email, we've sent a One Time Password (OTP) to %s", usr[0].Email)
	w.Write([]byte(text))
	w.WriteHeader(http.StatusOK)
	return
}

// AdminCategoriesNew handle requests on route /admin/categories/new
func (uh *UserHandler) UserNewValidity(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var otp string
	otp = r.FormValue("otpfield")

	usr := &entity.User{}
	usr.Name = name
	usr.Email = email
	usr.Phone = phone
	usr.Password = pass
	usr.ID = id

	if otp == "123456789"{
		user, errs := uh.userSrv.StoreUser(usr)
		if len(errs) > 0 {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		p := fmt.Sprintf("localhost:8383/v1/users/%d", user.ID)
		log.Println(user.Name)
		log.Println(p)
		w.Header().Set("Location", p)
		text := fmt.Sprintf("Successfully created account")
		w.Write([]byte(text))
		w.WriteHeader(http.StatusCreated)
		return
	}else {
		w.Header().Set("Content-Type", "application/json")
		//http.Error(w, "second", http.StatusNotFound)
		http.Error(w, "Wrong OTP", http.StatusNotFound)
		return
	}
}

func (uh *UserHandler) UserLogin(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//if req.Method != "POST" {
	//	http.Redirect(w, req, "/Loginpage", http.StatusSeeOther)
	//	return
	//}
	l := req.ContentLength
	body := make([]byte, l)
	req.Body.Read(body)
	user := &entity.User{}

	err := json.Unmarshal(body, user)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	usr, errs := uh.userSrv.Login(user.Email)

	log.Println(user.Email)

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(user.Password))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	p := fmt.Sprintf("/v1/users/%d", usr.ID)
	w.Header().Set("Location", p)
	w.WriteHeader(http.StatusCreated)
	return
}

// AdminCategoriesUpdate handle requests on /admin/categories/update
func (uh *UserHandler) UserUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	user, errs := uh.userSrv.User(uint(id))

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	l := r.ContentLength

	body := make([]byte, l)
	r.Body.Read(body)

	var usr []entity.User
	err = json.Unmarshal(body, &usr)
	if err != nil {
		panic(err.Error())
	}

	user, errs = uh.userSrv.UpdateUser(&usr[0])

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(user, "", "\t\t")

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
func (uh *UserHandler) UserDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, errs := uh.userSrv.DeleteUser(uint(id))

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}
//
//func (uh *UserHandler) UserChangePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	id, err := strconv.Atoi(ps.ByName("id"))
//	if err != nil {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//	user, errs := uh.userSrv.User(uint(id))
//
//	if len(errs) > 0 {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	l := r.ContentLength
//	body := make([]byte, l)
//	r.Body.Read(body)
//	usr := &entity.User{}
//
//	_ = json.Unmarshal(body, usr)
//
//	usr.Password =
//
//	user, errs = uh.userSrv.UpdateUser(user)
//
//	if len(errs) > 0 {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	output, err := json.MarshalIndent(user, "", "\t\t")
//
//	if err != nil {
//		w.Header().Set("Content-Type", "application/json")
//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(output)
//	return

	//if r.Method == http.MethodGet {
	//
	//	idRaw := r.URL.Query().Get("id")
	//	id, err := strconv.Atoi(idRaw)
	//
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	usr, errs := uh.userSrv.User(uint(id))
	//
	//	if len(errs) > 0 {
	//		panic(errs)
	//	}
	//
	//	uh.tmpl.ExecuteTemplate(w, "changepass.layout", usr)
	//
	//}
	//usr := &entity.User{}
	//id, _ := strconv.Atoi(r.FormValue("id"))
	////usr.ID = uint(id)
	//usr.ID = uint(id)
	//user, err1 := uh.userSrv.User(usr.ID)
	//if len(err1) > 0{
	//	panic(err1)
	//}
	//
	//log.Println(usr.ID)
	//usr.Password = r.FormValue("password")
	//var confp = r.FormValue("confpass")
	//var oldp = r.FormValue("oldpass")
	//
	//err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldp))
	//if err != nil {
	//	log.Println("This is not your old password")
	//	http.Redirect(w, r, "/users", 301)
	//	return
	//}
	//if usr.Password == confp {
	//
	//	hashedpass, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	//	if err != nil {
	//		//http.Error(w, "Server error, unable to create your account.", 500)
	//		//return errors.New("server error, unable to create your account")
	//		panic(err.Error())
	//	}
	//	usr.Password=string(hashedpass)
	//
	//	_, errs := uh.userSrv.ChangePassword(usr)
	//	log.Println("Succesfully changed")
	//	http.Redirect(w, r, "/Loginpage", 303)
	//	if len(errs) > 0{
	//		panic(errs)
	//	}
	//}else{
	//	log.Println("Passwords don't match")
	//	uh.tmpl.ExecuteTemplate(w, "user.update.layout", nil)
	//}



	//usr := &entity.User{}
	//id, _ := strconv.Atoi(r.FormValue("id"))
	//usr.ID = uint(id)
	//usr.Name = r.FormValue("name")
	//usr.Email = r.FormValue("email")
	//usr.Phone = r.FormValue("phone")
	//usr.Password = r.FormValue("newpass")
	//var confp = r.FormValue("confnewpass")
	//
	//if usr.Password == confp{
	//	hashedpass, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//	usr.Password=string(hashedpass)
	//
	//	_, errs := uh.userSrv.ChangePassword(usr)
	//
	//	if len(errs) > 0 {
	//		panic(errs)
	//	}
	//}else{
	//	errors.New("The passwords you entered dont match")
	//}
	//http.Redirect(w, r, "/users", http.StatusSeeOther)
//}
