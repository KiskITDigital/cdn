package auth

import (
	"context"

	"gitlab.ubrato.ru/ubrato/cdn/internal/lib/cerr"
	"gitlab.ubrato.ru/ubrato/cdn/internal/lib/token"
)

func (s *Service) ValidateAccessToken(ctx context.Context, accessToken string) (token.Claims, error) {
	claims, err := s.tokenAuthorizer.ValidateToken(accessToken)
	if err != nil {
		return token.Claims{}, cerr.Wrap(cerr.ErrAuthorize, cerr.CodeUnauthorized, "token not valid", nil)
	}

	return claims, nil
}
