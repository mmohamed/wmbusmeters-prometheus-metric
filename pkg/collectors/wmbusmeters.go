package collectors

import (
	"encoding/json"

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

var (
	waterCubicUsed = prometheus.NewDesc(
		waterCubicUsedKey,
		"Cubic meters of water used",
		[]string{"id", "name", "meter", "timestamp", "prefix", "serial_number", "current_alarms", "previous_alarms"}, nil,
	)
	heatKwhConsumption = prometheus.NewDesc(
		heatKwhConsumptionKey,
		"KWH Heat consumption",
		[]string{"id", "name", "meter", "timestamp"}, nil,
	)
	heatCubicConsumption = prometheus.NewDesc(
		heatCubicConsumptionKey,
		"Cubic Heat consumption",
		[]string{"id", "name", "meter", "timestamp"}, nil,
	)
	heatCubicFlow = prometheus.NewDesc(
		heatCubicFlowKey,
		"Cubic Heat flow",
		[]string{"id", "name", "meter", "timestamp"}, nil,
	)
	heatKwhFlow = prometheus.NewDesc(
		heatKwhFlowKey,
		"KWH Heat flow",
		[]string{"id", "name", "meter", "timestamp"}, nil,
	)
	heatTemperatureFlow = prometheus.NewDesc(
		heatTemperatureFlowKey,
		"Celsius Heat flow temperature",
		[]string{"id", "name", "meter", "timestamp"}, nil,
	)
	heatTemperatureReturn = prometheus.NewDesc(
		heatTemperatureReturnKey,
		"Celsius Heat return temperature",
		[]string{"id", "name", "meter", "timestamp"}, nil,
	)
	
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
	ch <- waterCubicUsed
}

// Collect implements the prometheus.Collector interface.
func (collector *wmbusmetersCollector) Collect(ch chan<- prometheus.Metric) {

	addWaterGauge := func(desc *prometheus.Desc, data map[string]interface{}, v float64, lv ...string) {
		lv = append([]string{data["id"].(string), data["name"].(string), data["meter"].(string), data["timestamp"].(string), data["prefix"].(string), data["serial_number"].(string), data["current_alarms"].(string), data["previous_alarms"].(string)}, lv...)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, v, lv...)
	}

	addHeatGauge := func(desc *prometheus.Desc, data map[string]interface{}, v float64, lv ...string) {
		lv = append([]string{data["id"].(string), data["name"].(string), data["meter"].(string), data["timestamp"].(string)}, lv...)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, v, lv...)
	}

	count := 0

	for id, metric := range collector.inMemoriesData {
		var parsedData map[string]interface{}
		json.Unmarshal([]byte(metric), &parsedData)

		if parsedData["media"].(string) == "water" {
			addWaterGauge(waterCubicUsed, parsedData, parsedData["total_m3"].(float64))
			delete(collector.inMemoriesData, id)
			count++
		}

		if parsedData["media"].(string) == "heat" {
			addHeatGauge(heatKwhConsumption, parsedData, parsedData["total_energy_consumption_kwh"].(float64))
			addHeatGauge(heatCubicConsumption, parsedData, parsedData["total_volume_m3"].(float64))
			addHeatGauge(heatCubicFlow, parsedData, parsedData["volume_flow_m3h"].(float64))
			addHeatGauge(heatKwhFlow, parsedData, parsedData["power_kw"].(float64))
			addHeatGauge(heatTemperatureFlow, parsedData, parsedData["flow_temperature_c"].(float64))
			addHeatGauge(heatTemperatureReturn, parsedData, parsedData["return_temperature_c"].(float64))
			delete(collector.inMemoriesData, id)
			count++
		}
	}
	glog.Infof("Found metrics %d ", count)
}
