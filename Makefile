

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

