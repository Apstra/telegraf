package aos_test

import (
    "testing"
    "github.com/influxdata/telegraf/plugins/inputs/aos"
    "github.com/influxdata/telegraf/plugins/inputs/aos/aos_streaming"
    "github.com/influxdata/telegraf/testutil"
    "github.com/stretchr/testify/assert"
)

func TestExtractProbeMessage(t *testing.T) {
    plugin := &aos.Aos{
        Port: 7777,
        Address: "blah",
        StreamingType: []string{"alerts"},
        AosServer: "blah",
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
    assert.NoError(t, plugin.Start(&acc))
    val := &aos_streaming.ProbeMessage_Int64Value{Int64Value:10}
    alert := &aos_streaming.ProbeMessage{
        Value: val,
    }
    ssl.ExtractProbeData(alert, "probe_msg_1")
    tag_value := acc.TagValue("probe_message", "device")
    plugin.Stop()
    assert.Equal(
        t, tag_value, "probe_msg_1", "The probe message was not added.")
    
}

func TestExtractAlertDataForProbeAlert(t *testing.T) {
    plugin := &aos.Aos{
        Port: 7778,
        Address: "blah",
        StreamingType: []string{"alerts"},
        AosServer: "blah",
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
    tags["device"] = "probe_alert"
    ssl.ExtractAlertData("probe_alert", tags, alert, true)
    tag_value := acc.TagValue("alert_probe", "device")
    plugin.Stop()
    assert.Equal(
        t, tag_value, "probe_alert", "The alert message was not added.")
}