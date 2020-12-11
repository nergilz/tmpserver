package server

import (
	"context"
	"net/http"

	"github.com/nergilz/tmpserver/utils"
)

// CtxKey это тип ключа для контекста запроса. используется для того, чтобы переменные,
// записанные в контекст в разных пакетах не пересекались
type CtxKey string

// Ключи для переменных, хранящихся в контексте запроса
const (
	СtxKeyUser CtxKey = "auth" // туть мы храним ключ для модели пользователя
)

// middleware
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken := r.Header.Get("Authorization")
		tokenMetadata, err := utils.VerifyJWTtoken(rawToken, s.us.GetSecret())
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(http.StatusText(http.StatusForbidden)))
			s.log.Warningf("not verify token : %v", err)
			return
		}
		userFromToken, err := utils.CheckJWTtoken(tokenMetadata)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			s.log.Warningf("ckeck token, not valid : %v", err)
			return
		}
		u, err := s.us.FindByID(userFromToken.ID)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(http.StatusText(http.StatusForbidden)))
			s.log.Warningf("user not find by id : %v", err)
			return
		}
		ctx := context.WithValue(r.Context(), СtxKeyUser, u) // add user по ключу в context
		next.ServeHTTP(w, r.WithContext(ctx))                // передали в обработчик
	})
}
