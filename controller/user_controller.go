package controller

import (
    "context"
    "encoding/json"
    "net/http"
    "os"
    "time"
    "todo-list/config"
    "todo-list/model"

    "github.com/golang-jwt/jwt/v5"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
    var userCollection *mongo.Collection = config.GetCollection("users")
    var user model.User

    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Check if user with the same email already exists
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var existingUser model.User
    err = userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
    if err == nil {
        http.Error(w, "User with this email already exists", http.StatusConflict)
        return
    }
    if err != mongo.ErrNoDocuments {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

    user.Password = string(hashedPassword)
    user.ID = primitive.NewObjectID()

    // Insert into DB
    _, err = userCollection.InsertOne(ctx, user)
    if err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    // Generate JWT Token
    jwtKey := []byte(os.Getenv("JWT_SECRET"))
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID.Hex(),
        "email":   user.Email,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    })

    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    // Response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "User created successfully",
        "token":   tokenString,
        "user": map[string]string{
            "id":       user.ID.Hex(),
            "username": user.Username,
            "email":    user.Email,
        },
    })
}


func GetAllUsers(w http.ResponseWriter, r *http.Request) {
    userCollection := config.GetCollection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := userCollection.Find(ctx, bson.M{})
    if err != nil {
        http.Error(w, "Error fetching users", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var users []model.User
    for cursor.Next(ctx) {
        var user model.User
        if err := cursor.Decode(&user); err == nil {
            user.Password = "" // Hide password
            users = append(users, user)
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func GetUserByID(w http.ResponseWriter, r *http.Request, userID string) {
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var user model.User
    userCollection := config.GetCollection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    user.Password = ""
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request, userID string) {
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var updatedUser model.User
    json.NewDecoder(r.Body).Decode(&updatedUser)

    update := bson.M{
        "$set": bson.M{
            "username": updatedUser.Username,
            "email":    updatedUser.Email,
        },
    }

    userCollection := config.GetCollection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := userCollection.UpdateByID(ctx, objID, update)
    if err != nil || result.MatchedCount == 0 {
        http.Error(w, "User not found or update failed", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request, userID string) {
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    userCollection := config.GetCollection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objID})
    if err != nil || result.DeletedCount == 0 {
        http.Error(w, "User not found or delete failed", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

func Login(w http.ResponseWriter, r *http.Request) {
    var userCollection *mongo.Collection = config.GetCollection("users")
    var loginData struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    // Decode incoming login JSON
    err := json.NewDecoder(r.Body).Decode(&loginData)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Find user in database
    var user model.User
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err = userCollection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&user)
    if err != nil {
        http.Error(w, "Email or password is incorrect", http.StatusUnauthorized)
        return
    }

    // Compare hashed password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
    if err != nil {
        http.Error(w, "Email or password is incorrect", http.StatusUnauthorized)
        return
    }

    // Generate JWT
    jwtKey := []byte(os.Getenv("JWT_SECRET"))
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID.Hex(),
        "email":   user.Email,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    })

    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    // Return response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Login successful",
        "token":   tokenString,
        "user": map[string]string{
            "id":       user.ID.Hex(),
            "username": user.Username,
            "email":    user.Email,
        },
    })
}

