package repo

import "common/database"

type Manager struct {
	Mongo *database.MongoManager
	Redis *database.RedisManager
}

func New() *Manager {
	return &Manager{
		Mongo: database.NewMongo(),
		Redis: database.NewRedis(),
	}
}

func (m *Manager) Close() {
	if m.Mongo != nil {
		m.Mongo.Close()
	}
	if m.Redis != nil {
		m.Redis.Close()
	}
}
