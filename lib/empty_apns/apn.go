package empty_apns

import (
	"log"

	"domliang.com/empty/lib/config"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
)

func SendNotify(notyMsg string) {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal("config error:", err)
	}
	authKey, err := token.AuthKeyFromFile(conf.Apns.AuthKey)
	if err != nil {
		log.Fatal("token error:", err)
	}

	token := &token.Token{
		AuthKey: authKey,
		// KeyID from developer account (Certificates, Identifiers & Profiles -> Keys)
		KeyID: conf.Apns.KeyID,
		// TeamID from developer account (View Account -> Membership)
		TeamID: conf.Apns.TeamID,
	}

	notification := &apns2.Notification{}
	notification.DeviceToken = conf.Apns.DeviceToken
	notification.Topic = "ai.pte.ios"
	notification.Payload = []byte(`{"aps":{"alert":"` + notyMsg + `","sound" : "chime.aiff"}}`)

	client := apns2.NewTokenClient(token)
	client.Push(notification)
	// print(res)
}
