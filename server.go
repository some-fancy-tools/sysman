package sysman

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	ui "development.thatwebsite.xyz/utils/sysman/sysman/dist"
)

func (s *SysManager) Routes() {
	s.mux.Handle("/", http.FileServer(http.FS(ui.FS)))

	s.mux.HandleFunc("/api/v1/sysman/governors", s.handleGovernors())
	s.mux.HandleFunc("/api/v1/sysman/frequencies", s.handleFrequencies())
}

func (s *SysManager) handleFrequencies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		val, err := s.GetScalingFrequencies()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(val); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err.Error())
			return
		}
	}
}

func (s *SysManager) handleGovernors() http.HandlerFunc {
	type resp struct {
		AvailableGovernors map[string][]string `json:"availableGovernors"`
		SelectedGovernor   map[string]string   `json:"selectedGovernor"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			m1, m2, err := s.GetScalingGovernors()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, err.Error())
				return
			}
			w.Header().Set("content-type", "application/json")
			if err := json.NewEncoder(w).Encode(resp{AvailableGovernors: m2, SelectedGovernor: m1}); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, err.Error())
				return
			}
		case http.MethodPost:
			govs := []Governor{}
			if err := json.NewDecoder(r.Body).Decode(&govs); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, err.Error())
				return
			}
			if err := s.SetScalingGovernor(govs); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, err.Error())
				return
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			// fmt.Fprintln(w, err.Error())
			return

		}
	}
}

func (s *SysManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
