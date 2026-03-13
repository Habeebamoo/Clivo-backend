package models

type Subscriber struct {
	SubscriberId  string  `json:"subscriberId"`
	Email         string  `json:"email"`
}

type SubscriberRequest struct {
	Email  string  `json:"email"`
}