package main

import (
	"context"
	"log"
	"net"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func beckend(user, realm string, hc bool) string {
	// ClientPassword := "12345"      //configuration.RadClientPassword
	// NASid := "store"               //configuration.RadNASid
	// sec := "secret"                //configuration.RadSecret
	// RadSvr := "192.168.1.172:1812" //configuration.RadServer

	ClientPassword := configuration.RadClientPassword
	// NASid := configuration.RadNASid
	sec := configuration.RadSecret
	RadSvr := configuration.RadServer
	var NASid string

	switch realm {
	case "store":
		NASid = configuration.RadNASidS
	case "partner":
		NASid = configuration.RadNASidP
	case "emp":
		NASid = configuration.RadNASidE
	default:
		NASid = configuration.RadNASidS
	}

	packet := radius.New(radius.CodeAccessRequest, []byte(sec))
	rfc2865.NASIdentifier_SetString(packet, NASid)
	rfc2865.UserName_SetString(packet, user)
	rfc2865.UserPassword_SetString(packet, ClientPassword)

	response, err := radius.Exchange(context.Background(), packet, RadSvr)
	if err != nil {
		log.Fatal(err)
	}
	IPAttribute := response.Get(8)

	if response.Code != 2 || IPAttribute == nil {
		return "Not Available"
	}

	FramedIP := net.IPv4(IPAttribute[0], IPAttribute[1], IPAttribute[2], IPAttribute[3])
	FramedIPstring := FramedIP.String()
	if !hc {
		log.Println("Framed-IP-string : ", FramedIPstring)
		log.Println("Code:", response.Code)
	}
	return FramedIPstring
}
