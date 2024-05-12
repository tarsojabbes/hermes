package auth

import (
	"hermes/database"
	"net/http"

	"github.com/google/uuid"
)

func IsAuthenticated(r *http.Request) bool {
	publisherUUID := r.Header.Get("HERMESUUID")
	return (publisherUUID != "" && database.FindPublisherByUUID(publisherUUID))
}

func RegisterAsPublisher() string {
	uuid := uuid.New()
	database.InsertPublisher(uuid)
	return uuid.String()
}