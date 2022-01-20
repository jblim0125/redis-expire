package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/jblim0125/redredis-expire/common"
	"github.com/jblim0125/redredis-expire/models"
)

// Manager Redis Manager
type Manager struct {
	log        *common.Logger
	client     *redis.Client
	serverInfo models.DataStorageInfo
}

// CreateInstance create instance
func (Manager) CreateInstance(log *common.Logger) *Manager {
	return &Manager{
		log: log,
	}
}

// ClientInit client init and connect test
func (manager *Manager) ClientInit(server models.DataStorageInfo) error {
	serverAddr := fmt.Sprintf("%s:%d", server.IP, server.Port)
	redisDB, _ := strconv.Atoi(server.Database)
	client := redis.NewClient(&redis.Options{
		Addr:     serverAddr,
		Password: server.PW,
		Username: server.ID,
		DB:       redisDB,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	manager.client = client
	manager.serverInfo = server
	return nil
}

// Destroy close and shutdown
func (manager *Manager) Destroy() {
	manager.client.Close()
}

// SetNotifyEvent set(enable/disable) notify event
func (manager *Manager) SetNotifyEvent(enable bool) error {
	/*
	   K     Keyspace events, published with __keyspace@<db>__ prefix.
	   E     Keyevent events, published with __keyevent@<db>__ prefix.
	   g     Generic commands (non-type specific) like DEL, EXPIRE, RENAME, ...
	   $     String commands
	   l     List commands
	   s     Set commands
	   h     Hash commands
	   z     Sorted set commands
	   t     Stream commands
	   d     Module key type events
	   x     Expired events (events generated every time a key expires)
	   e     Evicted events (events generated when a key is evicted for maxmemory)
	   m     Key miss events (events generated when a key that doesn't exist is accessed)
	   A     Alias for "g$lshztxed", so that the "AKE" string means all the events except "m".
	*/
	// This only expired events
	notiResult, err := manager.client.Do(context.Background(), "CONFIG", "SET", "notify-keyspace-events", "Ex").Result()
	if err != nil {
		manager.log.Errorf("unable to set keyspace events %v", err.Error())
		return err
	}
	manager.log.Infof("Notify Set Result[ %s ]", notiResult)
	return nil
}

// StartEventReceiver start redis event receiver
func (manager *Manager) StartEventReceiver() error {
	manager.log.Errorf("[ Redis ] Event Receiver Start .................................................... [ OK ]")

	channel := fmt.Sprintf("__keyevent@%s__:expired", manager.serverInfo.Database)
	manager.log.Infof("publish subscribe[ %s ]", channel)
	// this is telling redis to subscribe to events published in the keyevent channel, specifically for expired events
	pubSub := manager.client.PSubscribe(context.Background(), channel)
	// this is just to show publishing events and catching the expired events in the same codebase.

	// infinite loop
	for {
		message, err := pubSub.ReceiveMessage(context.Background())
		if err != nil {
			manager.log.Error(err)
		} else {
			manager.log.Errorf("MSG[ %s ]", message)
		}
	}
	return nil
}

// WriteData write data to redis server
func (manager *Manager) WriteData(data *models.SomeData) error {
	value, _ := json.Marshal(&data.SomeDataValue)
	result, err := manager.client.Set(context.Background(), data.Key, value, time.Second*time.Duration(3)).Result()
	if err != nil {
		return err
	}
	manager.log.Infof("Write Data Result[ %s ]", result)
	return nil
}

// DeleteData delete data to redis server
func (manager *Manager) DeleteData(data *models.SomeData) error {
	//value, err := json.Marshal(data.SomeDataValue)
	result, err := manager.client.Del(context.Background(), data.Key).Result()
	if err != nil {
		return err
	}
	manager.log.Infof("Delete Data Result[ %s ]", result)
	return nil
}
