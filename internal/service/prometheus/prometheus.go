package prometheus

import (
	"context"
	"go-xops/internal/response"
	"go-xops/pkg/common"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

func scalarV(v *model.Scalar) response.PmtRep {
	var p response.PmtRep
	p.TimeStap = v.Timestamp.String()
	p.Value = v.Value.String()
	return p
}

func vectorV(v model.Vector) []response.PmtRep {
	var (
		ps []response.PmtRep
		p  response.PmtRep
	)
	for _, k := range v {
		p.TimeStap = k.Timestamp.String()
		p.Metric = k.Metric.String()
		p.Value = k.Value.String()
		ps = append(ps, p)
	}
	return ps
}

func matrixV(v model.Matrix) []response.PmtRep {
	var (
		ps []response.PmtRep
		p  response.PmtRep
	)
	for _, i := range v {
		for _, j := range i.Values {
			p.TimeStap = j.Timestamp.String()
			p.Value = j.Value.String()
			ps = append(ps, p)
		}
	}
	return ps
}

func stringV(v *model.String) response.PmtRep {
	var p response.PmtRep
	p.TimeStap = v.Timestamp.String()
	p.Value = v.Value
	return p
}

func justType(v model.Value) interface{} {
	switch v.Type() {
	case model.ValScalar:
		v, _ := v.(*model.Scalar)
		return scalarV(v)
	case model.ValVector:
		v, _ := v.(model.Vector)
		return vectorV(v)
	case model.ValMatrix:
		v, _ := v.(model.Matrix)
		return matrixV(v)
	case model.ValString:
		v, _ := v.(*model.String)
		return stringV(v)
	default:
		logrus.Println("uknow type")
		return nil
	}
}

func PrometheusAPIQuery_Test(key, job []string) ([]response.PmtRep, error) {
	var (
		wg   sync.WaitGroup
		i    int
		data []response.PmtRep
	)
	add := common.Conf.PrometheusApiAddress.Address
	if add == "" {
		return nil, errors.New("prometheus api 不能够为空")
	}
	client, err := api.NewClient(api.Config{
		Address: add,
	})
	if err != nil {
		logrus.Println(err)
		os.Exit(1)
	}
	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, v := range job {
		i++
		wg.Add(1)
		go func(job string) {
			for _, j := range key {
				p := j + "{job='" + v + "'}"
				value, _, err := v1api.Query(ctx, p, time.Now())
				if err != nil {
					os.Exit(1)
				}
				h := justType(value)
				if value, ok := h.([]response.PmtRep); ok {
					data = append(data, value...)
				} else if value, ok := h.(response.PmtRep); ok {
					data = append(data, value)
				} else {
					errors.New("类型错误")
				}
			}
			wg.Done()
		}(v)
		if i == 30 {
			wg.Wait()
			i = 0
		}
		wg.Wait()
	}
	return data, nil
}
