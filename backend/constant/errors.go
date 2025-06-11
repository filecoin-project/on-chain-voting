package constant

import "errors"

var ErrAlreadySyncHeight= errors.New("already synchronized to the latest block height")

var ErrGistNotFound = errors.New("gist not found")

var ErrAlgorithmMismatch = errors.New("voted algorithm is not equal to constant.VotingAlgorithm")