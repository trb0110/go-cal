package gocal

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	_ "github.com/denisenkom/go-mssqldb"
	"net/http"
)

var server = flag.String("mssql", "127.0.0.1", "the database server")
var port  = flag.Int("port", 1433, "the database port")
var user = flag.String("user", "sa", "the database user")
var password = flag.String("password", "R13032andompassword", "the database password")
var database = flag.String("database", "countcal", "the database name")

type UserControl struct {
	myconnection *SQLConnection
}

func NewUserControl() {
	UC := &UserControl{
		myconnection: NewDBConnection(),
	}
	userControl=UC
}

type SQLConnection struct {
	originalSession *sql.DB
}

func NewDBConnection() (conn *SQLConnection) {
	conn = new(SQLConnection)
	conn.createLocalConnection()
	return
}


func (c *SQLConnection) createLocalConnection() (err error) {
	log.Println("Connecting to SQL DB server....")

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", *server, *user, *password, *port, *database)
	c.originalSession, err = sql.Open("mssql", connString)
	if err == nil {
		log.Println("Connection established to SQL DB server")
		return
	} else {
		log.Println("Error occured while creating SQL DB connection: %s", err.Error())
	}
	// Set maximum number of connections in idle connection pool.
	c.originalSession.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	c.originalSession.SetMaxOpenConns(15)
	return
}

func (c *SQLConnection)  DBUserLogin(creds User)(user *User){

	strSql := `SELECT  [user_id], username,[role_id] ,[password],[calorie_limit] FROM [countcal].dbo.[user] where [username] = '`+creds.Username+`'
																							and [password] = '`+creds.Password+`'`

	queryResult, rerr := c.originalSession.Query(strSql)
	defer queryResult.Close()
	if rerr != nil {
		log.Println(rerr)
	}

	user = &User{}
	for queryResult.Next() {
		err := queryResult.Scan(&user.UserID,&user.Username,&user.Role,&user.Password,&user.CalorieLimit)
		if err!=nil{
			//fmt.Println("Check UserName Error	" , err2)
			if err == sql.ErrNoRows{
				return nil
			}

			return nil
		}
	}
	defer queryResult.Close()
	return user
}
func (c *SQLConnection)  DBUserGet(id string)(user *User){

	strSql := `SELECT  [user_id], username,[role_id] ,[password] FROM [countcal].dbo.[user] where [user_id] = `+id+``

	queryResult, rerr := c.originalSession.Query(strSql)
	defer queryResult.Close()
	if rerr != nil {
		log.Println(rerr)
	}

	user = &User{}
	for queryResult.Next() {
		err := queryResult.Scan(&user.UserID,&user.Username,&user.Role,&user.Password)
		if err != nil {
			log.Println(err)
		}
	}
	defer queryResult.Close()
	return user
}


func (c *SQLConnection)  DBCalorieGetOne(calorieId string)(calorie *Calorie){

	strSql := `SELECT  [calorie_id], [user_id],[calorie_stamp] ,[food],[calorie_count] FROM [countcal].dbo.[calorie] where [calorie_id] = `+calorieId+``


	queryResult, rerr := c.originalSession.Query(strSql)
	defer queryResult.Close()
	if rerr != nil {
		log.Println(rerr)
	}

	calorie = &Calorie{}
	for queryResult.Next() {
		err := queryResult.Scan(&calorie.CalorieId,&calorie.UserID,&calorie.Timestamp,&calorie.Food,&calorie.CalorieCount)
		if err != nil {
			log.Println(err)
		}
	}
	defer queryResult.Close()
	return calorie
}

func (c *SQLConnection)  DBCaloriesGetBulk(startdate string, enddate string)(calorie []Calorie){

	strSql := `SELECT  [calorie_id], [user_id],[calorie_stamp] ,[food],[calorie_count],
 					(select [username]from [countcal].dbo.[user] where[user_id]=t.[user_id]) as [username] 
 
 				FROM [countcal].dbo.[calorie] t `
	if(len(startdate)>0&&len(enddate)==0){
		strSql+= ` where calorie_stamp>'`+startdate+`'`
	}
	if(len(enddate)>0&&len(startdate)==0){
		strSql+= ` where calorie_stamp<'`+enddate+`'`
	}
	if(len(enddate)>0&&len(startdate)>0){
		strSql+= ` where calorie_stamp>'`+startdate+`' and calorie_stamp<'`+enddate+`'`
	}
	strSql+=`order by calorie_stamp asc`

	queryResult, rerr := c.originalSession.Query(strSql)
	defer queryResult.Close()
	if rerr != nil {
		log.Println(rerr)
	}

	cals := make([]Calorie,0)
	for queryResult.Next(){
		calorie := Calorie{}
		err := queryResult.Scan(&calorie.CalorieId,&calorie.UserID,&calorie.Timestamp,&calorie.Food,&calorie.CalorieCount,&calorie.Username)
		if err!=nil{
			if err == sql.ErrNoRows{
				return nil
			}
			return nil
		}
		cals= append(cals,calorie)
	}

	defer queryResult.Close()
	return cals
}
func (c *SQLConnection)  DBCalorieUserGetBulk(userId string, startdate string, enddate string)(calorie []Calorie){

	strSql := `SELECT  [calorie_id], [user_id],[calorie_stamp] ,[food],[calorie_count] FROM [countcal].dbo.[calorie] where [user_id] = `+userId+``
	if(len(startdate)>0&&len(enddate)==0){
		strSql+= ` and calorie_stamp>'`+startdate+`'`
	}
	if(len(enddate)>0&&len(startdate)==0){
		strSql+= ` and calorie_stamp<'`+enddate+`'`
	}
	if(len(enddate)>0&&len(startdate)>0){
		strSql+= ` and calorie_stamp>'`+startdate+`' and calorie_stamp<'`+enddate+`'`
	}

	strSql+=`order by user_id,calorie_stamp asc`

	queryResult, rerr := c.originalSession.Query(strSql)
	defer queryResult.Close()
	if rerr != nil {
		log.Println(rerr)
	}

	cals := make([]Calorie,0)
	for queryResult.Next(){
		calorie := Calorie{}
		err := queryResult.Scan(&calorie.CalorieId,&calorie.UserID,&calorie.Timestamp,&calorie.Food,&calorie.CalorieCount)
		if err!=nil{
			if err == sql.ErrNoRows{
				return nil
			}
			return nil
		}
		cals= append(cals,calorie)
	}

	defer queryResult.Close()
	return cals
}


func (c *SQLConnection)  DBCalorieCreate(calorie Calorie)(response int, err error){

	timestamp := ""
	if(len(calorie.Timestamp)>0){
		timestamp="'"+calorie.Timestamp+"'"
	}else{
		timestamp="getdate()"
	}
	strSql := `
 				insert into [countcal].[dbo].[calorie] ([user_id],[food],[calorie_count],[calorie_stamp])
  				values ( `+ calorie.UserID+`,'`+ calorie.Food+`',`+ calorie.CalorieCount+`,`+timestamp+`)
				`

	queryResult, rerr := c.originalSession.Query(strSql)
	if rerr != nil {
		return http.StatusInternalServerError, rerr
	}

	defer queryResult.Close()
	return http.StatusOK , nil
}


func (c *SQLConnection)  DBCalorieUpdate(calorie Calorie)(response int, err error){

	strSql := `				update [countcal].[dbo].[calorie]
							set  
				`

	if calorie.Food!="" {
		strSql += `[food] = '`+calorie.Food+`',`
	}

	if calorie.CalorieCount!="" {
		strSql += `[calorie_count] = '`+calorie.CalorieCount+`',`
	}
	if calorie.Timestamp!="" {
		strSql += `[calorie_stamp] = '`+calorie.Timestamp+`',`
	}

	strSql = strSql[:len(strSql)-1]
	strSql+= ` where [calorie_id] = `+calorie.CalorieId+``

	queryResult, rerr := c.originalSession.Query(strSql)

	if rerr != nil {
		return http.StatusInternalServerError, rerr
	}
	defer queryResult.Close()
	return http.StatusOK , nil
}


func (c *SQLConnection)  DBCalorieDelete(id string)(response int, err error){

	strSql := `DELETE [countcal].[dbo].[calorie] where [calorie_id] = `+id+``

	queryResult, rerr := c.originalSession.Query(strSql)
	if rerr != nil {
		//log.Println("Error DBRegister", rerr)
		//printSqlToLog(strSql)
		return http.StatusInternalServerError, rerr
	}

	defer queryResult.Close()
	return http.StatusOK , nil
}
