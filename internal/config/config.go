package config

import (
	"log"

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
	viper.SetConfigFile(path + "/.env") // Asegura que leemos el archivo correcto
	viper.SetConfigType("env")          // Especifica el tipo
	viper.AutomaticEnv()                // Lee variables de sistema si el archivo falla

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Aviso: No se pudo leer el archivo: %v", err)
	}

	// DEBUG: Ver qué llaves encontró
	log.Printf("DEBUG: Llaves encontradas por Viper: %v", viper.AllKeys())

	err = viper.Unmarshal(&config)
	log.Printf("DEBUG: Estructura Config mapeada: %+v", config)

	return
}

/* func LoadConfig(path string) (config Config, err error) {
	// 1. Cargamos el .env explícitamente
	godotenv.Load(path + "/.env")

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Intentamos leer
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Aviso: Viper no pudo leer el archivo: %v", err)
	}

	// DEBUG: Ver dónde está buscando Viper
	log.Printf("DEBUG: Archivo usado por Viper: %s", viper.ConfigFileUsed())

	// 2. Mapeamos
	err = viper.Unmarshal(&config)

	// DEBUG: Imprimir qué cargó
	log.Printf("DEBUG: Config cargada: %+v", config)

	// 3. Validación crítica
	if config.DBSource == "" {
		log.Fatal("ERROR CRÍTICO: La variable DB_SOURCE está vacía. Verifica tu .env")
	}

	return
} */
