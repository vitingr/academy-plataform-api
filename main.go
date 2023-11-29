package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Estrutura do Servidor HTTP
type Server struct {
	client *mongo.Client
}

type User struct {
	Id                 int    `json:"id"`
	Name               string `json:"fullname"`
	FirstName          string `json:"firstname"`
	LastName           string `json:"lastname"`
	Email              string `json:"email"`
	Address            string `json:"address"`
	TrainingPreference string `json:"training"`
	Premium            bool   `json:"premium"`
	Admin              bool   `json:"admin"`
}

type Service struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Subscribers string  `json:"subscribers"`
}

func NewServer(c *mongo.Client) *Server {
	return &Server{
		client: c,
	}
}

func (s *Server) handleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	call := s.client.Database("academy").Collection("users")

	query := bson.M{}

	cursor, err := call.Find(context.TODO(), query)
	if err != nil {
		log.Fatal(err)
	}

	results := []bson.M{}

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	// Response do Servidor
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(results)
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	call := s.client.Database("academy").Collection("users")

	var reqData User
	err := json.NewDecoder(r.Body).Decode(&reqData)

	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
	}

	result, err := call.InsertOne(context.TODO(), bson.M{
		"name":             reqData.Name,
		"firstname":        reqData.FirstName,
		"lastname":         reqData.LastName,
		"email":            reqData.Email,
		"address":          "undefined",
		"trainingPrefence": "undefined",
		"premium":          false,
		"admin":            false,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Response do Servidor
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleGetOneUser(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	userEmail, ok := param["email"]

	if !ok {
		http.Error(w, "Email parameter is missing", http.StatusBadRequest)
		return
	}

	call := s.client.Database("academy").Collection("users")

	query := bson.M{"email": userEmail}
	user := call.FindOne(context.TODO(), query)

	if user.Err() != nil {
		if user.Err() == mongo.ErrNoDocuments {
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode("User not found")
		} else {
			http.Error(w, "Error while trying to find user", http.StatusInternalServerError)
		}
	} else {
		var userData bson.M

		if err := user.Decode(&userData); err != nil {
			http.Error(w, "Error while decoding user data", http.StatusInternalServerError)
		}

		// Response do Servidor
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userData)
	}
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	call := s.client.Database("academy").Collection("users")

	var reqData User
	err := json.NewDecoder(r.Body).Decode(&reqData)

	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
	}

	query := bson.M{"_id": reqData.Id}

	result, err := call.UpdateOne(context.TODO(), query, reqData)
	if err != nil {
		log.Fatal(err)
	}

	// Response do Servidor
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleGetAllServices(w http.ResponseWriter, r *http.Request) {
	call := s.client.Database("academy").Collection("services")

	query := bson.M{}

	cursor, err := call.Find(context.TODO(), query)
	if err != nil {
		log.Fatal(err)
	}

	results := []bson.M{}

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	// Response do Servidor
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(results)
}

func (s *Server) handleCreateService(w http.ResponseWriter, r *http.Request) {
	call := s.client.Database("academy").Collection("services")

	var reqData Service
	err := json.NewDecoder(r.Body).Decode(&reqData)

	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
	}

	result, err := call.InsertOne(context.TODO(), bson.M{
		"title": reqData.Title,
		"description": reqData.Description,
		"price": reqData.Price,
		"subscribers": reqData.Subscribers,  
	})
	if err != nil {
		log.Fatal(err)
	}

	// Response do Servidor
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleGetOneService(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	serviceId := param["id"]

	call := s.client.Database("academy").Collection("services")

	query := bson.M{"_id": serviceId}

	service := call.FindOne(context.TODO(), query)

	if service.Err() != nil {
		if service.Err() == mongo.ErrNoDocuments {
			http.Error(w, "Service not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error while tried to find service", http.StatusInternalServerError)
		}
	}

	var serviceData bson.M

	if err := service.Decode(&serviceData); err != nil {
		http.Error(w, "Error while decoding service data", http.StatusInternalServerError)
	}

	// Response do Servidor
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(serviceData)
}

func main() {
	// Configurações configuração Banco de Dados MongoDB
	uri := "mongodb+srv://vitorgabrielsbo1460:xXzMW9c0UljiQykk@aprendendo.cmdbthe.mongodb.net/academy?retryWrites=true&w=majority"
	if uri == "" {
		log.Fatal("ERROR! Connection with MongoDB failed...")
	}
	// Tratamento de Erros
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Servidor funcionando corretamente...")
	}

	server := NewServer(client)

	// Configurações do CORS
	router := mux.NewRouter()
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"*"}), // Permitir todos os métodos
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	router.Use(corsMiddleware)

	// FindManyUsers
	router.HandleFunc("/users", server.handleGetAllUsers).Methods("GET")

	// Create User
	router.HandleFunc("/users/create", server.handleCreateUser).Methods("POST")

	// FindOneUser
	router.HandleFunc("/users/{email}", server.handleGetOneUser).Methods("GET")

	// Update User
	router.HandleFunc("/updateUser/{id}", server.handleUpdateUser).Methods("POST")

	// FindServices
	router.HandleFunc("/services", server.handleGetAllServices).Methods("GET")

	// Create Service
	router.HandleFunc("/services/create", server.handleCreateService).Methods("POST")

	// FindOneService
	router.HandleFunc("/services/{id}", server.handleGetOneService).Methods("GET")

	http.ListenAndServe(":3030", router)
}
