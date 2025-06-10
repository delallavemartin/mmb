package main

import (
    "log"
    "net/http"

    "go.uber.org/zap"

    "mllave.com/mllave/mmb/mb/src/messengerservice/mail"
    "mllave.com/mllave/mmb/mb/src/messengerservice/postoffice"
    "mllave.com/mllave/mmb/mb/src/messengerservice/delivery"
    "mllave.com/mllave/mmb/mb/src/reader"
    "mllave.com/mllave/mmb/mb/src/model/subscribers"
)

// PublisherHandler returns an http.HandlerFunc with injected dependencies.
func PublisherHandler(subs *subscribers.SubscribersList, logger *zap.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        msg := reader.CustomReader{Reader: r.Body}
        defer r.Body.Close()

        aSubscribersPostOffice := postoffice.PostOffice{Channel: make(chan mail.Mail)}

        go aSubscribersPostOffice.OnMessageReceived(func(aSubscriberMail mail.Mail) {
            anHTTPDelivery := delivery.HttpDelivery{Mail: aSubscriberMail}
            anHTTPDelivery.Delivers()
        })

        subs.NotifySubscribers(aSubscribersPostOffice.NotificationAssistant(msg.ToString()))
        logger.Info("Notification sent to subscribers")
        w.WriteHeader(http.StatusAccepted)
    }
}

// SubscriberHandler returns an http.HandlerFunc with injected dependencies.
func SubscriberHandler(subs *subscribers.SubscribersList, logger *zap.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        portNumber := reader.CustomReader{Reader: r.Body}
        defer r.Body.Close()

        subs.Add(portNumber.ToString())
        logger.Info("Subscriber added", zap.String("port", portNumber.ToString()))
        w.WriteHeader(http.StatusCreated)
    }
}