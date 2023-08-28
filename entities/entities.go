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
		return j.Ltv7 / float64(j.Users)
	case j.Ltv6 != 0:
		return j.Ltv6 / float64(j.Users)
	case j.Ltv5 != 0:
		return j.Ltv5 / float64(j.Users)
	case j.Ltv4 != 0:
		return j.Ltv4 / float64(j.Users)
	case j.Ltv3 != 0:
		return j.Ltv3 / float64(j.Users)
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

func (j *JsonStruct) GetLtv(i int) float64 {
	switch i {
	case 1:
		if j.Ltv1 != 0 {
			return j.Ltv1
		}
	case 2:
		if j.Ltv2 != 0 {
			return j.Ltv2
		}
	case 3:
		if j.Ltv3 != 0 {
			return j.Ltv3
		}
	case 4:
		if j.Ltv4 != 0 {
			return j.Ltv4
		}
	case 5:
		if j.Ltv5 != 0 {
			return j.Ltv5
		}
	case 6:
		if j.Ltv6 != 0 {
			return j.Ltv6
		}
	case 7:
		if j.Ltv7 != 0 {
			return j.Ltv7
		}
	}

	return 0
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

func (c *CsvStruct) getFlippedLtvArr() []float64 {
	return []float64{c.Ltv7, c.Ltv6, c.Ltv5, c.Ltv4, c.Ltv3, c.Ltv2, c.Ltv1}
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

func (c *CsvStruct) GetLtv(i int) float64 {
	var ltv float64
	switch i {
	case 1:
		if c.Ltv1 != 0 {
			ltv = c.Ltv1
		}
	case 2:
		if c.Ltv2 != 0 {
			ltv = c.Ltv2
		}
	case 3:
		if c.Ltv3 != 0 {
			ltv = c.Ltv3
		}
	case 4:
		if c.Ltv4 != 0 {
			ltv = c.Ltv4
		}
	case 5:
		if c.Ltv5 != 0 {
			ltv = c.Ltv5
		}
	case 6:
		if c.Ltv6 != 0 {
			ltv = c.Ltv6
		}
	case 7:
		if c.Ltv7 != 0 {
			ltv = c.Ltv7
		}
	}

	if ltv == 0 {
		for _, prevLtv := range c.getFlippedLtvArr() {
			if prevLtv != 0 {
				ltv = prevLtv
			}
		}
	}

	return ltv
}
