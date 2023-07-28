package user

import "time"

// 客户端 grpc <--> 服务端 匿名网络+json <--> 云管
type UserInfo struct {
	User string    `json:"user"`
	Pass string    `json:"pass"`
	Mac  []string  `json:"mac"`
	Time time.Time `json:"time"`
}

type AnonyNode struct {
	Ip string `json:"ip"`
}
