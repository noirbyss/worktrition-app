package service

type CompleteWaterRequest struct {
	UserID   string
	AmountMl int32
}

func (cwr CompleteWaterRequest) validate() error {
	if cwr.UserID == "" {
		return ErrEmptyUserID
	}

	if cwr.AmountMl <= 0 {
		return ErrAmountMlTooLow
	}

	return nil
}
