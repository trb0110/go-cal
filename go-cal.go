package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	. "gocal/go-cal"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var db = make(map[string]string)
func main() {

	NewUserControl()
	InitializeTokens()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
	}))

	authorized.POST("/login", Login)
	authorized.Use(Auth())
	{
		authorized.GET("/calorie", CalorieGET)
		authorized.POST("/calorie", CaloriePost)
		authorized.GET("/user/:userid", UserGET)
	}

	authorized.Use(AdminAuth())
	{
		authorized.DELETE("/calorie/:calorieid", CalorieDelete)
		authorized.PUT("/calorie", CalorieUpdate)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer:=c.Request.Header["Bearer"]
		if(len(bearer)>0){
			accessExists:=LookupTokenKey(bearer[0])
			if(accessExists){
				session,_ :=TokensMap.Get(bearer[0])
				if(session.Role=="1"){
					c.Set("token",bearer[0])
					c.Next()
					return
				}else{
					c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
				}
			}else{
				c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
			}
		}else{
			c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
		}
	}
}
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		bearer:=c.Request.Header["Bearer"]
		if(len(bearer)>0){
			accessExists:=LookupTokenKey(bearer[0])
			if(accessExists){
				c.Set("token",bearer[0])
				c.Next()
				return
			}else{
				c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
			}
		}else{
			c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
		}
	}
}
