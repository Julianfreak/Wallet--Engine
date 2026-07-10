package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config guarda toda la configuración de la aplicación
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig carga la configuración desde un archivo .env o variables de entorno
func LoadConfig(path string) (config Config, err error) {
	// 1. Intentamos cargar el archivo .env (útil en desarrollo local)
	// Si no existe, no pasa nada, viper buscará en el entorno del sistema
	if err := godotenv.Load(path + "/.env"); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// 2. Configuración de Viper
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv() // Sobreescribe los valores del archivo con las variables del sistema si existen

	// Intentamos leer el archivo si existe
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Aviso: No se pudo leer archivo de configuración: %v", err)
	}

	// 3. Mapeamos los valores al struct
	err = viper.Unmarshal(&config)
	return
}
