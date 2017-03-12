package session

type Session struct {
	Values map[interface{}]interface{}
	Id     string
	store  *Store
}

func NewSession() {

}
