package video

// VideoService service placeholder
type VideoService struct{}

func NewVideoService() *VideoService {
    return &VideoService{}
}

func (v *VideoService) Generate(prompt string) (string, error) {
    return "video generated", nil
}
