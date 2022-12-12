package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/sessions"
)

var tmpl *template.Template
var tmpl2 *template.Template
var originaldb *sql.DB
var db *sql.DB

// dummy user data
// var users = map[string]string{"user1@12": "password", "user2": "password"}
var users = existsaccount()

// creating a cookie session store
var store = sessions.NewCookieStore([]byte("secret_key"))

func getMySqlDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@(127.0.0.1:3306)/customer_portal")
	if err != nil {
		panic(err)
	}
	return db
	// See "Important settings" section.
	// db.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)
}
func getOriginalDB() *sql.DB {
	originaldb, err := sql.Open("mysql", "root:@(127.0.0.1:3306)/rocket_development")
	if err != nil {
		panic(err)
	}
	return originaldb
}

// reading the HTML files
func init() {
	tmpl = template.Must(template.ParseFiles("login.html"))
	tmpl2 = template.Must(template.ParseFiles("forms.html"))
}

// ///////////////////////////////////first part//////////////////////////////////////////
// customer login credential
type credentialInfo struct {
	Username string
	Password string
}

func StudentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	log_in := credentialInfo{

		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	if Authentication(w, r, &log_in.Username, &log_in.Password) {
		ProductsDetails(w, r, &log_in.Username)
	}

	// tmpl.Execute(w, struct {
	// 	Success bool
	// 	Log_in  credentialInfo
	// }{true, log_in})
}
func existsaccount() map[string]string {

	db = getMySqlDB()
	users := make(map[string]string)

	rows, err := db.Query("Select username, password From credentials")

	if err != nil {
		fmt.Printf("something went wrong!")
	} else {
		var (
			username string
			password string
		)

		for rows.Next() {
			rows.Scan(&username, &password)
			users[username] = fmt.Sprint(password)
		}

		fmt.Println(users)

	}
	return users
}

func signup(w http.ResponseWriter, r *http.Request) {
	db = getMySqlDB()
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	} else {
		newuser := credentialInfo{
			Username: r.FormValue("new_username"),
			Password: r.FormValue("new_password"),
		}
		_, exists := users[newuser.Username]
		if exists {
			w.Write([]byte("Account already exists,Please sign in!"))
		}
		_, err := db.Exec("insert into credentials (username,password) values(?,?)", newuser.Username, newuser.Password)
		if err != nil {
			w.Write([]byte("Something went wrong!"))
		}
		tmpl.Execute(w, struct {
			Message string
		}{"MEOW! YOU made it a new account!"})
		users = existsaccount()
	}
}
func Authentication(w http.ResponseWriter, r *http.Request, e *string, p *string) bool {
	// check if user exists
	storedPassword, exists := users[*e]
	if exists {
		// 	// Get registers and returns a session for the given name and session store.
		// 	// session.id is the name of the cookie that will be stored in the client's browser
		session, _ := store.Get(r, "session.id")
		if storedPassword == *p {
			session.Values["authenticated"] = true
			// 		// saves all sessions used during the current request
			session.Save(r, w)
			fmt.Printf("Login successfully")
			return true
		} else {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return false
		}
	}
	w.Write([]byte("Account doesn't exist,Please sign up first!"))
	return false
}

// Collect input info and show on html
// func StudentHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		tmpl.Execute(w, nil)
// 		return
// 	}
// 	log_in := credentialInfo{

// 		Username: r.FormValue("username"),
// 		Password: r.FormValue("password"),
// 	}
// 	tmpl.Execute(w, struct {
// 		Success bool
// 		Log_in  credentialInfo
// 	}{true, log_in})
// }

// ///////////////////////////////
// ///////////////////////////////
// ///////////////////////////////
// ///////////////////////////////////second and third part//////////////////////////////////////////

type CustomerInfo []struct {
	AddressID                int         `json:"addressId"`
	UserID                   int         `json:"userId"`
	ID                       int         `json:"id"`
	CustomerCreationDate     string      `json:"customerCreationDate"`
	Date                     interface{} `json:"date"`
	CompanyName              string      `json:"companyName"`
	CompanyHqadress          interface{} `json:"companyHqadress"`
	FullNameOfCompanyContact string      `json:"fullNameOfCompanyContact"`
	CompanyContactPhone      string      `json:"companyContactPhone"`
	CompanyContactEmail      string      `json:"companyContactEmail"`
	CompanyDesc              interface{} `json:"companyDesc"`
	FullNameServiceTechAuth  string      `json:"fullNameServiceTechAuth"`
	TechAuthPhoneService     string      `json:"techAuthPhoneService"`
	TechManagerEmailService  string      `json:"techManagerEmailService"`
	CreatedAt                string      `json:"createdAt"`
	UpdatedAt                string      `json:"updatedAt"`
	Address                  interface{} `json:"address"`
	User                     interface{} `json:"user"`
	Buildings                Building    `json:"buildings"`
}
type Building []struct {
	CustomerID                       int           `json:"customerId"`
	AddressID                        int           `json:"addressId"`
	ID                               int           `json:"id"`
	FullNameOfBuildingAdmin          string        `json:"fullNameOfBuildingAdmin"`
	EmailOfAdminOfBuilding           string        `json:"emailOfAdminOfBuilding"`
	PhoneNumOfBuildingAdmin          int           `json:"phoneNumOfBuildingAdmin"`
	FullNameOfTechContactForBuilding string        `json:"fullNameOfTechContactForBuilding"`
	TechContactEmailForBuilding      string        `json:"techContactEmailForBuilding"`
	TechContactPhoneForBuilding      int           `json:"techContactPhoneForBuilding"`
	CreatedAt                        string        `json:"createdAt"`
	UpdatedAt                        string        `json:"updatedAt"`
	Address                          interface{}   `json:"address"`
	Batteries                        Battery       `json:"batteries"`
	BuildingDetails                  []interface{} `json:"buildingDetails"`
}
type Battery []struct {
	EmployeeID         int         `json:"employeeId"`
	BuildingID         int         `json:"buildingId"`
	ID                 int         `json:"id"`
	Type               string      `json:"type"`
	Status             string      `json:"status"`
	CommissionDate     string      `json:"commissionDate"`
	LastInspectionDate string      `json:"lastInspectionDate"`
	OperationsCert     string      `json:"operationsCert"`
	Information        string      `json:"information"`
	Notes              string      `json:"notes"`
	CreatedAt          string      `json:"createdAt"`
	UpdatedAt          string      `json:"updatedAt"`
	Employee           interface{} `json:"employee"`
	Columns            Column      `json:"columns"`
}
type Column []struct {
	BatteryID         int      `json:"batteryId"`
	ID                int      `json:"id"`
	Type              string   `json:"type"`
	NumOfFloorsServed int      `json:"numOfFloorsServed"`
	Status            string   `json:"status"`
	Information       string   `json:"information"`
	Notes             string   `json:"notes"`
	CreatedAt         string   `json:"createdAt"`
	UpdatedAt         string   `json:"updatedAt"`
	Elevators         Elevator `json:"elevators"`
}
type Elevator []struct {
	ColumnID           int    `json:"columnId"`
	ID                 int    `json:"id"`
	SerialNumber       int    `json:"serialNumber"`
	Model              string `json:"model"`
	Type               string `json:"type"`
	Status             string `json:"status"`
	CommisionDate      string `json:"commisionDate"`
	LastInspectionDate string `json:"lastInspectionDate"`
	InspectionCert     string `json:"inspectionCert"`
	Information        string `json:"information"`
	Notes              string `json:"notes"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

var todostruct CustomerInfo

func ProductsDetails(w http.ResponseWriter, r *http.Request, e *string) {
	// customer, err := http.Get("https://localhost:7189/api/customer/custemail/ernest.doyle@weber.net")
	customer, err := http.Get("https://localhost:7189/api/customer/custemail/" + *e)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer customer.Body.Close()

	responsedata, err := ioutil.ReadAll(customer.Body)

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(responsedata, &todostruct)
	// fmt.Println(todostruct[0].ID)

	data := []string{}
	// datatotable := make(map[string]interface{})
	for i := 0; i < len(todostruct[0].Buildings); i++ {
		mybuilding := todostruct[0].Buildings[i]
		// fmt.Println(todostruct[0].Buildings[i].Batteries)
		data = append(data, "Building", "id:", strconv.Itoa(mybuilding.ID), "\n")
		for j := 0; j < len(todostruct[0].Buildings[i].Batteries); j++ {
			mybattery := todostruct[0].Buildings[i].Batteries[j]
			// fmt.Println(mybattery.ID)

			data = append(data, "Battery", "id:", strconv.Itoa(mybattery.ID), "\n")
			for m := 0; m < len(todostruct[0].Buildings[i].Batteries[j].Columns); m++ {
				mycolumn := todostruct[0].Buildings[i].Batteries[j].Columns[m]
				// fmt.Println(mycolumn)
				data = append(data, "Column", "id:", strconv.Itoa(mycolumn.ID), "\n")

				for n := 0; n < len(todostruct[0].Buildings[i].Batteries[j].Columns[m].Elevators); n++ {
					myelevator := todostruct[0].Buildings[i].Batteries[j].Columns[m].Elevators[n]
					// fmt.Println(myelevator)
					data = append(data, "Elevator", "id:", strconv.Itoa(myelevator.ID))

				}
			}
		}

	}

	tmpl2.Execute(w, struct {
		// Log_in  credentialInfo
		Message string
	}{Message: strings.Trim(fmt.Sprint(data), "[]")})
	fmt.Println(data)
}

// func ProductsDetailshandler(w http.ResponseWriter, r *http.Request) {
// 	customer, err := http.Get("https://localhost:7189/api/customer/custemail/ernest.doyle@weber.net")
// 	if err != nil {
// 		fmt.Print(err.Error())
// 		os.Exit(1)
// 	}
// 	defer customer.Body.Close()

// 	responsedata, err := ioutil.ReadAll(customer.Body)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var todostruct CustomerInfo

// 	json.Unmarshal(responsedata, &todostruct)
// 	// fmt.Println(todostruct[0].ID)

// 	data := []string{}
// 	// datatotable := make(map[string]interface{})
// 	for i := 0; i < len(todostruct[0].Buildings); i++ {
// 		mybuilding := todostruct[0].Buildings[i]
// 		// fmt.Println(todostruct[0].Buildings[i].Batteries)
// 		data = append(data, "Building", "id:", strconv.Itoa(mybuilding.ID), "\n")
// 		for j := 0; j < len(todostruct[0].Buildings[i].Batteries); j++ {
// 			mybattery := todostruct[0].Buildings[i].Batteries[j]
// 			// fmt.Println(mybattery.ID)

// 			data = append(data, "Battery", "id:", strconv.Itoa(mybattery.ID), "\n")
// 			for m := 0; m < len(todostruct[0].Buildings[i].Batteries[j].Columns); m++ {
// 				mycolumn := todostruct[0].Buildings[i].Batteries[j].Columns[m]
// 				// fmt.Println(mycolumn)
// 				data = append(data, "Column", "id:", strconv.Itoa(mycolumn.ID), "\n")

// 				for n := 0; n < len(todostruct[0].Buildings[i].Batteries[j].Columns[m].Elevators); n++ {
// 					myelevator := todostruct[0].Buildings[i].Batteries[j].Columns[m].Elevators[n]
// 					// fmt.Println(myelevator)
// 					data = append(data, "Elevator", "id:", strconv.Itoa(myelevator.ID))

// 				}
// 			}
// 		}

// 	}

// 	tmpl2.Execute(w, struct {
// 		Message string
// 	}{Message: strings.Trim(fmt.Sprint(data), "[]")})
// 	fmt.Println(data)
// }

// customer_id := todostruct[0].ID
//nil will give back the template itself.
//tmpl2.Execute(w, nil)
// tmpl2.Execute(w,struct {
// 	Success bool
// 	Customer []struct
// }{})

// ///////////////////////////////
// ///////////////////////////////
// ///////////////////////////////
// ///////////////////////////////////third part//////////////////////////////////////////

type DropdownItem struct {
	Name   string
	Values string
}
type interventionInfo struct {
	author         string
	customer_id    string
	building_id    string
	battery_id     string
	column_id      string
	elevator_id    string
	employee_id    string
	start_datetime string
	end_datetime   string
	result         string
	report         string
	status         string
}

var buildingdropdown = make(map[string]interface{})
var batterydropdown = make(map[string]interface{})
var columndropdown = make(map[string]interface{})

var elevatordropdown = make(map[string]interface{})

func dropdownHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
	<html>
	<body>
	
	<select > // for loop in html template example
	  {{range $key, $value := .}}
		<option value="{{ $value }}">{{ $key }}</option>
	  {{end}}

	</select>
	
</form>
	</body>
	</html>`

	dropdownTemplate, err := template.New("dropdownexample").Parse(string(html))
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(todostruct[0].Buildings); i++ {
		mybuilding := todostruct[0].Buildings[i]
		// fmt.Println(todostruct[0].Buildings[i].Batteries)
		// data = append(data, "id:", string(mycolumn.))
		buildingdropdown["Building"+strconv.Itoa(mybuilding.ID)] = strconv.Itoa(mybuilding.ID)

		for j := 0; j < len(todostruct[0].Buildings[i].Batteries); j++ {
			mybattery := todostruct[0].Buildings[i].Batteries[j]
			// fmt.Println(mybattery.ID)

			batterydropdown["Battery"+strconv.Itoa(mybattery.ID)] = strconv.Itoa(mybattery.ID)

			for m := 0; m < len(todostruct[0].Buildings[i].Batteries[j].Columns); m++ {
				mycolumn := todostruct[0].Buildings[i].Batteries[j].Columns[m]
				// fmt.Println(mycolumn)
				columndropdown["Column"+strconv.Itoa(mycolumn.ID)] = strconv.Itoa(mycolumn.ID)

				for n := 0; n < len(todostruct[0].Buildings[i].Batteries[j].Columns[m].Elevators); n++ {
					myelevator := todostruct[0].Buildings[i].Batteries[j].Columns[m].Elevators[n]
					// fmt.Println(myelevator)
					elevatordropdown["Elevator"+strconv.Itoa(myelevator.ID)] = strconv.Itoa(myelevator.ID)

				}
			}
		}
		// tmpl.Execute(w, struct {
		// 	Message string
		// }{Message: strings.Trim(fmt.Sprint(data), "[]")})
		// // fmt.Println(data)
		for k, v := range batterydropdown {
			fmt.Printf(k, v)
		}

	}
	// populate dropdown with fruits

	dropdownTemplate.Execute(w, buildingdropdown)
	dropdownTemplate.Execute(w, batterydropdown)

	dropdownTemplate.Execute(w, columndropdown)
	dropdownTemplate.Execute(w, elevatordropdown)
}

func newinterventionhandler(w http.ResponseWriter, r *http.Request) {
	newintervention := interventionInfo{
		// 		author:         r.FormValue("author"),
		// 		customer_id:    r.FormValue("customer_id"),
		// 		building_id:    r.FormValue("building_id"),
		// 		battery_id:     r.FormValue("battery_id"),
		// 		column_id:      r.FormValue("column_id"),
		// 		elevator_id:    r.FormValue("elevator_id"),
		// 		employee_id:    r.FormValue("employee_id"),
		// 		start_datetime: "",
		// 		end_datetime:   "",
		// 		result:         "Incomplete",
		// 		report:         r.FormValue("sid"),
		// 		status:         "Pending",
		// 	}
		// todo add submit button put into database
	}
	_ = newintervention
}

// func interventionHandler(w http.ResponseWriter, r *http.Request) {
// 	originaldb = getOriginalDB()

// 	if r.Method != http.MethodPost {
// 		tmpl.Execute(w, nil)
// 		return
// 	}

// 	newintervention := interventionInfo{
// 		author:         r.FormValue("author"),
// 		customer_id:    r.FormValue("customer_id"),
// 		building_id:    r.FormValue("building_id"),
// 		battery_id:     r.FormValue("battery_id"),
// 		column_id:      r.FormValue("column_id"),
// 		elevator_id:    r.FormValue("elevator_id"),
// 		employee_id:    r.FormValue("employee_id"),
// 		start_datetime: "",
// 		end_datetime:   "",
// 		result:         "Incomplete",
// 		report:         r.FormValue("sid"),
// 		status:         "Pending",
// 	}

// 	// Sid, _ := strconv.Atoi(student.Sid)
// 	if r.FormValue("submit") == "Send" {
// 		_, err := originaldb.Exec("insert into interventions (author,customer_id,building_id,battery_id,column_id,elevator_id,employee_id,start_datetime,end_datetime,result,report,status) values(?,?,?,?,?,?,?,?,?,?,?,?)", newintervention.author, newintervention.customer_id, newintervention.building_id, newintervention.battery_id, newintervention.column_id, newintervention.column_id, newintervention.elevator_id, newintervention.employee_id, newintervention.start_datetime, newintervention.end_datetime, newintervention.report, newintervention.result, newintervention.status)
// 		if err != nil {
// 			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
// 		} else {
// 			w.Write([]byte("Login successfully!"))
// 		}
// 	}

// }

// func Ajaxhandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("method:", r.Method)
// 	if r.Method == "GET" {
// 		r.ParseForm()
// 		fmt.Println(r.Form)
// 		if r.Form["button"][0] == "Search" {
// 			fmt.Println(r.Form)
// 			fmt.Println(r.Form["email"][0])

// 			if err != nil {
// 				panic(err)
// 			}
// 			defer session.Close()
// 			session.SetMode(mgo.Monotonic, true)
// 			c := session.DB("user").C("profile")
// 			result := Users{}
// 			err = c.Find(bson.M{"email": r.Form["email"][0]}).All(&result)
// 			fmt.Println(result)
// 			b, err := json.MarshalIndent(result, "", "  ")
// 			if err != nil {
// 				panic(err)
// 			}
// 			fmt.Printf("%s\n", b)
// 			// set header to 'application/json'
// 			w.Header().Set("Content-Type", "application/json")
// 			// write the response
// 			w.Write(b)
// 		}
// 	}
// }

func main() {

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets"))

	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", StudentHandler)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/ajaxtest", dropdownHandler)

	// mux.HandleFunc("/profile", ProductsDetailshandler)
	mux.HandleFunc("/newintervention", newinterventionhandler)
	http.ListenAndServe(":8080", mux)
}

// func Gin() {
// 	r := gin.Default()
// 	r.GET("/ping", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"message": "pong",
// 		})
// 	})
// 	r.Run() // listen and serve on 0.0.0.0:8080
// }
