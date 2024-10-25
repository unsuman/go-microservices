package types

type Invoice struct {
	OBUid         int64   `json:"obuID"`
	TotalAmount   float64 `json:"totalAmount"`
	TotalDistance float64 `json:"totalDistance"`
}

type Distance struct {
	Value float64 `json:"value"`
	OBUID int64   `json:"obuID"`
	Unix  int64   `json:"unix"`
}

type OBUData struct {
	OBUid int64   `json:"obuID"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}
