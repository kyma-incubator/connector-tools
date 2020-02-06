---
title: Installation
type: Details
---

The Connector allows you to integrate an application by registering its APIs.  The integration uses the [extensibility features](https://help.sap.com/viewer/0815bc232f5140bba54a58ab15c82e99/Current/en-US/9ed15aa6eac34b948693955da0c90174.html) available in the SAP C/4HANA cockpit based on concepts and solutions from the open-source project "Kyma", so that you can easily develop Side-by-Side extensions. 

> **NOTE**: The integration is possible in the **preview** mode only. Do not use it in production scenarios.

### Set up the Connector 

1. Navigate to the **Runtimes** view under **Extensibility** in SAP C/4HANA cockpit and click on the cluster that you will be using
2. Register a new system providing an appropriate name (e.g. 'mysystem')
3. Click **Kyma Console**.
4. In the Console, choose the Namespace, go to **Catalog** and select **[Preview] oAuth Proxy - Connector**.
5. Click **Add once** to install and configure the Connector:

    | Parameter              | Description                                                                     |
    | ---------------------- | ------------------------------------------------------------------------------- |
    | **Application Name**   | The name of the registered system from step 2.                                  |
    | **Cluster Domain**     | The domain of the Kyma cluster. For example, `demo.cluster.extend.cx.cloud.sap` |
    | **Headers Parameters** | Headers to be sent in the oauth request                                         |
    | **System URL**         | The URL of your APIs. For example `https://xcore1.fgvms.com`                    |
    | **Form Parameters**    | The request body form parameters to be sent in the oauth request                |
    |                        |