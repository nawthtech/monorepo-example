package video

// أنواع الفيديو
type VideoProvider interface {
    GenerateVideo(prompt string, options VideoOptions) (*VideoResponse, error)
    Name() string
}

type VideoRequest struct {
    Prompt     string
    Duration   int
    Resolution string
    Aspect     string
    Style      string
}

type VideoResponse struct {
    URL       string
    Duration  int
    Width     int
    Height    int
    Format    string
    Provider  string
    Timestamp int64
}