package models

import "time"

type User struct {
	Name     string
	Password string
}

type Session struct {
	ID          string    `json:"id", bson:"id"`
	Title       string    `json:"title", bson:"title"`
	Description string    `json:"description", bson:"description"`
	Private     bool      `json:"private", bson:"private"`
	Creator     string    `json:"creator", bson:"creator"`
	CreatedDate time.Time `json:"createdDate"`
	Members     []string  `json:"members"`
}

type MQTTServer struct {
	IP       string `json:"ip", bson:"ip"`
	Port     uint16 `json:"port", bson:"port"`
	Username string `json:"username", bson:"username"`
	Password string `json:"password", bson:"password"`
}

type TURNServer struct {
	URI      string `json:"uri", bson:"uri"`
	Username string `json:"username", bson:"username"`
	Password string `json:"password", bson:"password"`
}
