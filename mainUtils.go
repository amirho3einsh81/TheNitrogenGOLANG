package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	ptime "github.com/yaa110/go-persian-calendar"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type App struct {
	DB     *sql.DB
	Port   string
	Router *mux.Router
}

func (a *App) Initialize() {
	db, err := sql.Open("sqlite3", "./nitrogen.db")
	if err != nil {
		return
	}
	a.DB = db
	a.Router = mux.NewRouter()
	a.initialRouters()
}
func errViewer(w http.ResponseWriter, code int, err string) {
	viewer(w, code, map[string]string{"Error": err})
}
func viewer(w http.ResponseWriter, code int, content interface{}) {
	js, err := json.Marshal(content)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(js)
		return
	}
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}
func logView(code int, content interface{}) {
	fmt.Println("code: ", code, " message: ", content)
}
func (a *App) initialRouters() {
	//initial get all users from db
	/* for register a new user */
	a.Router.HandleFunc("/rEgIsTeRuSeR", a.registerUser).Methods("POST")
	/* for login a user or check user exists in nitrogen */
	a.Router.HandleFunc("/user/{userName}", a.checkUsername).Methods("GET")
	a.Router.HandleFunc("/phone/{phone}", a.checkPhone).Methods("GET")
	a.Router.HandleFunc("/user/{userName}/{u}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/editusername/{userName}", a.editUsername).Methods("POST")
	a.Router.HandleFunc("/information", a.Information).Methods("OPTIONS")
	a.Router.HandleFunc("/login/{userName}", a.Login).Methods("GET")
	a.Router.HandleFunc("/challenges", a.Challenges).Methods("GET")
	a.Router.HandleFunc("/home", a.Home).Methods("GET")
	a.Router.HandleFunc("/challenge/{challengeID}", a.Challenge).Methods("GET")
	a.Router.HandleFunc("/banner", a.banners).Methods("GET")
	a.Router.HandleFunc("/location", a.Location).Methods("POST")
	/* for edit user informations */
	a.Router.HandleFunc("/editprofile", a.editProf).Methods("POST")

	a.Router.HandleFunc("/editimg", a.EditImage).Methods("PUT")
	a.Router.HandleFunc("/started/{challengeID}", a.started).Methods("POST")
	a.Router.HandleFunc("/suggestion", a.Suggestion).Methods("GET")
	a.Router.HandleFunc("/enabledch", a.EnabledChallenges).Methods("GET")
	a.Router.HandleFunc("/enabledc", a.EnabledChallenge).Methods("GET")
	a.Router.HandleFunc("/", a.ActionGoal).Methods("POST")
}
func (a *App) Run() {
	fmt.Println("service started")
	log.Fatal(http.ListenAndServe(a.Port, a.Router))
}
func (a *App) checkUsername(w http.ResponseWriter, r *http.Request) {
	path := strings.ToLower(strings.TrimPrefix(r.URL.Path, "/user/"))
	rows, err := a.DB.Query("SELECT userName,pass,phoneNumber FROM users WHERE userName = ? ", path)
	if err != nil {
		errViewer(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	var u User
	for rows.Next() {
		err := rows.Scan(&u.Username, &u.Pass, &u.PhoneNum)
		if err != nil {
			errViewer(w, http.StatusNotFound, err.Error())
			return
		}
		viewer(w, http.StatusOK, u)
	}
}
func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/user/")
	s := strings.Split(path, "/")
	rows, err := a.DB.Query("SELECT userName,phoneNumber FROM users WHERE userName = ?", strings.ToLower(s[0]))
	if err != nil {
		errViewer(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	var u User
	for rows.Next() {
		err := rows.Scan(&u.Username, &u.PhoneNum)
		if err != nil {
			errViewer(w, http.StatusNotFound, err.Error())
			return
		}
		//viewer(w, http.StatusOK, u)
	}
	if strings.Contains(path, "/user") {
		viewer(w, http.StatusOK, u)
	} else if strings.Contains(path, "/information") {
		w.Write([]byte(path))
	} else if strings.Contains(path, "/purchase") {
		w.Write([]byte(path))
	} else {
		w.Header().Set("Content-Type", "text")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("status not found"))
	}
}
func (a *App) checkPhone(w http.ResponseWriter, r *http.Request) {
	path := strings.ToLower(strings.TrimPrefix(r.URL.Path, "/phone/"))
	rows, err := a.DB.Query("SELECT pass,phoneNumber FROM users WHERE phoneNumber = ? ", path)
	if err != nil {
		errViewer(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	var u User
	for rows.Next() {
		err := rows.Scan(&u.Pass, &u.PhoneNum)
		if err != nil {
			errViewer(w, http.StatusNotFound, err.Error())
			return
		}
		viewer(w, http.StatusOK, u)
	}
}
func (a *App) registerUser(w http.ResponseWriter, r *http.Request) {
	js, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("u 1 : ", err.Error())
		return
	}
	u := Register{}
	err = json.Unmarshal(js, &u)
	if err != nil {
		fmt.Println(" u2 : ", err.Error())
		return
	}
	_, err = a.DB.Exec("INSERT INTO users(phoneNumber,pass,userName) VALUES (?,?,?)", u.PhoneNum, u.Pass, strings.ToLower(u.Username))
	if err != nil {
		fmt.Println(" u3 : ", err.Error())
		return
	}
	_, err = a.DB.Exec("INSERT INTO informations(userName,title,biography,family,gender,userIP,brand,device,model) VALUES (?,?,?,?,?,?,?,?,?)", u.Username, u.Title, u.Biography, "", 0, u.UserIP, u.Brand, u.Device, u.Model)
	if err != nil {
		fmt.Println("u4 : ", err.Error())
		return
	}
	fmt.Println("ip user: ", u.UserIP)
	viewer(w, http.StatusCreated, u)
}
func (a *App) editProf(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//id := vars["id"]
	userName := strings.ToLower(r.Header.Get("User-name"))
	rows, err := a.DB.Query("SELECT userName, phoneNumber FROM users WHERE userName = ?", userName)
	var uinfo UserInformation
	if err != nil {
		//errViewer(w, http.StatusInternalServerError, "userName : "+userName+" "+err.Error())
		return
	} else {
		//viewer(w, http.StatusOK, "userName : "+userName)
		var uBody, exists, u UserInformation

		for rows.Next() {
			err = rows.Scan(&u.PhoneNum, &u.Username)
		}
		if err != nil {
			//errViewer(w, http.StatusInternalServerError, ""+err.Error())
			return
		}
		js, err1 := ioutil.ReadAll(r.Body)

		if err1 != nil {
			//errViewer(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := json.Unmarshal(js, &uBody); err != nil {
			return
		}
		user := strings.ToLower(uBody.Username)
		//start replacement contents
		query, err := a.DB.Query("SELECT userName FROM users WHERE userName = ?", user)
		if err != nil {
			return
		}
		for query.Next() {
			query.Scan(&exists.Username)
		}
		if !strings.EqualFold(exists.Username, user) {
			if user != "" {
				if !strings.EqualFold(u.Username, user) {
					_, err := a.DB.Exec("UPDATE users SET userName = ? WHERE userName = ?", user, userName)
					if err != nil {
						uinfo.Username = "0"
						return
					}
					_, err = a.DB.Exec("UPDATE informations SET userName = ? WHERE userName = ?", user, userName)
					_, err = a.DB.Exec("UPDATE payments SET userName = ? WHERE userName = ?", user, userName)

					uinfo.Username = "1"
				}
			}
		} else {
			uinfo.Username = "2"
		}

		query, err = a.DB.Query("SELECT phoneNumber FROM users WHERE phoneNumber = ?", uBody.PhoneNum)
		if err != nil {
			return
		}
		for query.Next() {
			query.Scan(&exists.PhoneNum)
		}
		if !strings.EqualFold(exists.PhoneNum, uBody.PhoneNum) {
			if uBody.PhoneNum != "" {
				if !strings.EqualFold(u.PhoneNum, uBody.PhoneNum) {
					_, err := a.DB.Exec("UPDATE users SET phoneNumber = ? WHERE userName= ?", uBody.PhoneNum, userName)
					if err != nil {
						uinfo.PhoneNum = "0"
						return
					}
					uinfo.PhoneNum = "1"
				}
			}
		} else {
			uinfo.PhoneNum = "2"
		}

		rows, err = a.DB.Query("SELECT userName, title, biography, family FROM informations WHERE userName = ?", userName)
		for rows.Next() {
			rows.Scan(&u.Title, &u.Biography)
		}
		if uBody.Title != "" {
			_, err = a.DB.Exec("UPDATE informations SET title =? WHERE userName = ?", uBody.Title, userName)
			if err != nil {
				uinfo.Title = "0"
				return
			}
			uinfo.Title = "1"
		}
		if uBody.Biography != "" {
			_, err = a.DB.Exec("UPDATE informations SET biography = ? WHERE userName = ?", uBody.Biography, userName)
			if err != nil {
				uinfo.Biography = "0"
				return
			}
			uinfo.Biography = "1"
		}
		if uBody.Family != "" {
			_, err = a.DB.Exec("UPDATE informations SET family = ? WHERE userName = ?", uBody.Family, userName)
			if err != nil {
				uinfo.Family = "0"
				return
			}
			uinfo.Family = "1"
		}
		_, err = a.DB.Exec("UPDATE informations SET gender = ? WHERE userName = ?", uBody.Gender, userName)
		if err != nil {
			return
		}
		viewer(w, http.StatusOK, uinfo)
	}
}
func (a *App) editUsername(w http.ResponseWriter, r *http.Request) {
	link := strings.ToLower(strings.TrimPrefix(r.URL.Path, "/editusername/"))
	fmt.Println(r.Header.Get("name"))
	fmt.Println(r.Header.Get("tag"))
	//vars := mux.Vars(r)
	//id := vars["id"]
	userName := strings.ToLower(link)
	rows, err := a.DB.Query("SELECT userName, phoneNumber FROM users WHERE userName = ?", userName)
	var uinfo UserInformation
	if err != nil {
		//errViewer(w, http.StatusInternalServerError, "userName : "+userName+" "+err.Error())
		return
	} else {
		//viewer(w, http.StatusOK, "userName : "+userName)
		var uBody, exists, u UserInformation

		for rows.Next() {
			err = rows.Scan(&u.PhoneNum, &u.Username)
		}
		if err != nil {
			//errViewer(w, http.StatusInternalServerError, ""+err.Error())
			return
		}
		js, err1 := ioutil.ReadAll(r.Body)

		if err1 != nil {
			//errViewer(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := json.Unmarshal(js, &uBody); err != nil {
			return
		}
		user := strings.ToLower(uBody.Username)

		//start replacement contents
		query, err := a.DB.Query("SELECT userName FROM users WHERE userName = ?", user)
		if err != nil {
			return
		}
		for query.Next() {
			query.Scan(&exists.Username)
		}
		if !strings.EqualFold(exists.Username, user) {
			if user != "" {
				if !strings.EqualFold(u.Username, user) {
					_, err := a.DB.Exec("UPDATE users SET userName = ? WHERE userName = ?", user, userName)
					if err != nil {
						uinfo.Username = "0"
						return
					}
					_, err = a.DB.Exec("UPDATE informations SET userName = ? WHERE userName = ?", user, userName)
					if err != nil {
						uinfo.Username = "0"
						return
					}
					_, err = a.DB.Exec("UPDATE payments SET userName = ? WHERE userName = ?", user, userName)
					if err != nil {
						uinfo.Username = "0"
						return
					}
					uinfo.Username = "1"
				}
			}
		} else {
			uinfo.Username = "2"
		}
		viewer(w, http.StatusOK, uinfo)
	}
}
func (a *App) Login(writer http.ResponseWriter, request *http.Request) {
	path := strings.TrimPrefix(request.URL.Path, "/login/")
	u, err := getUser(a.DB, path)
	if err != nil {
		fmt.Println("4: ", err.Error())
		return
	}
	viewer(writer, http.StatusOK, u)
}
func (a *App) Information(writer http.ResponseWriter, request *http.Request) {
	path := request.Header.Get("User-name")
	rows, err := a.DB.Query("SELECT userName,title,family,biography,gender,profile FROM informations WHERE userName = ?", path)
	if err != nil {
		fmt.Println("ui 2", err.Error())
		return
	}
	defer rows.Close()
	var ui UserInformation
	for rows.Next() {
		err := rows.Scan(&ui.Username, &ui.Title, &ui.Family, &ui.Biography, &ui.Gender, &ui.Profile)
		if err != nil {
			fmt.Println("ui 3", err.Error())
			return
		}
	}
	body, err := ioutil.ReadFile("./nitrogen/images/" + ui.Profile)
	ui.Byte = body
	viewer(writer, http.StatusOK, ui)
	//f, err := os.Open("D://nitrogen/images/" + ui.Profile)
	//if err != nil {
	//	log.Fatal(err.Error())
	//return
	//}
	//fmt.Println(body)

}
func (a *App) Challenge(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	challengeID := vars["challengeID"]
	rows, err := a.DB.Query("SELECT * FROM challenges WHERE challengeID = ?", challengeID)
	if err != nil {
		fmt.Println("challenge 2", err.Error())
		return
	}
	var c Challenge
	for rows.Next() {
		err := rows.Scan(&c.Id, &c.ChallengeID, &c.Title, &c.SubTitle,
			&c.Description, &c.SubDescription, &c.Image, &c.DayCount,
			&c.BgColor, &c.TxtColor)
		if err != nil {
			fmt.Println("challenge 3", err.Error())
			return
		}
		err = c.getTags(a.DB)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = c.getCategory(a.DB)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = c.getInfluences(a.DB)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = c.getConclusion(a.DB)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = c.getPrerequisite(a.DB)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = c.getTool(a.DB)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		c.ChallengeByCategory(a.DB)
		err = c.getStarted(a.DB, request.Header.Get("User-name"))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		c.Price, _ = getPrice(a.DB, challengeID)
	}
	viewer(writer, http.StatusOK, c)
}
func (c *Challenge) ChallengeByCategory(db *sql.DB) {
	rows, _ := db.Query("SELECT title,challengeID FROM categories "+
		"WHERE challengeID = ?  GROUP BY title",
		c.ChallengeID)
	for rows.Next() {
		var category Category
		rows.Scan(&category.Title, &category.ChallengeID)
		ch, _ := getChallengeByCategory(db, category.Title)
		c.Suggestion = ch
	}
}
func (a *App) Challenges(writer http.ResponseWriter, request *http.Request) {
	var cArr []Challenge
	rows, err := a.DB.Query("SELECT * FROM challenges")
	if err != nil {
		fmt.Println("challenge 2", err.Error())
		return
	}
	for rows.Next() {
		var c Challenge
		err := rows.Scan(&c.Id, &c.ChallengeID, &c.Title, &c.SubTitle,
			&c.Description, &c.SubDescription, &c.Image, &c.DayCount,
			&c.BgColor, &c.TxtColor)
		if err != nil {
			fmt.Println("challenge 3", err.Error())
			return
		}
		err = c.getTags(a.DB)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = c.getCategory(a.DB)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = c.getInfluences(a.DB)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = c.getConclusion(a.DB)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		cArr = append(cArr, c)
	}
	viewer(writer, http.StatusOK, cArr)
}
func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	var h Home
	rows, err := a.DB.Query("SELECT * FROM banners")
	if err != nil {
		return
	}
	defer rows.Close()
	h.fillNewchallenges(a.DB)
	h.Success, _ = getChallengeByCategory(a.DB, "موفقیت")
	h.Psychology, _ = getChallengeByCategory(a.DB, "خودساخته")
	h.Selfdevelopment, _ = getChallengeByCategory(a.DB, "توسعه فردی")
	h.Suggestion5, _ = getChallengeByCategory(a.DB, "کارآفرینی")
	var c Challenge
	h.Banners, _ = c.readingByChallengeId(a.DB, rows)
	viewer(w, http.StatusOK, h)
}
func (h *Home) fillNewchallenges(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM challenges ORDER BY id DESC")
	if err != nil {
		fmt.Println("challenge 2: \n", err.Error())
		return err
	}
	defer rows.Close()
	var c Challenge
	var cs []Challenge
	for rows.Next() {
		rows.Scan(&c.Id, &c.ChallengeID, &c.Title, &c.SubTitle,
			&c.Description, &c.SubDescription, &c.Image, &c.DayCount,
			&c.BgColor, &c.TxtColor)
		cs = append(cs, c)
	}
	h.Newchallenges = cs
	return nil
}
func getChallengeByCategory(db *sql.DB, str string) ([]Challenge, error) {
	rows, err := db.Query("SELECT title,challengeID FROM categories WHERE title = ? order by challengeID DESC", str)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []Category
	for rows.Next() {
		var category Category
		rows.Scan(&category.Title, &category.ChallengeID)
		categories = append(categories, category)

	}
	var challenges []Challenge
	for i := 0; i <= (len(categories) - 1); i++ {
		c, _ := getChallengeByID(db, categories[i].ChallengeID)
		challenges = append(challenges, c)

	}
	return challenges, nil
}
func getChallengeByID(db *sql.DB, id string) (Challenge, error) {
	row, err := db.Query("SELECT * FROM challenges WHERE challengeID = ?", id)
	if err != nil {
		return Challenge{}, err
	}
	defer row.Close()
	var c Challenge
	for row.Next() {
		err := row.Scan(&c.Id, &c.ChallengeID, &c.Title, &c.SubTitle,
			&c.Description, &c.SubDescription, &c.Image, &c.DayCount,
			&c.BgColor, &c.TxtColor)
		if err != nil {
			fmt.Println("challenge 3", err.Error())
			return Challenge{}, err
		}
		c.Tags = []Tag{}
		c.Category = []Category{}
		c.Influence = []Influence{}
		c.Conclusion = []Conclusion{}
	}
	return c, nil
}
func getChallengeByID2(db *sql.DB, id string) (Challenge2, error) {
	row, err := db.Query("SELECT * FROM challenges WHERE challengeID = ?", id)
	if err != nil {
		return Challenge2{}, err
	}
	defer row.Close()
	var c Challenge2
	for row.Next() {
		err := row.Scan(&c.Id, &c.ChallengeID, &c.Title, &c.SubTitle,
			&c.Description, &c.SubDescription, &c.Image, &c.DayCount,
			&c.BgColor, &c.TxtColor)
		c.getCategory2(db)
		if err != nil {
			fmt.Println("challenge 3", err.Error())
			return Challenge2{}, err
		}
	}

	return c, nil
}
func (a *App) Suggestion(w http.ResponseWriter, r *http.Request) {
	c, err := getSuggestion(a.DB)
	if err != nil {
		return
	}
	viewer(w, http.StatusOK, c)
}
func (c *Challenge) getTags(db *sql.DB) error {
	path := c.ChallengeID
	rows, err := db.Query("SELECT  challengeID,title FROM tags WHERE challengeID = ?", path)
	if err != nil {
		fmt.Println("Tags 2", err.Error())
		return err
	}
	for rows.Next() {
		var t Tag
		err := rows.Scan(&t.ChallengeID, &t.Title)
		if err != nil {
			fmt.Println("Tags 3", err.Error())
			return err
		}
		c.Tags = append(c.Tags, t)
	}
	var i Tag
	if c.Tags == nil {
		c.Tags = append(c.Tags, i)
	}
	return nil
}
func (c *Challenge) getStarted(db *sql.DB, username string) error {
	path := c.ChallengeID
	rows, err := db.Query("SELECT userName, challengeID,day, month,year, dayOfYear,challengeState,priceType FROM startDate WHERE challengeID = ? AND userName = ?", path, username)
	if err != nil {
		fmt.Println("started 2", err.Error())
		return err
	}
	var s Started
	for rows.Next() {
		err := rows.Scan(&s.Username, &s.ChallengeID, &s.Day, &s.Month, &s.Year, &s.DayOfYear, &s.State, &s.PriceType)
		if err != nil {
			fmt.Println("started 3", err.Error())
			return nil
		}
	}
	c.Started = s
	if s.Username == "" {
		var s Started
		c.Started = s
	}
	return nil
}
func (c *Challenge) getCategory(db *sql.DB) error {
	path := c.ChallengeID
	rows, err := db.Query("SELECT  challengeID,title FROM categories WHERE challengeID = ?", path)
	if err != nil {
		fmt.Println("Category 2", err.Error())
		return err
	}
	for rows.Next() {
		var cy Category
		err := rows.Scan(&cy.ChallengeID, &cy.Title)
		if err != nil {
			fmt.Println("Category 3", err.Error())
			return err
		}
		c.Category = append(c.Category, cy)
	}
	var i Category
	if c.Category == nil {
		c.Category = append(c.Category, i)
	}
	return nil
}
func (c *Challenge2) getCategory2(db *sql.DB) error {
	path := c.ChallengeID
	rows, err := db.Query("SELECT  challengeID,title FROM categories WHERE challengeID = ?", path)
	if err != nil {
		fmt.Println("Category 2", err.Error())
		return err
	}
	for rows.Next() {
		var cy Category
		err := rows.Scan(&cy.ChallengeID, &cy.Title)
		if err != nil {
			fmt.Println("Category 3", err.Error())
			return err
		}
		c.Category = append(c.Category, cy)
	}
	var i Category
	if c.Category == nil {
		c.Category = append(c.Category, i)
	}
	return nil
}
func (c *Challenge) getInfluences(db *sql.DB) error {
	path := c.ChallengeID
	rows, err := db.Query("SELECT  challengeID,title FROM influences WHERE challengeID = ?", path)
	if err != nil {
		fmt.Println("Influences 2", err.Error())
		return err
	}
	for rows.Next() {
		var cy Influence

		err := rows.Scan(&cy.ChallengeID, &cy.Title)
		if err != nil {
			fmt.Println("Influences 3", err.Error())
			return err
		}
		c.Influence = append(c.Influence, cy)
	}
	var i Influence
	if c.Influence == nil {
		c.Influence = append(c.Influence, i)
	}

	return nil
}
func (c *Challenge) getConclusion(db *sql.DB) error {
	path := c.ChallengeID
	rows, err := db.Query("SELECT  challengeID,title,description FROM conclusions WHERE challengeID = ?", path)
	if err != nil {
		fmt.Println("Conclusion 2", err.Error())
		return err
	}
	for rows.Next() {
		var cy Conclusion
		err := rows.Scan(&cy.ChallengeID, &cy.Title, &cy.Description)
		if err != nil {
			fmt.Println("Conclusion 3", err.Error())
			return err
		}
		c.Conclusion = append(c.Conclusion, cy)
	}
	var i Conclusion
	if c.Conclusion == nil {
		c.Conclusion = append(c.Conclusion, i)
	}
	return nil
}
func (c *Challenge) getPrerequisite(db *sql.DB) error {
	path := c.ChallengeID
	rows, err := db.Query("SELECT  challengeID,prerequisiteID FROM prerequisites WHERE challengeID = ?", path)
	if err != nil {
		return err
	}
	for rows.Next() {
		var pe Prerequisite
		err := rows.Scan(&pe.ChallengeID, &pe.PrerequisiteID)
		if err != nil {
			fmt.Println("Prerequisite 3", err.Error())
			return err
		}
		row, err := db.Query("SELECT * FROM challenges WHERE challengeID = ?", pe.PrerequisiteID)
		var challenge Challenge
		for row.Next() {
			row.Scan(
				&challenge.Id, &challenge.ChallengeID, &challenge.Title, &challenge.SubTitle,
				&challenge.Description, &challenge.SubDescription, &challenge.Image, &challenge.DayCount,
				&challenge.BgColor, &challenge.TxtColor)
		}
		c.Prerequisite = append(c.Prerequisite, challenge)
	}
	var i Challenge
	if c.Prerequisite == nil {
		c.Prerequisite = append(c.Prerequisite, i)
	}
	return nil
}
func (c *Challenge) getTool(db *sql.DB) error {
	path := c.ChallengeID
	rows, err := db.Query("SELECT challengeID,toolID,title,description FROM tools WHERE challengeID = ?", path)
	if err != nil {
		fmt.Println("tool: 1", err.Error())
		return err
	}
	defer rows.Close()
	var alt []Tools
	for rows.Next() {
		var t Tool
		if err := rows.Scan(&t.ChallengeID, &t.ToolID, &t.Title, &t.Description); err != nil {
			fmt.Println("tool: 1", err.Error())
			return err
		}
		c.Tool.ChallengeID = t.ChallengeID
		c.Tool.Title = t.Title
		c.Tool.ToolID = t.ToolID
		c.Tool.Description = t.Description
		row, err := db.Query("SELECT toolID ,name FROM tool WHERE toolID = ?", c.Tool.ToolID)
		if err != nil {
			return err
		}
		for row.Next() {
			var tool Tools
			row.Scan(&tool.ToolID, &tool.Name)
			alt = append(alt, tool)
		}
	}
	c.Tool.Tools = alt
	//c.Tool.Tools = alt
	return nil
}
func getPrice(db *sql.DB, challengeID string) (Price, error) {
	rows, err := db.Query("SELECT priceID,challengeID,title,type,price,sku FROM prices WHERE challengeID = ?", challengeID)
	if err != nil {
		return Price{}, err
	}
	defer rows.Close()
	var p Price
	for rows.Next() {
		rows.Scan(&p.PriceID, &p.ChallengeID, &p.Title, &p.Type, &p.Price, &p.Sku)
	}
	return p, nil
}
func getSuggestion(db *sql.DB) ([]Challenge, error) {
	//path := c.ChallengeID
	var c []Challenge
	for i := 0; i <= 7; i++ {
		row, err := db.Query("SELECT * FROM challenges WHERE id = ?", i)
		if err != nil {
			return nil, err
		}
		var challenge Challenge
		for row.Next() {
			row.Scan(
				&challenge.Id, &challenge.ChallengeID, &challenge.Title, &challenge.SubTitle,
				&challenge.Description, &challenge.SubDescription, &challenge.Image, &challenge.DayCount,
				&challenge.BgColor, &challenge.TxtColor)

			c = append(c, challenge)
		}
	}
	var i Challenge
	if c == nil {
		c = append(c, i)
	}
	return c, nil
}
func (a *App) started(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/started/")
	username := r.Header.Get("User-name")
	var s Started

	rows, err := a.DB.Query("SELECT userName, challengeID,day, month,year, dayOfYear,challengeState,dayCount,priceType FROM startDate WHERE challengeID = ? AND userName = ?", path, username)
	challenge, _ := getChallengeByID(a.DB, path)
	if err != nil {
		fmt.Println("started 2", err.Error())
		return
	}
	for rows.Next() {
		err := rows.Scan(&s.Username, &s.ChallengeID, &s.Day, &s.Month, &s.Year, &s.DayOfYear, &s.State, &s.DayCount, &s.PriceType)
		if err != nil {
			fmt.Println("started 3", err.Error())
			return
		}
	}

	if r.Header.Get("State-value") == "1" && s.Username == r.Header.Get("User-name") {
		fmt.Println("this row is exists")
		return
	}
	if r.Header.Get("State-value") == "1" {
		currentTime := time.Now()
		// Create a new instance of time.Time
		var t time.Time = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second(), 0, ptime.Iran())
		// Create a new instance of time.Time
		// Get a new instance of ptime.Time using time.Time
		pt := ptime.New(t)
		// Get the date in Persian calendar
		fmt.Println(pt.Day())
		a.DB.Exec("INSERT INTO startDate(userName, challengeID,day, month,year, dayOfYear,challengeState,dayCount) VALUES(?,?,?,?,?,?,?,?)", username, path, pt.Day(), pt.Month(), pt.Year(), pt.YearDay(), r.Header.Get("State-value"), challenge.DayCount)
		viewer(w, http.StatusOK, "enabled")
	} else if r.Header.Get("State-value") == "0" {
		a.DB.Exec("DELETE FROM startDate WHERE userName = ? AND challengeID = ?", username, path)
		viewer(w, http.StatusOK, "deleted")
	}
}
func checkPrerequesite(db *sql.DB, username string, challengeID string) (string, error) {
	rows, err := db.Query("SELECT  challengeID,prerequisiteID FROM prerequisites WHERE challengeID = ? AND userName = ?", challengeID, username)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		var pe Prerequisite
		err := rows.Scan(&pe.ChallengeID, &pe.PrerequisiteID)
		if err != nil {
			fmt.Println("Prerequisite 3", err.Error())
			return "", err
		}
		row, err := db.Query("SELECT * FROM challenges WHERE challengeID = ?", pe.PrerequisiteID)
		var challenge Challenge
		for row.Next() {
			row.Scan(
				&challenge.Id, &challenge.ChallengeID, &challenge.Title, &challenge.SubTitle,
				&challenge.Description, &challenge.SubDescription, &challenge.Image, &challenge.DayCount,
				&challenge.BgColor, &challenge.TxtColor)
		}
		//c.Prerequisite = append(c.Prerequisite, challenge)
	}
	//var i Challenge
	//if c.Prerequisite == nil {
	//	c.Prerequisite = append(c.Prerequisite, i)
	//}
	return "", nil
}
func (a *App) banners(w http.ResponseWriter, r *http.Request) {
	rows, err := a.DB.Query("SELECT * FROM banners")
	if err != nil {
		return
	}
	defer rows.Close()
	var c Challenge
	get, _ := c.readingByChallengeId(a.DB, rows)
	viewer(w, http.StatusOK, get)
}
func (c *Challenge) readingByChallengeId(db *sql.DB, rows *sql.Rows) ([]Challenge, error) {
	var challenges []Challenge
	for rows.Next() {
		var b Banner
		rows.Scan(&b.Id, &b.ChallengeID)

		row, err := db.Query("SELECT * FROM challenges WHERE challengeID = ?", b.ChallengeID)
		if err != nil {
			return nil, err
		}
		for row.Next() {
			var c Challenge
			row.Scan(&c.Id, &c.ChallengeID, &c.Title, &c.SubTitle, &c.Description, &c.SubDescription, &c.Image, &c.DayCount, &c.BgColor,
				&c.TxtColor)
			challenges = append(challenges, c)
		}
	}
	return challenges, nil
}
func (a *App) EnabledChallenges(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("User-name")
	rows, err := a.DB.Query("SELECT * FROM startDate WHERE userName = ? ", username)
	if err != nil {
		return
	}
	defer rows.Close()
	var started []Started
	//var challenges []Challenge
	for rows.Next() {
		var s Started
		rows.Scan(&s.Id, &s.Username, &s.ChallengeID, &s.Day, &s.Month, &s.Year, &s.DayOfYear, &s.State, &s.DayCount, &s.PriceType)
		s.Progress = progress(s.DayOfYear)
		c, _ := getChallengeByID2(a.DB, s.ChallengeID)

		s.Challenge = c
		started = append(started, s)
	}
	viewer(w, http.StatusOK, started)
}
func (a *App) EnabledChallenge(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("User-name")
	challengeID := r.Header.Get("Challenge-id")
	rows, err := a.DB.Query("SELECT * FROM startDate WHERE userName=? AND challengeID=?", username, challengeID)
	if err != nil {
		return
	}
	defer rows.Close()
	var s Started
	for rows.Next() {
		rows.Scan(&s.Id, &s.Username, &s.ChallengeID, &s.Day, &s.Month, &s.Year, &s.DayOfYear, &s.State, &s.DayCount, &s.PriceType)
		s.Progress = progress(s.DayOfYear)
		c, _ := getChallengeByID2(a.DB, s.ChallengeID)
		s.Challenge = c
	}
	viewer(w, http.StatusOK, s)
}

type sendDate struct {
	Date     string `json:"date"`
	Year     string `json:"year"`
	Startday int    `json:"startDay"`
}

func (a *App) ActionGoal(writer http.ResponseWriter, request *http.Request) {
	username := request.Header.Get("User-name")
	goalid := request.Header.Get("Goal-id")
	item := request.Header.Get("Item")
	var g goal
	g.UserName = username

	if strings.Contains(item, "Add-goal") {
		js, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return
		}
		err = json.Unmarshal(js, &g)
		if err != nil {
			fmt.Println("err: ", err.Error())
			return
		}
		currentTime := time.Now()
		var t = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second(), 0, ptime.Iran())
		pt := ptime.New(t)
		_, err = a.DB.Exec("INSERT INTO goals(goalID,userName,title,description,priority,color,startDay,dayCount,life,updat) VALUES (?,?,?,?,?,?,?,?,?,?)",
			g.GoalID, g.UserName, g.Title, g.Description, g.Priority, g.Color, pt.YearDay(), g.DayCount, g.Life, g.Update)
		if err != nil {
			fmt.Println("err msg: ", err.Error())
			return
		}
		var s sendDate
		s.Startday = pt.YearDay()
		s.Date = strconv.Itoa(pt.Year()) + "/" + strconv.Itoa(int(pt.Month())) + "/" + strconv.Itoa(pt.Day())
		s.Year = strconv.Itoa(pt.Year()) + "/" + strconv.Itoa(int(pt.Month())) + "/" + strconv.Itoa(pt.Day())
		fmt.Println("it`s ok added")
		viewer(writer, http.StatusOK, s)
	} else if strings.Contains(item, "Edit-goal") {
		js, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return
		}
		err = json.Unmarshal(js, &g)
		if err != nil {
			fmt.Println("err: ", err.Error())
			return
		}
		_, err = a.DB.Exec("UPDATE goals SET title=?,description=?,priority=?,color=?,dayCount=?,life=? WHERE userName = ? AND goalID = ?",
			g.Title, g.Description, g.Priority, g.Color, g.DayCount, g.Life, g.UserName, goalid)
		if err != nil {
			fmt.Println("err msg: ", err.Error())
			return
		}
		fmt.Println("it`s ok edited")
	} else if strings.Contains(item, "Delete-goal") {
		_, err := a.DB.Exec("DELETE FROM goals WHERE userName = ? AND goalID = ?", username, goalid)
		if err != nil {
			fmt.Println("err msg: ", err.Error())
			return
		}
		fmt.Println("it`s ok deleted")
	}
}
func (a *App) EditImage(writer http.ResponseWriter, request *http.Request) {
	js, err := ioutil.ReadAll(request.Body)
	username := request.Header.Get("User-name")
	var img ImageEditor
	if err != nil {
		fmt.Println("msg: ", err.Error())
		return
	}
	err = json.Unmarshal(js, &img)
	if err != nil {
		fmt.Println("msg: ", err.Error())
		return
	}
	currentTime := time.Now()
	var t = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second(), 0, ptime.Iran())
	pt := ptime.New(t)
	fmt.Println(username)

	m := "IMAGE" + strconv.Itoa(rand.Int()) + strconv.Itoa(pt.Year()) + "" + strconv.Itoa(int(pt.Month())) + "" + strconv.Itoa(pt.Day()) + "" + strconv.Itoa(pt.Hour()) + "" + strconv.Itoa(pt.Minute()) + "" + strconv.Itoa(pt.Second()) + ".jpeg"

	a.DB.Exec("UPDATE informations SET profile = ? WHERE userName = ?", m, username)

	f, err := os.Create("./nitrogen/images/" + m)
	if err != nil {
		fmt.Println("err msg: ", err.Error())
	}
	textdecodeing, _ := base64.StdEncoding.DecodeString(img.Image)
	f.Write([]byte(textdecodeing))
	f.Close()
}
func (a *App) writeLocation(username string, l Location) {
	row, err := a.DB.Query("SELECT userName FROM location WHERE userName = ?", username)
	if err != nil {
		fmt.Println("nashod 1: ", err.Error())
		return
	}
	defer row.Close()
	var l1 Location

	fmt.Println("msg: ", l.CountryCode)
	for row.Next() {
		err := row.Scan(&l1.Username)
		if err != nil {
			fmt.Println("nashod 2: ", err.Error())
			return
		}
	}
	if l1.Username == username {
		_, err := a.DB.Exec("UPDATE location SET country=?,countryCode=?,region=?,regionName=?,city=?,lat=?,lon=?,timezone=?,isp=?,org=?,ass=?,query=? WHERE userName=?", l.Country, l.CountryCode, l.Region, l.RegionName, l.City, l.Lat, l.Lon, l.Timezone, l.Isp, l.Org, l.As, l.Query, username)
		if err != nil {
			fmt.Println("nashod 3: ", err.Error())
			return
		}
	} else {
		_, err = a.DB.Exec("INSERT INTO location(userName,country,countryCode,region,regionName,city,lat,lon,timezone,isp,org,ass,query) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)", username, l.Country, l.CountryCode, l.Region, l.RegionName, l.City, l.Lat, l.Lon, l.Timezone, l.Isp, l.Org, l.As, l.Query)
		if err != nil {
			fmt.Println("nashod 4: ", err.Error())
			return
		}
	}
}
func (a *App) Location(writer http.ResponseWriter, request *http.Request) {
	username := request.Header.Get("User-name")
	bt, _ := ioutil.ReadAll(request.Body)
	var l Location
	json.Unmarshal(bt, &l)
	a.writeLocation(username, l)
}
