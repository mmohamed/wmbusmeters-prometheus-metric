package collectors

import (
	"encoding/json"
	"fmt"
	"bytes"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	waterCubicUsedKey  = "water_cubic_used"
	heatKwhConsumptionKey  = "heat_kwh_consumption"
	heatCubicConsumptionKey  = "heat_cubic_consumption"
	heatKwhFlowKey = "heat_kwh_flow"
	heatCubicFlowKey = "heat_cubic_flow"
	heatTemperatureFlowKey = "heat_temperature_flow"
	heatTemperatureReturnKey = "heat_temperature_return"
)

type wmbusmetersCollector struct {
	inMemoriesData map[string]string
}

// NewWMBusmetersCollector creates a new wmbusmeters prometheus collector.
func NewWMBusmetersCollector(inMemoriesData map[string]string) prometheus.Collector {
	return &wmbusmetersCollector{inMemoriesData: inMemoriesData}
}

// Describe implements the prometheus.Collector interface.
func (collector *wmbusmetersCollector) Describe(ch chan<- *prometheus.Desc) {

}

// Helper to check element in slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// Collect implements the prometheus.Collector interface.
func (collector *wmbusmetersCollector) Collect(ch chan<- prometheus.Metric) {

	addGauge := func(key string, help string, data map[string]interface{}, v float64) {
		labelValues := []string{"total_m3", "total_energy_consumption_kwh", "total_volume_m3", "volume_flow_m3h", "power_kw", "flow_temperature_c", "return_temperature_c"}
		labels := []string{"id", "name", "meter", "timestamp"}
		values := []string{data["id"].(string), data["name"].(string), data["meter"].(string), data["timestamp"].(string)}

		for key, value := range data {
			if !contains(labels, key) && !contains(labelValues, key) {
				var temp_buff bytes.Buffer
				fmt.Fprintf(&temp_buff, "%v", value)
				labels = append(labels, key)
				values = append(values, temp_buff.String())
			} 
		}
		
		desc := prometheus.NewDesc(key, help, labels, nil,)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, v, values...)
	}

	count := 0

	for id, metric := range collector.inMemoriesData {
		var parsedData map[string]interface{}
		json.Unmarshal([]byte(metric), &parsedData)

		if parsedData["media"].(string) == "water" {
			addGauge(waterCubicUsedKey, "Cubic meters of water used", parsedData, parsedData["total_m3"].(float64))
			delete(collector.inMemoriesData, id)
			count++
		}

		if parsedData["media"].(string) == "heat" {
			addGauge(heatKwhConsumptionKey, "KWH Heat consumption", parsedData, parsedData["total_energy_consumption_kwh"].(float64))
			addGauge(heatCubicConsumptionKey, "Cubic Heat consumption", parsedData, parsedData["total_volume_m3"].(float64))
			addGauge(heatCubicFlowKey, "Cubic Heat flow", parsedData, parsedData["volume_flow_m3h"].(float64))
			addGauge(heatKwhFlowKey, "KWH Heat flow", parsedData, parsedData["power_kw"].(float64))
			addGauge(heatTemperatureFlowKey, "Celsius Heat flow temperature", parsedData, parsedData["flow_temperature_c"].(float64))
			addGauge(heatTemperatureReturnKey, "Celsius Heat return temperature", parsedData, parsedData["return_temperature_c"].(float64))
			delete(collector.inMemoriesData, id)
			count++
		}
	}
	glog.Infof("Found metrics %d ", count)
}
