## 1.0.0 / 2026-02-06

* Initial release based on [treydock/cgroup_exporter v1.0.1](https://github.com/treydock/cgroup_exporter/tree/v1.0.1).
* Updated to go v1.24.0
* Replaced legacy go-kit logging with log/slog.
* Added support for tracking job step and task needed by Jobstats.
* Added UID metric needed for Jobstats.
* Removed support for gathering process information.
* Removed support for cgroupv1.
* Fixed calculation of memory RSS for cgroupv2.
* Disabled exporter metrics by default.