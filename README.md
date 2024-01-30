# Prometheus transceiver exporter
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

This tool provides diagnostics for pluggable transceivers (SFP, SFP+, QSFP, etc.) by leveraging the [ethtool library](https://github.com/wobcom/go-ethtool).
You can use this tool with [CumulusLinux](https://cumulusnetworks.com/products/cumulus-linux/).

## Command line options
You might want to set `-collector.interface-features.enable` to false, because it may result in huge amounts of timeseries (especially on many port switches).

```
Usage of ./transceiver-exporter:
  -collector.interface-features.enable
        Collect interface features (default true)
  -collector.optical-power-in-dbm
        Report optical powers in dBm instead of mW (default false -> mW)
  -exclude.interfaces string
        Comma seperated list of interfaces to exclude
  -exclude.interfaces-regex string
        Regexp of interfaces to exclude
  -exclude.interfaces-down
        Don't report on interfaces being management DOWN
  -include.interfaces string
        Comma seperated list of interfaces to include
  -include.interfaces-regex string
        Regexp of interfaces to include
  -version
        Print version and exit
  -web.listen-address string
        Address to listen on (default "[::]:9458")
  -web.telemetry-path string
        Path under which to expose metrics (default "/metrics")
```

## Exported metrics

Note: Transmit / Receive power (and thresholds) are exported as milliwatts just as they are read from the module. If you wish to have decibel milliwatts, you'll have to do the conversion `10 * math.Log10(value_in_milliwatts)`. Please also note that, this might result `-Inf` for a value of 0 which might cause trouble with software / standards (e.g. JSON) not fully implementing the IEE754 floating point standard.
Starting in version 1.1.0 we added the runtime option `-collector.optical-power-in-dbm` to enable conversion to dBm in the exporter.

* `transceiver_exporter_date_code_unix_time`: Vendor supplied date code exported as unix epoch
* `transceiver_exporter_driver_name_info`: Driver name
* `transceiver_exporter_driver_version_info`: Driver version
* `transceiver_exporter_encoding_info`: Transceiver encoding information
* `transceiver_exporter_expansion_rom_version_info`: Expansion ROM Version
* `transceiver_exporter_firmware_version_info`: Firmware version
* `transceiver_exporter_interface_feature_active`: Interfaces features as reported by interface driver. 1 if active.
* `transceiver_exporter_interface_feature_available`: Interfaces features as reported by interface driver. 1 if available.
* `transceiver_exporter_identifier_info`: Type of transceiver information
* `transceiver_exporter_laser_bias_current_high_alarm_threshold_milliamperes`: High alarm threshold for the laser bias current in milliamperes
* `transceiver_exporter_laser_bias_current_high_warning_threshold_milliamperes`: High warning threshold for the laser bias current in milliamperes
* `transceiver_exporter_laser_bias_current_low_alarm_threshold_milliamperes`: Low alarm threshold for the laser bias current in milliamperes
* `transceiver_exporter_laser_bias_current_low_warning_threshold_milliamperes`: Low warning threshold for the laser bias current in milliamperes
* `transceiver_exporter_laser_bias_current_milliamperes`: Laser bias current in in milliamperes
* `transceiver_exporter_laser_bias_current_supports_thresholds_bool`: 1 if thresholds for the laser bias current are supported
* `transceiver_exporter_laser_rx_power_high_alarm_threshold_milliwatts`: High alarm threshold for the laser rx power in milliwatts
* `transceiver_exporter_laser_rx_power_high_warning_threshold_milliwatts`: High warning threshold for the laser rx power in milliwatts
* `transceiver_exporter_laser_rx_power_low_alarm_threshold_milliwatts`: Low alarm threshold for the laser rx power in milliwatts
* `transceiver_exporter_laser_rx_power_low_warning_threshold_milliwatts`: Low warning threshold for the laser rx power in milliwatts
* `transceiver_exporter_laser_rx_power_milliwatts`: Laser rx power in milliwatts
* `transceiver_exporter_laser_rx_power_supports_thresholds_bool`: 1 if thresholds for the laser rx power are supported
* `transceiver_exporter_laser_tx_power_high_alarm_threshold_milliwatts`: High alarm threshold for the laser tx power in milliwatts
* `transceiver_exporter_laser_tx_power_high_warning_threshold_milliwatts`: High warning threshold for the laser tx power in milliwatts
* `transceiver_exporter_laser_tx_power_low_alarm_threshold_milliwatts`: Low alarm threshold for the laser tx power in milliwatts
* `transceiver_exporter_laser_tx_power_low_warning_threshold_milliwatts`: Low warning threshold for the laser tx power in milliwatts
* `transceiver_exporter_laser_tx_power_milliwatts`: Laser tx power in milliwatts
* `transceiver_exporter_laser_tx_power_supports_thresholds_bool`: 1 if thresholds for the laser tx power are supported
* `transceiver_exporter_module_supports_monitoring_bool`: 1 if the module supports real time monitoring
* `transceiver_exporter_module_temperature_degrees_celsius`: Module temperature in degrees celsius
* `transceiver_exporter_module_temperature_high_alarm_threshold_degrees_celsius`: High alarm threshold for the module temperature in degrees celsius
* `transceiver_exporter_module_temperature_high_warning_threshold_degrees_celsius`: High warning threshold for the module temperature in degrees celsius
* `transceiver_exporter_module_temperature_low_alarm_threshold_degrees_celsius`: Low alarm threshold for the module temperature in degrees celsius
* `transceiver_exporter_module_temperature_low_warning_threshold_degrees_celsius`: Low warning threshold for the module temperature in degrees celsius
* `transceiver_exporter_module_temperature_supports_thresholds_bool`: 1 if thresholds for module temperature are supported
* `transceiver_exporter_module_voltage_high_alarm_threshold_voltage`: High alarm threshold for the module voltage in volts
* `transceiver_exporter_module_voltage_high_warning_threshold_voltage`: High warning threshold for the module voltage in volts
* `transceiver_exporter_module_voltage_low_alarm_threshold_voltage`: Low alarm threshold for the module voltage in volts
* `transceiver_exporter_module_voltage_low_warning_threshold_voltage`: Low warning threshold for the module voltage in volts
* `transceiver_exporter_module_voltage_supports_thresholds_bool`: 1 if thresholds for modue voltage are supported
* `transceiver_exporter_module_voltage_volts`: Module supply voltage in Volts
* `transceiver_exporter_powerclass_info`: Highest power class supported by the transceiver
* `transceiver_exporter_powerclass_watts`: Maximum wattage supported by the transceivers power class
* `transceiver_exporter_signalingrate_bauds_per_second`: Signaling rate in bauds per second supported by the transceiver
* `transceiver_exporter_supported_link_length_meter`: Maximum supported link length for different media in meters
* `transceiver_exporter_vendor_name_info`: Vendor name
* `transceiver_exporter_vendor_oui_info`: Vendor IEE company ID
* `transceiver_exporter_vendor_part_number_info`: Vendor part number
* `transceiver_exporter_vendor_revision_info`: Vendor revision
* `transceiver_exporter_vendor_serial_number_info`: Vendor serial number
* `transceiver_exporter_wavelength_nanometer`: Wavelength in nanometers

## Maintainer
* @vidister

## Authors
* @fluepke
* @BarbarossaTM
* @vidister
