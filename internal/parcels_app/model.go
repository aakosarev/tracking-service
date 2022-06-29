package parcels_app

func encodeTrackNumber(trackNumber string, encoder map[string]rune) string {
	rs := make([]rune, 0, len(trackNumber))
	for _, r := range trackNumber {
		rs = append(rs, encoder[string(r)])
	}
	return string(rs)
}

type TrackingResult struct {
	States []struct {
		Location string `json:"location"`
		Date     string `json:"date"`
		Carrier  int64  `json:"carrier"`
		Status   string `json:"status"`
	} `json,dynamodbav:"states"`
	Origin           string   `json:"origin"`
	Destination      string   `json:"destination"`
	From             string   `json:"from"`
	To               string   `json:"to"`
	Weight           string   `json:"weight"`
	Carriers         []string `json:"carriers"`
	ExternalTracking []struct {
		URL    string `json:"url"`
		Method string `json:"method"`
		Slug   string `json:"slug"`
	} `json:"external_tracking"`
	Attributes []struct {
		L   string `json:"l"`
		N   string `json:"n"`
		Val string `json:"val"`
	} `json:"attributes"`
	SubStatus       string `json:"sub_status"`
	Status          string `json:"status"`
	DestinationCode string `json:"destination_code"`
	OriginCode      string `json:"originCode"`
	Services        []struct {
		Slug string `json:"slug"`
		Name string `json:"name"`
	} `json:"services"`
}

type DatabaseData struct {
	TrackingNumber string   `json:"trackingNumber"`
	Status         string   `json:"status"`
	SubStatus      string   `json:"subStatus"`
	Origin         string   `json:"origin"`
	Destination    string   `json:"destination"`
	From           string   `json:"from"`
	To             string   `json:"to"`
	Weight         string   `json:"weight"`
	Carriers       []string `json:"carriers"`
	States         []struct {
		Location string `json:"location"`
		Date     string `json:"date"`
		Carrier  int64  `json:"carrier"`
		Status   string `json:"status"`
	} `json:"states"`
}

func (tr *TrackingResult) ConvertToDatabaseData(trackingNumber string) DatabaseData {
	databaseData := DatabaseData{
		TrackingNumber: trackingNumber,
		Status:         tr.Status,
		SubStatus:      tr.SubStatus,
		Origin:         tr.Origin,
		Destination:    tr.Destination,
		From:           tr.From,
		To:             tr.To,
		Weight:         tr.Weight,
		Carriers:       tr.Carriers,
		States:         tr.States,
	}
	return databaseData
}
