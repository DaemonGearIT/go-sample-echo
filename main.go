package main

import(
	"fmt"
	"io"
	"os"
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo/engine/standard"
)

//Define types :
// 	User model that contain user basic information
//	handler struct contain a map named DB that save user information
type (
 	User struct {
 		Email string `json:"email" form:"email"`
 		Password string `json:"passwd" form:"passwd"`
 	}
 	handler struct {
 		db map[string]*User
 	}
)

//Basic method that creates a new user
func (h *handler) createUser(c echo.Context) error {
	//Create a new user instance
	u := new(User)

	//User attributes bind with the request
	if err := c.Bind(u); err != nil {
		return err
	}

	//Add User to DB
	h.db[u.Email] = u
	
	fmt.Println("DB: ", len(h.db))	
	return c.JSON(http.StatusCreated, u)
}

//Basic method that retrieves a User by email
func (h *handler) getUser(c echo.Context) error {
	//Obtain email paramater from request
	email := c.Param("email")

	//find user by email
	user := h.db[email]

	//user not found, throws error
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, user)
}

//Basic method that retrieves all Users
func (h * handler) getAllUsers(c echo.Context) error {
	fmt.Println("DB: ", len(h.db))
	//Transform User map into Array
	list := h.transformMapToArray(h.db)

	fmt.Println("List: ", len(list))
	return c.JSON(http.StatusOK, list)
}

//Basic method that updates an existing User
func (h *handler) updateUser(c echo.Context) error {
	//Obtain email paramater from request
	email := c.Param("email")
	fmt.Println("Updating User: ", email)
	//Find user by email
	user := h.db[email]
	
	//user not found, throws error
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	//Define a new User that contains PUT information
	uUser := new(User)
	
	//User attributes bind with the request
	if err := c.Bind(uUser); err != nil {
		return err
	}

	//Validate attributes is not empty
	if uUser.Email != "" {
		user.Email = uUser.Email	
	}
	
	if uUser.Password != "" {
		user.Password = uUser.Password	
	}
	
	//Update database info
	h.db[email] = user

	return c.JSON(http.StatusOK, user)

}

//Basic method that deletes a User by email
func (h *handler) deleteUser(c echo.Context) error {
	//Obtain email paramater from request
	email := c.Param("email")
	fmt.Println("Delete User : " , email)
	//Find user by email
	user := h.db[email]

	//user not found, throws error
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	//Delete found user from the database
	delete(h.db, email)

	return c.JSON(http.StatusOK, h.db)
}

func (h *handler) fileUpload(c echo.Context) error {
	//Obtain file parameter from request
	fmt.Println("Uploading File")
	file, err := c.FormFile("file")

	if err != nil {
		fmt.Println("Error!! ", err)
		return err
	}

	fmt.Println("File", file)
	//Open file
	src, err := file.Open() 
	if err != nil {
		return err
	}

	//Create destination 
	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	//Copy file into destination 
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.String(http.StatusOK, "File uploaded!!")

}

//Transform User map to an Array
func (h *handler) transformMapToArray(m map[string]*User) []*User {
	//Define a User Array
	list := make([]*User, 0, len(m))
	//Append Users to Array
	for _, value := range m {
		list = append(list, value)
	}

	return list
}

//Create a new echo instance and runs on port 9090
func main() {
	h := new(handler)
	h.db = make(map[string]*User)
	
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Static("public"))


	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
	})

	//e.File("/upload", "public/index.html")

	//Routing User handler method
	e.GET("/users", h.getAllUsers)
	e.GET("/users/:email", h.getUser)
	e.POST("/users", h.createUser)
	e.PUT("/users/:email", h.updateUser)
	e.DELETE("/users/:email", h.deleteUser)
	e.POST("/upload", h.fileUpload)
	
	fmt.Println("Server running -----> localhost:9090")
	e.Run(standard.New(":9090"))
}


