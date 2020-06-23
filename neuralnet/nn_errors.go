package neuralnet

type NoMovesError struct {
	Err error
}

func (nme *NoMovesError) Error() string {
	return "There are no new moves for this brain"
}

type HitWallError struct {
	Err error
}

func (hwe *HitWallError) Error() string {
	return "The brain hit a wall"
}
