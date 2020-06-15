package web

import (
	"context"
	"database/sql"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

func NewSessionManager(dataSourceName string) (*scs.SessionManager, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	sessions := scs.New()
	sessions.Store = postgresstore.New(db)

	return sessions, nil
}

type SessionData struct {
	FlashMessage string
	// UserID uuid.UUID
}

func GetSessionData(session *scs.SessionManager, ctx context.Context) SessionData {
	var data SessionData

	data.FlashMessage = session.PopString(ctx, "flash")
	// data.UserID, _ = session.Get(ctx, "user_id").(uuid.UUID)

	return data
}
