# IEEE802Sniffer
> Proyecto académico para la captura, análisis y evaluación de tramas LAN y WLAN conforme a los estándares IEEE 802.3 y IEEE 802.11 a/b/g/n/ac/ax/ah.

## 📌 Descripción

Este programa fue desarrollado en **Go** y permite capturar, interpretar y evaluar tramas de red en entornos **LAN (Ethernet)** y **WLAN (Wi-Fi)**. Se identifican los diferentes tipos de tramas (datos, control, gestión) y se extraen campos clave según los estándares definidos.

También se analizan los **campos de seguridad y calidad de servicio (QoS)** presentes en las tramas IEEE 802.11, contextualizándolos dentro de un escenario de red WLAN real o simulado.

## 🧰 Funcionalidades principales

- 📡 Captura de tramas Ethernet (IEEE 802.3).
- 📶 Captura de tramas WLAN (IEEE 802.11 a/b/g/n/ac/ax/ah).
- 🧩 Identificación de tramas de datos, control y gestión.
- 🔐 Análisis de campos de **seguridad** (WEP, WPA, WPA2, etc.).
- 🚦 Evaluación de campos de **QoS** (Wi-Fi Multimedia - WMM).
