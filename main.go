package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"twitter_clone/models"
	"twitter_clone/services"
)

// Get User By User Id
func getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Verify http request type
	verifyHttpMethod(w, r, http.MethodGet)

    // Get the "id" parameter from the URL
	idStr := r.URL.Query().Get("id")

	// Convert the "id" parameter to an integer
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "invalid ID", http.StatusBadRequest)
        return
    }

	// Get the user by ID
    user, err := services.GetUserById(id)
    if err != nil {
        http.Error(w, "error retrieving user", http.StatusInternalServerError)
        return
    }

    // Serialize user as JSON and send the response
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(user); err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// Get User Tweets by UserId
func getUserTweetsHandler(w http.ResponseWriter, r *http.Request) {
	// Verify http request type
	verifyHttpMethod(w, r, http.MethodGet)

	// Get the "id" parameter from the URL
    idStr := r.URL.Query().Get("id")

	// Convert the "id" parameter to an integer
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "invalid ID", http.StatusBadRequest)
        return
    }

	// Get the tweets by the tweets
    tweets, err := services.GetTweetsByUserId(id)
    if err != nil {
        http.Error(w, "error retrieving user tweets", http.StatusInternalServerError)
        return
    }

    // Serialize tweets as JSON and send the response
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(tweets); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// Get a suggested tweet for the logged-in user (provided they don't follow them yet)
func getSuggestedTweet(w http.ResponseWriter, r *http.Request) {
	// Verify http request type
	verifyHttpMethod(w, r, http.MethodGet)

	// Get the "id" parameter from the URL
	idStr := r.URL.Query().Get("id")

	// Convert the "id" parameter to an integer
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	// Get the tweet by ID
	tweet, err := services.GetSuggestedTweet(id)
	noTweetsError:= "sql.ErrNoRows"
	if err != nil {
		if err.Error() == noTweetsError {
			// Handle the case when there are no suggested tweets
			// Return a custom response indicating no suggested tweets found
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "ErrNoRows"}`))
		} else {
			http.Error(w, "error retrieving suggested tweet", http.StatusInternalServerError)
		}
		return
	}

	// Serialize the tweet as JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tweet); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Get User Feed by User Id (with pagination)
func getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// Verify http request type
	verifyHttpMethod(w, r, http.MethodGet)

	// Get the "id," "page," and "pageSize" parameters from the URL
    idStr := r.URL.Query().Get("id")
    pageStr := r.URL.Query().Get("page")
    pageSizeStr := r.URL.Query().Get("pageSize")

	// Convert the "id" parameter to an integer
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "invalid ID", http.StatusBadRequest)
        return
    }

	// Convert the "page" parameter to an integer
    page, err := strconv.Atoi(pageStr)
    if err != nil {
        http.Error(w, "invalid page", http.StatusBadRequest)
        return
    }

	// Convert the "pageSize" parameter to an integer
    pageSize, err := strconv.Atoi(pageSizeStr)
    if err != nil {
        http.Error(w, "invalid pageSize", http.StatusBadRequest)
        return
    }

	// Fetch tweetFeed for the specified user with pagination.
    tweetFeed, err := services.GetTweetFeedByUserId(id, page, pageSize)
    if err != nil {
        http.Error(w, "error retrieving data", http.StatusInternalServerError)
        return
    }

    // Serialize tweet feed as JSON and send the response
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(tweetFeed); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// Create tweet (data from request body)
func createTweetHandler(w http.ResponseWriter, r *http.Request) {
	// Verify http request type
	verifyHttpMethod(w, r, http.MethodPost)

    // Decode the request body into the existing Tweet struct
    var newTweet models.Tweet

    if err := json.NewDecoder(r.Body).Decode(&newTweet); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    // Call the CreateTweet function with the parsed tweet
    if err := services.CreateTweet(newTweet); err != nil {
        http.Error(w, "error creating the tweet", http.StatusInternalServerError)
        return
    }

    // Return a success response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newTweet)
}

// Follow user (data from request body)
func createFollowerHandler(w http.ResponseWriter, r *http.Request) {
	// Verify http request type
	verifyHttpMethod(w, r, http.MethodPut)
	
    // Create a new instance of the Follower struct to hold the request body data
    var follower models.Follower

    // Parse the request body into the Follower struct
    err := json.NewDecoder(r.Body).Decode(&follower)
    if err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    // Call the CreateFollower function to store the Follower data
    err = services.CreateFollower(follower)
    if err != nil {
        http.Error(w, "error creating follow request", http.StatusInternalServerError)
        return
    }

    // Return a success response
    w.WriteHeader(http.StatusCreated)
}

// Unfollow user (data from request body)
func unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	// Verify http request type
	verifyHttpMethod(w, r, http.MethodDelete)
	
    // Create a new instance of the Follower struct to hold the request body data
    var follower models.Follower

    // Parse the request body into the Follower struct
    err := json.NewDecoder(r.Body).Decode(&follower)
    if err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    // Call the CreateFollower function to store the Follower data
    err = services.DeleteFollowerFromUser(follower.UserID, follower.FollowerUserID)
    if err != nil {
        http.Error(w, "error deleting followed user", http.StatusInternalServerError)
        return
    }

    // Return a success response
    w.WriteHeader(http.StatusAccepted)
}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {
    // Verify http request type
	verificationError := verifyHttpMethod(w, r, http.MethodPost)
	if verificationError != nil {
		http.Error(w, verificationError.Error(), http.StatusBadRequest)
		return
	}

    // Parse the request body
    var requestBody struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&requestBody); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    // Get the user by username and password
    user, err := services.GetUserByUsernameAndPassword(requestBody.Username, requestBody.Password)
    if err != nil {
        http.Error(w, "invalid username or password", http.StatusBadRequest)
        return
    }

    // Serialize user as JSON and send the response
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// Mediator function for determining if it's a create or remove follower request
func followerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		unfollowUserHandler(w, r)
	} else if r.Method == http.MethodPut {
		createFollowerHandler(w, r)
	} else {
		http.Error(w, "invalid request", http.StatusBadRequest)
	}
}

// Helper method for verifying http request types
func verifyHttpMethod(w http.ResponseWriter, r *http.Request, expectedMethod string) error {
	if r.Method != expectedMethod {
		return errors.New("invalid request type")
	}
	return nil
}

func initNetHttpApi(port int) {
	// Get
    http.HandleFunc("/api/user", getUserByIdHandler)
    // Get
	http.HandleFunc("/api/user/tweet", getUserTweetsHandler)
    // Post
    http.HandleFunc("/api/user/login", userLoginHandler)
    // Post / Delete
	http.HandleFunc("/api/user/follower", followerHandler)
    // Get
	http.HandleFunc("/api/feed", getUserFeedHandler)
	// Get
	http.HandleFunc("/api/feed/suggested", getSuggestedTweet)
	// Post
	http.HandleFunc("/api/tweet", createTweetHandler)
    // Host SPA app
    http.Handle("/", http.FileServer(http.Dir("views/")))

    // Start the HTTP server
    portNumber:= fmt.Sprintf(":%d", port)
    log.Fatal(http.ListenAndServe(portNumber , nil))
    log.Printf("Server is listening on port %d\n", port)
}

func main() {
	initNetHttpApi(3000)
}