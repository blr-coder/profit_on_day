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

func (j *JsonStruct) GetRevenue() float64 {
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

func (j *JsonStruct) GetCountry() string {
	return j.Country
}

func (j *JsonStruct) GetCampaignID() string {
	return j.CampaignId
}

type CsvStruct struct {
	UserId     int     `csv:"user_id"`
	CampaignId string  `csv:"CampaignId"`
	Country    string  `csv:"Country"`
	Ltv1       float64 `csv:"Ltv1"`
	Ltv2       float64 `csv:"Ltv2"`
	Ltv3       float64 `csv:"Ltv3"`
	Ltv4       float64 `csv:"Ltv4"`
	Ltv5       float64 `csv:"Ltv5"`
	Ltv6       float64 `csv:"Ltv6"`
	Ltv7       float64 `csv:"Ltv7"`
}

func (c *CsvStruct) GetRevenue() float64 {
	switch true {
	case c.Ltv7 != 0:
		return c.Ltv7
	case c.Ltv6 != 0:
		return c.Ltv6
	case c.Ltv5 != 0:
		return c.Ltv5
	case c.Ltv4 != 0:
		return c.Ltv4
	case c.Ltv3 != 0:
		return c.Ltv3
	default:
		// TODO: ???
		return 0
	}
}

func (c *CsvStruct) GetCountry() string {
	return c.Country
}

func (c *CsvStruct) GetCampaignID() string {
	return c.CampaignId
}
