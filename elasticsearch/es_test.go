package elasticsearch_test

import (
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/webliupeng/gin-tonic/elasticsearch"
)

type mockServer struct {
}

func (d *mockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/idx_name/_count" {
		w.Write([]byte(`{
			"count": 16,
			"_shards": {
				"total": 2,
				"successful": 2,
				"failed": 0
			}
		}`))
	}

	if r.RequestURI == "/idx_name/_search" {
		w.Write([]byte(`
		{
			"took": 4,
			"timed_out": false,
			"_shards": {
				"total": 2,
				"successful": 2,
				"failed": 0
			},
			"hits": {
				"total": 16,
				"max_score": 3.7568405,
				"hits": [
					{
						"_index": "wx_bot_contacts",
						"_type": "contacts",
						"_id": "wxid_sv994slu64n122_7645904467@chatroom",
						"_score": 3.7568405,
						"_source": {
							"openid": "wxid_sv994slu64n122"
						}
					},
					{
						"_index": "wx_bot_contacts",
						"_type": "contacts",
						"_id": "wxid_sv994slu64n122_weixin",
						"_score": 3.7568405,
						"_source": {
							"openid": "wxid_sv994slu64n122"
						}
					},
					{
						"_index": "wx_bot_contacts",
						"_type": "contacts",
						"_id": "wxid_sv994slu64n122_fmessage",
						"_score": 3.7568405,
						"_source": {
							"openid": "wxid_sv994slu64n122"
						}
					},
					{
						"_index": "wx_bot_contacts",
						"_type": "contacts",
						"_id": "wxid_sv994slu64n122_gh_93769de9c356",
						"_score": 3.7568405,
						"_source": {
							"openid": "wxid_sv994slu64n122"
						}
					},
					{
						"_index": "wx_bot_contacts",
						"_type": "contacts",
						"_id": "wxid_sv994slu64n122_gh_20f7ecf66011",
						"_score": 3.7568405,
						"_source": {
							"openid": "wxid_sv994slu64n122"
						}
					},
					{
						"_index": "wx_bot_contacts",
						"_type": "contacts",
						"_id": "wxid_sv994slu64n122_anderslau",
						"_score": 3.7568405,
						"_source": {
							"openid": "wxid_sv994slu64n122"
						}
					},
					{
						"_index": "wx_bot_contacts",
						"_type": "contacts",
						"_id": "wxid_sv994slu64n122_5341696798@chatroom",
						"_score": 3.7568405,
						"_source": {
							"openid": "wxid_sv994slu64n122"
						}
					},
					{
						"_index": "wx_bot_contacts",
						"_type": "contacts",
						"_id": "wxid_sv994slu64n122_medianote",
						"_score": 3.541602,
						"_source": {
							"openid": "wxid_sv994slu64n122"
						}
					},
					{
						"_index": "wx_bot_contacts",
						"_type": "contacts",
						"_id": "wxid_sv994slu64n122_floatbottle",
						"_score": 3.541602,
						"_source": {
							"openid": "wxid_sv994slu64n122"
						}
					},
					{
						"_index": "wx_bot_contacts",
						"_type": "contacts",
						"_id": "wxid_sv994slu64n122_filehelper",
						"_score": 3.541602,
						"_source": {
							"openid": "wxid_sv994slu64n122"
						}
					}
				]
			}
		}`))
	}

}
func TestEs(t *testing.T) {

	ts := httptest.NewServer(&mockServer{})
	defer ts.Close()
	client := elasticsearch.NewClient(ts.URL + "/")

	res := client.Count("idx_name", gin.H{})
	assert.Equal(t, 16, res.Body().Count)

	res2 := client.Search("idx_name", gin.H{
		"size": 10,
	})
	assert.Equal(t, 10, len(res2.Body().Hits.Hits))
}
