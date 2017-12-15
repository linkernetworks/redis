package redis

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"strings"
)

const (
	KeyModelVersion     = "model-version"
	KeyModelUrl         = "model-url"
	ClientsKeyPattern   = "clients:*"
	KeepAliveKeyPattern = "keepalive:*"
)

func (service *Service) RemoveExpiredClients() {
	c := service.Pool.Get()
	defer c.Close()

	// Get all rooms, ex. clients:room1
	rooms, err := redis.Strings(c.Do("KEYS", ClientsKeyPattern))
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, room := range rooms {
		// Get all clients in room for each room, ex. clients:room1
		clients, err := redis.Strings(c.Do("SMEMBERS", room))
		if err != nil {
			log.Fatal(err)
			continue
		}

		// Get live clients
		keepalives, err := redis.Strings(c.Do("KEYS", KeepAliveKeyPattern))
		if err != nil {
			log.Fatal(err)
			continue
		}

		// For each client
		for _, client := range clients {
			isDead := true
			for _, liveclient := range keepalives {
				if strings.Contains(liveclient, client) {
					isDead = false
				}
			}
			// If client is dead, append to dead clients
			if isDead {
				log.Printf("client %s is expired \n", client)
				_, err = c.Do("SREM", room, client)
				if err != nil {
					log.Fatal(err)
					continue
				}
			}
		}
	}
}
