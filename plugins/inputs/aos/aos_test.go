package aos_test

import (
    "testing"
    "os"

    "github.com/influxdata/telegraf/plugins/inputs/aos"
    "github.com/influxdata/telegraf/plugins/inputs/aos/aos_streaming"
    "github.com/influxdata/telegraf/testutil"
    "github.com/stretchr/testify/assert"
)

var aos_server string = os.Getenv("AOS_SERVER")

func TestExtractAlertDataForProbeAlert(t *testing.T) {
    plugin := &aos.Aos{
        Port: 7777,
        Address: "blah",
        StreamingType: []string{"alerts"},
        AosServer: aos_server,
        AosPort: 443,
        AosLogin: "admin",
        AosPassword: "admin",
        AosProtocol: "https",
        RefreshInterval: 1000,
    }

    ssl := &aos.StreamAos{
        Listener: nil,
        Aos: plugin,
    }

    acc := testutil.Accumulator{}
    //plugin.Accumulator = acc
    assert.NoError(t, plugin.Start(&acc))

    pi := "p"
    sn := "s"
    alert := &aos_streaming.ProbeAlert{
        ExpectedInt: new(int64),
        ActualInt: new(int64),
        ProbeId: &pi,
        StageName: &sn,
        //KeyValuePairs: nil,
    }
    tags := make(map[string]string)
    ssl.ExtractAlertData("probe_alert", tags, alert, true)
}
