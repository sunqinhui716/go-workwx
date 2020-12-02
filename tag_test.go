package workwx

import (
	"encoding/json"
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)



func TestAddTag(t *testing.T) {
	c.Convey("构造一个 reqTag", t, func() {
		tag := Tag{
			TagID:     1,
			TagName: "bar",
		}

		a := reqTag{
			tag,
		}

		c.Convey("执行序列化", func() {
			result, err := a.intoBody()

			c.Convey("序列化应该成功", func() {
				c.So(err, c.ShouldBeNil)

				c.Convey("序列化结果应该符合预期", func() {
					expectedPayload := []byte(`{
								"tagid": 1,
								"tagname": "bar"
								}`)
					var expected map[string]interface{}
					err := json.Unmarshal(expectedPayload, &expected)
					c.So(err, c.ShouldBeNil)

					var actual map[string]interface{}
					err = json.Unmarshal(result, &actual)
					c.So(err, c.ShouldBeNil)

					c.So(actual, c.ShouldResemble, expected)
				})
			})
		})
	})
}


func TestUpdateTag(t *testing.T) {
	c.Convey("构造一个 reqTag", t, func() {
		tag := Tag{
			TagID:     1,
			TagName: "bar",
		}

		a := reqTag{
			tag,
		}

		c.Convey("执行序列化", func() {
			result, err := a.intoBody()

			c.Convey("序列化应该成功", func() {
				c.So(err, c.ShouldBeNil)

				c.Convey("序列化结果应该符合预期", func() {
					expectedPayload := []byte(`{
								"tagid": 1,
								"tagname": "bar"
								}`)
					var expected map[string]interface{}
					err := json.Unmarshal(expectedPayload, &expected)
					c.So(err, c.ShouldBeNil)

					var actual map[string]interface{}
					err = json.Unmarshal(result, &actual)
					c.So(err, c.ShouldBeNil)

					c.So(actual, c.ShouldResemble, expected)
				})
			})
		})
	})
}
