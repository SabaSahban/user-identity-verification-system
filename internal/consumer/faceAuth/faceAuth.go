package faceAuth

import (
	"bank-authentication-system/pkg/model"
	"bank-authentication-system/pkg/mqtt"
	"bank-authentication-system/pkg/service/imagga"
	"bank-authentication-system/pkg/service/mail"
	"bank-authentication-system/pkg/state"
	"bank-authentication-system/pkg/storage/s3"
	"fmt"

	"github.com/sirupsen/logrus"
)

const emailSubject = "YOUR MEETING"

type FaceAuth struct {
	Imagga   *imagga.Imagga
	UserRepo model.UserRepo
	MailGun  *mail.Mailgun
	MQTT     *mqtt.MQTT
	S3       *s3.S3
}

func (f *FaceAuth) Process() {
	events, err := f.MQTT.Channel.Consume(
		f.MQTT.Queue,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		logrus.Errorf("failed to consume messages: %v\n", err)
	}

	logrus.Infof("consumer started on queue: %s", f.MQTT.Queue)

	for event := range events {
		userID := string(event.Body)

		imageOneUrl := s3.GetImageFromS3("image_one" + userID)
		imageTwoUrl := s3.GetImageFromS3("image_two" + userID)

		imageOneResp, err := f.Imagga.DetectFace(imageOneUrl)
		if err != nil {
			logrus.Errorf("face detection failed for face one %s", err.Error())
		}

		imageTwoResp, err := f.Imagga.DetectFace(imageTwoUrl)
		if err != nil {
			logrus.Errorf("face detection failed for face two %s", err.Error())
		}

		isImageOneFace := f.isFace(imageOneResp)
		isImageTwoFace := f.isFace(imageTwoResp)

		if !isImageOneFace || !isImageTwoFace {
			logrus.Info("can't detect face in one of the pictures")

			err := f.UserRepo.UpdateStateByUserID(userID, state.RejectState)
			if err != nil {
				logrus.Errorf("error updating user state: %s", err.Error())
			}
		} else {
			similarity, err := f.Imagga.FindFaceSimilarity(imageOneResp.Result.Faces[0].FaceId, imageTwoResp.Result.Faces[0].FaceId)
			if err != nil {
				logrus.Info("can't find similarity: %s", err.Error())
			}

			similarityScore := similarity.Result.Score

			f.checkSimilarity(similarityScore, userID)
		}

		user, err := f.UserRepo.FindUserByUserID(userID)
		if err != nil {
			logrus.Errorf("failed to find state by userID %s", err)
		}

		if user.State == state.AcceptState {
			emailBody := fmt.Sprintf("your authentication state is: %s", user.State, "https://meet.google.com/vwe-corw-fbz")
			logrus.Info("authentication state email is successfully sent:\"https://meet.google.com/vwe-corw-fbz\"")

			err = f.MailGun.Send(emailBody, emailSubject, user.Email)
		} else {
			emailBody := fmt.Sprintf("your authentication state is: %s")

			err = f.MailGun.Send(emailBody, emailSubject, user.Email)
			logrus.Info("authentication state email is successfully sent")

		}

		if err != nil {
			logrus.Errorf("failed to send state with mailgun %s", err)
		}
	}
}

func (f *FaceAuth) checkSimilarity(similarityScore float64, userID string) {
	logrus.Infof("similiary is %f percent", similarityScore)

	if similarityScore > 80 {
		logrus.Info("pictures are similar")

		err := f.UserRepo.UpdateStateByUserID(userID, state.AcceptState)
		if err != nil {
			logrus.Errorf("can't update user's state: %s", err.Error())
		}
	} else {
		logrus.Info("pictures are not similar")

		err := f.UserRepo.UpdateStateByUserID(userID, state.RejectState)
		if err != nil {
			logrus.Errorf("can't update user's state: %s", err.Error())
		}
	}
}

func (f *FaceAuth) isFace(imageResp *imagga.FaceDetectionResponse) bool {
	if len(imageResp.Result.Faces) == 0 {
		return false
	}

	return true
}
