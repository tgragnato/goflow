# GoFlow

[![Go](https://github.com/tgragnato/goflow/actions/workflows/go.yml/badge.svg)](https://github.com/tgragnato/goflow/actions/workflows/go.yml)
[![CodeQL](https://github.com/tgragnato/goflow/actions/workflows/codeql.yml/badge.svg)](https://github.com/tgragnato/goflow/actions/workflows/codeql.yml)

This application is a NetFlow/IPFIX/sFlow collector in Go.

It gathers network information (IP, interfaces, routers) from different flow protocols, serializes it in a common format.

You will want to use GoFlow if:
* You receive a decent amount of network samples and need horizontal scalability
* Have protocol diversity and need a consistent format
* Require raw samples and build aggregation and custom enrichment

This software is the entry point of a pipeline. The storage, transport, enrichment, graphing, alerting are not provided.

![GoFlow System diagram](/docs/diagram.png)

## Origins

This work is a fork of a previous [open-source GoFlow code](https://github.com/cloudflare/goflow) built and used at Cloudflare.
It lives in its own GitHub organization to be maintained more easily.

Among the differences with the original code:
The serializer and transport options have been revamped to make this program more user-friendly and target new use-cases like logging providers.
Minimal changes in the decoding libraries.

## Modularity

In order to enable load-balancing and optimizations, the GoFlow library has a `decoder` which converts the payload of a flow packet into a structure.

The `producer` converts the samples into another format.
Out of the box, this repository provides a protobuf producer (`pb/flow.pb`) and a raw producer.
In the case of the protobuf producer, the records in a single flow packet are extracted and made in their own protobuf. Custom mapping allows to add new fields without rebuilding the proto.

The `format` directory offers various utilities to format a message. It calls specific functions to marshal as JSON or text for instance.

The `transport` provides different way of processing the message. Either sending it via syslog or send it to a file (or stdout).

GoFlow is a wrapper of all the functions and chains them.

You can build your own collector using this base and replace parts:
* Use different transport (e.g: RabbitMQ instead of Syslog)
* Convert to another format (e.g: Cap'n Proto, Avro, instead of protobuf)
* Decode different samples (e.g: not only IP networks, add MPLS)
* Different metrics system (e.g: [OpenTelemetry](https://opentelemetry.io/))

### Protocol difference

The sampling protocols have distinct features:

**sFlow** is a stateless protocol which sends the full header of a packet with router information (interfaces, destination AS) while **NetFlow/IPFIX** rely on templates that contain fields (e.g: source IPv6).

The sampling rate in NetFlow/IPFIX is provided by **Option Data Sets**. This is why it can take a few minutes for the packets to be decoded until all the templates are received (**Option Template** and **Data Template**).

Both of these protocols bundle multiple samples (**Data Set** in NetFlow/IPFIX and **Flow Sample** in sFlow) in one packet.

The advantages of using an abstract network flow format, such as protobuf, is it enables summing over the protocols (e.g: per ASN or per port, rather than per (ASN, router) and (port, router)).

To read more about the protocols and how they are mapped inside, check out [page](/docs/protocols.md)

### Features of GoFlow

Collection:
* NetFlow v5
* IPFIX/NetFlow v9 (sampling rate provided by the Option Data Set)
* sFlow v5

(adding NetFlow v1,7,8 is being evaluated)

Production:
* Convert to protobuf or json
* Prints to the console/file

Monitoring via Prometheus metrics

## Get started

To read about agents that samples network traffic, check this [page](/docs/agents.md).

To set up the collector, download the latest release corresponding to your OS and run the following command (the binaries have a suffix with the version):

```bash
$ ./goflow
```

By default, this command will launch an sFlow collector on port `:6343` and a NetFlowV9/IPFIX collector on port `:2055`.

By default, the samples received will be printed in JSON format on the stdout.

```json
{
    "type": "SFLOW_5",
    "time_received_ns": 1681583295157626000,
    "sequence_num": 2999,
    "sampling_rate": 100,
    "sampler_address": "192.168.0.1",
    "time_flow_start_ns": 1681583295157626000,
    "time_flow_end_ns": 1681583295157626000,
    "bytes": 1500,
    "packets": 1,
    "src_addr": "fd01::1",
    "dst_addr": "fd01::2",
    "etype": "IPv6",
    "proto": "TCP",
    "src_port": 443,
    "dst_port": 50001
}
```

If you are using a log integration (e.g: Loki with Promtail, Splunk, Fluentd, Google Cloud Logs, etc.), just send the output into a file.

```bash
$ ./goflow -transport.file /var/logs/goflow.log
```

To enable syslog, use the following arguments:

```bash
$ ./goflow -transport=syslog -transport.syslog.address 'localhost:514' -transport.syslog.protocol 'udp'
```

By default, the collector will listen for IPFIX/NetFlow V9 on port 2055 and sFlow on port 6343.
To change the sockets binding, you can set the `-listen` argument and a URI for each protocol (`netflow`, `sflow` and `nfl` as scheme) separated by a comma.
For instance, to create 4 parallel sockets of sFlow and one of NetFlow V5, you can use:

```bash
$ ./goflow -listen 'sflow://:6343?count=4,nfl://:2055'
```

More information about workers and resource usage is avaialble on the [Performance page](/docs/performance.md).

### Docker

You can also run directly with a container:
```
$ sudo docker run -p 6343:6343/udp -p 2055:2055/udp -ti tgragnato/goflow:latest [...]
```

### Mapping extra fields

In the case of exotic template fields or extra payload not supported by GoFlow of out the box, it is possible to pass a mapping file using `-mapping mapping.yaml`.
A [sample file](mapping.yaml) is available in the main directory.

For instance, certain devices producing IPFIX use `ingressPhysicalInterface` (id: 252) and do not use `ingressInterface` (id: 10). Using the following you can have the interface mapped in the InIf protobuf field without changing the code.

```yaml
ipfix:
  mapping:
    - field: 252
      destination: in_if
    - field: 253
      destination: out_if
```

### Output format considerations

The JSON format is advised only when consuming a small amount of data directly.
For bigger workloads, the protobuf output format provides a binary representation and is preferred.
It can also be extended with enrichment as long as the user keep the same IDs.

If you want to develop applications, build `pb/flow.proto` into the language you want:
When adding custom fields, picking a field ID ≥ 1000 is suggested.

Check the docs for more information about [compiling protobuf](/docs/protobuf.md).

## Security notes and assumptions

By default, the buffer for UDP is 9000 bytes.
Protections were added to avoid DOS on sFlow since the various length fields are 32 bits.
There are assumptions on how many records and list items a sample can have (eg: AS-Path).

### Contributions

This project welcomes pull-requests, whether it's documentation, instrumentation (metrics), internals (protocol libraries),
integration (new CLI feature) or else!
Just make sure to check for the use-cases via an issue.

This software would not exist without the testing and commits from its users and [contributors](/docs/contributors.md).

## License

Licensed under the BSD-3 License.
