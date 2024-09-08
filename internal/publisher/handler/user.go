package handler

import (
	"bank-authentication-system/pkg/model"
	"bank-authentication-system/pkg/mqtt"
	"bank-authentication-system/pkg/state"
	"bank-authentication-system/pkg/storage/s3"
	"bank-authentication-system/pkg/util"
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	UserRepo model.UserRepo
	S3       *s3.S3
	MQTT     *mqtt.MQTT
}

func (h *UserHandler) RegisterRequestHandler(c echo.Context) error {
	userID := uuid.New().String()
	logrus.Infof("user ID is %s", userID)

	email := c.FormValue("email")
	lastName := c.FormValue("last_name")
	nationalID := c.FormValue("national_id")

	err := h.UserRepo.Save(model.User{
		UUID:       userID,
		Email:      email,
		LastName:   lastName,
		NationalID: util.HashString(nationalID),
		IP:         c.RealIP(),
		State:      state.PendingState,
	})
	if err != nil {
		logrus.Errorf("failed to save user to database: %s", err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	err = h.uploadTos3(userID, c)
	if err != nil {
		logrus.Errorf("failed to upload image to s3: %s", err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	err = h.MQTT.Channel.PublishWithContext(
		context.Background(),
		"",
		h.MQTT.Queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(userID),
		},
	)
	if err != nil {
		logrus.Errorf("failed to send userID with rabbitMQ: %s", err.Error())

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "Your identity verification request has been registered")
}

func (h *UserHandler) uploadTos3(userID string, c echo.Context) error {
	//uploading image one to s3
	image, err := c.FormFile("image_one")
	if err != nil {
		return err
	}

	file, err := image.Open()
	if err != nil {
		logrus.Errorf("error opening image one")

		return err
	}

	uploader := s3manager.NewUploader(h.S3.Session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(h.S3.Cfg.Bucket),
		Key:    aws.String("image_one" + userID),
		Body:   file,
		ACL:    aws.String("public-read"), // Set ACL to public-read
	})
	if err != nil {
		logrus.Errorf("error uploading image one to S3")

		return err
	}

	logrus.Infof("image one uploaded to s3 successfuly")

	//uploading image two to s3
	image, err = c.FormFile("image_two")
	if err != nil {
		return err
	}

	file, err = image.Open()
	if err != nil {
		logrus.Errorf("error opening image two")

		return err
	}

	uploader = s3manager.NewUploader(h.S3.Session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(h.S3.Cfg.Bucket),
		Key:    aws.String("image_two" + userID),
		Body:   file,
		ACL:    aws.String("public-read"), // Set ACL to public-read
	})
	if err != nil {
		logrus.Errorf("error uploading image two to S3")

		return err
	}

	return nil
}

func (h *UserHandler) CheckRequestStatusHandler(c echo.Context) error {
	nationalID := c.Param("id")

	user, err := h.UserRepo.FindByNationalID(util.HashString(nationalID))
	if err != nil {
		logrus.Errorf("Failed to retrieve user from the database: %s", err.Error())

		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	if c.RealIP() != user.IP {
		return c.JSON(http.StatusForbidden, "ip mismatch")
	}

	switch user.State {
	case state.PendingState:
		return c.String(http.StatusOK, "Authentication in progress")
	case state.RejectState:
		return c.String(http.StatusOK, "Your identity verification request has been rejected, please try again later.")
	case state.AcceptState:
		return c.String(http.StatusOK, fmt.Sprintf("Identity verification successful, your username is %s. The Session's Link is: https://meet.google.com/vwe-corw-fbz", user.UUID))
	}

	return c.String(http.StatusNotFound, "Authentication state not found")
}
