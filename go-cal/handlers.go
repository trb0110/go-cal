package gocal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var userControl *UserControl

func UserGET(c *gin.Context) {
	token := c.MustGet("token").(string)
	fmt.Println(token)
	userId := c.Params.ByName("userid")
	user := userControl.myconnection.DBUserGet(userId)

	c.JSON(http.StatusOK, gin.H{"user": user})
}
func Login(c *gin.Context) {
	var creds User
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := userControl.myconnection.DBUserLogin(creds)
	if(len(user.UserID)==0){
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong username or password"})
	}else{
		tokenK,err := GenerateToken(user)
		if err !=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		//w.Header().Set("Token", tokenK)
		user.Password=""
		user.Token=tokenK

		c.JSON(http.StatusOK, gin.H{"user": user})
	}


}


func CalorieGET(c *gin.Context) {
	userId := c.Query("userid")
	calorieId := c.Query("calorieid")
	startdate := c.Query("startdate")
	enddate := c.Query("enddate")

	token := c.MustGet("token").(string)
	session,_:= TokensMap.Get(token)
	//fmt.Println(session)
	if(session.Role=="1"){
		calories := userControl.myconnection.DBCaloriesGetBulk(startdate,enddate)
		c.JSON(http.StatusOK, gin.H{"calories": calories})
	}else{
		if(len(userId)>0 && len(calorieId)==0){
			if _, err := strconv.Atoi(userId); err == nil {
				if(session.userid==userId){
					calories := userControl.myconnection.DBCalorieUserGetBulk(userId,startdate,enddate)
					c.JSON(http.StatusOK, gin.H{"calories": calories})
				}else{
					c.String(http.StatusUnauthorized, "unauthorized")
				}
			}else{
				c.String(http.StatusBadRequest, "Please pass a valid param")
			}
		}else if(len(calorieId)>0 && len(userId)>0){
			if _, err := strconv.Atoi(calorieId); err == nil {
				if(session.userid==userId) {
					calorie := userControl.myconnection.DBCalorieGetOne(calorieId)
					c.JSON(http.StatusOK, gin.H{"calories": calorie})
				}else{
					c.String(http.StatusUnauthorized, "unauthorized")
				}
			}else{
				c.String(http.StatusBadRequest, "Please pass a valid param")
			}
		}else{
			c.String(http.StatusBadRequest, "Please pass a valid param")
		}
	}

}
func CaloriePost(c *gin.Context) {
	var calorie Calorie
	if err := c.ShouldBindJSON(&calorie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp,_ := userControl.myconnection.DBCalorieCreate(calorie)

	c.JSON(resp, gin.H{"calories": calorie})
}
func CalorieDelete(c *gin.Context) {
	calorieid := c.Params.ByName("calorieid")
	resp,_ := userControl.myconnection.DBCalorieDelete(calorieid)
	c.String(resp, "%d",resp)
}
func CalorieUpdate(c *gin.Context) {
	var calorie Calorie
	if err := c.ShouldBindJSON(&calorie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp,_ := userControl.myconnection.DBCalorieUpdate(calorie)
	c.JSON(resp, gin.H{"calories": calorie})
}