package redis

import (
	"log"
	"strings"

	"github.com/gomodule/redigo/redis"
)

const (
	KeyModelVersion     = "model-version"
	KeyModelUrl         = "model-url"
	ClientsKeyPattern   = "clients:*"
	KeepAliveKeyPattern = "keepalive:*"
)

func (service *Service) RemoveExpiredClients() error {
	c := service.Pool.Get()
	defer c.Close()

	// Get all rooms, ex. clients:room1
	rooms, err := redis.Strings(c.Do("KEYS", ClientsKeyPattern))
	if err != nil {
		return err
	}

	for _, room := range rooms {
		// Get all clients in room for each room, ex. clients:room1
		clients, err := redis.Strings(c.Do("SMEMBERS", room))
		if err != nil {
			log.Printf("redis SMEMBERS: %v\n", err)
			continue
		}

		// Get live clients
		keepalives, err := redis.Strings(c.Do("KEYS", KeepAliveKeyPattern))
		if err != nil {
			log.Printf("redis do KEYS: %v\n", err)
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
					log.Printf("redis SREM: %v\n", err)
					continue
				}
			}
		}
	}
	return nil
}
