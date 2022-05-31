package params

type IPAddr struct {
	Data    IPAddrData `json:"data"`
	Msg     string     `json:"msg"`
	Success bool       `json:"success"`
	CodNo   string     `json:"codNo"`
	Charge  bool       `json:"charge"`
}
type IPAddrData struct {
	Country   string `json:"country"`
	CountryID string `json:"country_id"`
	Region    string `json:"region"`
	RegionID  string `json:"region_id"`
	City      string `json:"city"`
	CityID    string `json:"city_id"`
	IP        string `json:"ip"`
	LongIP    string `json:"long_ip"`
	Isp       string `json:"isp"`
}
