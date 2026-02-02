package validation

type dataSigner interface {
	Sign(data any) (string, error)
}

type validator struct {
	signer dataSigner
}

func NewValidator(signer dataSigner) *validator {
	return &validator{
		signer: signer,
	}
}

func (v *validator) Validate(data any) (string, error) {
	// TODO: actually validate the data :p
	return v.signer.Sign(data)
}
