variable "count" { default = 2 }
variable "region" { default = "australia-southeast1" }
variable "region_zone" { default = "australia-southeast1-c" }
variable "dns_zone" { default = "ranchermeetup.baillie.cloud." }
variable "public_key_path" { default = "~/.ssh/id_ed25519.pub" }
variable "rancher_token" { default = "" }

provider "google" {
    credentials = "${file("secrets/account.json")}"
    project     = "martinbaillie-rancher-meetup"
    region      = "${var.region}"
}

resource "google_compute_network" "rancher-meetup" {
  name = "rancher-meetup-network"
}

resource "google_compute_subnetwork" "rancher-meetup" {
  name          = "rancher-meetup-subnetwork"
  network       = "${google_compute_network.rancher-meetup.name}"
  ip_cidr_range = "10.240.0.0/24"
}

resource "google_compute_firewall" "rancher-meetup-allow-internal" {
  name          = "rancher-meetup-firewall-internal"
  network       = "${google_compute_network.rancher-meetup.name}"
  source_ranges = ["10.240.0.0/24", "10.200.0.0/16"]
  allow {
    protocol = "tcp"
  }
  allow {
    protocol = "udp"
  }
  allow {
    protocol = "icmp"
  }
}

resource "google_compute_firewall" "rancher-meetup-allow-external" {
  name          = "rancher-meetup-firewall-external"
  network       = "${google_compute_network.rancher-meetup.name}"
  source_ranges = ["0.0.0.0/0"]
  allow {
    protocol    = "tcp"
    ports       = ["22", "80", "443", "8080"]
  }
  allow {
    protocol    = "udp"
    ports       = ["22"]
  }
  allow {
    protocol    = "icmp"
  }
}

resource "google_compute_instance" "rancher-meetup" {
    count = "${var.count}"
    machine_type = "n1-standard-2"
    zone = "${var.region_zone}"
    name = "rancher-meetup-${count.index + 1}"
    tags = ["rancher-meetup-cluster","rancher-meetup-${count.index + 1}"]

    disk {
        image = "rancheros"
    }

    network_interface {
        subnetwork   = "${google_compute_subnetwork.rancher-meetup.name}"
        address      = "10.240.0.${count.index+10}"
        access_config { }
    }

    metadata {
        "ssh-keys" = "rancher:${file("${var.public_key_path}")}"
        "google-username" = "rancher"
        "user-data" = <<CLOUDCONFIG
#cloud-config
hostname: rancher-meetup-${count.index + 1}
locale: en_AU.UTF-8
write_files:
  - path: /etc/rc.local
    permissions: "0755"
    owner: root
    content: |
        #!/bin/bash

        # Wait for Docker and t'internet
        wait-for-docker
        while [ -e $${OK} ]; do
            ping google.com -c 4 && OK=true && break; sleep 3;
        done
        
        # Prep for LetsEncrypt
        touch /acme.json
        chmod 600 /acme.json

        # Some meetup demo pre-caching
        docker pull martinbaillie/rancher-cowsay-goproverb-ui
        docker pull martinbaillie/rancher-cowsay-goproverb-api

        docker pull martinbaillie/grafana-conf
        docker pull martinbaillie/prometheus-conf
        docker pull martinbaillie/prometheus-alertmanager-conf

        docker pull martinbaillie/prom-aggregation-gateway
        docker pull martinbaillie/cadvisor
        docker pull martinbaillie/apache-bench

        docker pull prom/prometheus
        docker pull prom/alertmanager
        docker pull prom/node-exporter

        docker pull grafana/grafana

        docker pull containous/traefik:experimental
        docker pull infinityworks/prometheus-rancher-exporter
        docker pull gmehta3/prometheus-container-monitor

        # Rancher agent
        docker run --rm --privileged \
            -e CATTLE_HOST_LABELS='rancher-meetup=${count.index + 1}' \
            -e CATTLE_AGENT_IP=$(wget -qO - --header 'Metadata-Flavor: Google' http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip) \
            -v /var/run/docker.sock:/var/run/docker.sock \
            -v /var/lib/rancher:/var/lib/rancher rancher/agent:v1.2.2 \
            https://try-api.rancher.com/v1/scripts/${var.rancher_token}
CLOUDCONFIG
    }
}

# NOTE: can't do GCE network load balancing with RancherOS as need the guest
# environment (specifically the `ip_forwarding_daemon`). `rancher/#260` on
# their Github is asking for this. All it needs is a new system service in ros.
#
# Using CloudDNS instead of network load balancing.

#resource "google_compute_address" "rancher-meetup" {
  #name = "rancher-meetup-address"
#}

#resource "google_compute_http_health_check" "rancher-meetup" {
  #name         = "rancher-meetup-http-health-check"
  #request_path = "/ping"
  #port         = "8080"
#}

#resource "google_compute_target_pool" "rancher-meetup" {
  #name          = "rancher-meetup-target-pool"
  #instances = [ "${formatlist("%s/%s", google_compute_instance.rancher-meetup.*.zone, google_compute_instance.rancher-meetup.*.name)}" ]
  #health_checks = [ "${google_compute_http_health_check.rancher-meetup.name}" ]
#}

#resource "google_compute_forwarding_rule" "rancher-meetup-80" {
  #name       = "rancher-meetup-forwarding-rule-80"
  #target     = "${google_compute_target_pool.rancher-meetup.self_link}"
  #port_range = "80"
  #ip_address = "${google_compute_address.rancher-meetup.address}"
#}

#resource "google_compute_forwarding_rule" "rancher-meetup-443" {
  #name       = "rancher-meetup-forwarding-rule-443"
  #target     = "${google_compute_target_pool.rancher-meetup.self_link}"
  #port_range = "443"
  #ip_address = "${google_compute_address.rancher-meetup.address}"
#}

#resource "google_compute_forwarding_rule" "rancher-meetup-8080" {
  #name       = "rancher-meetup-forwarding-rule-8080"
  #target     = "${google_compute_target_pool.rancher-meetup.self_link}"
  #port_range = "8080"
  #ip_address = "${google_compute_address.rancher-meetup.address}"
#}

#resource "google_compute_firewall" "rancher-meetup-allow-lb-ping" {
  #name          = "rancher-meetup-firewall-lb-ping"
  #network       = "${google_compute_network.rancher-meetup.name}"
  #source_ranges = ["130.211.0.0/22"]
  #allow {
    #protocol    = "tcp"
    #ports       = ["8080"]
  #}
#}

#output "rancher_meetup_lb" {
  #value = "${google_compute_address.rancher-meetup.address}"
#}

# NOTE: using GCE CloudDNS instead of network load balancing
#resource "google_dns_managed_zone" "rancher-meetup" {
  #name     = "rancher-meetup-dns-managed-zone"
  #dns_name = "ranchermeetup.baillie.cloud."
#}

resource "google_dns_record_set" "rancher-meetup" {
  name = "${var.dns_zone}"
  type = "A"
  ttl  = 60
  managed_zone = "rancher-meetup-dns-managed-zone"
  rrdatas = ["${google_compute_instance.rancher-meetup.*.network_interface.0.access_config.0.assigned_nat_ip}"]
}

resource "google_dns_record_set" "rancher-meetup-wildcard" {
  name = "*.${var.dns_zone}"
  type = "A"
  ttl  = 60
  managed_zone = "rancher-meetup-dns-managed-zone"
  rrdatas = ["${google_compute_instance.rancher-meetup.*.network_interface.0.access_config.0.assigned_nat_ip}"]
}

resource "google_dns_record_set" "rancher-meetup-instances" {
  count = "${var.count}"
  name = "rancher-meetup-${count.index + 1}.${var.dns_zone}"
  type = "A"
  ttl = 60
  managed_zone = "rancher-meetup-dns-managed-zone"
  rrdatas = ["${element(google_compute_instance.rancher-meetup.*.network_interface.0.access_config.0.assigned_nat_ip, count.index)}"]
}

output "rancher_meetup_dns_zone" {
  value = "${google_dns_record_set.rancher-meetup-wildcard.name}"
}

output "rancher_meetup_instance_dns" {
  value = "${join(" ", google_dns_record_set.rancher-meetup-instances.*.name)}"
}

output "rancher_meetup_instance_ips" {
  value = "${join(" ", google_compute_instance.rancher-meetup.*.network_interface.0.access_config.0.assigned_nat_ip)}"
}
