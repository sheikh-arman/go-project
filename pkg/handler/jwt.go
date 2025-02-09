package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	chiro "github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/sheikh-arman/api-server/newsfeed"
)

var (
	jwtkey    = []byte("adsads")
	TokenAuth *jwtauth.JWTAuth
	// tokenString string
	// token       jwt.Token
)

var (
	ID            int
	feeds         []newsfeed.Item
	feeds2        map[int]newsfeed.Item
	Credslist     map[string]string
	refreshTokens = make(map[string]string) // Refresh Token স্টোর করার জন্য
)

func InitCred() {
	TokenAuth = jwtauth.New(string(jwa.HS256), jwtkey, nil)
	Credslist = make(map[string]string)

	creds := []newsfeed.Credentials{
		{
			Username: "arman",
			Password: "123",
		},
	}

	for _, cred := range creds {
		Credslist[cred.Username] = cred.Password
	}
}

func InitDB() {
	ID = 1
	var feed newsfeed.Item
	feed = newsfeed.Item{
		Id:    ID,
		Title: "Nothing",
		Post:  "Lorem Ipsum Doller Site",
	}
	// feeds2[ID] = feed
	ID++
	feeds = append(feeds, feed)

	feed = newsfeed.Item{
		Id:    ID,
		Title: "Nothing2",
		Post:  "Lorem Ipsum Doller Site2",
	}
	// feeds2[ID] = feed
	ID++
	feeds = append(feeds, feed)

	feed = newsfeed.Item{
		Id:    ID,
		Title: "Nothing3",
		Post:  "Lorem Ipsum Doller Site3",
	}
	// feeds2[ID] = feed
	ID++
	feeds = append(feeds, feed)
}

func WriteJsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	log.Println(err)
}

func GetNewsFeeds(w http.ResponseWriter, r *http.Request) {
	log.Println("test")
	sort.SliceStable(feeds, func(i, j int) bool {
		return feeds[i].Id < feeds[j].Id
	})
	WriteJsonResponse(w, http.StatusOK, feeds2)
}

func GetNewsFeed(w http.ResponseWriter, r *http.Request) {
	param := chiro.URLParam(r, "id")
	paramsID, _ := strconv.Atoi(param)

	for _, curFeed := range feeds {
		if curFeed.Id == paramsID {
			WriteJsonResponse(w, http.StatusOK, curFeed)
			return
		}
	}
	WriteJsonResponse(w, http.StatusNotFound, "Newsfeed doesn't exist")
}

func CreateNewsFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newFeed newsfeed.Item
	err := json.NewDecoder(r.Body).Decode(&newFeed)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	newFeed.Id = ID
	feeds = append(feeds, newFeed)
	ID++
}

func DeleteNewsFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := chiro.URLParam(r, "id")
	paramsID, _ := strconv.Atoi(param)
	for index, curFeed := range feeds {
		if curFeed.Id == paramsID {
			feeds = append(feeds[:index], feeds[index+1:]...)
			break
		}
	}
	// fmt.Println(feeds)
	WriteJsonResponse(w, http.StatusOK, feeds)
}

func UpdateNewsFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := chiro.URLParam(r, "id")
	paramID, _ := strconv.Atoi(param)
	var newFeed newsfeed.Item
	err := json.NewDecoder(r.Body).Decode(&newFeed)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	for index, curFeed := range feeds {
		if curFeed.Id == paramID {
			newFeed.Id = paramID
			feeds[index] = newFeed
			err := json.NewEncoder(w).Encode(feeds)
			log.Println(err)
			return
		}
	}
	err = json.NewEncoder(w).Encode("No data")
	log.Println(err)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds newsfeed.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	correctPassword, ok := Credslist[creds.Username]
	if !ok || creds.Password != correctPassword {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Access Token (10 মিনিট)
	accessExpireTime := time.Now().Add(10 * time.Minute)
	_, accessToken, err := TokenAuth.Encode(map[string]interface{}{
		"aud": creds.Username,
		"exp": accessExpireTime.Unix(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Refresh Token (7 দিন)
	refreshExpireTime := time.Now().Add(7 * 24 * time.Hour)
	_, refreshToken, err := TokenAuth.Encode(map[string]interface{}{
		"aud": creds.Username,
		"exp": refreshExpireTime.Unix(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Refresh Token স্টোর করা হচ্ছে
	refreshTokens[refreshToken] = creds.Username

	// Access Token রেসপন্স হিসাবে পাঠানো
	w.Header().Set("Authorization", "Bearer "+accessToken)

	// Refresh Token HTTP-Only Cookie তে পাঠানো
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  refreshExpireTime,
		HttpOnly: true, // নিরাপত্তার জন্য
	})

	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Refresh Token Cookie মুছে ফেলা
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	})

	// Refresh Token মেমোরি থেকে মুছে ফেলা
	for token := range refreshTokens {
		delete(refreshTokens, token)
	}

	w.WriteHeader(http.StatusOK)
}

func StartServer(port int) {
	InitCred()
	InitDB()

	r := chiro.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/login", Login)
	r.Post("/refresh-token", RefreshToken) // Refresh Token API
	r.Group(func(r chiro.Router) {
		// jwtauth-> will learn later
		r.Use(jwtauth.Verifier(TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Route("/newsfeeds", func(r chiro.Router) {
			r.Get("/", GetNewsFeeds)
			r.Get("/{id}", GetNewsFeed)
			r.Post("/", CreateNewsFeed)
			r.Delete("/{id}", DeleteNewsFeed)
			r.Put("/{id}", UpdateNewsFeed)
		})
		r.Post("/logout", Logout)
	})
	// port := 5050
	Server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: r,
	}
	fmt.Println("Serving on " + strconv.Itoa(port))
	// http.ListenAndServe(strconv.Itoa(port), r)
	fmt.Println(Server.ListenAndServe())
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	refreshToken := cookie.Value
	username, exists := refreshTokens[refreshToken]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// new Access Token
	accessExpireTime := time.Now().Add(10 * time.Minute)
	_, newAccessToken, err := TokenAuth.Encode(map[string]interface{}{
		"aud": username,
		"exp": accessExpireTime.Unix(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// নAccess Token
	w.Header().Set("Authorization", "Bearer "+newAccessToken)
	w.WriteHeader(http.StatusOK)
}
