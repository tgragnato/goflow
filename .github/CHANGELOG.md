## [0.5.0](https://github.com/tgragnato/goflow/compare/v0.4.0...v0.5.0) (2024-12-04)


### Features

* **protobuf:** add and populate the sampler_hostname field ([e3ef7a1](https://github.com/tgragnato/goflow/commit/e3ef7a1e8eb3ff8d02b6c499feeb769d97572722))
* **protobuf:** introduce additional as path fields ([e735e5a](https://github.com/tgragnato/goflow/commit/e735e5a9764c52def45b05d54d76a9417f41f43a))


### Bug Fixes

* **producer:** replace direct string handling with json.Marshal ([3c18550](https://github.com/tgragnato/goflow/commit/3c185506325ee1b4f9024f3f58cc8580eb80665f))

## [0.4.0](https://github.com/tgragnato/goflow/compare/v0.3.0...v0.4.0) (2024-11-19)


### Features

* populate the ip and asn fields ([b2af7e4](https://github.com/tgragnato/goflow/commit/b2af7e4be306197add35c1d9e23dd1d05fa000c3))


### Bug Fixes

* **gosec:** integer overflow conversion uint -> uint32 ([52b89ab](https://github.com/tgragnato/goflow/commit/52b89abea762495487cb54d43f6a3208b760503a))
* **protobuf:** remove SrcAddrIp and DstAddrIp ([cb09310](https://github.com/tgragnato/goflow/commit/cb093100a2de215978835845a4733b9a9ec45a35))

## [0.3.0](https://github.com/tgragnato/goflow/compare/v0.2.0...v0.3.0) (2024-11-11)


### Features

* add support to the lms_target_index target field ([8c12dc3](https://github.com/tgragnato/goflow/commit/8c12dc3e3f3b9d2e66c67d28f61933f059484195))
* **build:** fix upload-artifact GitHub action ([#353](https://github.com/tgragnato/goflow/issues/353)) ([844e345](https://github.com/tgragnato/goflow/commit/844e3457555e796adafd22a26d84580049fb4f7d))
* **geoip:** merge pb_ext into pb and embed geoip ([1badcc1](https://github.com/tgragnato/goflow/commit/1badcc1c13bc3f8cde1587201b58f3d8d8ca352e))

## [0.2.0](https://github.com/tgragnato/goflow/compare/v0.1.0...v0.2.0) (2024-09-09)


### Features

* **producer:** refactor encapsulation and configuration ([#342](https://github.com/tgragnato/goflow/issues/342)) ([ec96100](https://github.com/tgragnato/goflow/commit/ec961004f02c405ffd644103de4b3aa50d11c3fb))
* **syslog:** remove the timestamp from the log message ([e500ea0](https://github.com/tgragnato/goflow/commit/e500ea0e1c4a73b88f4a37427ba1097a5e0d5176))


### Bug Fixes

* **producer:** mpls labels in NetFlow ([#345](https://github.com/tgragnato/goflow/issues/345)) ([b01b1aa](https://github.com/tgragnato/goflow/commit/b01b1aacd2e450f9db47c8086c22253e5cfc9bc1))
* **producer:** should parse with empty config ([#347](https://github.com/tgragnato/goflow/issues/347)) ([127c1a3](https://github.com/tgragnato/goflow/commit/127c1a3272125fa41c22af7e93ac14700711793d))
* udp receiver not passing errors ([#323](https://github.com/tgragnato/goflow/issues/323)) ([f8e113a](https://github.com/tgragnato/goflow/commit/f8e113a38469c40c0d1824e534ea06eee193d59e))

## [0.1.0](https://github.com/tgragnato/goflow/compare/v0.0.1...v0.1.0) (2024-08-20)


### Features

* **enricher:** move the protobuf definition in pb_ext ([9c6d5c2](https://github.com/tgragnato/goflow/commit/9c6d5c2907d58316d883b6cfed89ae2ebfdbb40f))
* **transport:** replace kafka with syslog ([3a441b4](https://github.com/tgragnato/goflow/commit/3a441b4db4fb09f0893193dded4d1290515bc88a))


### Bug Fixes

* **codeql:** conversion of uint64 to int32 without upper bound checks ([c9ca2bc](https://github.com/tgragnato/goflow/commit/c9ca2bc5eabfe396c475e1f7a49c8c32c5cdee40))

