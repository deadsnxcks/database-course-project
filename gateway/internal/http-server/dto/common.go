package dto
 
type Created struct {
	ID int64 `json:"id"`
}
 
type OK struct {
	Status string `json:"status"`
}
 
func Ok() OK { return OK{Status: "ok"} }