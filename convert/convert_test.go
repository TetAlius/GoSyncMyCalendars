package convert_test

import (
	"testing"
	"time"

	"reflect"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/convert"
	"github.com/TetAlius/GoSyncMyCalendars/logger"
)

type first struct {
	Field1    string `convert:"field_1"`
	Field2    int    `convert:"field_2"`
	Field3    string `convert:"field_3"`
	Field4    int    `convert:"field_4"`
	Field5    string `convert:"field_5"`
	Field6    int    `convert:"field_6"`
	Something *third `convert:"something"`
}

type second struct {
	Field7    string `convert:"field_1"`
	Field8    int    `convert:"field_2"`
	Field9    string `convert:"field_3"`
	Field10   int    `convert:"field_4"`
	Field11   string `convert:"field_5"`
	Field12   int    `convert:"field_6"`
	Something *third `convert:"something"`
}

type third struct {
	Field123 string `convert:"field_123"`
}

func (*third) Convert(m interface{}, tag string, opts string) (convert.Converter, error) {
	logger.Debugf("Convert: %s", m)
	return &third{Field123: m.(map[string]interface{})["field_123"].(string)}, nil
}
func (this *third) Deconvert() interface{} {
	logger.Debugln("Deconverting")
	return map[string]interface{}{"field_123": this.Field123}
}

func TestConvert(t *testing.T) {
	from := first{Field1: "ASd11", Field2: 2, Field3: "ASD3", Field4: 4, Field5: "ASd5", Field6: 6, Something: &third{Field123: "ASDASD"}}
	to := new(second)

	logger.Debugf("%s", reflect.TypeOf(from.Field1))
	err := convert.Convert(from, to)
	if err != nil {
		t.Fatalf("error converting: %s", err.Error())
	}
	logger.Debugf("%s", to.Field7)
	logger.Debugf("%d", to.Field8)
	logger.Debugf("%s", to.Field9)
	logger.Debugf("%d", to.Field10)
	logger.Debugf("%s", to.Field11)
	logger.Debugf("%d", to.Field12)
	logger.Debugf("%s", to.Something)
}

func TestConvertEvent(t *testing.T) {
	var event api.GoogleEvent
	start := new(api.GoogleTime)
	end := new(api.GoogleTime)
	now := time.Now()
	more := now.Add(time.Hour * 2)
	start.DateTime = now
	end.DateTime = more
	event.Start = start
	event.End = end
	eventOut := new(api.OutlookEvent)

	t.Logf("start: %s", now.UTC().Format(time.RFC3339))
	t.Logf("end: %s", more.UTC().Format(time.RFC3339))
	err := convert.Convert(event, eventOut)
	if err != nil {
		t.Fatalf("error converting event: %s", err.Error())
	}
	if event.Start.DateTime.UTC().Format(time.RFC3339) != eventOut.Start.DateTime.UTC().Format(time.RFC3339) {
		t.Fatalf("convertion of start went wrong g: %s, o: %s", event.Start.DateTime.UTC().Format(time.RFC3339), eventOut.Start.DateTime.UTC().Format(time.RFC3339))
	}
	if event.End.DateTime.UTC().Format(time.RFC3339) != eventOut.End.DateTime.UTC().Format(time.RFC3339) {
		t.Fatalf("convertion of end went wrong g: %s, o: %s", event.End.DateTime.UTC().Format(time.RFC3339), eventOut.End.DateTime.UTC().Format(time.RFC3339))
	}
}
