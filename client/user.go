package client

type User struct {
	Mid       int64
	Uname     string
	LoginInfo *LoginInfo
}

type LoginInfo struct {
	AccessKey string
}
