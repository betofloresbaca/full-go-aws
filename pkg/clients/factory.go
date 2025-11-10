package clients

import (
	"context"
	"reflect"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// clientFactory es una fábrica de clientes de AWS con cache en memoria (Singleton)
type clientFactory struct {
	mu    sync.RWMutex
	cache map[string]interface{}
	cfg   aws.Config
}

var (
	instance *clientFactory
	once     sync.Once
)

// getFactoryInstance retorna la instancia singleton de ClientFactory
func getFactoryInstance() *clientFactory {
	once.Do(func() {
		ctx := context.Background()
		cfg, _ := config.LoadDefaultConfig(ctx)
		instance = &clientFactory{
			cache: make(map[string]interface{}),
			cfg:   cfg,
		}
	})
	return instance
}

// GetClient obtiene un cliente del cache o lo construye si no existe
// La clave se genera automáticamente del nombre del tipo T
// factory es una función que recibe la configuración de AWS y retorna el cliente del tipo T
// Retorna el cliente ya casteado al tipo correcto
func GetClient[T any](factory func(cfg aws.Config) T) T {
	// Generar la clave automáticamente del nombre del tipo
	var zero T
	key := reflect.TypeOf(zero).String()

	cf := getFactoryInstance()

	cf.mu.RLock()
	if cached, exists := cf.cache[key]; exists {
		cf.mu.RUnlock()
		return cached.(T)
	}
	cf.mu.RUnlock()

	// Si no existe, lo creamos
	client := factory(cf.cfg)

	// Guardamos en cache
	cf.mu.Lock()
	cf.cache[key] = client
	cf.mu.Unlock()

	return client
}
