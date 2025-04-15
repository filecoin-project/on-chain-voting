// Copyright (C) 2023-2024 StorSwift Inc.
// This file is part of the PowerVoting library.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import "time"

type DeveloperPower struct {
	Account   string `json:"account"`
	Power     int    `json:"power"`
	TimeStamp int    `json:"time_stamp"`
}

// Developer represents the developer structure containing information about the developer's repository.
type Developer struct {
	Data struct {
		Repository struct {
			DefaultBranchRef struct {
				Target struct {
					History struct {
						Nodes    []Nodes `json:"nodes"`
						PageInfo struct {
							EndCursor   string `json:"endCursor"`
							HasNextPage bool   `json:"hasNextPage"`
						} `json:"pageInfo"`
					} `json:"history"`
				} `json:"target"`
			} `json:"defaultBranchRef"`
		} `json:"repository"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}
type PaginationVariables struct {
	Cursor *string `json:"cursor"`
}

// Nodes represents the nodes structure containing details about committed date, author, and committer.
type Nodes struct {
	CommittedDate time.Time `json:"committedDate"` // Date of commit.
	Author        Author    `json:"author"`        // Author information.
	Committer     Committer `json:"committer"`     // Committer information.
}

// Author represents the author structure containing details about the user.
type Author struct {
	User User `json:"user"` // User information.
}

// Committer represents the committer structure containing details about the user.
type Committer struct {
	User User `json:"user"` // User information.
}

// User represents the user structure containing details like login.
type User struct {
	Login string `json:"login"` // User login information.
}
