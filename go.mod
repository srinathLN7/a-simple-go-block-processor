module github.com/srinathLN7/flow/processor

go 1.18

replace (
	github.com/srinathLN7/flow/bc => ./src/bc
	github.com/srinathLN7/flow/consensus => ./src/consensus
	github.com/srinathLN7/flow/util => ./lib/util
)

require (
	github.com/TwiN/go-color v1.1.0
	github.com/srinathLN7/flow/bc v0.0.0-00010101000000-000000000000
	github.com/srinathLN7/flow/consensus v0.0.0-00010101000000-000000000000
	github.com/srinathLN7/flow/util v0.0.0-00010101000000-000000000000
)

require (
	github.com/joho/godotenv v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
)
