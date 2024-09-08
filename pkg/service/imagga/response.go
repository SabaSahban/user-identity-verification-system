package imagga

type FaceDetectionResponse struct {
	Result struct {
		Faces []struct {
			Confidence  float64 `json:"confidence"`
			Coordinates struct {
				Height int `json:"height"`
				Width  int `json:"width"`
				Xmax   int `json:"xmax"`
				Xmin   int `json:"xmin"`
				Ymax   int `json:"ymax"`
				Ymin   int `json:"ymin"`
			} `json:"coordinates"`
			FaceId string `json:"face_id"`
		} `json:"faces"`
	} `json:"result"`
	Status struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"status"`
}

type FaceSimilarityResponse struct {
	Result struct {
		Score float64 `json:"score"`
	} `json:"result"`
	Status struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"status"`
}
