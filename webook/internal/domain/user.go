package domain

type User struct {
	Id       int64
	Email    string
	Password string

	NickName string
	Birth    string
	Synopsis string

	Ctime int64
	Utime int64
}
