
# WMBusMeters Prometheus Metrics (exporter)

## Deploy & Configure
### Deploy the exporter
```bash
# deploy the exporter
kubectl apply -f sample-deploy.yaml
# test get metrics
curl wmbusmeters-prometheus-metric.kube-system.svc.cluster.local:9001/metrics
# test push metrics
curl wmbusmeters-prometheus-metric.kube-system.svc.cluster.local:9001/push -d "$METER_JSON"
```
### Configure the WMBUS
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: wmbus
  namespace: wmbus
  labels:
    app: wmbus
spec:
    ...
    spec:
      containers:
      - name: wmbus
        image: wmbusmeters/wmbusmeters
        command: 
        - "sh"
        - "-c"
        - /wmbusmeters/wmbusmeters --format=json --listento=t1 --shell='/usr/bin/curl wmbusmeters-prometheus-metric.kube-system.svc.cluster.local:9001/push -s --data "$METER_JSON"' ...
        securityContext:
          privileged: true
        volumeMounts:
        - mountPath: /dev/
          name: dev
      volumes:
        - name: dev
          hostPath:
            path: /dev/
```

## Metrics

| Metric name | Metric type | Labels |
|-------------|-------------|-------------|
|water_cubic_used*|Gauge|id=\<meter-id\> <br/> name=\<meter-name\> <br/> meter=\<meter-brand\> <br/> timestamp=\<meter-value-timestamp\> <br/> prefix=\<meter-prefix-serial-number\> <br/> serial_number=\<meter-serial-number\> <br/> current_alarms=\<meter-current-alarms\> <br/> previous_alarms=\<meter-previous-alarms\>| 
|heat_kwh_consumption|Gauge|id=\<meter-id\> <br/> name=\<meter-name\> <br/> meter=\<meter-brand\> <br/> timestamp=\<meter-value-timestamp\>| 
|heat_cubic_consumption|Gauge|id=\<meter-id\> <br/> name=\<meter-name\> <br/> meter=\<meter-brand\> <br/> timestamp=\<meter-value-timestamp\>| 
|heat_kwh_flow|Gauge|id=\<meter-id\> <br/> name=\<meter-name\> <br/> meter=\<meter-brand\> <br/> timestamp=\<meter-value-timestamp\>| 
|heat_cubic_flow|Gauge|id=\<meter-id\> <br/> name=\<meter-name\> <br/> meter=\<meter-brand\> <br/> timestamp=\<meter-value-timestamp\>| 
|heat_temperature_flow|Gauge|id=\<meter-id\> <br/> name=\<meter-name\> <br/> meter=\<meter-brand\> <br/> timestamp=\<meter-value-timestamp\>| 
|heat_temperature_return|Gauge|id=\<meter-id\> <br/> name=\<meter-name\> <br/> meter=\<meter-brand\> <br/> timestamp=\<meter-value-timestamp\>| 

## Build
```bash
docker build --tag wmbusmeters-prometheus-metric:local . -f Dockerfile
# For multi plateform 
# docker buildx build --push --platform linux/arm/v7,linux/arm64,linux/amd64 --tag medinvention/wmbusmeters-prometheus-metric:0.0.1 . -f Dockerfile
```

**water metric label was set up from izar frame deconding, maybe not working with another brand of water meter*

### References

- https://github.com/wmbusmeters/wmbusmeters
- https://github.com/fstab/grok_exporter
- https://wmbusmeters.github.io/wmbusmeters-wiki/PROMETHEUS.html