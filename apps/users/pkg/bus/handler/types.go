package handler

type BusHandler interface {
	Load(key any) (value any, ok bool)
	Store(key any, value any) (duplicated bool)
	Delete(key any)
}
