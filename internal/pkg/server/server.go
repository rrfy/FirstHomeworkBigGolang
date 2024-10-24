package server

import (
	"encoding/json"
	"homework1/internal/pkg/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host    string
	storage *storage.Storage
}

type Entry struct {
	Value interface{} `json:"value"`
}

func New(host string, st *storage.Storage) *Server {
	s := &Server{
		host:    host,
		storage: st,
	}

	return s
}

func (r *Server) newAPI() *gin.Engine {
	engine := gin.New()

	engine.GET("/hello-world", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Hello world!")
	})

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	// Маршруты для работы со скалярными значениями
	engine.PUT("/scalar/set/:key", r.handlerSet)
	engine.GET("/scalar/get/:key", r.handlerGet)

	return engine
}

func (r *Server) handlerSet(ctx *gin.Context) {
	key := ctx.Param("key")

	var v Entry

	// Декодируем значение, которое может быть любого типа
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.storage.Set(key, v.Value)

	ctx.Status(http.StatusOK)
}

func (r *Server) handlerGet(ctx *gin.Context) {
	key := ctx.Param("key")

	// Получаем значение из хранилища
	value, err := r.storage.Get(key)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Отправляем значение в JSON-формате
	ctx.JSON(http.StatusOK, Entry{Value: value})
}

// Start запускает сервер
func (r *Server) Start() {
	r.newAPI().Run(r.host)
}
