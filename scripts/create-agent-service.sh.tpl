#!/bin/bash

#
# SPDX-License-Identifier: Apache-2.0
#  SPDX-FileCopyrightText: Huawei Inc.
#

wget -O /tmp/create-agent-service.sh https://raw.githubusercontent.com/eclipse-xpanse/xpanse-agent/refs/heads/main/scripts/create-agent-service.sh

chmod 775 /tmp/create-agent-service.sh

/tmp/create-agent-service.sh --serviceId ${serviceId} --resourceName ${resourceName} --pollingInterval ${pollingInterval} --xpanseApiEndpoint ${xpanseApiEndpoint} --agentVersion ${agentVersion}