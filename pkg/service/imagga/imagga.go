package imagga

import (
	"bank-authentication-system/pkg/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Imagga struct {
	Cfg config.Imagga
}

func (i *Imagga) DetectFace(imageURL string) (*FaceDetectionResponse, error) {
	fmt.Println(imageURL)
	client := &http.Client{}

	baseUrl := "https://api.imagga.com/v2/faces/detections?image_url="
	req, _ := http.NewRequest("GET", baseUrl+imageURL+"&return_face_id=1", nil)
	req.SetBasicAuth(i.Cfg.ApiKey, i.Cfg.ApiSecret)

	resp, err := client.Do(req)

	if err != nil {
		logrus.Errorf("Error when sending request to the server")
		return nil, err
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))

	var response FaceDetectionResponse
	err = json.Unmarshal(resp_body, &response)
	if err != nil {
		logrus.Errorf("Error unmarshaling JSON response: %v", err)
		return nil, err
	}

	return &response, nil
}

func (i *Imagga) FindFaceSimilarity(faceIDOne string, faceIDTwo string) (*FaceSimilarityResponse, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://api.imagga.com/v2/faces/similarity?face_id="+faceIDOne+"&second_face_id="+faceIDTwo, nil)
	req.SetBasicAuth(i.Cfg.ApiKey, i.Cfg.ApiSecret)

	resp, err := client.Do(req)

	if err != nil {
		logrus.Errorf("Error when sending request to the server")

		return nil, err
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))

	var response FaceSimilarityResponse
	err = json.Unmarshal(resp_body, &response)
	if err != nil {
		logrus.Errorf("Error unmarshaling JSON response: %v", err)

		return nil, err
	}

	return &response, nil
}
