package types

type OBUData struct {
	OBUid int64   `json:"obuID"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}
