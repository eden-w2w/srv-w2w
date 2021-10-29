package v0

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-w2w/lib-modules/clients/gaode"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/mozillazg/go-pinyin"
	"strings"
)

func init() {
	Router.Register(courier.NewRouter(GetDistricts{}))
}

// GetDistricts
type GetDistricts struct {
	httpx.MethodGet

	Keyword string `in:"query" name:"keyword" default:""`
}

func (req GetDistricts) Path() string {
	return "/districts"
}

type GetDistrictsResponse struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func (req GetDistricts) Output(ctx context.Context) (result interface{}, err error) {
	request := gaode.DistrictRequest{
		Keywords:    req.Keyword,
		SubDistrict: 1,
		Extensions:  "base",
		Page:        1,
		Offset:      20,
	}
	resp, err := global.Config.ClientGaode.District(request)
	if err != nil {
		return
	}

	if len(resp.Districts) == 0 {
		return nil, nil
	}

	response := make(map[string][]GetDistrictsResponse)
	root := resp.Districts[0]
	for _, district := range root.Districts {
		py := pinyin.LazyConvert(district.Name, nil)
		c := strings.ToUpper(string(py[0][0]))
		if _, ok := response[c]; !ok {
			response[c] = make([]GetDistrictsResponse, 0)
		}
		item := GetDistrictsResponse{
			Value: district.ADCode,
			Label: district.Name,
		}
		if district.Level == "street" {
			item.Value = district.Name
		}
		response[c] = append(response[c], item)
	}
	return response, nil
}
