package tracking_more

type InputData struct {
	TrackingNumber string `json:"tracking_number"`
	CourierCode    string `json:"courier_code"`
}

type TrackingResult struct {
	Code int64 `json:"code"`
	Data struct {
		ScheduledAddress string `json:"Scheduled_Address"`
		Archived         bool   `json:"archived"`
		Consignee        string `json:"consignee"`
		CourierCode      string `json:"courier_code"`
		CreatedAt        string `json:"created_at"`
		CustomerEmail    string `json:"customer_email"`
		CustomerName     string `json:"customer_name"`
		CustomerPhone    string `json:"customer_phone"`
		DeliveryStatus   string `json:"delivery_status"`
		Destination      string `json:"destination"`
		DestinationInfo  struct {
			ArrivedAbroadDate      string `json:"arrived_abroad_date"`
			ArrivedDestinationDate string `json:"arrived_destination_date"`
			CourierCode            string `json:"courier_code"`
			CourierPhone           string `json:"courier_phone"`
			CustomsReceivedDate    string `json:"customs_received_date"`
			DepartedAirportDate    string `json:"departed_airport_date"`
			DispatchedDate         string `json:"dispatched_date"`
			ReceivedDate           string `json:"received_date"`
			ReferenceNumber        string `json:"reference_number"`
			Trackinfo              []struct {
				CheckpointDate              string `json:"checkpoint_date"`
				CheckpointDeliveryStatus    string `json:"checkpoint_delivery_status"`
				CheckpointDeliverySubstatus string `json:"checkpoint_delivery_substatus"`
				Location                    string `json:"location"`
				TrackingDetail              string `json:"tracking_detail"`
			} `json:"trackinfo"`
			Weblink string `json:"weblink"`
		} `json:"destination_info"`
		DestinationTrackNumber string `json:"destination_track_number"`
		ExchangeNumber         string `json:"exchangeNumber"`
		ID                     string `json:"id"`
		LastestCheckpointTime  string `json:"lastest_checkpoint_time"`
		LatestEvent            string `json:"latest_event"`
		LogisticsChannel       string `json:"logistics_channel"`
		Note                   string `json:"note"`
		OrderNumber            string `json:"order_number"`
		OriginInfo             struct {
			ArrivedAbroadDate      string `json:"arrived_abroad_date"`
			ArrivedDestinationDate string `json:"arrived_destination_date"`
			CourierCode            string `json:"courier_code"`
			CourierPhone           string `json:"courier_phone"`
			CustomsReceivedDate    string `json:"customs_received_date"`
			DepartedAirportDate    string `json:"departed_airport_date"`
			DispatchedDate         string `json:"dispatched_date"`
			ReceivedDate           string `json:"received_date"`
			ReferenceNumber        string `json:"reference_number"`
			Trackinfo              []struct {
				CheckpointDate              string `json:"checkpoint_date"`
				CheckpointDeliveryStatus    string `json:"checkpoint_delivery_status"`
				CheckpointDeliverySubstatus string `json:"checkpoint_delivery_substatus"`
				Location                    string `json:"location"`
				TrackingDetail              string `json:"tracking_detail"`
			} `json:"trackinfo"`
			Weblink string `json:"weblink"`
		} `json:"origin_info"`
		Original              string `json:"original"`
		Previously            string `json:"previously"`
		ScheduledDeliveryDate string `json:"scheduled_delivery_date"`
		ServiceCode           string `json:"service_code"`
		ShippingDate          string `json:"shipping_date"`
		SignSupported         string `json:"sign_supported"`
		StatusInfo            string `json:"status_info"`
		StayTime              int64  `json:"stay_time"`
		Substatus             string `json:"substatus"`
		Title                 string `json:"title"`
		TrackingNumber        string `json:"tracking_number"`
		TransitTime           int64  `json:"transit_time"`
		UpdateDate            string `json:"update_date"`
		Updating              bool   `json:"updating"`
		Weight                string `json:"weight"`
	} `json:"data"`
	Message string `json:"message"`
}

type DatabaseData struct {
	TrackingNumber        string           `json:"tracking_number"`
	CourierCode           string           `json:"courier_code"`
	LastestCheckpointTime string           `json:"lastest_checkpoint_time"`
	LatestEvent           string           `json:"latest_event"`
	Trackinfo             []TrackinfoPoint `json:"trackinfo"`
}

type TrackinfoPoint struct {
	CheckpointDate              string `json:"checkpoint_date"`
	CheckpointDeliveryStatus    string `json:"checkpoint_delivery_status"`
	CheckpointDeliverySubstatus string `json:"checkpoint_delivery_substatus"`
	Location                    string `json:"location"`
	TrackingDetail              string `json:"tracking_detail"`
}

func (tr *TrackingResult) ConvertToDatabaseData() DatabaseData {
	trackinfo := []TrackinfoPoint{}
	var checkpointDeliveryStatus string
	var checkpointDeliverySubstatus string
	for _, v := range tr.Data.OriginInfo.Trackinfo {

		switch v.CheckpointDeliveryStatus {
		case "transit":
			checkpointDeliveryStatus = "TRANSIT : Courier has picked up package from shipper, the package is on the way to destination"
		case "delivered":
			checkpointDeliveryStatus = "DELIVERED : The package was delivered successfully"
		case "pending":
			checkpointDeliveryStatus = "PENDING : New package added that are pending to track"
		case "pickup":
			checkpointDeliveryStatus = "PICKUP : Also known as \"Out For Delivery\", courier is about to deliver the package, or the package is wating for addressee to pick up"
		case "expired":
			checkpointDeliveryStatus = "EXPIRED : No tracking_more information for 30days for express service, or no tracking_more information for 60 days for postal service since the package added"
		case "undelivered":
			checkpointDeliveryStatus = "UNDELIVERED : Also known as \"Failed Attempt\", courier attempted to deliver but failded, usually left a notice and will try to delivery again"
		case "exception":
			checkpointDeliveryStatus = "EXCEPTION : Package missed, addressee returned package to sender or other exceptions"
		case "InfoReceived":
			checkpointDeliveryStatus = "INFO RECEIVED : Carrier has received request from shipper and is about to pick up the shipment"
		}

		switch v.CheckpointDeliverySubstatus {
		case "transit001":
			checkpointDeliverySubstatus = "TRANSIT 1 : Package is on the way to destination"
		case "transit002":
			checkpointDeliverySubstatus = "TRANSIT 2 : Package arrived at a hub or sorting center"
		case "transit003":
			checkpointDeliverySubstatus = "TRANSIT 3 : Package arrived at delivery facility"
		case "transit004":
			checkpointDeliverySubstatus = "TRANSIT 4 : Package arrived at destination country"
		case "transit005":
			checkpointDeliverySubstatus = "TRANSIT 5 : Customs clearance completed"
		case "transit006":
			checkpointDeliverySubstatus = "TRANSIT 6 : Item Dispatched"
		case "transit007":
			checkpointDeliverySubstatus = "TRANSIT 7 : Depart from Airport"

		case "delivered001":
			checkpointDeliverySubstatus = "DELIVERED 1 : Package delivered successfully"
		case "delivered002":
			checkpointDeliverySubstatus = "DELIVERED 2 : Package picked up by the addressee"
		case "delivered003":
			checkpointDeliverySubstatus = "DELIVERED 3 : Package received and signed by addressee"
		case "delivered004":
			checkpointDeliverySubstatus = "DELIVERED 4 : Package was left at the front door or left with your neighbour"

		case "undelivered001":
			checkpointDeliverySubstatus = "UNDELIVERED 1 : Address-related issues"
		case "undelivered002":
			checkpointDeliverySubstatus = "UNDELIVERED 2 : Receiver not home"
		case "undelivered003":
			checkpointDeliverySubstatus = "UNDELIVERED 3 : Impossible to locate the addressee"
		case "undelivered004":
			checkpointDeliverySubstatus = "UNDELIVERED 4 : Undelivered due to other reasons"

		case "pickup001":
			checkpointDeliverySubstatus = "PICKUP 1 : The package is out for delivery"
		case "pickup002":
			checkpointDeliverySubstatus = "PICKUP 2 : The package is ready for collection"
		case "pickup003":
			checkpointDeliverySubstatus = "PICKUP 3 : The customer is contacted before the final delivery"

		case "notfound001":
			checkpointDeliverySubstatus = "NOT FOUND 1 : The package is waiting for courier to pick up"
		case "notfound002":
			checkpointDeliverySubstatus = "NOT FOUND 2 : No tracking_more information found"

		case "exception004":
			checkpointDeliverySubstatus = "EXCEPTION 4 : The package is unclaimed"
		case "exception005":
			checkpointDeliverySubstatus = "EXCEPTION 5 : Other exceptions"
		case "exception006":
			checkpointDeliverySubstatus = "EXCEPTION 6 : Package was detained by customs"
		case "exception007":
			checkpointDeliverySubstatus = "EXCEPTION 7 : Package was lost or damaged during delivery"
		case "exception008":
			checkpointDeliverySubstatus = "EXCEPTION 8 : Logistics order was cancelled before courier pick up the package"
		case "exception009":
			checkpointDeliverySubstatus = "EXCEPTION 9 : Package was refused by addressee"
		case "exception0010":
			checkpointDeliverySubstatus = "EXCEPTION 10 : Package has been returned to sender"
		case "exception0011":
			checkpointDeliverySubstatus = "EXCEPTION 11 : Package is beening sent to sender"
		}

		trackinfoPoint := TrackinfoPoint{
			CheckpointDate:              v.CheckpointDate,
			CheckpointDeliveryStatus:    checkpointDeliveryStatus,
			CheckpointDeliverySubstatus: checkpointDeliverySubstatus,
			Location:                    v.Location,
			TrackingDetail:              v.TrackingDetail,
		}
		trackinfo = append(trackinfo, trackinfoPoint)
	}

	databaseData := DatabaseData{
		TrackingNumber:        tr.Data.TrackingNumber,
		CourierCode:           tr.Data.CourierCode,
		LastestCheckpointTime: tr.Data.LastestCheckpointTime,
		LatestEvent:           tr.Data.LatestEvent,
		Trackinfo:             trackinfo,
	}

	return databaseData
}
