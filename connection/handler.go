package connection

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/ladon"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
)

const (
	connectionsResource = "rn:hydra:connections"
	connectionResource  = "rn:hydra:connection:%s"
	scope               = "hydra.connections"
)

type Handler struct {
	Manager Manager
	H       herodot.Herodot
	W       warden.Warden
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST("/oauth2/connections", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r.URL.Query().Get("local") != "" {
			h.FindLocal(w, r, ps)
			return
		}

		if r.URL.Query().Get("remote") != "" && r.URL.Query().Get("provider") != "" {
			h.FindRemote(w, r, ps)
			return
		}

		h.Get(w, r, ps)
	})

	r.GET("/oauth2/connections/:id", h.Get)
	r.DELETE("/oauth2/connections/:id", h.Delete)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var conn Connection
	var ctx = context.Background()

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: connectionsResource,
		Action:   "create",
	}, scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&conn); err != nil {
		h.H.WriteErrorCode(ctx, w, r, http.StatusBadRequest, err)
		return
	}

	if v, err := govalidator.ValidateStruct(conn); err != nil {
		h.H.WriteErrorCode(ctx, w, r, http.StatusBadRequest, err)
		return
	} else if !v {
		h.H.WriteErrorCode(ctx, w, r, http.StatusBadRequest, errors.New("Payload did not validate."))
		return
	}

	conn.ID = uuid.New()
	if err := h.Manager.Create(&conn); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.WriteCreated(ctx, w, r, "/oauth2/connections/"+conn.ID, &conn)
}

func (h *Handler) FindLocal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: connectionsResource,
		Action:   "find",
	}, scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	conns, err := h.Manager.FindAllByLocalSubject(r.URL.Query().Get("local"))
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.Write(ctx, w, r, conns)
}

func (h *Handler) FindRemote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: connectionsResource,
		Action:   "find",
	}, scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	conns, err := h.Manager.FindByRemoteSubject(r.URL.Query().Get("provider"), r.URL.Query().Get("local"))
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.Write(ctx, w, r, conns)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()
	var id = ps.ByName("id")

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: fmt.Sprintf(connectionResource, id),
		Action:   "get",
	}, scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	conn, err := h.Manager.Get(id)
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.Write(ctx, w, r, conn)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()
	var id = ps.ByName("id")

	if _, err := h.W.HTTPActionAllowed(ctx, r, &ladon.Request{
		Resource: fmt.Sprintf(connectionResource, id),
		Action:   "delete",
	}, scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := h.Manager.Delete(ps.ByName("id")); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
