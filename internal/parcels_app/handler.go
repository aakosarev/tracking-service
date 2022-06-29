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

	if trackingNumber == "" {
		w.WriteHeader(400)
		w.Write([]byte(`{"message":"Request type error. Please check the API documentation for the request type of this API"}`))
		return
	}

	reqURL := "https://parcels_app.com/api/v2/parcels"

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

	req, err := http.NewRequest(http.MethodPost, reqURL, strings.NewReader(fmt.Sprintf("trackingId=%s&carrier=Auto-Detect&language=en&country=Unknown&platform=web-desktop&wd=false&c=false&p=5&l=3&se=1440x791%%2C1440x791%%2C829x680%%2Cno%%2CMacIntel%%2CGecko%%2CMozilla%%2CNetscape%%2CGoogle+Inc.%%2Ctrue%%2Ctrue%%2CGoogle+Inc.+(Intel+Inc.)%%2CANGLE+(Intel+Inc.%%2C+Intel(R)+Iris(TM)+Plus+Graphics+OpenGL+Engine%%2C+OpenGL+4.1)%%2Ctrue%%2Ctrue%%2C195%%2C12%%2C2203383760&extra%%5BdetectedSlugs%%5D=%%5B%%22ups%%22%%2C%%22moscow-cargo%%22%%2C%%22emirates-airlines-cargo%%22%%5D", encodedTrackingNumber)))
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

	resp, _ := httpClient.Do(req)
	defer resp.Body.Close()

	trackingResult := TrackingResult{}

	if err = json.NewDecoder(resp.Body).Decode(&trackingResult); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}

	status := resp.StatusCode
	h.logger.Debug(status)
	if status != http.StatusOK {
		w.WriteHeader(status)
		w.Write([]byte(fmt.Sprintf(`{"message":"Message from ParcelsApp: %s"}`, "some error has occured")))
		return
	}

	databaseData := trackingResult.ConvertToDatabaseData(trackingNumber)
	h.logger.Debug(databaseData.TrackingNumber)
	if err = h.service.CreateOrUpdate(databaseData); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Request response is successful"}`))
	return
}
