package autproductcontroller

import (
	"errors"
	"net/http"
	autproductmodel "product/aut-product/aut-product-model"
	"product/entities"
	"product/helper"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func ViewLogin(c echo.Context) error {
	helper.Lock.Lock()
	defer helper.Lock.Unlock()

	helper.Template(c, "view/login.html", nil)

	return c.NoContent(http.StatusOK)
}

func Login(c echo.Context) error {
	helper.Lock.Lock()
	defer helper.Lock.Unlock()

	username := c.FormValue("username")
	password := c.FormValue("password")

	user := entities.Users{}

	err := autproductmodel.Login(&user, username, password)
	if err != nil {
		err = errors.New("Username atau Password salah")
		data := map[string]interface{}{
			"error": err,
		}
		helper.Template(c, "view/login.html", data)
	}

	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["loggedIn"] = true
	sess.Values["email"] = user.Email
	sess.Values["username"] = user.Username
	sess.Values["name"] = user.Name

	sess.Save(c.Request(), c.Response())

	c.Redirect(http.StatusSeeOther, "/pasien")

	return c.NoContent(http.StatusOK)
}

func ViewRegister(c echo.Context) error {
	helper.Lock.Lock()
	defer helper.Lock.Unlock()

	helper.Template(c, "view/register.html", nil)

	return c.NoContent(http.StatusOK)
}

func Register(c echo.Context) error {
	helper.Lock.Lock()
	defer helper.Lock.Unlock()

	user := entities.Users{
		Name:      c.FormValue("name"),
		Email:     c.FormValue("email"),
		Username:  c.FormValue("username"),
		Password:  c.FormValue("password"),
		Cpassword: c.FormValue("cpassword"),
	}

	errorMessage := make(map[string]interface{})

	if err := c.Validate(user); err != nil {
		arr := helper.ConvertErr(err)

		errorMessage["validation"] = arr
		errorMessage["user"] = user

		helper.Template(c, "view/register.html", errorMessage)
	} else {
		email := autproductmodel.Unic(user, user.Email, "email")
		username := autproductmodel.Unic(user, user.Username, "username")

		if email || username {
			unic := make(map[string]interface{})
			if email {
				unic["Email"] = "Email Sudah Di Gunakan"
			} else {
				unic["Username"] = "Username Sudah Di Gunakan"
			}
			errorMessage["validation"] = unic
			errorMessage["user"] = user

			helper.Template(c, "view/register.html", errorMessage)
		} else {

			user.Password, _ = helper.HashPassword(user.Cpassword)
			err := autproductmodel.Register(&user)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			errorMessage["pesan"] = "Registrasi Berhasil, Silahkan Login"

			helper.Template(c, "view/register.html", errorMessage)
		}
	}

	return c.NoContent(http.StatusOK)
}

func Logout(c echo.Context) error {
	helper.Lock.Lock()
	defer helper.Lock.Unlock()

	sess, _ := session.Get("session", c)

	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	c.Redirect(http.StatusSeeOther, "/")

	return c.NoContent(http.StatusOK)
}
