package secondary

type TextEncoder interface {
    Encode(text string) ([]float32, error)
}