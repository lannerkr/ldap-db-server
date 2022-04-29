package main

import (
	"log"
	"strings"
)

func userFromCn(oname string) (user, realm string) {

	oname = strings.Replace(oname, ")", "", 1)
	oname = strings.Replace(oname, "(", "", 1)
	newcn := strings.Split(oname, ",")
	if len(newcn) < 2 {
		log.Printf("object name is invalid: %v\n", oname)
		return "", ""
	}
	userN := strings.Split(newcn[0], "=")
	realmN := strings.Split(newcn[1], "=")
	if len(userN) < 2 || len(realmN) < 2 {
		log.Printf("cn is invalid, user: %v , realm: %v\n", userN, realmN)
		return "", ""
	}
	user = userN[1]
	realm = realmN[1]
	log.Printf("user: %v , realm: %v\n", user, realm)

	return user, realm

}
