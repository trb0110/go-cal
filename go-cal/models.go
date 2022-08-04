package gocal


type User struct {
	UserID       		string `json:"UserID"									xml:"UserID"`
	Username     		string `json:"Username"									xml:"Username"`
	Password     		string `json:"Password"									xml:"Password"`
	Role         		string `json:"Role"										xml:"Role"`
	CalorieLimit 		string `json:"CalorieLimit"								xml:"CalorieLimit"`
	Token        		string `json:"Token"									xml:"Token"`
}

type Calorie struct {
	CalorieId		string 								`json:"CalorieId"							xml:"CalorieId"`
	Timestamp 		string 								`json:"Timestamp"							xml:"Timestamp"`
	Food	 		string								`json:"Food"								xml:"Food"`
	CalorieCount	string								`json:"CalorieCount"						xml:"CalorieCount"`
	UserID			string								`json:"UserID"								xml:"UserID"`
	Username		string								`json:"Username"							xml:"Username"`
}