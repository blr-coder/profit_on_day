package entities

type JsonStruct struct {
	CampaignId string  `json:"CampaignId"`
	Country    string  `json:"Country"`
	Ltv1       float64 `json:"Ltv1"`
	Ltv2       float64 `json:"Ltv2"`
	Ltv3       float64 `json:"Ltv3"`
	Ltv4       float64 `json:"Ltv4"`
	Ltv5       float64 `json:"Ltv5"`
	Ltv6       float64 `json:"Ltv6"`
	Ltv7       float64 `json:"Ltv7"`
	Users      int     `json:"Users"`
}

// Revenue общая прибыль
func (j *JsonStruct) Revenue() float64 {
	switch true {
	case j.Ltv7 != 0:
		return j.Ltv7 * float64(j.Users)
	case j.Ltv6 != 0:
		return j.Ltv6 * float64(j.Users)
	case j.Ltv5 != 0:
		return j.Ltv5 * float64(j.Users)
	case j.Ltv4 != 0:
		return j.Ltv4 * float64(j.Users)
	case j.Ltv3 != 0:
		return j.Ltv3 * float64(j.Users)
	default:
		// TODO: ???
		return 0
	}
}

type CsvStruct struct {
	UserId     int     `json:"user_id"`
	CampaignId string  `json:"CampaignId"`
	Country    string  `json:"Country"`
	Ltv1       float64 `json:"Ltv1"`
	Ltv2       float64 `json:"Ltv2"`
	Ltv3       float64 `json:"Ltv3"`
	Ltv4       float64 `json:"Ltv4"`
	Ltv5       float64 `json:"Ltv5"`
	Ltv6       float64 `json:"Ltv6"`
	Ltv7       float64 `json:"Ltv7"`
}
