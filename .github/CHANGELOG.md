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

