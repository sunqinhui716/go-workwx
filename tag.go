package workwx

// AddTag 添加标签
func (c *WorkwxApp) AddTag(req Tag) (int, error) {
	resp, err := c.execTag(reqTag{
		Tag: req,
	})
	if err != nil {
		return 0, err
	}
	return resp.TagID, nil
}

