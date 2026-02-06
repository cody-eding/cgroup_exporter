## 1.0.0 / 2026-02-06

* Initial release based on [treydock/cgroup_exporter v1.0.1](https://github.com/treydock/cgroup_exporter/tree/v1.0.1).
* Added support for tracking job step and task needed by Jobstats.
* Added UID metric needed for Jobstats.
* Removed support for external process tracking.
* Removed support for cgroupv1.
* Fixed calculation of memory RSS for cgroupv2.
* Disabled exporter metric by default.