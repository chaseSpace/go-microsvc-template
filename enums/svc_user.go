package enums

type Sex int32

const (
	SexUnknown Sex = iota
	SexMale
	SexFemale
)

func (s Sex) Int32() int32 {
	return int32(s)
}
