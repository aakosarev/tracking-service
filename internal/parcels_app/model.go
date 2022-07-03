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
		Location string `json:"location,omitempty"`
		Date     string `json:"date,omitempty"`
		Carrier  int64  `json:"carrier,omitempty"`
		Status   string `json:"status,omitempty"`
	} `json:"states,omitempty"`
	Origin           string   `json:"origin,omitempty"`
	Destination      string   `json:"destination,omitempty"`
	From             string   `json:"from,omitempty"`
	To               string   `json:"to,omitempty"`
	Weight           string   `json:"weight,omitempty"`
	Carriers         []string `json:"carriers,omitempty"`
	ExternalTracking []struct {
		URL        string `json:"url,omitempty"`
		Method     string `json:"method,omitempty"`
		Slug       string `json:"slug,omitempty"`
		Copy       bool   `json:"copy,omitempty"`
		TrackingID string `json:"tracking_id,omitempty"`
	} `json:"external_tracking,omitempty"`
	Attributes []struct {
		L   string `json:"l,omitempty"`
		N   string `json:"n,omitempty"`
		Val string `json:"val,omitempty"`
	} `json:"attributes,omitempty"`
	SubStatus       string `json:"sub_status,omitempty"`
	Status          string `json:"status,omitempty"`
	DestinationCode string `json:"destination_code,omitempty"`
	OriginCode      string `json:"originCode,omitempty"`
	Services        []struct {
		Slug       string `json:"slug,omitempty"`
		Name       string `json:"name,omitempty"`
		IsFinished bool   `json:"is_finished,omitempty"`
	} `json:"services,omitempty"`
}

type DatabaseData struct {
	TrackingNumber string   `json:"trackingNumber,omitempty"`
	Status         string   `json:"status,omitempty"`
	SubStatus      string   `json:"subStatus,omitempty"`
	Origin         string   `json:"origin,omitempty"`
	Destination    string   `json:"destination,omitempty"`
	From           string   `json:"from,omitempty"`
	To             string   `json:"to,omitempty"`
	Weight         string   `json:"weight,omitempty"`
	Carriers       []string `json:"carriers,omitempty"`
	States         []struct {
		Location string `json:"location,omitempty"`
		Date     string `json:"date,omitempty"`
		Carrier  int64  `json:"carrier,omitempty"`
		Status   string `json:"status,omitempty"`
	} `json:"states,omitempty"`
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
