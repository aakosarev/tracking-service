package parcels_app

import (
	"encoding/json"
	"fmt"
	"github.com/aakosarev/tracking-service/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type Service interface {
	CreateOrUpdate(databaseData DatabaseData) error
}

type handler struct {
	service Service
	logger  *logging.Logger
}

func NewHandler(service Service, logger *logging.Logger) *handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/parcels_app", h.Track)
}

func (h *handler) Track(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Put parcels_app data in database")
	w.Header().Set("Content-Type", "application/json")
	trackingNumber := r.URL.Query().Get("tracking_number")
	h.logger.Debug(trackingNumber)
	if trackingNumber == "" {
		w.WriteHeader(400)
		w.Write([]byte(`{"message":"Request type error. Please check the API documentation for the request type of this API"}`))
		return
	}

	reqURL := "https://parcelsapp.com/api/v2/parcels"

	encoder := map[string]rune{
		"1": 'o',
		"2": 'p',
		"3": 'q',
		"4": 'r',
		"5": 's',
		"6": 't',
		"7": 'u',
		"8": 'v',
		"9": 'w',
		"0": 'n',
		"-": 'k',
	}

	encodedTrackingNumber := encodeTrackNumber(trackingNumber, encoder)
	h.logger.Debug(encodedTrackingNumber)

	req, err := http.NewRequest(http.MethodPost, reqURL, strings.NewReader(fmt.Sprintf("trackingId=%s&carrier=Auto-Detect&language=Unknown&country=Unknown&platform=web-desktop&wd=Unknown&c=Unknown&p=Unknown&l=Unknown&se=1920x1040%%2C1920x1040%%2C771x969%%2Cno%%2CWin32%%2CGecko%%2CMozilla%%2CNetscape%%2CGoogle+Inc.%%2Ctrue%%2Ctrue%%2CGoogle+Inc.+(Intel)%%2CANGLE+(Intel%%2C+Intel(R)+UHD+Graphics+620+Direct3D11+vs_5_0+ps_5_0%%2C+D3D11)%%2Ctrue%%2Ctrue%%2C184%%2C10%%2C4124653317", encodedTrackingNumber)))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Referer", fmt.Sprintf("https://parcels_app.com/en/tracking/%s", trackingNumber))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	defer req.Body.Close()

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}
	defer resp.Body.Close()

	trackingResult := TrackingResult{}

	if err = json.NewDecoder(resp.Body).Decode(&trackingResult); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}
	h.logger.Debug("res:", trackingResult)

	status := resp.StatusCode
	h.logger.Debug(status)
	if status != http.StatusOK {
		w.WriteHeader(status)
		w.Write([]byte(fmt.Sprintf(`{"message":"Message from ParcelsApp: %s"}`, "some error has occured")))
		return
	}

	databaseData := trackingResult.ConvertToDatabaseData(trackingNumber)
	h.logger.Debug(databaseData)

	if err = h.service.CreateOrUpdate(databaseData); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Request response is successful"}`))
	return
}
