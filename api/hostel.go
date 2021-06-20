package api

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Hostel struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	DateOfOrigin time.Time `json:"dateoforigin"`
	NumOfRooms   int       `json:"numofrooms"`
}

type HostelHandler struct {
	sync.Mutex
	store map[string]*Hostel
}

func NewHostelHandler() *HostelHandler {
	return &HostelHandler{
		store: map[string]*Hostel{
			"h1": &Hostel{
				ID:           "h1",
				DateOfOrigin: time.Time(time.Now()),
				Name:         "Tagore Bhavan",
				NumOfRooms:   122,
			},
		},
	}
}

func hostelHandler(r *mux.Router) {

	h := NewHostelHandler()

	r.HandleFunc("", h.listHostels).Methods("GET")
	r.HandleFunc("/{id}", h.getHostel).Methods("GET")
	r.HandleFunc("", h.createHostel).Methods("POST")
	r.HandleFunc("/{id}", h.updateHostel).Methods("PUT")
	r.HandleFunc("/{id}", h.deleteHostel).Methods("DELETE")
}

// Get the list of Hostels
func (h *HostelHandler) listHostels(w http.ResponseWriter, r *http.Request) {
	hostels := make([]*Hostel, len(h.store))

	h.Lock()
	i := 0
	for _, coaster := range h.store {
		hostels[i] = coaster
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(hostels)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// Get a single Hostel object
func (h *HostelHandler) getHostel(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	h.Lock()
	ok := h.isValidId(id)
	hostel := h.store[id]
	h.Unlock()

	if ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("This id is not valid!"))
		return
	}

	jsonBytes, err := json.Marshal(hostel)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// Create a Hostel object
func (h *HostelHandler) createHostel(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("need Content-Type 'application/json'"))
		return
	}

	var hostel Hostel
	err = json.Unmarshal(bodyBytes, &hostel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	id := fmt.Sprintf("I%v", hash(hostel.Name))
	hostel.ID = id

	if ok, err := h.isValidData(&hostel); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	h.Lock()
	h.store[hostel.ID] = &hostel
	defer h.Unlock()
}

// Update the Hostel object
func (h *HostelHandler) updateHostel(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("need Content-Type 'application/json'"))
		return
	}

	var hostel Hostel
	err = json.Unmarshal(bodyBytes, &hostel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	newId := fmt.Sprintf("I%v", hash(hostel.Name))

	h.Lock()
	defer h.Unlock()

	if id == newId {
		if ok, err := h.isValidData(&hostel); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		hostel.ID = newId

		if ok, err := h.isValidData(&hostel); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		delete(h.store, id)
	}

	h.store[newId] = &hostel
}

// Delete the Hostel object
func (h *HostelHandler) deleteHostel(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	h.Lock()
	if _, ok := h.store[id]; ok {
		delete(h.store, id)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No such hostel exists"))
	}
	h.Unlock()
}

// Used for hashing ID to uint32
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// Can be used to check whether proper input was provided
func (h *HostelHandler) isValidId(id string) bool {
	_, ok := h.store[id]
	if ok {
		return false
	}
	return true
}

// Can be used to validate the data entered by User
// To be used while updating or creating objects
func (h *HostelHandler) isValidData(info *Hostel) (bool, error) {

	if ok := h.isValidId(info.ID); !ok {
		return false, fmt.Errorf("Not a valid ID!")
	}
	// h.isUniqueName(info.Name)

	return true, nil
}
