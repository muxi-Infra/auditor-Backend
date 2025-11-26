package internal

type Contents struct {
	Topic       Topics  `json:"topic"`
	LastComment Comment `json:"last_comment"`
	NextComment Comment `json:"next_comment"`
}

type Topics struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Pictures []string `json:"pictures"`
}

type Comment struct {
	Content  string   `json:"content"`
	Pictures []string `json:"pictures"`
}

type ContentOption func(*Contents)
type TopicOption func(*Topics)
type CommentOption func(*Comment)

// NewContents 你应当始终使用此函数来创建对象
func NewContents(opts ...ContentOption) *Contents {
	c := &Contents{}
	if len(opts) == 0 {
		return nil
	}
	for _, opt := range opts {
		opt(c)
	}

	if c.LastComment.Pictures == nil {
		c.LastComment.Pictures = []string{}
	}
	if c.NextComment.Pictures == nil {
		c.NextComment.Pictures = []string{}
	}

	return c
}

func WithTopicText(title, content string) ContentOption {
	return func(c *Contents) {
		c.Topic.Title = title
		c.Topic.Content = content
	}
}

func WithTopicPictures(pics []string) ContentOption {
	return func(c *Contents) {
		c.Topic.Pictures = pics
	}
}

func WithLastCommentText(content string) ContentOption {
	return func(c *Contents) {
		c.LastComment.Content = content
	}
}

func WithLastCommentPictures(pics []string) ContentOption {
	return func(c *Contents) {
		c.LastComment.Pictures = pics
	}
}

func WithNextCommentText(content string) ContentOption {
	return func(c *Contents) {
		c.NextComment.Content = content
	}
}

func WithNextCommentPictures(pics []string) ContentOption {
	return func(c *Contents) {
		c.NextComment.Pictures = pics
	}
}

func WithCommentPictures(pics []string) CommentOption {
	return func(c *Comment) {
		c.Pictures = pics
	}
}
