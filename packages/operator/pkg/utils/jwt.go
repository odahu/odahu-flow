/*
 * Copyright 2020 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/user"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
)

// Errors
var (
	ErrMalformedJWT    = errors.New("malformed JWT")
	ErrClaimExtraction = errors.New("claim extraction failed")
)

func ExtractUserInfoFromToken(token string, config config.Claims) (*user.UserInfo, error) {
	// Parse a token, but don't verify it
	// Ignore the error, but make sure token isn't nil in case of there were parsing errors
	parsedToken, _ := jwt.Parse(token, nil)
	if parsedToken == nil {
		return nil, ErrMalformedJWT
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrClaimExtraction
	}

	userName, err := extractClaim(claims, config.Name)
	if err != nil {
		return nil, err
	}

	userEmail, err := extractClaim(claims, config.Email)
	if err != nil {
		return nil, err
	}

	return &user.UserInfo{
		Username: userName,
		Email:    userEmail,
	}, nil
}

func extractClaim(claims jwt.MapClaims, claimName string) (string, error) {
	claimValue, ok := claims[claimName]
	if !ok {
		return "", fmt.Errorf("%s claim is missing", claimName)
	}

	claimValueStr, ok := claimValue.(string)
	if !ok {
		return "", fmt.Errorf("%s claim is not the string type", claimName)
	}

	return claimValueStr, nil
}
