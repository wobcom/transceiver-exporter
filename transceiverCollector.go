package main

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gitlab.com/wobcom/ethtool"
	"gitlab.com/wobcom/ethtool/eeprom"
	"net"
	"strconv"
	"strings"
)

var (
	driverDesc              *prometheus.Desc
	driverVersionDesc       *prometheus.Desc
	firmwareVersionDesc     *prometheus.Desc
	busInfoDesc             *prometheus.Desc
	expansionRomVersionDesc *prometheus.Desc

	interfaceFeatureActiveDesc    *prometheus.Desc
	interfaceFeatureAvailableDesc *prometheus.Desc

	identifierDesc                            *prometheus.Desc
	encodingDesc                              *prometheus.Desc
	powerClassDesc                            *prometheus.Desc
	powerClassWattageDesc                     *prometheus.Desc
	signalingRateDesc                         *prometheus.Desc
	supportedLinkLengthsDesc                  *prometheus.Desc
	vendorNameDesc                            *prometheus.Desc
	vendorPNDesc                              *prometheus.Desc
	vendorRevDesc                             *prometheus.Desc
	vendorSNDesc                              *prometheus.Desc
	vendorOUIDesc                             *prometheus.Desc
	dateCodeDesc                              *prometheus.Desc
	wavelengthDesc                            *prometheus.Desc
	moduleSupportsMonitoringDesc              *prometheus.Desc
	moduleTemperatureDesc                     *prometheus.Desc
	moduleTemperatureThresholdsSupportedDesc  *prometheus.Desc
	moduleTemperatureHighAlarmThresholdDesc   *prometheus.Desc
	moduleTemperatureHighWarningThresholdDesc *prometheus.Desc
	moduleTemperatureLowAlarmThresholdDesc    *prometheus.Desc
	moduleTemperatureLowWarningThresholdDesc  *prometheus.Desc
	moduleVoltageDesc                         *prometheus.Desc
	moduleVoltageThresholdsSupportedDesc      *prometheus.Desc
	moduleVoltageHighAlarmThresholdDesc       *prometheus.Desc
	moduleVoltageHighWarningThresholdDesc     *prometheus.Desc
	moduleVoltageLowAlarmThresholdDesc        *prometheus.Desc
	moduleVoltageLowWarningThresholdDesc      *prometheus.Desc

	/* Laser monitoring information */
	laserSupportsMonitoringDesc *prometheus.Desc

	laserBiasDesc                     *prometheus.Desc
	laserBiasThresholdsSupportedDesc  *prometheus.Desc
	laserBiasHighAlarmThresholdDesc   *prometheus.Desc
	laserBiasHighWarningThresholdDesc *prometheus.Desc
	laserBiasLowAlarmThresholdDesc    *prometheus.Desc
	laserBiasLowWarningThresholdDesc  *prometheus.Desc

	laserTxPowerDesc                     *prometheus.Desc
	laserTxPowerThresholdsSupportedDesc  *prometheus.Desc
	laserTxPowerHighAlarmThresholdDesc   *prometheus.Desc
	laserTxPowerHighWarningThresholdDesc *prometheus.Desc
	laserTxPowerLowAlarmThresholdDesc    *prometheus.Desc
	laserTxPowerLowWarningThresholdDesc  *prometheus.Desc

	laserRxPowerDesc                     *prometheus.Desc
	laserRxPowerThresholdsSupportedDesc  *prometheus.Desc
	laserRxPowerHighAlarmThresholdDesc   *prometheus.Desc
	laserRxPowerHighWarningThresholdDesc *prometheus.Desc
	laserRxPowerLowAlarmThresholdDesc    *prometheus.Desc
	laserRxPowerLowWarningThresholdDesc  *prometheus.Desc
)

type TransceiverCollector struct{}

type MeasurementDesc struct {
	ValueDesc                 *prometheus.Desc
	ThresholdsSupportedDesc   *prometheus.Desc
	ThresholdsHighAlarmDesc   *prometheus.Desc
	ThresholdsHighWarningDesc *prometheus.Desc
	ThresholdsLowAlarmDesc    *prometheus.Desc
	ThresholdsLowWarningDesc  *prometheus.Desc
}

func init() {
	interfaceLabels := []string{"interface"}

	driverDesc = prometheus.NewDesc(prefix+"driver_name_info", "Driver name", []string{"interface", "driver_name"}, nil)
	driverVersionDesc = prometheus.NewDesc(prefix+"driver_version_info", "Driver version", []string{"interface", "driver_version"}, nil)
	firmwareVersionDesc = prometheus.NewDesc(prefix+"firmware_version_info", "Firmware version", []string{"interface", "firmware_version"}, nil)
	busInfoDesc = prometheus.NewDesc(prefix+"bus_info", "Bus information", []string{"interface", "bus_information"}, nil)
	expansionRomVersionDesc = prometheus.NewDesc(prefix+"expansion_rom_version_info", "Expansion ROM Version", []string{"interface", "expansion_rom_version"}, nil)

	interfaceFeatureActiveDesc = prometheus.NewDesc(prefix+"interface_feature_active", "Interfaces features as reported by interface driver. 1 if active.", []string{"interface", "feature_name"}, nil)
	interfaceFeatureAvailableDesc = prometheus.NewDesc(prefix+"interface_feature_available", "Interfaces features as reported by interface driver. 1 if available.", []string{"interface", "feature_name"}, nil)

	identifierDesc = prometheus.NewDesc(prefix+"identifier_info", "Type of transceiver information", []string{"interface", "identifier"}, nil)
	encodingDesc = prometheus.NewDesc(prefix+"encoding_info", "Transceiver encoding information", []string{"interface", "encoding"}, nil)
	powerClassDesc = prometheus.NewDesc(prefix+"powerclass_info", "Highest power class supported by the transceiver", interfaceLabels, nil)
	powerClassWattageDesc = prometheus.NewDesc(prefix+"powerclass_watts", "Maximum wattage supported by the transceivers power class", interfaceLabels, nil)
	signalingRateDesc = prometheus.NewDesc(prefix+"signalingrate_bauds_per_second", "Signaling rate in bauds per second supported by the transceiver", interfaceLabels, nil)
	supportedLinkLengthsDesc = prometheus.NewDesc(prefix+"supported_link_length_meter", "Maximum supported link length for different media in meters", []string{"interface", "media"}, nil)
	vendorNameDesc = prometheus.NewDesc(prefix+"vendor_name_info", "Vendor name", []string{"interface", "vendor_name"}, nil)
	vendorPNDesc = prometheus.NewDesc(prefix+"vendor_part_number_info", "Vendor part number", []string{"interface", "vendor_part_number"}, nil)
	vendorRevDesc = prometheus.NewDesc(prefix+"vendor_revision_info", "Vendor revision", []string{"interface", "vendor_revision"}, nil)
	vendorSNDesc = prometheus.NewDesc(prefix+"vendor_serial_number_info", "Vendor serial number", []string{"interface", "vendor_serial_number"}, nil)
	vendorOUIDesc = prometheus.NewDesc(prefix+"vendor_oui_info", "Vendor IEE company ID", []string{"interface", "vendor_oui"}, nil)
	dateCodeDesc = prometheus.NewDesc(prefix+"date_code_unix_time", "Vendor supplied date code exported as unix epoch", interfaceLabels, nil)
	wavelengthDesc = prometheus.NewDesc(prefix+"wavelength_nanometer", "Wavelength in nanometers", interfaceLabels, nil)
	moduleSupportsMonitoringDesc = prometheus.NewDesc(prefix+"module_supports_monitoring_bool", "1 if the module supports real time monitoring", interfaceLabels, nil)

	moduleTemperatureDesc = prometheus.NewDesc(prefix+"module_temperature_degrees_celsius", "Module temperature in degrees celsius", interfaceLabels, nil)
	moduleTemperatureThresholdsSupportedDesc = prometheus.NewDesc(prefix+"module_temperature_supports_thresholds_bool", "1 if thresholds for module temperature are supported", interfaceLabels, nil)
	moduleTemperatureHighAlarmThresholdDesc = prometheus.NewDesc(prefix+"module_temperature_high_alarm_threshold_degrees_celsius", "High alarm threshold for the module temperature in degrees celsius", interfaceLabels, nil)
	moduleTemperatureHighWarningThresholdDesc = prometheus.NewDesc(prefix+"module_temperature_high_warning_threshold_degrees_celsius", "High warning threshold for the module temperature in degrees celsius", interfaceLabels, nil)
	moduleTemperatureLowAlarmThresholdDesc = prometheus.NewDesc(prefix+"module_temperature_low_alarm_threshold_degrees_celsius", "Low alarm threshold for the module temperature in degrees celsius", interfaceLabels, nil)
	moduleTemperatureLowWarningThresholdDesc = prometheus.NewDesc(prefix+"module_temperature_low_warning_threshold_degrees_celsius", "Low warning threshold for the module temperature in degrees celsius", interfaceLabels, nil)

	moduleVoltageDesc = prometheus.NewDesc(prefix+"module_voltage_volts", "Module supply voltage in Volts", interfaceLabels, nil)
	moduleVoltageThresholdsSupportedDesc = prometheus.NewDesc(prefix+"module_voltage_supports_thresholds_bool", "1 if thresholds for modue voltage are supported", interfaceLabels, nil)
	moduleVoltageHighAlarmThresholdDesc = prometheus.NewDesc(prefix+"module_voltage_high_alarm_threshold_voltage", "High alarm threshold for the module voltage in volts", interfaceLabels, nil)
	moduleVoltageHighWarningThresholdDesc = prometheus.NewDesc(prefix+"module_voltage_high_warning_threshold_voltage", "High warning threshold for the module voltage in volts", interfaceLabels, nil)
	moduleVoltageLowAlarmThresholdDesc = prometheus.NewDesc(prefix+"module_voltage_low_alarm_threshold_voltage", "Low alarm threshold for the module voltage in volts", interfaceLabels, nil)
	moduleVoltageLowWarningThresholdDesc = prometheus.NewDesc(prefix+"module_voltage_low_warning_threshold_voltage", "Low warning threshold for the module voltage in volts", interfaceLabels, nil)

	/* Laser monitoring information */
	laserLabels := []string{"interface", "laser_index"}
	laserSupportsMonitoringDesc = prometheus.NewDesc(prefix+"laser_supports_monitoring_bool", "1 if the laser supports real time monitoring", laserLabels, nil)
	laserBiasDesc = prometheus.NewDesc(prefix+"laser_bias_current_milliamperes", "Laser bias current in in milliamperes", laserLabels, nil)
	laserBiasThresholdsSupportedDesc = prometheus.NewDesc(prefix+"laser_bias_current_supports_thresholds_bool", "1 if thresholds for the laser bias current are supported", laserLabels, nil)
	laserBiasHighAlarmThresholdDesc = prometheus.NewDesc(prefix+"laser_bias_current_high_alarm_threshold_milliamperes", "High alarm threshold for the laser bias current in milliamperes", laserLabels, nil)
	laserBiasHighWarningThresholdDesc = prometheus.NewDesc(prefix+"laser_bias_current_high_warning_threshold_milliamperes", "High warning threshold for the laser bias current in milliamperes", laserLabels, nil)
	laserBiasLowAlarmThresholdDesc = prometheus.NewDesc(prefix+"laser_bias_current_low_alarm_threshold_milliamperes", "Low alarm threshold for the laser bias current in milliamperes", laserLabels, nil)
	laserBiasLowWarningThresholdDesc = prometheus.NewDesc(prefix+"laser_bias_current_low_warning_threshold_milliamperes", "Low warning threshold for the laser bias current in milliamperes", laserLabels, nil)

	laserTxPowerDesc = prometheus.NewDesc(prefix+"laser_tx_power_milliwatts", "Laser tx power in milliwatts", laserLabels, nil)
	laserTxPowerThresholdsSupportedDesc = prometheus.NewDesc(prefix+"laser_tx_power_supports_thresholds_bool", "1 if thresholds for the laser tx power are supported", laserLabels, nil)
	laserTxPowerHighAlarmThresholdDesc = prometheus.NewDesc(prefix+"laser_tx_power_high_alarm_threshold_milliwatts", "High alarm threshold for the laser tx power in milliwatts", laserLabels, nil)
	laserTxPowerHighWarningThresholdDesc = prometheus.NewDesc(prefix+"laser_tx_power_high_warning_threshold_milliwatts", "High warning threshold for the laser tx power in milliwatts", laserLabels, nil)
	laserTxPowerLowAlarmThresholdDesc = prometheus.NewDesc(prefix+"laser_tx_power_low_alarm_threshold_milliwatts", "Low alarm threshold for the laser tx power in milliwatts", laserLabels, nil)
	laserTxPowerLowWarningThresholdDesc = prometheus.NewDesc(prefix+"laser_tx_power_low_warning_threshold_milliwatts", "Low warning threshold for the laser tx power in milliwatts", laserLabels, nil)

	laserRxPowerDesc = prometheus.NewDesc(prefix+"laser_rx_power_milliwatts", "Laser rx power in milliwatts", laserLabels, nil)
	laserRxPowerThresholdsSupportedDesc = prometheus.NewDesc(prefix+"laser_rx_power_supports_thresholds_bool", "1 if thresholds for the laser rx power are supported", laserLabels, nil)
	laserRxPowerHighAlarmThresholdDesc = prometheus.NewDesc(prefix+"laser_rx_power_high_alarm_threshold_milliwatts", "High alarm threshold for the laser rx power in milliwatts", laserLabels, nil)
	laserRxPowerHighWarningThresholdDesc = prometheus.NewDesc(prefix+"laser_rx_power_high_warning_threshold_milliwatts", "High warning threshold for the laser rx power in milliwatts", laserLabels, nil)
	laserRxPowerLowAlarmThresholdDesc = prometheus.NewDesc(prefix+"laser_rx_power_low_alarm_threshold_milliwatts", "Low alarm threshold for the laser rx power in milliwatts", laserLabels, nil)
	laserRxPowerLowWarningThresholdDesc = prometheus.NewDesc(prefix+"laser_rx_power_low_warning_threshold_milliwatts", "Low warning threshold for the laser rx power in milliwatts", laserLabels, nil)
}

func NewTransceiverCollector() *TransceiverCollector {
	return &TransceiverCollector{}
}

func (t *TransceiverCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- driverDesc
	ch <- driverVersionDesc
	ch <- firmwareVersionDesc
	ch <- busInfoDesc
	ch <- expansionRomVersionDesc

	ch <- identifierDesc
	ch <- encodingDesc
	ch <- powerClassDesc
	ch <- powerClassWattageDesc
	ch <- signalingRateDesc
	ch <- supportedLinkLengthsDesc
	ch <- vendorNameDesc
	ch <- vendorPNDesc
	ch <- vendorRevDesc
	ch <- vendorSNDesc
	ch <- vendorOUIDesc
	ch <- dateCodeDesc
	ch <- wavelengthDesc
	ch <- moduleSupportsMonitoringDesc
	ch <- moduleTemperatureDesc
	ch <- moduleTemperatureThresholdsSupportedDesc
	ch <- moduleTemperatureHighAlarmThresholdDesc
	ch <- moduleTemperatureHighWarningThresholdDesc
	ch <- moduleTemperatureLowAlarmThresholdDesc
	ch <- moduleTemperatureLowWarningThresholdDesc
	ch <- moduleVoltageDesc
	ch <- moduleVoltageThresholdsSupportedDesc
	ch <- moduleVoltageHighAlarmThresholdDesc
	ch <- moduleVoltageHighWarningThresholdDesc
	ch <- moduleVoltageLowAlarmThresholdDesc
	ch <- moduleVoltageLowWarningThresholdDesc

	ch <- laserSupportsMonitoringDesc

	ch <- laserBiasDesc
	ch <- laserBiasThresholdsSupportedDesc
	ch <- laserBiasHighAlarmThresholdDesc
	ch <- laserBiasHighWarningThresholdDesc
	ch <- laserBiasLowAlarmThresholdDesc
	ch <- laserBiasLowWarningThresholdDesc

	ch <- laserTxPowerDesc
	ch <- laserTxPowerThresholdsSupportedDesc
	ch <- laserTxPowerHighAlarmThresholdDesc
	ch <- laserTxPowerHighWarningThresholdDesc
	ch <- laserTxPowerLowAlarmThresholdDesc
	ch <- laserTxPowerLowWarningThresholdDesc

	ch <- laserRxPowerDesc
	ch <- laserRxPowerThresholdsSupportedDesc
	ch <- laserRxPowerHighAlarmThresholdDesc
	ch <- laserRxPowerHighWarningThresholdDesc
	ch <- laserRxPowerLowAlarmThresholdDesc
	ch <- laserRxPowerLowWarningThresholdDesc
}

// blacklistedIfaceNames enumerates the system's network interfaces and filters out command line blacklisted ones as well as loopback interfaces
func getMonitoredInterfaces() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return []string{}, errors.Wrapf(err, "Could not enumerate system's interfaces")
	}

	ifaceNames := []string{}
	blacklistedIfaceNames := strings.Split(*excludeInterfaces, ",")
	for index, blacklistedIfaceName := range blacklistedIfaceNames {
		blacklistedIfaceNames[index] = strings.Trim(blacklistedIfaceName, " ")
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback > 0 {
			continue
		}
		if contains(blacklistedIfaceNames, iface.Name) {
			continue
		}
		ifaceNames = append(ifaceNames, iface.Name)
	}
	return ifaceNames, nil
}

func (t *TransceiverCollector) Collect(ch chan<- prometheus.Metric) {
	ifaceNames, err := getMonitoredInterfaces()
	if err != nil {
		log.Error(err.Error())
		return
	}
	tool, err := ethtool.NewEthtool()
	if err != nil {
		log.Errorf("Could not instanciate ethtool: %v", err)
		return
	}

	for _, ifaceName := range ifaceNames {
		iface, err := tool.NewInterface(ifaceName, true)
		if err != nil {
			log.Errorf("Error fetching information for interface %s: %v", ifaceName, err)
			// continue
		}
		if iface != nil {
			exportMetricsForInterface(iface, ch)
		}
	}
}

func exportMetricsForInterface(iface *ethtool.Interface, ch chan<- prometheus.Metric) {
	if *collectInterfaceFeatures {
		features, err := iface.GetFeatures()
		if err == nil {
			for name, status := range features {
				ch <- prometheus.MustNewConstMetric(interfaceFeatureAvailableDesc, prometheus.GaugeValue, boolToFloat64(status.Available), iface.Name, name)
				ch <- prometheus.MustNewConstMetric(interfaceFeatureActiveDesc, prometheus.GaugeValue, boolToFloat64(status.Active), iface.Name, name)
			}
		}
	}
	if iface.DriverInfo != nil {
		exportDriverInfoMetricsForInterface(iface.Name, iface.DriverInfo, ch)
	}
	if iface.Eeprom != nil {
		exportEEPROMMetricsForInterface(iface.Name, iface.Eeprom, ch)
	}
}

func exportDriverInfoMetricsForInterface(ifaceName string, driverInfo *ethtool.DriverInfo, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(driverDesc, prometheus.GaugeValue, 1, ifaceName, driverInfo.DriverName)
	ch <- prometheus.MustNewConstMetric(driverVersionDesc, prometheus.GaugeValue, 1, ifaceName, driverInfo.DriverVersion)
	ch <- prometheus.MustNewConstMetric(firmwareVersionDesc, prometheus.GaugeValue, 1, ifaceName, driverInfo.FirmwareVersion)
	ch <- prometheus.MustNewConstMetric(busInfoDesc, prometheus.GaugeValue, 1, ifaceName, driverInfo.BusInfo)
	ch <- prometheus.MustNewConstMetric(expansionRomVersionDesc, prometheus.GaugeValue, 1, ifaceName, driverInfo.ExpansionRomVersion)
}

func exportEEPROMMetricsForInterface(ifaceName string, rom eeprom.EEPROM, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(identifierDesc, prometheus.GaugeValue, 1, ifaceName, rom.GetIdentifier().String())
	ch <- prometheus.MustNewConstMetric(encodingDesc, prometheus.GaugeValue, 1, ifaceName, rom.GetEncoding().String())
	ch <- prometheus.MustNewConstMetric(powerClassDesc, prometheus.GaugeValue, float64(byte(rom.GetPowerClass())), ifaceName)
	ch <- prometheus.MustNewConstMetric(powerClassWattageDesc, prometheus.GaugeValue, rom.GetPowerClass().GetMaxPower(), ifaceName)
	ch <- prometheus.MustNewConstMetric(signalingRateDesc, prometheus.GaugeValue, rom.GetSignalingRate(), ifaceName)
	for mediaName, supportedLength := range rom.GetSupportedLinkLengths() {
		ch <- prometheus.MustNewConstMetric(supportedLinkLengthsDesc, prometheus.GaugeValue, supportedLength, ifaceName, mediaName)
	}
	ch <- prometheus.MustNewConstMetric(vendorNameDesc, prometheus.GaugeValue, 1, ifaceName, rom.GetVendorName())
	ch <- prometheus.MustNewConstMetric(vendorPNDesc, prometheus.GaugeValue, 1, ifaceName, rom.GetVendorPN())
	ch <- prometheus.MustNewConstMetric(vendorRevDesc, prometheus.GaugeValue, 1, ifaceName, rom.GetVendorRev())
	ch <- prometheus.MustNewConstMetric(vendorSNDesc, prometheus.GaugeValue, 1, ifaceName, rom.GetVendorSN())
	ch <- prometheus.MustNewConstMetric(vendorOUIDesc, prometheus.GaugeValue, 1, ifaceName, rom.GetVendorOUI().String())
	ch <- prometheus.MustNewConstMetric(dateCodeDesc, prometheus.GaugeValue, float64(rom.GetDateCode().Unix()), ifaceName)
	ch <- prometheus.MustNewConstMetric(wavelengthDesc, prometheus.GaugeValue, rom.GetWavelength(), ifaceName)
	ch <- prometheus.MustNewConstMetric(moduleSupportsMonitoringDesc, prometheus.GaugeValue, boolToFloat64(rom.SupportsMonitoring()), ifaceName)

	if rom.SupportsMonitoring() {
		temperature, err := rom.GetModuleTemperature()
		if err == nil {
			exportMeasurement([]string{ifaceName}, temperature, &MeasurementDesc{
				moduleTemperatureDesc,
				moduleTemperatureThresholdsSupportedDesc,
				moduleTemperatureHighAlarmThresholdDesc,
				moduleTemperatureHighWarningThresholdDesc,
				moduleTemperatureLowAlarmThresholdDesc,
				moduleTemperatureLowWarningThresholdDesc,
			}, ch)
		}
		voltage, err := rom.GetModuleVoltage()
		if err == nil {
			exportMeasurement([]string{ifaceName}, voltage, &MeasurementDesc{
				moduleVoltageDesc,
				moduleVoltageThresholdsSupportedDesc,
				moduleVoltageHighAlarmThresholdDesc,
				moduleVoltageHighWarningThresholdDesc,
				moduleVoltageLowAlarmThresholdDesc,
				moduleVoltageLowWarningThresholdDesc,
			}, ch)
		}
		for index, laser := range rom.GetLasers() {
			if !laser.SupportsMonitoring() {
				continue
			}
			laserLabels := []string{ifaceName, strconv.Itoa(index)}

			bias, err := laser.GetBias()
			if err == nil {
				exportMeasurement(laserLabels, bias, &MeasurementDesc{
					laserBiasDesc,
					laserBiasThresholdsSupportedDesc,
					laserBiasHighAlarmThresholdDesc,
					laserBiasHighWarningThresholdDesc,
					laserBiasLowAlarmThresholdDesc,
					laserBiasLowWarningThresholdDesc,
				}, ch)
			}
			txPower, err := laser.GetTxPower()
			if err == nil {
				exportMeasurement(laserLabels, txPower, &MeasurementDesc{
					laserTxPowerDesc,
					laserTxPowerThresholdsSupportedDesc,
					laserTxPowerHighAlarmThresholdDesc,
					laserTxPowerHighWarningThresholdDesc,
					laserTxPowerLowAlarmThresholdDesc,
					laserTxPowerLowWarningThresholdDesc,
				}, ch)
			}
			rxPower, err := laser.GetRxPower()
			if err == nil {
				exportMeasurement(laserLabels, rxPower, &MeasurementDesc{
					laserRxPowerDesc,
					laserRxPowerThresholdsSupportedDesc,
					laserRxPowerHighAlarmThresholdDesc,
					laserRxPowerHighWarningThresholdDesc,
					laserRxPowerLowAlarmThresholdDesc,
					laserRxPowerLowWarningThresholdDesc,
				}, ch)
			}
		}
	}
}

func exportMeasurement(labels []string, measurement eeprom.Measurement, measurementDesc *MeasurementDesc, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(measurementDesc.ValueDesc, prometheus.GaugeValue, measurement.GetValue(), labels...)
	thresholdsSupported := measurement.SupportsThresholds()
	ch <- prometheus.MustNewConstMetric(measurementDesc.ThresholdsSupportedDesc, prometheus.GaugeValue, boolToFloat64(thresholdsSupported), labels...)
	if thresholdsSupported {
		thresholds, err := measurement.GetAlarmThresholds()
		if err != nil {
			return
		}
		ch <- prometheus.MustNewConstMetric(measurementDesc.ThresholdsHighAlarmDesc, prometheus.GaugeValue, thresholds.GetHighAlarm(), labels...)
		ch <- prometheus.MustNewConstMetric(measurementDesc.ThresholdsHighWarningDesc, prometheus.GaugeValue, thresholds.GetHighWarning(), labels...)
		ch <- prometheus.MustNewConstMetric(measurementDesc.ThresholdsLowAlarmDesc, prometheus.GaugeValue, thresholds.GetLowAlarm(), labels...)
		ch <- prometheus.MustNewConstMetric(measurementDesc.ThresholdsLowWarningDesc, prometheus.GaugeValue, thresholds.GetLowWarning(), labels...)
	}
}
