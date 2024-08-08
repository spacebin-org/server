/*
 * Copyright 2020-2024 Luke Whritenour

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 *     http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server_test

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/lukewhrit/spacebin/internal/database"
	"github.com/lukewhrit/spacebin/internal/database/databasefakes"
	"github.com/lukewhrit/spacebin/internal/server"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DocumentResponse struct {
	Payload database.Document
	Error   string
}

type FetchDocumentSuite struct {
	suite.Suite

	srv *server.Server
}

func (s *FetchDocumentSuite) SetupTest() {
	mockDB := &databasefakes.FakeDatabase{}

	mockDB.GetDocumentReturns(database.Document{
		ID:        "12345678",
		Content:   "test",
		CreatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
		UpdatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
	}, nil)

	s.srv = server.NewServer(&mockConfig, mockDB)
	s.srv.MountHandlers()
}

func (s *FetchDocumentSuite) TestFetchDocument() {
	req, _ := http.NewRequest(http.MethodGet, "/api/12345678", nil)

	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s.srv)

	checkResponseCode(s.T(), http.StatusOK, res.Result().StatusCode)

	x, _ := io.ReadAll(res.Result().Body)
	var body DocumentResponse
	json.Unmarshal(x, &body)

	expectedResponse := DocumentResponse{
		Payload: database.Document{
			ID:        "12345678",
			Content:   "test",
			CreatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
			UpdatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
		},
	}

	require.Equal(s.T(), expectedResponse.Payload, body.Payload)
}
