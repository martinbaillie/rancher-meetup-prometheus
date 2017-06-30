# Sydney Rancher Meetup June '17
### About 

> Inclusive monitoring of a microservice architecture with Rancher and the Prometheus ecosystem – (Martin Baillie @ IAG) 
>
> In this talk and demo, we will look at how you can use the Prometheus ecosystem to monitor, alert, graph and forecast on the internal state of a microservice architecture.  We’ll focus on how you can scale this solution out as your platform grows using techniques like zero-conf discoverability and telemetry, and also how you can better encourage an inclusive monitoring culture within your engineering teams.

- [Meetup](https://www.meetup.com/Rancher-Sydney/events/240631132)
- [Deck](https://speakerdeck.com/martinbaillie/syd-rancher-meetup-june-17)

### Demo components
#### Deck
- [Marp](https://github.com/yhatt/marp)

#### Ingress 
- [Traefik](https://traefik.io)  with the new Rancher metadata service provider and Let's Encrypt

#### Monitoring stack
Official upstream images with [confd](https://github.com/martinbaillie/confd) sidekicks to provide custom configuration discovered from the Rancher metadata service. Grafana dashboards, Prometheus alerts stored as code beside the images.
- [Prometheus](https://prometheus.io)
- [Grafana](https://grafana.com)
- [Alertmanager](https://github.com/prometheus/alertmanager)
- [Node Exporter](https://github.com/prometheus/node_exporter) (Host metrics)
- [CAdvisor](https://github.com/google/cadvisor) (Docker metrics)
- [Rancher Exporter](https://github.com/infinityworksltd/prometheus-rancher-exporter) (Rancher metrics)
- [Container Monitor](https://github.com/ibrokethecloud/prometheus-container-monitor) (IPSec metrics)

#### Rancher Cowsay GoProverb API
Rob Pike's talk on [Go Proverbs](https://www.google.com.au/url?sa=t&rct=j&q=&esrc=s&source=web&cd=1&cad=rja&uact=8&ved=0ahUKEwjRxqfHuOTUAhVILpQKHb6uD38QtwIIKDAA&url=https%3A%2F%2Fwww.youtube.com%2Fwatch%3Fv%3DPAAkCSZUG1c&usg=AFQjCNF7fuC6vP6xQ2W_s2R1K_A-Q6LeDA)
- [Go-kit](https://github.com/go-kit/kit)
- [Prometheus Golang Client](https://github.com/prometheus/client_golang)

#### Rancher Cowsay GoProverb UI
- [Weaveworks Prom JS](https://github.com/weaveworks/promjs)
- [Weaveworks Prom Aggregation Gateway](https://github.com/weaveworks/prom-aggregation-gateway)

#### Load generator
- [Apache Bench](https://httpd.apache.org/docs/2.4/programs/ab.html)
