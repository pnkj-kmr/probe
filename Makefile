

# // @Summary Get user by ID
# // @Description Get user information by user ID
# // @Tags users
# // @Accept json
# // @Produce json
# // @Param id path string true "User ID"
# // @Param verbose query bool false "Verbose mode"
# // @Param input body UserRequest true "User request body"
# // @Success 200 {object} UserResponse
# // @Failure 400 {object} ErrorResponse
# // @Router /users/{id} [get]

# IncludeDetails bool `json:"includeDetails" example:"true"`

docgen:
	/Users/pankaj/Workspace/golang/bin/swag init

docgen:
	/Users/pankaj/Workspace/golang/bin/goweight probe

run:
	go run main.go

air:
	air

test:
	go run tmp/producer/main.go \
	-brokers pankaj.local:9092 \
	-version "3.0.1" \
	-username "infraon-producer" -passwd "infraon@123" \
	-algorithm "sha512" \
	-ca "/Users/pankaj/Workspace/golang/src/probe/tmp/ca-cert" \
	-certificate "/Users/pankaj/Workspace/golang/src/probe/tmp/ca-cert" \
	-tls \
	-tls-skip-verify



test2:
	go run tmp/producer2/main.go \
	-brokers pankaj.local:9092 \
	-version "3.0.1" \
	-username "infraon-producer" -passwd "infraon@123" \
	-algorithm "sha512" \
	-ca "/Users/pankaj/Workspace/golang/src/probe/tmp/ca-cert" \
	-certificate "/Users/pankaj/Workspace/golang/src/probe/tmp/ca-cert" \
	-tls \
	-tls-skip-verify \
	-c 100

