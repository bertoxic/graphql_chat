package posts

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePostInput_Sanitize(t *testing.T) {
	postData :=
		CreatePostInput{
			Title:   "title ",
			Content: "body   ",
		}
	postData.Sanitize()
	wantData :=
		CreatePostInput{
			Title:   "title",
			Content: "body",
		}
	require.Equal(t, wantData, postData)
}
func TestCreatePostInput_Validate(t *testing.T) {

}
