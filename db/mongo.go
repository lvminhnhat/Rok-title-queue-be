package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
}

var (
	instance *Mongo
	once     sync.Once
	mu       sync.Mutex
)

// ConnectMongo tạo một kết nối MongoDB và trả về một instance của Mongo struct
func ConnectMongo(uri string) (*Mongo, error) {
	// Tạo context với timeout để kết nối không bị treo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Cấu hình các tùy chọn cho client MongoDB
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(20).                        // Số kết nối tối đa
		SetMinPoolSize(5).                         // Số kết nối tối thiểu
		SetMaxConnIdleTime(30 * time.Second).      // Thời gian chờ kết nối rảnh
		SetServerSelectionTimeout(5 * time.Second) // Thời gian chọn server

	// Kết nối tới MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Kiểm tra kết nối
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB")

	// Trả về một struct Mongo chứa client
	return &Mongo{
		Client: client,
	}, nil
}

// GetMongoInstance trả về instance duy nhất của MongoDB
func GetMongoInstance(uri string) (*Mongo, error) {
	mu.Lock()
	defer mu.Unlock()

	var err error
	once.Do(func() {
		instance, err = ConnectMongo(uri)
	})

	return instance, err
}

// DisconnectMongo dùng để ngắt kết nối MongoDB khi không sử dụng nữa
func (m *Mongo) DisconnectMongo() error {
	if err := m.Client.Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("failed to disconnect MongoDB: %w", err)
	}
	log.Println("Successfully disconnected from MongoDB")
	return nil
}

// filepath: d:\Code\rok\Backend-2\db\mongo.go
func (m *Mongo) ExecuteWithRetry(ctx context.Context, operation func(context.Context) error) error {
	retries := 3
	var err error

	for i := 0; i < retries; i++ {
		err = operation(ctx)
		if err == nil {
			return nil
		}

		// Chỉ retry cho một số lỗi nhất định
		if cmdErr, ok := err.(mongo.CommandError); ok {
			if cmdErr.HasErrorLabel("TransientTransactionError") ||
				cmdErr.HasErrorLabel("NetworkError") {
				time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
				continue
			}
		}

		// Lỗi khác không retry
		return err
	}

	return err
}
