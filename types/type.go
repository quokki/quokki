package types

type QuokkiPower struct {
	Available int64 `json:"available"`
	Used      int64 `json:"used"`
	Reserved  int64 `json:"reserved"`
}
