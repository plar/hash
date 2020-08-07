package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/plar/hash/domain"
	"github.com/plar/hash/domain/repository"
	"github.com/plar/hash/service/hasher"
)

type hasherHandler struct {
	svc hasher.Service
}

func (h hasherHandler) createHash(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	hashID, err := h.svc.Create(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(hashID.Bytes())
}

func (h hasherHandler) getHash(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(requestField(r, 0), 10, 64)
	hash, err := h.svc.Get(domain.HashID(id))
	if err != nil {
		if errors.Is(err, repository.ErrHashNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	w.Write(hash.Base64())
}
