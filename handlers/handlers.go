package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"htmx-learn/db"
	"htmx-learn/templates/components"
	"htmx-learn/templates/pages"
)

type Handlers struct {
	counterStore *db.CounterStore
	userStore    *db.UserStore
}

func New(database *db.DB) *Handlers {
	return &Handlers{
		counterStore: db.NewCounterStore(database),
		userStore:    db.NewUserStore(database),
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if err := pages.Home().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) CounterPage(w http.ResponseWriter, r *http.Request) {
	count, err := h.counterStore.Get()
	if err != nil {
		log.Printf("Error getting counter: %v", err)
		count = 0
	}
	if err := pages.CounterPage(count).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) DynamicPage(w http.ResponseWriter, r *http.Request) {
	if err := pages.DynamicPage().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) CounterIncrement(w http.ResponseWriter, r *http.Request) {
	count, err := h.counterStore.Increment()
	if err != nil {
		log.Printf("Error incrementing counter: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := components.CountDisplay(count).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) CounterDecrement(w http.ResponseWriter, r *http.Request) {
	count, err := h.counterStore.Decrement()
	if err != nil {
		log.Printf("Error decrementing counter: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := components.CountDisplay(count).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) CounterReset(w http.ResponseWriter, r *http.Request) {
	count, err := h.counterStore.Reset()
	if err != nil {
		log.Printf("Error resetting counter: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := components.CountDisplay(count).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetTime(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	if err := components.TimeDisplay(currentTime).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userStore.GetAll()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	convertedUsers := make([]components.User, len(users))
	for i, user := range users {
		convertedUsers[i] = components.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	}
	
	for _, user := range convertedUsers {
		if err := components.UserCard(user).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	
	name := r.FormValue("user-name")
	email := r.FormValue("user-email")
	
	if name == "" || email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}
	
	user, err := h.userStore.Add(name, email)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	convertedUser := components.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	
	if err := components.UserCard(convertedUser).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	err = h.userStore.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Error deleting user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) SearchUsers(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	
	query := r.FormValue("search")
	users, err := h.userStore.Search(query)
	if err != nil {
		log.Printf("Error searching users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	convertedUsers := make([]components.User, len(users))
	for i, user := range users {
		convertedUsers[i] = components.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	}
	
	if err := components.SearchResults(convertedUsers).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}