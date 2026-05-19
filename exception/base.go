package exception

type ILocation interface {
	GetStart() (int, int)
	GetEnd() (int, int)
	GetFile() string
}

type IException interface {
	error
	GetLoc() ILocation
	GetMsg() string
}
