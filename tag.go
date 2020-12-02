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

// UpdateTag 编辑标签
func (c *WorkwxApp) UpdateTag(req Tag) error {
	_, err := c.execEditTag(reqEditTag{
		TagName: req.TagName,
		TagID:   req.TagID,
	})
	return err
}
